package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/claudiu/gocron"
	"github.com/denisbakhtin/blog/routes"
	"github.com/denisbakhtin/blog/shared"
	"github.com/gorilla/context"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	migrate := flag.String("migrate", "skip", "Run DB migrations: up, down, redo, new [MIGRATION_NAME] and then os.Exit(0)")
	mode := flag.String("mode", "debug", "Application mode: debug, release, test")
	flag.Parse()

	shared.SetMode(mode)
	shared.Init()
	shared.RunMigrations(migrate)

	//Periodic tasks
	gocron.Every(1).Day().Do(shared.CreateXMLSitemap)
	gocron.Start()

	routes.Init()

	log.Printf("Listening on :%d port in %q mode\n", 8080, *mode)
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux)))
}
