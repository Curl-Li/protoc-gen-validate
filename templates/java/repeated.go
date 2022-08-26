package java

const repeatedConstTpl = `{{ renderConstants (.Elem "" "") }}
{{ $ctx := . }}{{ range $index, $r := .Rules.Rules -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const repeatedTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetIgnoreEmpty }}
		if ( !{{ accessor $ctx }}.isEmpty() ) {
{{- end -}}
{{- if $r.GetMinItems }}
			cn.spaceli.pgv.RepeatedValidation.minItems({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMinItems }});
{{- end -}}
{{- if $r.GetMaxItems }}
			cn.spaceli.pgv.RepeatedValidation.maxItems({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMaxItems }});
{{- end -}}
{{- if $r.GetUnique }}
			cn.spaceli.pgv.RepeatedValidation.unique({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end }}
{{- if $r.GetItems }}
			cn.spaceli.pgv.RepeatedValidation.forEach({{ accessor $ctx }}, item -> {
				{{ render ($ctx.ElemWithErrIndex "item" "" $index) }}
			});
{{- end }}
{{- if $r.GetIgnoreEmpty }}
		}
{{- end -}}
{{- end -}}
`
