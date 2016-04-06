package controllers

import (
	"fmt"
	"github.com/denisbakhtin/blog/models"
	"github.com/denisbakhtin/blog/shared"
	"github.com/gorilla/feeds"
	"log"
	"net/http"
	"time"
)

//RssXML handles GET /rss route
func RssXML(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	if r.Method == "GET" {

		now := time.Now()
		domain := shared.GetConfig().Domain
		feed := &feeds.Feed{
			Title:       "Blog boilerplate",
			Link:        &feeds.Link{Href: domain},
			Description: "Basic blog boilerplate in Go",
			Author:      &feeds.Author{Name: "Blog Author"},
			Created:     now,
			Copyright:   fmt.Sprintf("© %s", "Blog"),
		}

		feed.Items = make([]*feeds.Item, 0)
		posts, err := models.GetPublishedPosts()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		for i := range posts {
			feed.Items = append(feed.Items, &feeds.Item{
				Id:          fmt.Sprintf("%s/posts/%d", domain, posts[i].ID),
				Title:       posts[i].Name,
				Link:        &feeds.Link{Href: fmt.Sprintf("%s/posts/%d", domain, posts[i].ID)},
				Description: string(posts[i].Excerpt()),
				Author:      &feeds.Author{Name: posts[i].Author.Name},
				Created:     now,
			})
		}

		rss, err := feed.ToRss()
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		fmt.Fprintln(w, rss)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}
