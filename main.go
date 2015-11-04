package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/blog/controllers"
	"github.com/denisbakhtin/blog/system"
	"github.com/gorilla/csrf"
)

var csrfMiddleware func(http.Handler) http.Handler

func main() {
	migrate := flag.String("migrate", "skip", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	system.SetMode(mode)
	system.Init()
	system.RunMigrations(migrate)

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

	//TODO: set csrf.Secure(true) in release mode if http is protected by ssl
	csrfMiddleware = csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))
	{
		http.Handle("/admin", RestrictedCsrf(controllers.Dashboard))

		http.Handle("/admin/users", RestrictedCsrf(controllers.UserIndex))
		http.Handle("/admin/new_user", RestrictedCsrf(controllers.UserCreate))
		http.Handle("/admin/edit_user/", RestrictedCsrf(controllers.UserUpdate))
		http.Handle("/admin/delete_user", RestrictedCsrf(controllers.UserDelete))

		http.Handle("/admin/pages", RestrictedCsrf(controllers.PageIndex))
		http.Handle("/admin/new_page", RestrictedCsrf(controllers.PageCreate))
		http.Handle("/admin/edit_page/", RestrictedCsrf(controllers.PageUpdate))
		http.Handle("/admin/delete_page", RestrictedCsrf(controllers.PageDelete))

		http.Handle("/admin/posts", RestrictedCsrf(controllers.PostIndex))
		http.Handle("/admin/new_post", RestrictedCsrf(controllers.PostCreate))
		http.Handle("/admin/edit_post/", RestrictedCsrf(controllers.PostUpdate))
		http.Handle("/admin/delete_post", RestrictedCsrf(controllers.PostDelete))

		http.Handle("/admin/tags", RestrictedCsrf(controllers.TagIndex))
		http.Handle("/admin/new_tag", RestrictedCsrf(controllers.TagCreate))
		http.Handle("/admin/delete_tag", RestrictedCsrf(controllers.TagDelete))

		http.Handle("/admin/upload", Restricted(controllers.Upload)) //without csrf, since markdown editor does not support csrf :/
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public")))) //styles, js, images

	log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}

//Default executes default middleware chain for a HandlerFunc
func Default(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return system.SessionMiddleware(
		system.TemplateMiddleware(
			system.DataMiddleware(http.HandlerFunc(fn)),
		),
	)
}

//Restricted executes default + restriced middleware chain for a HandlerFunc
func Restricted(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return system.SessionMiddleware(
		system.TemplateMiddleware(
			system.DataMiddleware(
				system.RestrictedMiddleware(http.HandlerFunc(fn)),
			),
		),
	)
}

//RestrictedCsrf executes default + restriced + csrf middleware chain for a HandlerFunc
func RestrictedCsrf(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return system.SessionMiddleware(
		system.TemplateMiddleware(
			system.DataMiddleware(
				system.RestrictedMiddleware(
					csrfMiddleware(http.HandlerFunc(fn)),
				),
			),
		),
	)
}
