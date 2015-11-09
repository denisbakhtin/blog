package oauth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/system"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

//state should be regenerated per auth request
var (
	oauthState = "random_oauth_state_string"
)

//FacebookLogin handles /facebook_login route
func FacebookLogin(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)
	session.Values["oauth_redirect"] = r.Referer()
	session.Save(r, w)
	http.Redirect(w, r, fbConfig().AuthCodeURL(oauthState), 303)
}

//FacebookCallback handlers /facebook_callback route
func FacebookCallback(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	session := context.Get(r, "session").(*sessions.Session)
	state := r.FormValue("state")
	if oauthState != state {
		err := fmt.Errorf("Wrong state string: Expected %s, got %s. Please, try again", oauthState, state)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
		return
	}

	token, err := fbConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
		return
	}

	client := fbConfig().Client(oauth2.NoContext, token)
	response, err := client.Get(fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=name,email,birthday", token.AccessToken))

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
		return
	}

	defer response.Body.Close()
	str, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(500)
		tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
		return
	}

	var fbuser struct {
		ID       string
		Name     string
		Email    string
		Birthday string
	}

	err = json.Unmarshal([]byte(str), &fbuser)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(500)
		tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
		return
	}

	redirectURL := session.Values["oauth_redirect"]
	delete(session.Values, "oauth_redirect")
	session.Values["oauth_email"] = fbuser.Email
	session.Values["oauth_name"] = fbuser.Name
	session.Save(r, w)
	if url, ok := redirectURL.(string); ok {
		http.Redirect(w, r, url, 303)
	} else {
		http.Redirect(w, r, "/", 303)
	}
}

//facebook config
func fbConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     system.GetConfig().Oauth.Facebook.ClientID,
		ClientSecret: system.GetConfig().Oauth.Facebook.ClientSecret,
		RedirectURL:  system.GetConfig().Oauth.Facebook.RedirectURL,
		Scopes:       []string{"email", "user_about_me"},
		Endpoint:     facebook.Endpoint,
	}
}
