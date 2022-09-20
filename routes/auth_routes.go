package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"mojo-auth-test-1/cookie_access"
	"mojo-auth-test-1/views"

	"github.com/gin-gonic/gin"
)

type AuthCodeResult int

const (
	AuthExpired AuthCodeResult = iota
	AuthDbError
	AuthSuccessSameBrowser
	AuthSuccessDifferentBrowser
)

// Name of the email.
const emailCookie = "given-email"

var signInTemplate *views.View
var waitSignInTemplate *views.View
var emailPattern *regexp.Regexp

func init() {
	emailPattern = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
}

// validate that the given string looks like an email address
func isValidEmail(email string) bool {
	return emailPattern.MatchString(email)
}

func getSignInService(context *gin.Context) {
	email := cookie_access.GetCookie(context, emailCookie)
	if len(email) > 0 {
		context.Redirect(http.StatusFound, "/auth/wait-sign-in")
		return
	}

	err := signInTemplate.Render(context, gin.H{
		"greeting": "Hello (sign in) world!",
	})
	if err != nil {
		fmt.Printf("auth signInTemplate render failed with error %v\n", err)
	}
	//emailValue := cookie_access.GetSessionValue(context, emailCookie)
	//if emailValue != "" {
	//	context.Redirect(http.StatusFound, "/auth/waiting")
	//	return
	//}
}

func postSignInService(context *gin.Context) {
	authEmail := strings.Trim(context.PostForm("auth_info_email"), " \n\r\t")
	if isValidEmail(authEmail) {
		cookie_access.SetSessionCookie(context, emailCookie, authEmail)
		context.Redirect(http.StatusFound, "/auth/wait-sign-in")
		return
	}

	//	if len(authEmail) == db_access.TokenLength {
	//		code := finishSignInWithCode(context, authEmail)
	//		if code == AuthSuccessSameBrowser {
	//			context.Redirect(http.StatusFound, "/")
	//			return
	//		} else if code == AuthSuccessDifferentBrowser {
	//			context.Redirect(http.StatusFound, "/auth/different-browser")
	//			return
	//		}
	//	}
	//	cookie_access.SetSessionCookie(context, submittedEmailCookie, authEmail)
	//
	//	messages.AddFlashMessage(context, "error", "That is not a valid email address or sign-in code.")
	//	context.Redirect(http.StatusFound, "/auth/sign-in")
	//	return
	//}
	//
	//if len(authEmail) > 0 {
	//	setupForSignIn(context, authEmail)
	//} else {
	//	context.Redirect(http.StatusFound, "/auth/sign-in")
	//}
	context.Redirect(http.StatusFound, "/auth/sign-in")
}

func postSignOutService(context *gin.Context) {
	cookie_access.SetSessionValue(context, emailCookie, "")
	cookie_access.SetSessionValue(context, cookie_access.IsAuthorized, "")
	context.Redirect(http.StatusFound, "/")
}

func getWaitSignInService(context *gin.Context) {
	email := cookie_access.GetCookie(context, emailCookie)
	if len(email) == 0 {
		context.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	err := waitSignInTemplate.Render(context, gin.H{
		"greeting": "Hello (wait sign in) world!",
	})
	if err != nil {
		fmt.Printf("auth waitSignInTemplate render failed with error %v\n", err)
	}
	//emailValue := cookie_access.GetSessionValue(context, emailCookie)
	//if emailValue != "" {
	//	context.Redirect(http.StatusFound, "/auth/waiting")
	//	return
	//}
}

func postCancelSignInService(context *gin.Context) {
	cookie_access.RemoveCookie(context, emailCookie)
	cookie_access.SetSessionValue(context, emailCookie, "")
	cookie_access.SetSessionValue(context, cookie_access.IsAuthorized, "")
	context.Redirect(http.StatusFound, "/")
}

func postWaitSignInService(context *gin.Context) {
	authCode := strings.Trim(context.PostForm("auth_code_input"), " \n\r\t")
	if authCode == "1234" {
		cookie_access.SetSessionValue(context, cookie_access.IsAuthorized, "true")
		email := cookie_access.GetCookie(context, emailCookie)
		cookie_access.SetSessionValue(context, emailCookie, email)
		cookie_access.RemoveCookie(context, emailCookie)
		context.Redirect(http.StatusFound, "/show")
		return
	}

	//authEmail := strings.Trim(context.PostForm("auth_info_email"), " \n\r\t")
	//if !isValidEmail(authEmail) {
	//	if len(authEmail) == db_access.TokenLength {
	//		code := finishSignInWithCode(context, authEmail)
	//		if code == AuthSuccessSameBrowser {
	//			context.Redirect(http.StatusFound, "/")
	//			return
	//		} else if code == AuthSuccessDifferentBrowser {
	//			context.Redirect(http.StatusFound, "/auth/different-browser")
	//			return
	//		}
	//	}
	//	cookie_access.SetSessionCookie(context, submittedEmailCookie, authEmail)
	//
	//	messages.AddFlashMessage(context, "error", "That is not a valid email address or sign-in code.")
	//	context.Redirect(http.StatusFound, "/auth/sign-in")
	//	return
	//}
	//
	//if len(authEmail) > 0 {
	//	setupForSignIn(context, authEmail)
	//} else {
	//	context.Redirect(http.StatusFound, "/auth/sign-in")
	//}
	context.Redirect(http.StatusFound, "/auth/wait-sign-in")
}

func SkipAuthorizer() gin.HandlerFunc {
	return func(context *gin.Context) {
		emailValue := cookie_access.GetSessionValue(context, emailCookie)
		isAuthorized := cookie_access.GetSessionValue(context, cookie_access.IsAuthorized)

		if len(emailValue) > 0 && isAuthorized == "true" {
			context.Redirect(http.StatusTemporaryRedirect, "/")
			context.AbortWithStatus(http.StatusTemporaryRedirect)
		}
		context.Next()
	}
}

func InitializeAuthRoutes(router *gin.Engine) {
	signInTemplate = views.NewView("layout.html", "templates/views/auth/signIn.html")
	waitSignInTemplate = views.NewView("layout.html", "templates/views/auth/waiting.html")

	skipAuth := SkipAuthorizer()
	router.GET("/auth/sign-in", skipAuth, getSignInService)
	router.POST("/auth/sign-in", postSignInService)
	router.GET("/auth/wait-sign-in", skipAuth, getWaitSignInService)
	router.POST("/auth/wait-sign-in", postWaitSignInService)
	router.POST("/auth/cancel-sign-in", postCancelSignInService)
	router.POST("/auth/sign-out", postSignOutService)
}
