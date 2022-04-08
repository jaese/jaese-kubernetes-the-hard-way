package utils

import (
	"bytes"
	"text/template"
	"os"

	"github.com/codeskyblue/go-sh"
)

func MustRun(session *sh.Session) {
	session.ShowCMD = true
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(); err != nil {
		panic(err)
	}
}

func MustRunWithStringInput(input string, session *sh.Session) {
	session.ShowCMD = true
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.SetInput(input)

	if err := session.Run(); err != nil {
		panic(err)
	}
}

func TextTemplateExecuteString(tmpl *template.Template, data any) (string, error) {
	var b bytes.Buffer
	if err := tmpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
