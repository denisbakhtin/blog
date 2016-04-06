package routes

import (
	"github.com/denisbakhtin/blog/controllers"
	"github.com/denisbakhtin/blog/controllers/oauth"
	"github.com/denisbakhtin/blog/shared"
	"net/http"
)

func Init() {
	shared.InitCsrf()

	http.Handle("/", shared.Default(controllers.Home))

	if shared.GetConfig().SignupEnabled {
		http.Handle("/signup", shared.Default(controllers.SignUp))
	}
	http.Handle("/signin", shared.Default(controllers.SignIn))
	http.Handle("/logout", shared.Default(controllers.Logout))

	http.Handle("/pages/", shared.Default(controllers.PageShow))
	http.Handle("/posts/", shared.Default(controllers.PostShow))
	http.Handle("/tags/", shared.Default(controllers.TagShow))
	http.Handle("/archives/", shared.Default(controllers.ArchiveShow))
	http.Handle("/rss", shared.Default(controllers.RssXML))
	http.Handle("/search", shared.Default(controllers.Search))
	http.Handle("/new_comment", shared.Default(controllers.CommentCreate))

	//comment oauth login
	http.Handle("/facebook_login", shared.Default(oauth.FacebookLogin))
	http.Handle("/facebook_callback", shared.Default(oauth.FacebookCallback))
	http.Handle("/google_login", shared.Default(oauth.GoogleLogin))
	http.Handle("/google_callback", shared.Default(oauth.GoogleCallback))
	http.Handle("/linkedin_login", shared.Default(oauth.LinkedinLogin))
	http.Handle("/linkedin_callback", shared.Default(oauth.LinkedinCallback))
	http.Handle("/vk_login", shared.Default(oauth.VkLogin))
	http.Handle("/vk_callback", shared.Default(oauth.VkCallback))

	{
		http.Handle("/admin", shared.Restricted(controllers.Dashboard))

		http.Handle("/admin/users", shared.Restricted(controllers.UserIndex))
		http.Handle("/admin/new_user", shared.Restricted(controllers.UserCreate))
		http.Handle("/admin/edit_user/", shared.Restricted(controllers.UserUpdate))
		http.Handle("/admin/delete_user", shared.Restricted(controllers.UserDelete))

		http.Handle("/admin/pages", shared.Restricted(controllers.PageIndex))
		http.Handle("/admin/new_page", shared.Restricted(controllers.PageCreate))
		http.Handle("/admin/edit_page/", shared.Restricted(controllers.PageUpdate))
		http.Handle("/admin/delete_page", shared.Restricted(controllers.PageDelete))

		http.Handle("/admin/posts", shared.Restricted(controllers.PostIndex))
		http.Handle("/admin/new_post", shared.Restricted(controllers.PostCreate))
		http.Handle("/admin/edit_post/", shared.Restricted(controllers.PostUpdate))
		http.Handle("/admin/delete_post", shared.Restricted(controllers.PostDelete))
		http.Handle("/admin/post_on_facebook", shared.RestrictedWithoutCSRF(controllers.PostOnFacebook))

		http.Handle("/admin/tags", shared.Restricted(controllers.TagIndex))
		http.Handle("/admin/new_tag", shared.Restricted(controllers.TagCreate))
		http.Handle("/admin/delete_tag", shared.Restricted(controllers.TagDelete))

		http.Handle("/admin/comments", shared.Restricted(controllers.CommentIndex))
		http.Handle("/admin/new_comment", shared.Restricted(controllers.CommentReply))
		http.Handle("/admin/edit_comment/", shared.Restricted(controllers.CommentUpdate))
		http.Handle("/admin/delete_comment", shared.Restricted(controllers.CommentDelete))

		//markdown editor does not support csrf when uploading images, so I have to apply CSRF middleware manually per route, sigh :/
		http.Handle("/admin/upload", shared.RestrictedWithoutCSRF(controllers.Upload))
	}

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("static/public")))) //styles, js, images
}
