package templates

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

type WrapTemplate struct {
	pgs.Template
	*pgs.ModuleBase
	Format bool
}

func (wt *WrapTemplate) Execute(w io.Writer, data interface{}) error {
	if !wt.Format {
		return wt.Template.Execute(w, data)
	}
	out := &bytes.Buffer{}
	if err := wt.Template.Execute(out, data); err != nil {
		return err
	}
	if err := wt.formatCode(out, w); err != nil {
		wt.Logf("format code failed, %v", err)
		// fallback
		return wt.Template.Execute(w, data)
	}
	return nil
}

func (wt *WrapTemplate) formatCode(r io.Reader, w io.Writer) error {
	tpl, ok := wt.Template.(*template.Template)
	if !ok {
		return errors.New("pgs template is not *template.Template")
	}
	formatter := CodeFormatterFor(tpl)
	if formatter == nil {
		return fmt.Errorf("unsupported formatter for code: %s", tpl.Name())
	}
	return formatter(r, w)
}
