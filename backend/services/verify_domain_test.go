package services

import (
	"context"
	"testing"

	"github.com/mujhtech/s3ase/database/models"
	"github.com/mujhtech/s3ase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type fakeDNSResolver struct {
	cname string
	txt   []string
}

func (f fakeDNSResolver) LookupCNAME(context.Context, string) (string, error) {
	return f.cname, nil
}

func (f fakeDNSResolver) LookupTXT(context.Context, string) ([]string, error) {
	return f.txt, nil
}

func TestVerifyDomainServiceRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	members := mocks.NewMockAppMemberRepository(ctrl)
	domains := mocks.NewMockDomainRepository(ctrl)
	members.EXPECT().FindAppMemberByAppIDAndUserId(gomock.Any(), "app-id", "owner-id").
		Return(&models.AppMember{Role: models.AppMemberRoleOwner}, nil)
	domain := &models.Domain{
		ID: "domain-id", AppID: "app-id", Domain: "cdn.example.com",
		CnameRecord: "files.s3ase.dev", TxtRecord: "s3ase-verification=token",
	}
	domains.EXPECT().FindDomainByAppID(gomock.Any(), "app-id").Return(domain, nil)
	domains.EXPECT().UpdateDomain(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, updated *models.Domain) error {
			require.Equal(t, models.DomainStatusVerified, updated.Status)
			return nil
		},
	)

	verified, err := (&VerifyDomainService{
		App: &models.App{ID: "app-id"}, User: &models.User{ID: "owner-id"},
		AppMemberRepo: members, DomainRepo: domains,
		Resolver: fakeDNSResolver{cname: "files.s3ase.dev.", txt: []string{"s3ase-verification=token"}},
	}).Run(context.Background())

	require.NoError(t, err)
	require.Equal(t, models.DomainStatusVerified, verified.Status)
}
