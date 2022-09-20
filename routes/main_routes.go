package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mojo-auth-test-1/cookie_access"
	"mojo-auth-test-1/views"

	"github.com/gin-gonic/gin"
	go_mojoauth "github.com/mojoauth/go-sdk"
	"github.com/mojoauth/go-sdk/api"
	"github.com/mojoauth/go-sdk/mojoerror"
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
		jwt := cookie_access.GetSessionValue(context, jwtToken)
		fmt.Printf("Got jwt of length %d\n", len(jwt))
		if len(jwt) == 0 {
			context.Redirect(http.StatusTemporaryRedirect, "/auth/sign-in")
			context.AbortWithStatus(http.StatusTemporaryRedirect)
			return
		}

		cfg := go_mojoauth.Config{
			ApiKey: os.Getenv("MOJO_APP_ID"),
		}
		errors := ""
		mojoClient, err := go_mojoauth.NewMojoAuth(&cfg)
		res, err := api.Mojoauth{Client: mojoClient}.VerifyToken(jwt)
		fmt.Printf("res from VerifyToken is %#v\n", res)
		if err != nil {
			errors += err.(mojoerror.Error).OrigErr().Error()
			//		respCode = 500
		} else if res.IsValid {
			context.Next()
			return
		} else {
			errors += "res.IsValid is false?"
		}

		if errors != "" {
			log.Printf(errors)
			context.Redirect(http.StatusTemporaryRedirect, "/auth/sign-in")
			context.AbortWithStatus(http.StatusTemporaryRedirect)

			return
		}
		fmt.Println("Didn't get errors, but somehow wound up here. res.Invalid is", res.IsValid)

		context.Redirect(http.StatusTemporaryRedirect, "/")
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
