package java

const msgTpl = `
{{ if not (ignored .) -}}
	/**
	 * Validates {@code {{ simpleName . }}} protobuf objects.
	 */
	public static class {{ simpleName . }}Validator implements cn.spaceli.pgv.ValidatorImpl<{{ qualifiedName . }}> {
		{{- template "msgInner" . -}}
	}
{{- end -}}
`

const msgInnerTpl = `
	{{ $ctx := . }}
	{{- range .NonOneOfFields }}
		{{ renderConstants (context $ctx .) }}
	{{ end }}
	{{ template "oneOfConst" . }}

	public void assertValid({{ qualifiedName . }} proto, cn.spaceli.pgv.ValidatorIndex index) throws RuntimeException {
	{{ if disabled . }}
		// Validate is disabled for {{ simpleName . }}
		return;
	{{- else -}}
	{{ range .NonOneOfFields -}}
		{{ render (context $ctx .) }}
	{{ end -}}
	{{ template "oneOf" . }}
	{{- end }}
	}
`
