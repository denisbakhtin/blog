package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/blog/controllers"
	"github.com/denisbakhtin/blog/controllers/oauth"
	"github.com/denisbakhtin/blog/system"
	"github.com/gorilla/csrf"
)

//gorilla/csrf middleware
var CSRF func(http.Handler) http.Handler

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	migrate := flag.String("migrate", "skip", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	system.SetMode(mode)
	system.Init()
	system.RunMigrations(migrate)
	CSRF = csrf.Protect([]byte(system.GetConfig().CsrfSecret), csrf.Secure(system.GetConfig().Ssl), csrf.Path("/"), csrf.Domain(system.GetConfig().Domain))

	//Periodic tasks
	gocron.Every(1).Day().Do(system.CreateXMLSitemap)
	gocron.Start()

	http.Handle("/", Default(controllers.Home))
	if system.GetConfig().SignupEnabled {
		http.Handle("/signup", Default(controllers.SignUp))
	}
	http.Handle("/signin", Default(controllers.SignIn))
	http.Handle("/logout", Default(controllers.Logout))

	http.Handle("/pages/", Default(controllers.PageShow))
	http.Handle("/posts/", Default(controllers.PostShow))
	http.Handle("/tags/", Default(controllers.TagShow))
	http.Handle("/archives/", Default(controllers.ArchiveShow))
	http.Handle("/rss", Default(controllers.RssXML))
	http.Handle("/search", Default(controllers.Search))
	http.Handle("/new_comment", Default(controllers.CommentCreate))

	//comment oauth login
	http.Handle("/facebook_login", Default(oauth.FacebookLogin))
	http.Handle("/facebook_callback", Default(oauth.FacebookCallback))
	http.Handle("/google_login", Default(oauth.GoogleLogin))
	http.Handle("/google_callback", Default(oauth.GoogleCallback))
	http.Handle("/linkedin_login", Default(oauth.LinkedinLogin))
	http.Handle("/linkedin_callback", Default(oauth.LinkedinCallback))
	http.Handle("/vk_login", Default(oauth.VkLogin))
	http.Handle("/vk_callback", Default(oauth.VkCallback))

	{
		http.Handle("/admin", Restricted(controllers.Dashboard))

		http.Handle("/admin/users", Restricted(controllers.UserIndex))
		http.Handle("/admin/new_user", Restricted(controllers.UserCreate))
		http.Handle("/admin/edit_user/", Restricted(controllers.UserUpdate))
		http.Handle("/admin/delete_user", Restricted(controllers.UserDelete))

		http.Handle("/admin/pages", Restricted(controllers.PageIndex))
		http.Handle("/admin/new_page", Restricted(controllers.PageCreate))
		http.Handle("/admin/edit_page/", Restricted(controllers.PageUpdate))
		http.Handle("/admin/delete_page", Restricted(controllers.PageDelete))

		http.Handle("/admin/posts", Restricted(controllers.PostIndex))
		http.Handle("/admin/new_post", Restricted(controllers.PostCreate))
		http.Handle("/admin/edit_post/", Restricted(controllers.PostUpdate))
		http.Handle("/admin/delete_post", Restricted(controllers.PostDelete))

		http.Handle("/admin/tags", Restricted(controllers.TagIndex))
		http.Handle("/admin/new_tag", Restricted(controllers.TagCreate))
		http.Handle("/admin/delete_tag", Restricted(controllers.TagDelete))

		http.Handle("/admin/comments", Restricted(controllers.CommentIndex))
		http.Handle("/admin/new_comment", Restricted(controllers.CommentReply))
		http.Handle("/admin/edit_comment/", Restricted(controllers.CommentUpdate))
		http.Handle("/admin/delete_comment", Restricted(controllers.CommentDelete))

		//markdown editor does not support csrf when uploading images, so I have to apply CSRF middleware manually per route, sigh :/
		http.Handle("/admin/upload", RestrictedWithoutCSRF(controllers.Upload))
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public")))) //styles, js, images

	log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}

//Default executes default middleware chain for a HandlerFunc
func Default(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return CSRF(
		system.SessionMiddleware(
			system.TemplateMiddleware(
				system.DataMiddleware(
					http.HandlerFunc(fn),
				),
			),
		),
	)
}

//Restricted executes default + restriced middleware chain for a HandlerFunc
func Restricted(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return CSRF(
		RestrictedWithoutCSRF(fn),
	)
}

//RestrictedWithoutCSRF executes default + restriced middleware chain without CSRF middleware
func RestrictedWithoutCSRF(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return system.SessionMiddleware(
		system.TemplateMiddleware(
			system.DataMiddleware(
				system.RestrictedMiddleware(http.HandlerFunc(fn)),
			),
		),
	)
}
