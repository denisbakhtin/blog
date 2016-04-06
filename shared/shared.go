package shared

import (
	"fmt"
	"github.com/denisbakhtin/blog/models"
	"os"
	"strings"
)

//Application mode
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

var mode string //application mode: debug, release, test

//Init initializes core system elements (DB, sessions, templates, et al)
func Init() {
	loadConfig()
	CreateSessionStore(config.SessionSecret, config.Domain, config.Ssl)
	loadTemplates()
	connection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name)
	models.InitDB(connection)
}

//SetMode sets application mode
func SetMode(flagmode *string) {
	switch *flagmode {
	case DebugMode:
		mode = DebugMode
	case ReleaseMode:
		mode = ReleaseMode
	default:
		mode = TestMode
	}
}

//GetMode returns application mode
func GetMode() string {
	return mode
}

//Getwd returns app root directory (even for tests)
func Getwd() string {
	wd, _ := os.Getwd()
	wd = strings.Replace(wd, "/controllers", "", 1)
	wd = strings.Replace(wd, "/views", "", 1)
	wd = strings.Replace(wd, "/config", "", 1)
	return wd
}
