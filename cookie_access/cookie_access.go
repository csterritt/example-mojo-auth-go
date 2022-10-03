package cookie_access

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetTempCookie - Get the value of the named cookie, or return ""
func GetTempCookie(context *gin.Context, name string) string {
	var value string
	holder, err := context.Request.Cookie(name)
	if err != nil {
		value = ""
	} else {
		value = holder.Value
	}

	return strings.Trim(value, " \n\r\t")
}

// SetTempCookie - Add a new browser-session (that is, it lives until the browser is closed)
// cookie into the cookie storage.
func SetTempCookie(context *gin.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		SameSite: 1,
	}

	http.SetCookie(context.Writer, cookie)
}

// RemoveCookie - Remove a cookie from the cookie storage.
func RemoveCookie(context *gin.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
		SameSite: 1,
	}

	http.SetCookie(context.Writer, cookie)
}
