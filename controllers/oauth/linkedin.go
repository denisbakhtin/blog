package oauth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/shared"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

//LinkedinLogin handles /linkedin_login route
func LinkedinLogin(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)
	session.Values["oauth_redirect"] = r.Referer()
	session.Save(r, w)
	http.Redirect(w, r, inConfig().AuthCodeURL(oauthState), 303)
}

//LinkedinCallback handles /linkedin_callback route
func LinkedinCallback(w http.ResponseWriter, r *http.Request) {
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

	token, err := inConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	client := inConfig().Client(oauth2.NoContext, token)
	req, err := http.NewRequest(
		"GET",
		"https://api.linkedin.com/v1/people/~:(email-address,first-name,last-name,id,headline)?format=json",
		nil,
	)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}
	req.Header.Set("Bearer", token.AccessToken)
	response, err := client.Do(req)

	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	defer response.Body.Close()
	str, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(500)
		tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
		return
	}

	var inuser struct {
		ID        string
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Headline  string
		Email     string `json:"emailAddress"`
	}

	err = json.Unmarshal(str, &inuser)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(500)
		tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
		return
	}

	redirectURL := session.Values["oauth_redirect"]
	delete(session.Values, "oauth_redirect")
	session.Values["oauth_email"] = inuser.Email
	session.Values["oauth_name"] = inuser.FirstName + " " + inuser.LastName
	session.Save(r, w)
	if url, ok := redirectURL.(string); ok {
		http.Redirect(w, r, url, 303)
	} else {
		http.Redirect(w, r, "/", 303)
	}
}

func inConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     shared.GetConfig().Oauth.Linkedin.ClientID,
		ClientSecret: shared.GetConfig().Oauth.Linkedin.ClientSecret,
		RedirectURL:  shared.GetConfig().Oauth.Linkedin.RedirectURL,
		Scopes:       []string{"r_basicprofile", "r_emailaddress"},
		Endpoint:     linkedin.Endpoint,
	}
}
