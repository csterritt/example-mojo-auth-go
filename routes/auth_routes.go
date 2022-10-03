package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"mojo-auth-test-1/cookie_access"
	"mojo-auth-test-1/messages"
	"mojo-auth-test-1/views"

	"github.com/gin-gonic/gin"
	go_mojoauth "github.com/mojoauth/go-sdk"
	"github.com/mojoauth/go-sdk/api"
	"github.com/mojoauth/go-sdk/httprutils"
	"github.com/mojoauth/go-sdk/mojoerror"
)

/*
	{
		"authenticated":true,
		"oauth":{
			"access_token":"<long-string>",
			"id_token":"<long-string>",
			"refresh_token":"3bf1c663-5226-4b35-8710-b04beb77cbb4",
			"expires_in":"2022-09-23T20:26:27Z",
			"token_type":"Bearer"
		},
		"user":{
			"created_at":"2022-09-20T20:26:27Z",
			"updated_at":"2022-09-20T20:26:27Z",
			"issuer":"https://www.mojoauth.com",
			"user_id":"6324d91c3b65af019d250ac7",
			"identifier":"csterritt@gmail.com"
		}
	}
*/
type AuthCodeResult int
type MojoAuthState struct {
	StateId string `json:"state_id"`
}
type MojoAuthResult struct {
	Authenticated bool `json:"authenticated"`
	Oauth         struct {
		AccessToken  string `json:"access_token"`
		IdToken      string `json:"id_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	User struct {
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		Issuer     string `json:"issuer"`
		UserId     string `json:"user_id"`
		Identifier string `json:"identifier"`
	}
}

// Name of the email.
const emailCookie = "given-email"
const stateIdCookie = "state-id"
const isAuthCookie = "is-auth"

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
}

func postSignInService(context *gin.Context) {
	authEmail := strings.Trim(context.PostForm("auth_info_email"), " \n\r\t")
	if len(authEmail) == 0 {
		messages.AddFlashMessage(context, "error", "That is not a valid email.")
		context.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	if !isValidEmail(authEmail) {
		messages.AddFlashMessage(context, "error", "That is not a valid email.")
		context.Redirect(http.StatusFound, "/auth/sign-in")
		return
	} else {
		cookie_access.SetSessionCookie(context, emailCookie, authEmail)

		errors := ""

		cfg := go_mojoauth.Config{
			ApiKey: os.Getenv("MOJO_APP_ID"),
		}
		mojoClient, err := go_mojoauth.NewMojoAuth(&cfg)
		var res *httprutils.Response
		if err != nil {
			errors += err.(mojoerror.Error).OrigErr().Error()
			//      respCode = 500
		} else {
			body := map[string]string{
				"email": authEmail,
			}
			queryParams := map[string]string{
				"language": "en",
			}
			res, err = api.Mojoauth{Client: mojoClient}.SigninWithEmailOTP(body, queryParams)
			if err != nil {
				errors += "SWEOTP: " + err.(mojoerror.Error).OrigErr().Error()
				//		respCode = 500
			}
		}

		if errors != "" {
			log.Printf(errors)
		} else {
			fmt.Printf("Raw body: '%s'\n", res.Body)
			var data MojoAuthState
			err = json.Unmarshal([]byte(res.Body), &data)
			if err == nil {
				fmt.Printf("Raw body: '%s', struct %#v\n", res.Body, data)
				cookie_access.SetSessionValue(context, stateIdCookie, data.StateId)
				messages.AddFlashMessage(context, "info", "An email with the validation code has been sent.")
				context.Redirect(http.StatusFound, "/auth/wait-sign-in")
				return
			} else {
				fmt.Println("Error on JSON unmarshall:", err)
			}
		}
	}

	messages.AddFlashMessage(context, "error", "An internal error occured, please try again.")
	context.Redirect(http.StatusFound, "/auth/sign-in")
}

func postSignOutService(context *gin.Context) {
	cookie_access.SetSessionValue(context, emailCookie, "")
	cookie_access.SetSessionValue(context, stateIdCookie, "")
	cookie_access.SetSessionValue(context, isAuthCookie, "")
	messages.AddFlashMessage(context, "info", "Signed out successfully.")
	context.Redirect(http.StatusFound, "/")
}

func getWaitSignInService(context *gin.Context) {
	email := cookie_access.GetCookie(context, emailCookie)
	if len(email) == 0 {
		context.Redirect(http.StatusFound, "/auth/sign-in")
		return
	}

	err := waitSignInTemplate.Render(context, gin.H{})
	if err != nil {
		fmt.Printf("auth waitSignInTemplate render failed with error %v\n", err)
	}
}

func postCancelSignInService(context *gin.Context) {
	cookie_access.RemoveCookie(context, emailCookie)
	cookie_access.SetSessionValue(context, emailCookie, "")
	cookie_access.SetSessionValue(context, stateIdCookie, "")
	cookie_access.SetSessionValue(context, isAuthCookie, "")
	messages.AddFlashMessage(context, "info", "Sign in cancelled.")
	context.Redirect(http.StatusFound, "/")
}

func postWaitSignInService(context *gin.Context) {
	stateIdValue := cookie_access.GetSessionValue(context, stateIdCookie)
	if len(stateIdValue) == 0 {
		cookie_access.RemoveCookie(context, emailCookie)
		cookie_access.SetSessionValue(context, stateIdCookie, "")
		cookie_access.SetSessionValue(context, isAuthCookie, "")
		context.Redirect(http.StatusFound, "/")
	}

	authCode := strings.Trim(context.PostForm("auth_code_input"), " \n\r\t")
	if len(authCode) == 0 {
		messages.AddFlashMessage(context, "error", "You must enter the sign-in code from the email.")
		context.Redirect(http.StatusFound, "/auth/wait-sign-in")
		return
	}

	body := map[string]string{
		"state_id": stateIdValue,
		"otp":      authCode,
	}
	cfg := go_mojoauth.Config{
		ApiKey: os.Getenv("MOJO_APP_ID"),
	}
	errors := ""
	mojoClient, err := go_mojoauth.NewMojoAuth(&cfg)
	res, err := api.Mojoauth{Client: mojoClient}.VerifyEmailOTP(body)
	if err != nil {
		errors += err.(mojoerror.Error).OrigErr().Error()
		//		respCode = 500
	}

	if errors != "" {
		log.Printf(errors)

		messages.AddFlashMessage(context, "error", "That is an invalid or out-of-date code.")
		context.Redirect(http.StatusFound, "/auth/wait-sign-in")
		return
	}
	fmt.Println(res.Body)
	var info MojoAuthResult
	err = json.Unmarshal([]byte(res.Body), &info)
	if err == nil {
		fmt.Printf("Saving jwtToken of length %d, refreshToken of length %d\n", len(info.Oauth.AccessToken), len(info.Oauth.RefreshToken))
		cookie_access.RemoveCookie(context, emailCookie)
		cookie_access.SetSessionValue(context, stateIdCookie, "")
		cookie_access.SetSessionValue(context, isAuthCookie, "true")
		messages.AddFlashMessage(context, "info", "Signed in successfully.")
		context.Redirect(http.StatusFound, "/")
		return
	} else {
		fmt.Printf("Got error unmarshalling MojoAuthResult: %v\n", err)
	}

	context.Redirect(http.StatusFound, "/auth/wait-sign-in")
}

func SkipAuthorizer() gin.HandlerFunc {
	return func(context *gin.Context) {
		isAuth := cookie_access.GetSessionValue(context, isAuthCookie)
		if isAuth == "true" {
			context.Redirect(http.StatusFound, "/")
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
