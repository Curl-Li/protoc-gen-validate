package shared

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

func RegisterFunctions(tpl *template.Template, params pgs.Parameters) {
	tpl.Funcs(map[string]interface{}{
		"disabled":  Disabled,
		"ignored":   Ignored,
		"required":  RequiredOneOf,
		"oneofRule": OneOfRule,
		"context":   rulesContext,
		"render":    Render(tpl),
		"has":       Has,
		"needs":     Needs,
		"fileneeds": FileNeeds,
	})
}
