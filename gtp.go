package gtp

import (
	"fmt"
	"github.com/taskie/jc"
	htmltemplate "html/template"
	"io"
	"text/template"
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

func (gtp *Gtp) Run(r io.Reader, w io.Writer) error {
	if len(gtp.TemplateFilePaths) == 0 {
		return fmt.Errorf("Please specify template file path")
	}
	jco := jc.Jc{
		FromType: gtp.DataType,
	}
	var data interface{}
	err := jco.Decode(r, &data)
	if err != nil {
		return err
	}
	err = gtp.Execute(w, data)
	return err
}
