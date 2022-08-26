package java

const timestampConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.Const }}
		private final com.google.protobuf.Timestamp {{ constantName $ctx "Const" }} = {{ tsLit $r.GetConst }};
{{- end -}}
{{- if $r.Lt }}
		private final com.google.protobuf.Timestamp {{ constantName $ctx "Lt" }} = {{ tsLit $r.GetLt }};
{{- end -}}
{{- if $r.Lte }}
		private final com.google.protobuf.Timestamp {{ constantName $ctx "Lte" }} = {{ tsLit $r.Lte }};
{{- end -}}
{{- if $r.Gt }}
		private final com.google.protobuf.Timestamp {{ constantName $ctx "Gt" }} = {{ tsLit $r.GetGt }};
{{- end -}}
{{- if $r.Gte }}
		private final com.google.protobuf.Timestamp {{ constantName $ctx "Gte" }} = {{ tsLit $r.GetGte }};
{{- end -}}
{{- if $r.Within }}
		private final com.google.protobuf.Duration {{ constantName $ctx "Within" }} = {{ durLit $r.GetWithin }};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const timestampTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
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
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.range({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ if $r.Lt }}{{ constantName $ctx "Lt" }}{{ else }}null{{ end }}, {{ if $r.Lte }}{{ constantName $ctx "Lte" }}{{ else }}null{{ end }}, {{ if $r.Gt }}{{ constantName $ctx "Gt" }}{{ else }}null{{ end }}, {{ if $r.Gte }}{{ constantName $ctx "Gte" }}{{ else }}null{{ end }}, com.google.protobuf.util.Timestamps.comparator());
{{- else -}}
{{- if $r.Lt }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.lessThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lt" }}, com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- if $r.Lte }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.lessThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Lte" }}, com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- if $r.Gt }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.greaterThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gt" }}, com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- if $r.Gte }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.greaterThanOrEqual({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Gte" }}, com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- end -}}
{{- if $r.LtNow }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.lessThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, cn.spaceli.pgv.TimestampValidation.currentTimestamp(), com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- if $r.GtNow }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.ComparativeValidation.greaterThan({{ errorName $ctx $index }}, {{ accessor $ctx }}, cn.spaceli.pgv.TimestampValidation.currentTimestamp(), com.google.protobuf.util.Timestamps.comparator());
{{- end -}}
{{- if $r.Within }}
			if ({{ hasAccessor $ctx }}) cn.spaceli.pgv.TimestampValidation.within({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Within" }}, cn.spaceli.pgv.TimestampValidation.currentTimestamp());
{{- end -}}
{{- end -}}
`
