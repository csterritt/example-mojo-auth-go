package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"mojo-auth-test-1/cookie_access"
	"mojo-auth-test-1/messages"
	"mojo-auth-test-1/views"

	"github.com/gin-gonic/gin"
	mojoauth "github.com/mojoauth/go-sdk"
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
const wantedLocationCookie = "wanted-loc"

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

func testIsValidPhoneAndClean(phoneNumber string) (bool, string) {
	cleanNumber := make([]rune, 0)

	for _, ch := range phoneNumber {
		if unicode.IsDigit(ch) {
			cleanNumber = append(cleanNumber, ch)
		}
	}

	res := string(cleanNumber)

	return len(res) == 10, res
}

func getSignInService(context *gin.Context) {
	email := cookie_access.GetTempCookie(context, emailCookie)
	if len(email) > 0 {
		RedirectTo(context, http.StatusFound, WaitSignInPath)
		return
	}

	err := signInTemplate.Render(context, gin.H{
		"greeting": "Hello (sign in) world!",
	})
	if err != nil {
		fmt.Printf("auth signInTemplate render failed with error %v\n", err)
	}
}

func signInWithEmail(context *gin.Context, authEmail string) {
	cookie_access.SetTempCookie(context, emailCookie, authEmail)

	errors := ""

	cfg := mojoauth.Config{
		ApiKey: os.Getenv("MOJO_APP_ID"),
	}
	mojoClient, err := mojoauth.NewMojoAuth(&cfg)
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
		var data MojoAuthState
		err = json.Unmarshal([]byte(res.Body), &data)
		if err == nil {
			cookie_access.SetSessionValue(context, stateIdCookie, data.StateId)
			RedirectToWithInfo(context, http.StatusFound, WaitSignInPath, "An email with the validation code has been sent.")
			return
		} else {
			fmt.Println("Error on JSON unmarshall:", err)
		}
	}

	RedirectToWithError(context, http.StatusFound, SignInPath, "An internal error occurred, please try again.")
}

func signInWithPhone(context *gin.Context, phoneNumber string) {
	cookie_access.SetTempCookie(context, emailCookie, phoneNumber)

	errors := ""

	cfg := mojoauth.Config{
		ApiKey: os.Getenv("MOJO_APP_ID"),
	}
	mojoClient, err := mojoauth.NewMojoAuth(&cfg)
	var res *httprutils.Response
	if err != nil {
		errors += err.(mojoerror.Error).OrigErr().Error()
		//      respCode = 500
	} else {
		body := map[string]string{
			"phone": "+1" + phoneNumber,
		}
		queryParams := map[string]string{
			"language": "en",
		}
		res, err = api.Mojoauth{Client: mojoClient}.SigninWithPhoneOTP(body, queryParams)
		if err != nil {
			errors += "SWPOTP: " + err.(mojoerror.Error).OrigErr().Error()
			//		respCode = 500
		}
	}

	if errors != "" {
		log.Printf(errors)
	} else {
		var data MojoAuthState
		err = json.Unmarshal([]byte(res.Body), &data)
		if err == nil {
			cookie_access.SetSessionValue(context, stateIdCookie, data.StateId)
			RedirectToWithInfo(context, http.StatusFound, WaitSignInPath, "An text with the validation code has been sent.")
			return
		} else {
			fmt.Println("Error on JSON unmarshall:", err)
		}
	}

	RedirectToWithError(context, http.StatusFound, SignInPath, "An internal error occurred, please try again.")
}

func postSignInService(context *gin.Context) {
	authEmail := strings.Trim(context.PostForm("auth_info_email"), " \n\r\t")
	if len(authEmail) != 0 && isValidEmail(authEmail) {
		signInWithEmail(context, authEmail)
		return
	}

	phoneNumber := strings.Trim(context.PostForm("auth_info_phone"), " \n\r\t")
	if len(phoneNumber) != 0 {
		isValid, cleanNumber := testIsValidPhoneAndClean(phoneNumber)

		if isValid {
			signInWithPhone(context, cleanNumber)
			return
		}
	}

	RedirectToWithError(context, http.StatusFound, SignInPath, "You must enter a valid email or phone number to sign in.")
}

