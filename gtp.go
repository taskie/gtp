package gtp

import (
	"fmt"
	htmltemplate "html/template"
	"io"
	"text/template"

	"github.com/taskie/jc"
)

var (
	Version  = "0.1.0-beta"
	Revision = ""
)

type Gtp struct {
	TemplateFilePaths []string
	TemplateType      string
	DataType          string
}

func (gtp *Gtp) Execute(w io.Writer, data interface{}) error {
	switch gtp.TemplateType {
	case "text":
		tpl, err := template.ParseFiles(gtp.TemplateFilePaths...)
		if err != nil {
			return err
		}
		err = tpl.Execute(w, data)
		return err
	case "html":
		tpl, err := htmltemplate.ParseFiles(gtp.TemplateFilePaths...)
		if err != nil {
			return err
		}
		err = tpl.Execute(w, data)
		return err
	default:
		return fmt.Errorf("Invalid template type: %s", gtp.TemplateType)
	}
}

func (gtp *Gtp) Run(w io.Writer, r io.Reader) error {
	if len(gtp.TemplateFilePaths) == 0 {
		return fmt.Errorf("Please specify template file path")
	}
	dec := jc.NewDecoder(r, gtp.DataType)
	var data interface{}
	err := dec.Decode(&data)
	if err != nil {
		return err
	}
	err = gtp.Execute(w, data)
	return err
}
