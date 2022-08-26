package java

const durationConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.Const }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Const" }} = {{ durLit $r.GetConst }};
{{- end -}}
{{- if $r.Lt }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Lt" }} = {{ durLit $r.GetLt }};
{{- end -}}
{{- if $r.Lte }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Lte" }} = {{ durLit $r.GetLte }};
{{- end -}}
{{- if $r.Gt }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Gt" }} = {{ durLit $r.GetGt }};
{{- end -}}
{{- if $r.Gte }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Gte" }} = {{ durLit $r.GetGte }};
{{- end -}}
{{- if $r.In }}
		private final com.google.protobuf.Duration[] {{ constantName $ctx "In" }} = new com.google.protobuf.Duration[]{
			{{- range $r.In }}
			{{ durLit . }},
			{{- end }}
		};
{{- end -}}
{{- if $r.NotIn }}
		private final com.google.protobuf.Duration[] {{ constantName $ctx "NotIn" }} = new com.google.protobuf.Duration[]{
			{{- range $r.NotIn }}
			{{ durLit . }},
			{{- end }}
		};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const durationTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetRequired }}
		if ({{ hasAccessor $ctx }}) {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName $ctx $index }}, {{ accessor $ctx }});
		} else {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName $ctx $index }}, null);
		};
{{- end -}}
{{- if $r.Const }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ConstantValidation.constant({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Const" }});
{{- end -}}
{{- if and (or $r.Lt $r.Lte) (or $r.Gt $r.Gte)}}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.range({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ if $r.Lt }}{{ constantName $ctx "Lt" }}{{ else }}null{{ end }}, {{ if $r.Lte }}{{ constantName $ctx "Lte" }}{{ else }}null{{ end }}, {{ if $r.Gt }}{{ constantName $ctx "Gt" }}{{ else }}null{{ end }}, {{ if $r.Gte }}{{ constantName $ctx "Gte" }}{{ else }}null{{ end }}, com.google.protobuf.util.Durations.comparator());
{{- else -}}
{{- if $r.Lt }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.lessThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lt" }}, com.google.protobuf.util.Durations.comparator());
{{- end -}}
{{- if $r.Lte }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.lessThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lte" }}, com.google.protobuf.util.Durations.comparator());
{{- end -}}
{{- if $r.Gt }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.greaterThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gt" }}, com.google.protobuf.util.Durations.comparator());
{{- end -}}
{{- if $r.Gte }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.greaterThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gte" }}, com.google.protobuf.util.Durations.comparator());
{{- end -}}
{{- end -}}
{{- if $r.In }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.CollectiveValidation.in({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "In" }});
{{- end -}}
{{- if $r.NotIn }}
		if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.CollectiveValidation.notIn({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "NotIn" }});
{{- end -}}
{{- end -}}
`
