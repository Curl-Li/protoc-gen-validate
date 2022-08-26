package java

const mapConstTpl = `
{{- if or (ne (.Elem "" "").Typ "none") (ne (.Key "" "").Typ "none") }}
		{{ renderConstants (.Key "key" "Key") }}
		{{ renderConstants (.Elem "value" "Value") }}
{{- end -}}
{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const mapTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetIgnoreEmpty }}
			if ( !{{ accessor $ctx }}.isEmpty() ) {
{{- end -}}
{{- if $r.GetMinPairs }}
			cn.spaceli.pgv.MapValidation.min({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMinPairs }});
{{- end -}}
{{- if $r.GetMaxPairs }}
			cn.spaceli.pgv.MapValidation.max({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMaxPairs }});
{{- end -}}
{{- if $r.GetNoSparse }}
			cn.spaceli.pgv.MapValidation.noSparse({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if and $r.GetKeys (ne ($ctx.Key "" "").Typ "none") }}
			cn.spaceli.pgv.MapValidation.validateParts({{ accessor $ctx }}.keySet(), key -> {
				{{ render ($ctx.KeyWithErrIndex "key" "Key" $index) }}
			});
{{- end -}}
{{ if and $r.GetValues (ne ($ctx.Key "" "").Typ "none") }}
			cn.spaceli.pgv.MapValidation.validateParts({{ accessor $ctx }}.values(), value -> {
				{{ render ($ctx.ElemWithErrIndex "value" "Value" $index) }}
			});
{{- end -}}
{{- if $r.GetIgnoreEmpty }}
			}
{{- end -}}
{{- end -}}
`
