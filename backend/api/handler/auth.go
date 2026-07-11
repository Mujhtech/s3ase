package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	jwtAuth "github.com/mujhtech/s3ase/auth"
	"github.com/mujhtech/s3ase/config"
	"github.com/mujhtech/s3ase/internal/pkg/auth"
	"github.com/mujhtech/s3ase/internal/pkg/response"
	"github.com/mujhtech/s3ase/services"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

const (
	AuthProviderKey = "auth_provider_key"
	AuthRedirectUrl = "redirect_url"
	AuthStateKey    = "auth_state_key"
)

func getProviderFromPath(r *http.Request) (string, error) {
	rawRef, err := pathParamOrError(r, AuthProviderKey)
	if err != nil {
		return "", err
	}

	return url.PathUnescape(rawRef)
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {

	//ctx := r.Context()

	provider, redirectUrl, err := getProvider(h.cfg, r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	state, err := auth.GenerateState(h.cfg, provider.Name(), redirectUrl)
	if err != nil {
		_ = response.InternalServerError(w, r, err)
		return
	}

	authUrl := provider.AuthCodeURL(string(state), oauth2.SetAuthURLParam("prompt", "consent"))

	// set state cookie
	createCookie(w, CookieOptions{
		Name:     AuthStateKey,
		State:    string(state),
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   300,
		Secure:   h.cfg.Server.SSL || strings.HasPrefix(h.cfg.Auth.RedirectUrl, "https://"),
		HttpOnly: true,
	})

	if err = response.Redirect(w, r, authUrl, http.StatusTemporaryRedirect, true); err != nil {
		log.Error().Err(err).Msg("failed to redirect")
		return
	}
}

func (h *Handler) AuthenticateCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	provider, _, err := getProvider(h.cfg, r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	authState, err := queryParamOrError(r, "state")

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	code, err := queryParamOrError(r, "code")

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	cookie, ok := getCookie(r, AuthStateKey)

	if !ok {
		_ = response.Unauthorized(w, r, fmt.Errorf("state cookie not found"))
		return
	}

	_, err = auth.VerifyState(h.cfg, cookie, authState)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	token, err := provider.GetOAuthToken(code)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	authUser, err := provider.GetUser(token)
	// create new user or link account

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	// create or link user
	createOrLinkUserService := &services.CreateOrLinkUserService{
		UserRepo: h.store.UserRepo,
		AuthUser: authUser,
	}

	user, err := createOrLinkUserService.Run(ctx)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	// create token
	tokenManager := jwtAuth.NewJWTAuth(h.cfg, h.store.UserRepo, h.store.TokenRepo, h.cache)

	_, authToken, err := tokenManager.CreateToken(ctx, user.ID, "bearer")

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	redirectTo, err := url.Parse(h.cfg.Auth.UIRedirectUrl)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	query := url.Values{}
	query.Add("token", authToken)
	redirectTo.RawQuery = query.Encode()

	// reset cookie
	createCookie(w, CookieOptions{
		Name:     AuthStateKey,
		State:    "",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
		Secure:   h.cfg.Server.SSL || strings.HasPrefix(h.cfg.Auth.RedirectUrl, "https://"),
		HttpOnly: true,
	})

	if err = response.Redirect(w, r, redirectTo.String(), http.StatusTemporaryRedirect, true); err != nil {
		log.Error().Err(err).Msg("failed to redirect")
		return
	}
}

func (h *Handler) AuthenticateCallbackPost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provider, redirectUrl, err := getProvider(h.cfg, r)

	if err != nil {
		_ = response.BadRequest(w, r, err)
		return
	}

	authState := queryParamOrDefault(r, "state", "")

	code := queryParamOrDefault(r, "code", "")

	var dst interface{}

	query := url.Values{}

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {

		query.Add("error", err.Error())

		if err = response.Redirect(w, r, fmt.Sprintf("%s?%s", redirectUrl, query.Encode()), http.StatusTemporaryRedirect, true); err != nil {
			log.Error().Err(err).Msg("failed to redirect")
			return
		}
	}

	log.Ctx(ctx).Info().Msgf("dst: %+v", dst)

	query.Add("state", authState)
	query.Add("code", code)

	if err = response.Redirect(w, r, fmt.Sprintf("/%s/callback?%s", provider.Name(), query.Encode()), http.StatusSeeOther, true); err != nil {
		log.Error().Err(err).Msg("failed to redirect")
		return
	}
}

func getProvider(cfg *config.Config, r *http.Request) (auth.AuthProvider, string, error) {
	authProvider, err := getProviderFromPath(r)
	if err != nil {
		return nil, "", err
	}

	redirectUrl := queryParamOrDefault(r, AuthRedirectUrl, "")

	if redirectUrl == "" {
		redirectUrl = fmt.Sprintf("%s/api/ui/auth/%s/callback", cfg.Auth.RedirectUrl, authProvider)
	}

	provider, err := auth.GetAuthProvider(cfg, authProvider, redirectUrl)

	if err != nil {
		return nil, "", err
	}

	return provider, redirectUrl, nil
}
