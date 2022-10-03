package routes

import (
	"mojo-auth-test-1/messages"

	"github.com/gin-gonic/gin"
)

type Path int

const (
	RootPath Path = iota
	ShowPath
	SignInPath
	WaitSignInPath
	CancelSignInPath
	SignOutPath
)

var Paths = []string{
	"/",
	"/show",
	"/auth/sign-in",
	"/auth/wait-sign-in",
	"/auth/cancel-sign-in",
	"/auth/sign-out",
}

func RedirectToWithInfo(context *gin.Context, status int, path Path, message string) {
	messages.AddFlashMessage(context, messages.InfoMessage, message)
	context.Redirect(status, Paths[path])
}

func RedirectToWithError(context *gin.Context, status int, path Path, message string) {
	messages.AddFlashMessage(context, messages.ErrorMessage, message)
	context.Redirect(status, Paths[path])
}

func RedirectTo(context *gin.Context, status int, path Path) {
	context.Redirect(status, Paths[path])
}
