package system

import (
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/nicksnyder/go-i18n/i18n"
	"html/template"
	"log"
	"net/http"
)

var (
	store *sessions.FilesystemStore
	//store *sessions.CookieStore
)

func createSession() {
	store = sessions.NewFilesystemStore("", []byte(config.SessionSecret))
	//store = sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options = &sessions.Options{Domain: config.Domain, Path: "/", Secure: config.Ssl, HttpOnly: true, MaxAge: 7 * 86400}
}

//SessionMiddleware creates gorilla session and stores it in context
func SessionMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		session, _ := store.Get(r, "session") //ignore unrecoverable error if file storage has been removed from /tmp dir after server reboot. Instead check session == nil
		if session == nil {
			var err error
			session, err = store.New(r, "session")
			if err != nil {
				log.Printf("ERROR: can't get session: %s", err)
				http.Error(w, err.Error(), 500)
				return //abort chain
			}
		}
		context.Set(r, "session", session)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//LocaleMiddleware stores current locale
func LocaleMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookieLang := ""
		if cookieL, _ := r.Cookie("lang"); cookieL != nil {
			cookieLang = cookieL.Value
		}
		acceptLang := r.Header.Get("Accept-Language")
		defaultLang := config.Language // known valid language
		T, lang, err := i18n.TfuncAndLanguage(cookieLang, acceptLang, defaultLang)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			log.Printf("ERROR: %s\n", err)
			return
		}
		context.Set(r, "T", T)
		context.Set(r, "language", lang.Tag)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//TemplateMiddleware stores parsed templates in context. Must be preceded by LocaleMiddleware
func TemplateMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t := tmpl.Funcs(map[string]interface{}{
			"T": context.Get(r, "T").(i18n.TranslateFunc), //translation func for current locale, see LocaleMiddleware
		})
		context.Set(r, "template", t)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//DataMiddleware inits common request data (active user, et al). Must be preceded by SessionMiddleware
func DataMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//set active user
		session := context.Get(r, "session").(*sessions.Session)
		if uID, ok := session.Values["user_id"]; ok {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				context.Set(r, "user", user)
			}
		}
		//enable signup link
		if config.SignupEnabled {
			context.Set(r, "signup_enabled", true)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//RestrictedMiddleware verifies presence on 'user' in context (which is set by DataMiddleware, if user has signed in
func RestrictedMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if user := context.Get(r, "user"); user != nil {
			//access granted
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(403)
			context.Get(r, "template").(*template.Template).Lookup("errors/403").Execute(w, nil)
			log.Printf("ERROR: unauthorized access to %s\n", r.RequestURI)
		}
	}
	return http.HandlerFunc(fn)
}
