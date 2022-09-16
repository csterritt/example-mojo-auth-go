package messages

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Name of the cookie.
const sessionName = "flash-messages"

func getCookieStore() *sessions.CookieStore {
	// In real-world applications, use env variables to store the session key.
	sessionKey := "test-session-key"
	return sessions.NewCookieStore([]byte(sessionKey))
}

// AddFlashMessage -- Add a new message into the cookie storage.
func AddFlashMessage(context *gin.Context, name, value string) {
	session, _ := getCookieStore().Get(context.Request, sessionName)
	session.AddFlash(value, name)

	_ = session.Save(context.Request, context.Writer)
}

// GetFlashMessages -- Get flash messages from the cookie storage.
func GetFlashMessages(context *gin.Context, name string) []string {
	session, _ := getCookieStore().Get(context.Request, sessionName)
	fm := session.Flashes(name)
	// If we have some messages.
	if len(fm) > 0 {
		_ = session.Save(context.Request, context.Writer)
		// Initiate a strings slice to return messages.
		var flashes []string
		for _, fl := range fm {
			// Add message to the slice.
			flashes = append(flashes, fl.(string))
		}

		return flashes
	}

	return nil
}
