package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/mujhtech/s3ase/api/dto"
	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/database/store"
	errs "github.com/mujhtech/s3ase/errors"
)

type CreateOrUpdateDomainService struct {
	App           *models.App
	AppMemberRepo store.AppMemberRepository
	DomainRepo    store.DomainRepository
	User          *models.User
	Body          *dto.CreateOrUpdateDomainRequestDto
	CnameTarget   string
}

func (c *CreateOrUpdateDomainService) Run(ctx context.Context) (*models.Domain, error) {
	if err := requireAppOwner(ctx, c.AppMemberRepo, c.App.ID, c.User.ID); err != nil {
		return nil, err
	}

	domain, err := c.DomainRepo.FindDomainByAppID(ctx, c.App.ID)

	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, err
	}

	userDomain, err := normalizeDomain(c.Body.Domain)
	if err != nil {
		return nil, err
	}
	target, err := normalizeDomain(c.CnameTarget)
	if err != nil {
		return nil, fmt.Errorf("domain CNAME target is not configured: %w", err)
	}

	if domain == nil {
		verification := make([]byte, 24)
		if _, err := rand.Read(verification); err != nil {
			return nil, err
		}
		domain := &models.Domain{
			ID:          uuid.New().String(),
			AppID:       c.App.ID,
			Domain:      userDomain,
			CnameRecord: target,
			TxtRecord:   "s3ase-verification=" + base64.RawURLEncoding.EncodeToString(verification),
			Status:      models.DomainStatusPending,
			CreatedBy:   c.User.ID,
			Metadata:    map[string]interface{}{},
		}

		if err := c.DomainRepo.CreateDomain(ctx, domain); err != nil {
			return nil, err
		}
	} else {
		domainChanged := domain.Domain != userDomain
		domain.Domain = userDomain
		domain.CnameRecord = target
		domain.Status = models.DomainStatusPending
		if domainChanged {
			verification := make([]byte, 24)
			if _, err := rand.Read(verification); err != nil {
				return nil, err
			}
			domain.TxtRecord = "s3ase-verification=" + base64.RawURLEncoding.EncodeToString(verification)
		}

		if err := c.DomainRepo.UpdateDomain(ctx, domain); err != nil {
			return nil, err
		}

	}

	return domain, nil
}

func normalizeDomain(value string) (string, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Hostname() == "" || parsed.User != nil || parsed.Port() != "" {
		return "", fmt.Errorf("%w: invalid domain", errs.ErrInvalidInput)
	}
	host := strings.TrimSuffix(parsed.Hostname(), ".")
	if ip := net.ParseIP(host); ip != nil || !strings.Contains(host, ".") {
		return "", fmt.Errorf("%w: a public domain name is required", errs.ErrInvalidInput)
	}
	return host, nil
}

type DNSResolver interface {
	LookupCNAME(ctx context.Context, host string) (string, error)
	LookupTXT(ctx context.Context, host string) ([]string, error)
}

type VerifyDomainService struct {
	App           *models.App
	User          *models.User
	AppMemberRepo store.AppMemberRepository
	DomainRepo    store.DomainRepository
	Resolver      DNSResolver
}

func (s *VerifyDomainService) Run(ctx context.Context) (*models.Domain, error) {
	if err := requireAppOwner(ctx, s.AppMemberRepo, s.App.ID, s.User.ID); err != nil {
		return nil, err
	}
	domain, err := s.DomainRepo.FindDomainByAppID(ctx, s.App.ID)
	if err != nil {
		return nil, err
	}
	resolver := s.Resolver
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	cname, cnameErr := resolver.LookupCNAME(ctx, domain.Domain)
	txtValues, txtErr := resolver.LookupTXT(ctx, "_s3ase."+domain.Domain)
	cnameMatches := cnameErr == nil && strings.EqualFold(strings.TrimSuffix(cname, "."), strings.TrimSuffix(domain.CnameRecord, "."))
	txtMatches := false
	if txtErr == nil {
		for _, value := range txtValues {
			if value == domain.TxtRecord {
				txtMatches = true
				break
			}
		}
	}
	if cnameMatches && txtMatches {
		domain.Status = models.DomainStatusVerified
	} else {
		domain.Status = models.DomainStatusReview
	}
	if err := s.DomainRepo.UpdateDomain(ctx, domain); err != nil {
		return nil, err
	}
	return domain, nil
}
