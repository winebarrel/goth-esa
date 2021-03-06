package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/winebarrel/goth-esa/esa"
)

var indexTemplate = `
<p><a href="/auth/esa">Log in with esa</a></p>
`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}}</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>AvatarURL: {{.AvatarURL}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
`

func main() {
	goth.UseProviders(
		esa.New(os.Getenv("ESA_KEY"), os.Getenv("ESA_SECRET"), "http://localhost:3000/auth/esa/callback", "read"),
	)

	// NOTE: `_gothic_session` only exists to handle the OAuth2 state parameter handling
	// cf. https://github.com/markbates/goth/issues/181#issuecomment-590391070
	store := sessions.NewCookieStore([]byte("secret"))
	r := mux.NewRouter()

	r.Path("/auth/{provider}/callback").Methods("GET").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)

		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		sess, _ := store.Get(req, "mysession")
		sess.Values["user"] = user
		sess.Save(req, res)

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Path("/logout/{provider}").Methods("GET").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)

		sess, _ := store.Get(req, "mysession")
		delete(sess.Values, "user")
		sess.Save(req, res)

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	r.Path("/auth/{provider}").Methods("GET").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	r.Path("/").Methods("GET").HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		sess, _ := store.Get(req, "mysession")
		user := sess.Values["user"]

		if user == nil {
			t, _ := template.New("foo").Parse(indexTemplate)
			t.Execute(res, nil)
		} else {
			t, _ := template.New("foo").Parse(userTemplate)
			t.Execute(res, user)
		}
	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe("localhost:3000", r))
}
