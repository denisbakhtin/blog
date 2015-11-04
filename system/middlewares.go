package system

import (
	"github.com/denisbakhtin/blog/models"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"golang.org/x/net/context"
	"html/template"
	"log"
	"net/http"
)

var (
	store *sessions.FilesystemStore
)

//CtxHandler is practically http.Handler but with context param
type CtxHandler interface {
	ServeHTTPCtx(context.Context, http.ResponseWriter, *http.Request)
}

//CtxHandlerFunc is practically http.HandlerFunc but with context param
type CtxHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func createSession() {
	store = sessions.NewFilesystemStore("", []byte(config.SessionSecret))
	store.Options = &sessions.Options{HttpOnly: true, MaxAge: 7 * 86400} //Also set Secure: true if using SSL, you should though
}

//ServeHTTPCtx makes CtxHandlerFunc implement CtxHandler interface
func (h CtxHandlerFunc) ServeHTTPCtx(ctx context.Context, rw http.ResponseWriter, req *http.Request) {
	h(ctx, rw, req)
}

//SessionMiddleware creates gorilla session and stores it in context
func SessionMiddleware(next CtxHandler) CtxHandler {
	return CtxHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		defer gcontext.Clear(r) //one day all golang apps will rely on std context... dreams
		session, err := store.Get(r, "session")
		if err != nil {
			log.Printf("ERROR: can't get session: %s", err)
			return //abort chain
		}
		next.ServeHTTPCtx(context.WithValue(ctx, "session", session), w, r)
	})
}

//TemplateMiddleware stores parsed templates in context
func TemplateMiddleware(next CtxHandler) CtxHandler {
	return CtxHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		next.ServeHTTPCtx(context.WithValue(ctx, "template", tmpl), w, r)
	})
}

//DataMiddleware inits common request data (active user, et al). Must be preceded by SessionMiddleware
func DataMiddleware(next CtxHandler) CtxHandler {
	return CtxHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		session := ctx.Value("session").(*sessions.Session)
		if uID, ok := session.Values["user_id"]; ok {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				ctx = context.WithValue(ctx, "user", user)
			}
		}
		if config.SignupEnabled {
			ctx = context.WithValue(ctx, "signup_enabled", true)
		}

		next.ServeHTTPCtx(ctx, w, r)
	})
}

//RestrictedMiddleware verifies presence on 'user' in context (which is set by DataMiddleware, if user has signed in
func RestrictedMiddleware(next CtxHandler) CtxHandler {
	return CtxHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if user := ctx.Value("user"); user != nil {
			//access granted
			next.ServeHTTPCtx(ctx, w, r)
		} else {
			w.WriteHeader(403)
			ctx.Value("template").(*template.Template).Lookup("errors/403").Execute(w, nil)
			log.Printf("ERROR: unauthorized access to %s\n", r.RequestURI)
		}
	})
}
