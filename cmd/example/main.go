package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/pat"
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

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	p := pat.New()

	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)

		if err != nil {
			fmt.Fprintln(res, err)
			return
		}

		sess, _ := store.Get(req, "mysqssion")
		sess.Values["user"] = user
		sess.Save(req, res)

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)

		sess, _ := store.Get(req, "mysqssion")
		delete(sess.Values, "user")
		sess.Save(req, res)

		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		sess, _ := store.Get(req, "mysqssion")
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
	log.Fatal(http.ListenAndServe("localhost:3000", p))
}
