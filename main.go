package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/blog/controllers"
	"github.com/denisbakhtin/blog/system"
	_ "github.com/gorilla/csrf"
)

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

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware

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

	{
		http.Handle("/admin", DefaultRestricted(controllers.Dashboard))

		http.Handle("/admin/users", DefaultRestricted(controllers.UserIndex))
		http.Handle("/admin/new_user", DefaultRestricted(controllers.UserCreate))
		http.Handle("/admin/edit_user/", DefaultRestricted(controllers.UserUpdate))
		http.Handle("/admin/delete_user", DefaultRestricted(controllers.UserDelete))

		http.Handle("/admin/pages", DefaultRestricted(controllers.PageIndex))
		http.Handle("/admin/new_page", DefaultRestricted(controllers.PageCreate))
		http.Handle("/admin/edit_page/", DefaultRestricted(controllers.PageUpdate))
		http.Handle("/admin/delete_page", DefaultRestricted(controllers.PageDelete))

		http.Handle("/admin/posts", DefaultRestricted(controllers.PostIndex))
		http.Handle("/admin/new_post", DefaultRestricted(controllers.PostCreate))
		http.Handle("/admin/edit_post/", DefaultRestricted(controllers.PostUpdate))
		http.Handle("/admin/delete_post", DefaultRestricted(controllers.PostDelete))

		http.Handle("/admin/tags", DefaultRestricted(controllers.TagIndex))
		http.Handle("/admin/new_tag", DefaultRestricted(controllers.TagCreate))
		http.Handle("/admin/delete_tag", DefaultRestricted(controllers.TagDelete))

		http.Handle("/admin/upload", DefaultRestricted(controllers.Upload))
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public")))) //styles, js, images

	//TODO: set csrf.Secure(true) for release mode if ssl is available
	//CSRF := csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))
	//log.Fatal(http.ListenAndServe(":8080", CSRF(http.DefaultServeMux)))
	log.Fatal(http.ListenAndServe(":8080", http.DefaultServeMux))
}

//Default executes default middleware chain for a HandlerFunc
func Default(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		system.SessionMiddleware(
			system.TemplateMiddleware(
				system.DataMiddleware(http.HandlerFunc(fn)),
			),
		).ServeHTTP(w, r)
	})
}

//DefaultRestricted executes default + restriced middleware chain for a HandlerFunc
func DefaultRestricted(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		system.SessionMiddleware(
			system.TemplateMiddleware(
				system.DataMiddleware(
					system.RestrictedMiddleware(http.HandlerFunc(fn)),
				),
			),
		).ServeHTTP(w, r)
	})
}
