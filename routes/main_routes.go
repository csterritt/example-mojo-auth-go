package routes

import (
	"fmt"
	"net/http"

	"mojo-auth-test-1/cookie_access"
	"mojo-auth-test-1/messages"
	"mojo-auth-test-1/views"

	"github.com/gin-gonic/gin"
)

var indexTemplate *views.View
var showTemplate *views.View

func getShowService(context *gin.Context) {
	err := showTemplate.Render(context, gin.H{
		"greeting": "Hello (show) world!",
	})
	if err != nil {
		fmt.Printf("main showTemplate render failed with error %v\n", err)
	}
}

func getIndexService(context *gin.Context) {
	err := indexTemplate.Render(context, gin.H{
		"greeting": "Hello (index) world!",
	})
	if err != nil {
		fmt.Printf("main indexTemplate render failed with error %v\n", err)
	}
}

func Authorizer() gin.HandlerFunc {
	return func(context *gin.Context) {
		isAuth := cookie_access.GetSessionValue(context, isAuthCookie)
		fmt.Printf("Got isAuth %s\n", isAuth)
		if len(isAuth) == 0 || isAuth != "true" {
			cookie_access.SetSessionValue(context, wantedLocationCookie, context.FullPath())
			messages.AddFlashMessage(context, "error", "You must sign in to continue.")
			context.Redirect(http.StatusTemporaryRedirect, "/auth/sign-in")
			context.AbortWithStatus(http.StatusTemporaryRedirect)
			return
		}

		context.Next()
	}
}

func InitializeMainRoutes(router *gin.Engine) {
	indexTemplate = views.NewView("layout.html", "templates/views/index.html")
	showTemplate = views.NewView("layout.html", "templates/views/show.html")

	authorized := Authorizer()
	//router.GET("/", authorized, getIndexService)
	router.GET("/", getIndexService)
	router.GET("/show", authorized, getShowService)
}
