package java

const boolConstTpl = `
{{- if .Rules.Error }}
	private final RuntimeException {{ errorName . 0 }} = {{ error . .Rules.Error }};
{{- end -}}
`

const boolTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
			cn.spaceli.pgv.ConstantValidation.constant({{ errorName . 0 }}, {{ accessor . }}, {{ $r.GetConst }});
{{- end }}`
