package java

const numConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.Const }}
	private final {{ javaTypeFor $ctx}} {{ constantName $ctx "Const" }} = {{ $r.GetConst }}{{ javaTypeLiteralSuffixFor $ctx }};
{{- end -}}
{{- if $r.Lt }}
	private final {{ javaTypeFor $ctx}} {{ constantName $ctx "Lt" }} = {{ $r.GetLt }}{{ javaTypeLiteralSuffixFor $ctx }};
{{- end -}}
{{- if $r.Lte }}
	private final {{ javaTypeFor $ctx}} {{ constantName $ctx "Lte" }} = {{ $r.GetLte }}{{ javaTypeLiteralSuffixFor $ctx }};
{{- end -}}
{{- if $r.Gt }}
	private final {{ javaTypeFor $ctx}} {{ constantName $ctx "Gt" }} = {{ $r.GetGt }}{{ javaTypeLiteralSuffixFor $ctx }};
{{- end -}}
{{- if $r.Gte }}
	private final {{ javaTypeFor $ctx}} {{ constantName $ctx "Gte" }} = {{ $r.GetGte }}{{ javaTypeLiteralSuffixFor $ctx }};
{{- end -}}
{{- if $r.In }}
	private final {{ javaTypeFor $ctx }}[] {{ constantName $ctx "In" }} = new {{ javaTypeFor $ctx }}[]{
		{{- range $r.In -}}
			{{- sprintf "%v" . -}}{{ javaTypeLiteralSuffixFor $ }},
		{{- end -}}
	};
{{- end -}}
{{- if $r.NotIn }}
	private final {{ javaTypeFor $ctx }}[] {{ constantName $ctx "NotIn" }} = new {{ javaTypeFor $ctx }}[]{
		{{- range $r.NotIn -}}
			{{- sprintf "%v" . -}}{{ javaTypeLiteralSuffixFor $ }},
		{{- end -}}
	};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}`

const numTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetIgnoreEmpty }}
		if ( {{ accessor $ctx }} != 0 ) {
{{- end -}}
{{- if $r.Const }}
			cn.spaceli.pgv.ConstantValidation.constant({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Const" }});
{{- end -}}
{{- if and (or $r.Lt $r.Lte) (or $r.Gt $r.Gte)}}
			cn.spaceli.pgv.ComparativeValidation.range({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ if $r.Lt }}{{ constantName $ctx "Lt" }}{{ else }}null{{ end }}, {{ if $r.Lte }}{{ constantName $ctx "Lte" }}{{ else }}null{{ end }}, {{ if $r.Gt }}{{ constantName $ctx "Gt" }}{{ else }}null{{ end }}, {{ if $r.Gte }}{{ constantName $ctx "Gte" }}{{ else }}null{{ end }}, java.util.Comparator.naturalOrder());
{{- else -}}
{{- if $r.Lt }}
			cn.spaceli.pgv.ComparativeValidation.lessThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lt" }}, java.util.Comparator.naturalOrder());
{{- end -}}
{{- if $r.Lte }}
			cn.spaceli.pgv.ComparativeValidation.lessThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lte" }}, java.util.Comparator.naturalOrder());
{{- end -}}
{{- if $r.Gt }}
			cn.spaceli.pgv.ComparativeValidation.greaterThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gt" }}, java.util.Comparator.naturalOrder());
{{- end -}}
{{- if $r.Gte }}
			cn.spaceli.pgv.ComparativeValidation.greaterThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gte" }}, java.util.Comparator.naturalOrder());
{{- end -}}
{{- end -}}
{{- if $r.In }}
			cn.spaceli.pgv.CollectiveValidation.in({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "In" }});
{{- end -}}
{{- if $r.NotIn }}
			cn.spaceli.pgv.CollectiveValidation.notIn({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "NotIn" }});
{{- end -}}
{{- if $r.GetIgnoreEmpty }}
		}
{{- end -}}
{{- end -}}
`
