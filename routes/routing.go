package routes

import "github.com/gin-gonic/gin"

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

func RedirectTo(context *gin.Context, status int, path Path) {
	context.Redirect(status, Paths[path])
}
