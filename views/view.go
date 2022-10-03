package views

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"

	"mojo-auth-test-1/messages"

	"github.com/gin-gonic/gin"
)

func NewView(layout string, files ...string) *View {
	files = append(layoutFiles(), files...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Render(context *gin.Context, data interface{}) error {
	infoFlashes := messages.GetFlashMessages(context, messages.InfoMessage)
	hasInfoFlashes := infoFlashes != nil && len(infoFlashes) > 0
	if hasInfoFlashes {
		for i := range infoFlashes {
			infoFlashes[i] = fmt.Sprintf("<div>%s</div>", infoFlashes[i])
		}
	}
	errorFlashes := messages.GetFlashMessages(context, messages.ErrorMessage)
	hasErrorFlashes := errorFlashes != nil && len(errorFlashes) > 0
	if hasErrorFlashes {
		for i := range errorFlashes {
			errorFlashes[i] = fmt.Sprintf("<div>%s</div>", errorFlashes[i])
		}
	}

	var dataMap gin.H
	if data == nil {
		dataMap = gin.H{
			"hasInfoFlashes":  hasInfoFlashes,
			"infoFlashes":     template.HTML(strings.Join(infoFlashes, "\n")),
			"hasErrorFlashes": hasErrorFlashes,
			"errorFlashes":    template.HTML(strings.Join(errorFlashes, "\n")),
		}
		data = dataMap
	} else {
		dataMap = data.(gin.H)
		dataMap["hasInfoFlashes"] = hasInfoFlashes
		dataMap["infoFlashes"] = template.HTML(strings.Join(infoFlashes, "\n"))
		dataMap["hasErrorFlashes"] = hasErrorFlashes
		dataMap["errorFlashes"] = template.HTML(strings.Join(errorFlashes, "\n"))
	}

	return v.Template.ExecuteTemplate(context.Writer, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob("templates/layouts/*.html")
	if err != nil {
		panic(err)
	}
	return files
}