func postSignOutService(context *gin.Context) {
	cookie_access.SetSessionValue(context, emailCookie, "")
	cookie_access.SetSessionValue(context, stateIdCookie, "")
	cookie_access.SetSessionValue(context, isAuthCookie, "")
	RedirectToWithInfo(context, http.StatusFound, RootPath, "Signed out successfully.")
}

func getWaitSignInService(context *gin.Context) {
	email := cookie_access.GetTempCookie(context, emailCookie)
	if len(email) == 0 {
		RedirectTo(context, http.StatusFound, SignInPath)
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
	RedirectToWithInfo(context, http.StatusFound, RootPath, "Sign in cancelled.")
}

func postWaitSignInService(context *gin.Context) {
	stateIdValue := cookie_access.GetSessionValue(context, stateIdCookie)
	if len(stateIdValue) == 0 {
		cookie_access.RemoveCookie(context, emailCookie)
		cookie_access.SetSessionValue(context, stateIdCookie, "")
		cookie_access.SetSessionValue(context, isAuthCookie, "")
		RedirectTo(context, http.StatusFound, RootPath)
	}

	authCode := strings.Trim(context.PostForm("auth_code_input"), " \n\r\t")
	if len(authCode) == 0 {
		RedirectToWithError(context, http.StatusFound, WaitSignInPath, "You must enter the sign-in code from the email.")
		return
	}

	body := map[string]string{
		"state_id": stateIdValue,
		"otp":      authCode,
	}
	cfg := mojoauth.Config{
		ApiKey: os.Getenv("MOJO_APP_ID"),
	}
	errors := ""
	mojoClient, err := mojoauth.NewMojoAuth(&cfg)
	_, err = api.Mojoauth{Client: mojoClient}.VerifyEmailOTP(body)
	if err != nil {
		errors += err.(mojoerror.Error).OrigErr().Error()
		log.Printf(errors)

		RedirectToWithError(context, http.StatusFound, WaitSignInPath, "That is an invalid or out-of-date code.")
		return
	}

	cookie_access.RemoveCookie(context, emailCookie)
	cookie_access.SetSessionValue(context, stateIdCookie, "")
	cookie_access.SetSessionValue(context, isAuthCookie, "true")

	wanted := cookie_access.GetSessionValue(context, wantedLocationCookie)
	if len(wanted) > 0 {
		cookie_access.SetSessionValue(context, wantedLocationCookie, "")
		messages.AddFlashMessage(context, messages.InfoMessage, "Signed in successfully.")
		context.Redirect(http.StatusFound, wanted)
	} else {
		RedirectToWithInfo(context, http.StatusFound, RootPath, "Signed in successfully.")
	}
}

func SkipAuthorizer() gin.HandlerFunc {
	return func(context *gin.Context) {
		isAuth := cookie_access.GetSessionValue(context, isAuthCookie)
		if isAuth == "true" {
			RedirectTo(context, http.StatusFound, RootPath)
		}

		context.Next()
	}
}

func InitializeAuthRoutes(router *gin.Engine) {
	signInTemplate = views.NewView("layout.html", "templates/views/auth/signIn.html")
	waitSignInTemplate = views.NewView("layout.html", "templates/views/auth/waiting.html")

	skipAuth := SkipAuthorizer()
	router.GET(Paths[SignInPath], skipAuth, getSignInService)
	router.POST(Paths[SignInPath], postSignInService)
	router.GET(Paths[WaitSignInPath], skipAuth, getWaitSignInService)
	router.POST(Paths[WaitSignInPath], postWaitSignInService)
	router.POST(Paths[CancelSignInPath], postCancelSignInService)
	router.POST(Paths[SignOutPath], postSignOutService)
}
