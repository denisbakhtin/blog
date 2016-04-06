package shared

import (
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
)

var (
	Store *sessions.FilesystemStore
)

func CreateSessionStore(secret string, domain string, secure bool) {
	Store = sessions.NewFilesystemStore("", []byte(secret))
	Store.Options = &sessions.Options{Domain: domain, Path: "/", Secure: secure, HttpOnly: true, MaxAge: 7 * 86400}
}

//DefaultData returns common to all pages template data
func DefaultData(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"ActiveUser":     context.Get(r, "user"),           //signed in models.User
		"Active":         "",                               //active uri shortening for menu item highlight
		"Title":          "",                               //page title:w
		"SignupEnabled":  context.Get(r, "signup_enabled"), //signup route is enabled (otherwise everyone can signup ;)
		csrf.TemplateTag: csrf.TemplateField(r),
	}
}

//ErrorData returns template data for error
func ErrorData(err error) map[string]interface{} {
	return map[string]interface{}{
		"Title": err.Error(),
		"Error": err.Error(),
	}
}

//Template returns parsed *html/template.Template
func Template(r *http.Request) *template.Template {
	return tmpl
}

//Session returns current session
func Session(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, "session") //ignore unrecoverable error if file storage has been removed from /tmp dir after server reboot. Instead check session == nil
	if session == nil {
		session, _ = Store.New(r, "session")
	}
	return session
}

//Atoi64 converts string to int64, returns 0 if error
func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

//Atob converts string to bool
func Atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
