package esa

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/markbates/goth"
)

// Session stores data during the auth process with esa.
type Session struct {
	AuthURL     string
	AccessToken string
}

var _ goth.Session = &Session{}

// GetAuthURL will return the URL set by calling the `BeginAuth` function on the esa provider.
func (s Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New(goth.NoAuthUrlErrorMessage)
	}
	return s.AuthURL, nil
}

// Authorize the session with esa and return the access token to be stored for future use.
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	token, err := p.config.Exchange(goth.ContextForClient(p.Client()), params.Get("code"))
	if err != nil {
		return "", err
	}

	if !token.Valid() {
		return "", errors.New("invalid token received from provider")
	}

	s.AccessToken = token.AccessToken
	return token.AccessToken, err
}

// Marshal the session into a string
func (s Session) Marshal() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s Session) String() string {
	return s.Marshal()
}

// CondVal convert string in string array
func CondVal(v string) []string {
	if v == "" {
		return nil
	}
	return []string{v}
}

// UnmarshalSession wil unmarshal a JSON string into a session.
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}
	err := json.NewDecoder(strings.NewReader(data)).Decode(s)
	return s, err
}
