package system

//go:generate rice embed-go

import (
	"github.com/GeertJohan/go.rice"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"os"
	"regexp"
)

func loadI18n() {
	box := rice.MustFindBox("../config")

	fn := func(path string, f os.FileInfo, err error) error {
		if matched, _ := regexp.MatchString("[a-z]{2}-[a-z]{2}\\.all\\.json", f.Name()); matched {
			err := i18n.ParseTranslationFileBytes(f.Name(), []byte(box.MustString(path)))
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := box.Walk("", fn)
	if err != nil {
		log.Fatal(err)
	}
}
