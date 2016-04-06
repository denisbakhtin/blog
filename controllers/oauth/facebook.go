package oauth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/denisbakhtin/blog/shared"
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
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	token, err := fbConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	client := fbConfig().Client(oauth2.NoContext, token)
	response, err := client.Get(fmt.Sprintf("https://graph.facebook.com/me?access_token=%s&fields=name,email,birthday", token.AccessToken))

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
		tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
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

//PostOnFacebook creates new post on facebook page wall
func PostOnFacebook(link, message string) error {
	//see http://stackoverflow.com/questions/17197970/facebook-permanent-page-access-token
	//for info on obtaining upexpirable page access token
	//also https://developers.facebook.com/docs/graph-api/reference/v2.5/page/feed for api description

	token := &oauth2.Token{
		AccessToken: shared.GetConfig().Oauth.Facebook.Token, //page access token
	}
	client := fbConfig().Client(oauth2.NoContext, token)
	response, err := client.Post(
		fmt.Sprintf(
			"https://graph.facebook.com/v2.5/%s/feed?access_token=%s&link=%s&message=%s",
			shared.GetConfig().Oauth.Facebook.Page,
			token.AccessToken,
			url.QueryEscape(link),
			url.QueryEscape(message),
		),
		"application/json",
		nil,
	)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if response.StatusCode != 200 {
		err := fmt.Errorf("ERROR: while posting on facebook: %s\n", body)
		return err
	}
	return nil
}

//facebook config
func fbConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     shared.GetConfig().Oauth.Facebook.ClientID,
		ClientSecret: shared.GetConfig().Oauth.Facebook.ClientSecret,
		RedirectURL:  shared.GetConfig().Oauth.Facebook.RedirectURL,
		Scopes:       []string{"email", "user_about_me"},
		Endpoint:     facebook.Endpoint,
	}
}
