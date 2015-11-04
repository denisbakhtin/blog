package helpers

import (
	"golang.org/x/net/context"
	"strconv"
)

//DefaultData returns common to all pages template data
func DefaultData(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"ActiveUser":    ctx.Value("user"),           //signed in models.User
		"Active":        "",                          //active uri shortening for menu item highlight
		"Title":         "",                          //page title:w
		"SignupEnabled": ctx.Value("signup_enabled"), //signup route is enabled (otherwise everyone can signup ;)
	}
}

//ErrorData returns template data for error
func ErrorData(err error) map[string]interface{} {
	return map[string]interface{}{
		"Title": err.Error(),
		"Error": err.Error(),
	}
}

//Atoi64 converts string to int64, returns 0 if error
func Atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

//Atob converts string to bool
func Atob(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}
