package java

const oneOfConstTpl = `
{{ $msg := . }}
{{ range .RealOneOfs }}
{{ $r := oneofRule .}}
{{- if $r.GetError }}
	private final RuntimeException {{ errorOneOfRequiredName $msg . }} = {{ error (context $msg (index .Fields 0)) $r.GetError }};
{{- end -}}
{{ range .Fields }}{{ renderConstants (context $msg .) }}{{- end -}}
{{- end -}}

`

const oneOfTpl = `
{{ $msg := . }}
{{ range .RealOneOfs }}
{{ $field := . }}
{{ $r := oneofRule .}}
switch (proto.get{{camelCase .Name }}Case()) {
	{{ range .Fields -}}
	case {{ oneof . }}:
		{{ render (context $msg .) }}
		break;
	{{ end -}}
	{{- if $r.GetRequired }}
	default: 
		cn.spaceli.pgv.RequiredValidation.required({{ errorOneOfRequiredName $msg $field }}, null);
	{{- end }}
}
{{- end -}}
			
`
