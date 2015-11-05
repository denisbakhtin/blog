package system

//go:generate rice embed-go

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/denisbakhtin/blog/models"
	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
)

//RunMigrations applies database migrations, where box - rice box for "migrations" dir, command:
//new - creates new blank migration in "migrations" directory. Edit that file as needed.
//"up", "down"- apply all pending migrations, or undo the last one
//"redo" - rollback last migration, then reapply it
//db - database handler
func RunMigrations(command *string) {
	box := rice.MustFindBox("../migrations")
	switch *command {
	case "new":
		migrateNew(box)
		os.Exit(0)
	case "up":
		migrateUp(models.GetDB(), box, 0)
		os.Exit(0)
	case "down":
		migrateDown(models.GetDB(), box, 1)
		os.Exit(0)
	case "redo":
		migrateDown(models.GetDB(), box, 1)
		migrateUp(models.GetDB(), box, 1)
		os.Exit(0)
	case "skip":
	default:
		log.Fatalf("Wrong migration flag %q, acceptable values: up, down\n", *command)
	}
}

//migrateNew creates new blank migration
func migrateNew(box *rice.Box) {
	if len(flag.Args()) == 0 {
		log.Fatalf("ERROR: Migration's name not specified\n")
		return
	}
	name := path.Join("migrations", fmt.Sprintf("%d_%s.sql", time.Now().Unix(), flag.Arg(0)))
	file, err := os.Create(name)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return
	}
	fmt.Fprintf(file, "-- +migrate Up\n")
	fmt.Fprintf(file, "-- SQL in section 'Up' is executed when this migration is applied\n\n\n")
	fmt.Fprintf(file, "-- +migrate Down\n")
	fmt.Fprintf(file, "-- SQL in section 'Down' is executed when this migration is rolled back\n\n\n")
	err = file.Close()
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	} else {
		log.Printf("INFO: File %s has been successfully created\n", name)
	}
}

//migrateUp applies {{max}} pending db migrations. If max == 0, it applies all
func migrateUp(db *sqlx.DB, box *rice.Box, max int) {
	migrations := getRiceMigrations(box)
	n, err := migrate.ExecMax(db.DB, "postgres", migrations, migrate.Up, max)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	} else {
		log.Printf("INFO: %d migration(s) applied\n", n)
	}
}

//migrateDown rolls back {{max}} db migrations. If max == 0, it rolles back all of them
func migrateDown(db *sqlx.DB, box *rice.Box, max int) {
	migrations := getRiceMigrations(box)
	n, err := migrate.ExecMax(db.DB, "postgres", migrations, migrate.Down, max)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	} else {
		log.Printf("INFO: %d migration(s) rolled back\n", n)
	}
}

//getRiceMigrations builds migration source from go.rice storage
func getRiceMigrations(box *rice.Box) *migrate.MemoryMigrationSource {
	source := &migrate.MemoryMigrationSource{}
	fn := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migFile, err := box.Open(path)
			if err != nil {
				return err
			}
			mig, err := migrate.ParseMigration(path, migFile)
			migFile.Close()
			if err != nil {
				return err
			}
			source.Migrations = append(source.Migrations, mig)
		}
		return nil
	}
	err := box.Walk("", fn)
	if err != nil {
		log.Panic(err)
		return nil
	}
	return source
}
