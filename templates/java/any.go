package java

const anyConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules }}
{{- if $r.In }}
	private final String[] {{ constantName $ctx "In" }} = new String[]{
		{{- range $r.In }}
		"{{ . }}",
		{{- end }}
	};
{{- end -}}
{{- if $r.NotIn }}
	private final String[] {{ constantName $ctx "NotIn" }} = new String[]{
		{{- range $r.NotIn }}
		"{{ . }}",
		{{- end }}
	};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const anyTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules }}
	{{- if $r.GetRequired }}
		if ({{ hasAccessor $ctx }}) {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName $ctx $index }}, {{ accessor $ctx }});
		} else {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName $ctx $index }}, null);
		};
	{{- end -}}
	{{- if $r.In }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.CollectiveValidation.in({{ errorName $ctx $index }}, {{ accessor $ctx }}.getTypeUrl(), {{ constantName $ctx "In" }});
	{{- end -}}
	{{- if $r.NotIn }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.CollectiveValidation.notIn({{ errorName $ctx $index }}, {{ accessor $ctx }}.getTypeUrl(), {{ constantName $ctx "NotIn" }});
	{{- end -}}
	{{- end -}}
`
