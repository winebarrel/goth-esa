// Package esa implements the OAuth2 protocol for authenticating users through esa.
// This package can be used as a reference implementation of an OAuth2 provider for Goth.
package esa

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL         string = "https://api.esa.io/oauth/authorize"
	tokenURL        string = "https://api.esa.io/oauth/token"
	endpointProfile string = "https://api.esa.io/v1/user"
)

// Provider is the implementation of `goth.Provider` for accessing esa.
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// New creates a new esa provider and sets up important connection details.
// You should always call `esa.New` to get a new provider.  Never try to
// create one manually.
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "esa",
	}
	p.config = newConfig(p, scopes)
	return p
}

// Name is the name used to retrieve this provider later.
func (p *Provider) Name() string {
	return p.providerName
}

// SetName is to update the name of the provider (needed in case of multiple providers of 1 type)
func (p *Provider) SetName(name string) {
	p.providerName = name
}

func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug is a no-op for the esa package.
func (p *Provider) Debug(debug bool) {}

// BeginAuth asks esa for an authentication end-point.
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	return &Session{
		AuthURL: p.config.AuthCodeURL(state),
	}, nil
}

// FetchUser will go to esa and access basic information about the user.
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	sess := session.(*Session)
	user := goth.User{
		AccessToken: sess.AccessToken,
		Provider:    p.Name(),
	}

	if user.AccessToken == "" {
		// data is not yet retrieved since accessToken is still empty
		return user, fmt.Errorf("%s cannot get user information without accessToken", p.providerName)
	}

	req, err := http.NewRequest("GET", endpointProfile, nil)
	if err != nil {
		return user, err
	}

	req.Header.Set("Authorization", "Bearer "+sess.AccessToken)

	response, err := p.Client().Do(req)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("%s responded with a %d trying to fetch user information", p.providerName, response.StatusCode)
	}

	err = userFromReader(response.Body, &user)
	return user, err
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	if len(scopes) > 0 {
		c.Scopes = append(c.Scopes, scopes...)
	}
	return c
}

func userFromReader(r io.Reader, user *goth.User) error {
	u := struct {
		Email      string `json:"email"`
		Icon       string `json:"icon"`
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
	}{}
	err := json.NewDecoder(r).Decode(&u)
	if err != nil {
		return err
	}
	user.AvatarURL = u.Icon
	user.Email = u.Email
	user.Name = u.Name
	user.NickName = u.Name
	user.UserID = strconv.Itoa(u.ID)
	return nil
}

//RefreshToken refresh token is not provided by esa
func (p *Provider) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	return nil, errors.New("Refresh token is not provided by yammer")
}

//RefreshTokenAvailable refresh token is not provided by esa
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}
