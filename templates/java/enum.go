package java

const enumConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.In }}
	private final {{ javaTypeFor . }}[] {{ constantName . "In" }} = new {{ javaTypeFor . }}[]{
		{{- range $r.In }}
		{{ javaTypeFor $ctx }}.forNumber({{- sprintf "%v" . -}}),
		{{- end }}
	};
{{- end -}}
{{- if $r.NotIn }}
	private final {{ javaTypeFor . }}[] {{ constantName . "NotIn" }} = new {{ javaTypeFor . }}[]{
		{{- range $r.NotIn }}
		{{ javaTypeFor $ctx }}.forNumber({{- sprintf "%v" . -}}),
		{{- end }}
	};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName . 0 }} = {{ error . $r.Error }};
{{- end -}}
`

const enumTpl = `{{ $f := .Field }}{{ $r := .Rules -}}
{{- if $r.Const }}
			cn.spaceli.pgv.ConstantValidation.constant({{ errorName . 0 }}, {{ accessor . }}, 
				{{ javaTypeFor . }}.forNumber({{ $r.GetConst }}));
{{- end -}}
{{- if $r.GetDefinedOnly }}
			cn.spaceli.pgv.EnumValidation.definedOnly({{ errorName . 0 }}, {{ accessor . }});
{{- end -}}
{{- if $r.In }}
			cn.spaceli.pgv.CollectiveValidation.in({{ errorName . 0 }}, {{ accessor . }}, {{ constantName . "In" }});
{{- end -}}
{{- if $r.NotIn }}
			cn.spaceli.pgv.CollectiveValidation.notIn({{ errorName . 0 }}, {{ accessor . }}, {{ constantName . "NotIn" }});
{{- end -}}
`
