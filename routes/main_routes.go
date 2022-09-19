package routes

import (
	"fmt"
	"net/http"

	"mojo-auth-test-1/cookie_access"
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
	//isAuthorized := cookie_access.GetSessionValue(context, cookie_access.IsAuthorized)
	//err := indexTemplate.Render(context, gin.H{"isAuthorized": isAuthorized})
	//if err != nil {
	//	fmt.Printf("main indexTemplate render failed with error %v\n", err)
	//}
}

func Authorizer() gin.HandlerFunc {
	return func(context *gin.Context) {
		emailValue := cookie_access.GetSessionValue(context, emailCookie)
		isAuthorized := cookie_access.GetSessionValue(context, cookie_access.IsAuthorized)

		if len(emailValue) > 0 && isAuthorized == "true" {
			context.Next()
			return
		}

		if len(emailValue) == 0 && len(isAuthorized) != 0 {
			cookie_access.SetSessionValue(context, cookie_access.IsAuthorized, "")
			context.Redirect(http.StatusTemporaryRedirect, "/")
			context.AbortWithStatus(http.StatusTemporaryRedirect)
			return
		}

		cookie_access.SetSessionValue(context, cookie_access.IsAuthorized, "")
		context.Redirect(http.StatusTemporaryRedirect, "/auth/sign-in")
		context.AbortWithStatus(http.StatusTemporaryRedirect)
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
