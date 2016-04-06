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
	"golang.org/x/oauth2/vk"
)

//VkLogin handles /vk_login route
func VkLogin(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)
	session.Values["oauth_redirect"] = r.Referer()
	session.Save(r, w)
	http.Redirect(w, r, vkConfig().AuthCodeURL(oauthState), 303)
}

//VkCallback handles /vk_callback route
func VkCallback(w http.ResponseWriter, r *http.Request) {
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

	token, err := vkConfig().Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	email := token.Extra("email")
	userID := int64(token.Extra("user_id").(float64))
	//if you need to invoke vk api, create a client, make your requests
	client := vkConfig().Client(oauth2.NoContext, token)
	response, err := client.Get(fmt.Sprintf("https://api.vk.com/method/users.get?access_token=%s&user_id=%d", token.AccessToken, userID))
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(400)
		tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
		return
	}

	type Response struct {
		UID       int64
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	var vkuser struct {
		Response []Response
	}
	if err := json.Unmarshal(buf, &vkuser); err != nil {
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(500)
		tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
		return
	}

	redirectURL := session.Values["oauth_redirect"]
	delete(session.Values, "oauth_redirect")
	session.Values["oauth_email"] = email
	session.Values["oauth_name"] = vkuser.Response[0].FirstName + " " + vkuser.Response[0].LastName
	session.Save(r, w)
	if url, ok := redirectURL.(string); ok {
		http.Redirect(w, r, url, 303)
	} else {
		http.Redirect(w, r, "/", 303)
	}
}

func vkConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     shared.GetConfig().Oauth.Vk.ClientID,
		ClientSecret: shared.GetConfig().Oauth.Vk.ClientSecret,
		RedirectURL:  shared.GetConfig().Oauth.Vk.RedirectURL,
		Scopes:       []string{"email"},
		Endpoint:     vk.Endpoint,
	}
}
