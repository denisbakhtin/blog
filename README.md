Blog boilerplate
===============

Provides essentials that most web blogs need - MVC pattern, user authorisation, SQL db migration, admin dashboard, javascript form validation, rss feeds, etc.

It consists of the following core components:
- Standard http multiplexer (simplicity is king, no more microframeworks)
- context - gorilla context (for native HandlerFuncs) - https://github.com/gorilla/context
- sessions - gorilla sessions - https://github.com/gorilla/sessions
- csrf - gorilla csrf - https://github.com/gorilla/csrf
- pq - Postgres driver for the Go language - https://github.com/lib/pq
- sqlx - Relational database access interface - https://github.com/jmoiron/sqlx
- go.rice - Asset packaging tool for easy deployment - https://github.com/GeertJohan/go.rice
- Twitter Bootstrap - popular HTML, CSS, JS framework for developing responsive, mobile first web projects - http://getbootstrap.com
- Parsley JS - form validation - http://parsleyjs.org
- Bootstrap Markdown Editor with image upload - https://github.com/inacho/bootstrap-markdown-editor
- blackfriday - markdown processor - https://github.com/russross/blackfriday
- bluemonday - html sanitizer (for excerpts, etc) - https://github.com/microcosm-cc/bluemonday
- RSS feeds - https://github.com/gorilla/feeds
- sitemap - XML sitemap for search engines - https://github.com/denisbakhtin/sitemap
- gocron - periodic task launcher (for sitemap generation, etc) - https://github.com/jasonlvhit/gocron
- sql-migrate - SQL schema migration tool - https://github.com/rubenv/sql-migrate
- search - Postgresql full text search on posts table, can be extended to multiple tables + rankings - http://www.postgresql.org/docs/9.4/static/textsearch-intro.html
- Post comments with oauth2 authentication
- Post blog entries on facebook page
- go-i18n - Multi language interface - https://github.com/nicksnyder/go-i18n

# TODO
- Social plugins (share, like buttons)
- Post updates on twitter, google +, linkedin

# Usage
```
git clone https://github.com/denisbakhtin/blog.git
cd blog
go get .
```
Copy sample config `cp config/config.json.example config/config.json`, create postgresql database, modify config/config.json accordingly.

Use `go run main.go -migrate=up` to apply initial database migrations.
`go run main.go` to launch web server. Or if you have https://github.com/cespare/reflex installed `make debug`.

# Deployment
```
make build
```
Upload `ginblog` binary and `public` directory to your server. If you find `rice embed-go` is running slow on your system, consider using other [go.rice packing options](https://github.com/GeertJohan/go.rice#tool-usage) with `go generate` command.

# Project structure

`/config`

Contains application configuration file.

`/controllers`

All your controllers that serve defined routes.

`/helpers`

Helper functions.

`/migrations`

Database schema migrations

`/models`

You database models.

`/public`

It has all your static files

`/system`

Core functions and structs.

`/views`

Your views using standard `Go` template system.

`main.go`

This file starts your web application, contains routes definition & some custom middlewares.

# Make it your own

I assume you have followed installation instructions and you have `ginblog` installed in your `GOPATH` location.

Let's say I want to create `Amazing Blog`. I create new `GitHub` repository `https://github.com/denisbakhtin/blog` (of course replace that with your own repository).

Now I have to prepare `blog`. First thing is that I have to delete its `.git` directory.

I issue:

```
rm -rf src/github.com/denisbakhtin/blog/.git
```

Then I want to replace all references from `github.com/denisbakhtin/blog` to `github.com/denisbakhtin/amazingblog`:

```
grep -rl 'github.com/denisbakhtin/blog' ./ | xargs sed -i 's/github.com\/denisbakhtin\/blog/github.com\/denisbakhtin\/amazingblog/g'
```

Now I have to move all `blog` files to the new location:

```
mv src/github.com/denisbakhtin/blog/ src/github.com/denisbakhtin/amazingblog
```

And push it to my new repository at `GitHub`:

```
cd src/github.com/denisbakhtin/amazingblog
git init
git add --all .
git commit -m "Amazing Blog First Commit"
git remote add origin https://github.com/denisbakhtin/amazingblog.git
git push -u origin master
```

You can now go back to your `GOPATH` and check if everything is ok:

```
go install github.com/denisbakhtin/amazingblog
```

And that's it. 

# Continuous Development

For Continuous Development I recommend using `Reflex` - https://github.com/cespare/reflex

You can install `Reflex` by issuing:

```
go get github.com/cespare/reflex
```

Then create a config file `reflex.conf` in your `GOPATH`:

```
# Restart server when .go, .html files change
-sr '(\.go|\.html)$' go run main.go
```

Now if you run:

```
reflex -c reflex.conf
```

Project will automatically rebuild itself when a change to *.go, *.html files occurs. For more options read https://github.com/cespare/reflex#usage

