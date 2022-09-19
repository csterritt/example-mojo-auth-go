package routes

import (
	"fmt"
	"net/http"

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

func getSignInService(context *gin.Context) {
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
}

func getWaitSignInService(context *gin.Context) {
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

func postWaitSignInService(context *gin.Context) {
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
}
