package cookie_access

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Name of the cookie.
const sessionName = "recipe-db-session"

// Maximum age of the session (about six months)
const sessionMaxAge = 3600 * 24 * 30 * 6

func getSessionStore() *sessions.CookieStore {
	// In real-world applications, use env variables to store the session key.
	store := sessions.NewCookieStore([]byte(os.Getenv("MAIN_AUTH_SECRET")), []byte(os.Getenv("MAIN_ENC_SECRET")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	return store
}

// GetSessionValue - Get the value of the named item from the session, or return ""
func GetSessionValue(context *gin.Context, name string) string {
	session, _ := getSessionStore().Get(context.Request, sessionName)
	valueHolder := session.Values[name]
	value := ""
	if valueHolder != nil {
		switch valueHolder.(type) {
		case string:
			value = valueHolder.(string)
		}
	}

	return strings.Trim(value, " \n\r\t")
}

// SetSessionValue - Get the value of the named item from the session, or return ""
func SetSessionValue(context *gin.Context, name string, value string) {
	session, _ := getSessionStore().Get(context.Request, sessionName)
	session.Values[name] = strings.Trim(value, " \n\r\t")

	err := session.Save(context.Request, context.Writer)
	if err != nil {
		// TODO: should probably add:
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("Got error trying to do a session.Save: %v\n", err)
	}
}
