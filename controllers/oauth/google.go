package oauth

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/shared"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	goauth2 "google.golang.org/api/oauth2/v2"
)

//GoogleLogin handles /google_login route
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)
	session.Values["oauth_redirect"] = r.Referer()
	session.Save(r, w)
	http.Redirect(w, r, goConfig().AuthCodeURL(oauthState), 303)
}

//GoogleCallback handles /google_callback route
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	session := context.Get(r, "session").(*sessions.Session)

	state := r.FormValue("state")
	if oauthState != state {
		err := fmt.Errorf("Wrong state string: Expected %s, got %s. Please, try again", oauthState, state)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	token, err := goConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	client := goConfig().Client(oauth2.NoContext, token)
	service, err := goauth2.New(client)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}
	uService := goauth2.NewUserinfoService(service)
	gouser, err := uService.Get().Do()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	redirectURL := session.Values["oauth_redirect"]
	delete(session.Values, "oauth_redirect")
	session.Values["oauth_email"] = gouser.Email
	session.Values["oauth_name"] = gouser.GivenName + " " + gouser.FamilyName
	session.Save(r, w)
	if url, ok := redirectURL.(string); ok {
		http.Redirect(w, r, url, 303)
	} else {
		http.Redirect(w, r, "/", 303)
	}
}

//google config
func goConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     shared.GetConfig().Oauth.Google.ClientID,
		ClientSecret: shared.GetConfig().Oauth.Google.ClientSecret,
		RedirectURL:  shared.GetConfig().Oauth.Google.RedirectURL,
		Scopes:       []string{goauth2.PlusLoginScope, goauth2.PlusMeScope, goauth2.UserinfoEmailScope, goauth2.UserinfoProfileScope},
		Endpoint:     google.Endpoint,
	}
}
