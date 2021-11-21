package esa_test

import (
	"os"
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/goth-esa/esa"
)

func Test_New(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := provider()

	a.Equal(p.ClientKey, os.Getenv("ESA_KEY"))
	a.Equal(p.Secret, os.Getenv("ESA_SECRET"))
	a.Equal(p.CallbackURL, "/foo")
}

func Test_Implements_Provider(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	a.Implements((*goth.Provider)(nil), provider())
}

func Test_BeginAuth(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := provider()
	session, err := p.BeginAuth("test_state")
	s := session.(*esa.Session)
	a.NoError(err)
	a.Contains(s.AuthURL, "api.esa.io/oauth/authorize")
}

func Test_SessionFromJSON(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	p := provider()
	session, err := p.UnmarshalSession(`{"AuthURL":"https://api.esa.io/oauth/authorize","AccessToken":"1234567890"}`)
	a.NoError(err)

	s := session.(*esa.Session)
	a.Equal(s.AuthURL, "https://api.esa.io/oauth/authorize")
	a.Equal(s.AccessToken, "1234567890")
}

func provider() *esa.Provider {
	return esa.New(os.Getenv("ESA_KEY"), os.Getenv("ESA_SECRET"), "/foo")
}
