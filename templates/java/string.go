package java

const stringConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.In }}
	private final {{ javaTypeFor $ctx }}[] {{ constantName $ctx "In" }} = new {{ javaTypeFor $ctx }}[]{
		{{- range $r.In -}}
			"{{- sprintf "%v" . -}}",
		{{- end -}}
	};
{{- end -}}
{{- if $r.NotIn }}
	private final {{ javaTypeFor $ctx }}[] {{ constantName $ctx "NotIn" }} = new {{ javaTypeFor $ctx }}[]{
		{{- range $r.NotIn -}}
			"{{- sprintf "%v" . -}}",
		{{- end -}}
	};
{{- end -}}
{{- if $r.Pattern }}
	com.google.re2j.Pattern {{ constantName $ctx "Pattern" }} = com.google.re2j.Pattern.compile({{ javaStringEscape $r.GetPattern }});
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const stringTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetIgnoreEmpty }}
		if ( !{{ accessor $ctx }}.isEmpty() ) {
{{- end -}}
{{- if $r.Const }}
			cn.spaceli.pgv.ConstantValidation.constant({{ errorName $ctx $index }}, {{ accessor $ctx }}, "{{ $r.GetConst }}");
{{- end -}}
{{- if $r.In }}
			cn.spaceli.pgv.CollectiveValidation.in({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "In" }});
{{- end -}}
{{- if $r.NotIn }}
			cn.spaceli.pgv.CollectiveValidation.notIn({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "NotIn" }});
{{- end -}}
{{- if $r.Len }}
			cn.spaceli.pgv.StringValidation.length({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetLen }});
{{- end -}}
{{- if $r.MinLen }}
			cn.spaceli.pgv.StringValidation.minLength({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMinLen }});
{{- end -}}
{{- if $r.MaxLen }}
			cn.spaceli.pgv.StringValidation.maxLength({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMaxLen }});
{{- end -}}
{{- if $r.LenBytes }}
			cn.spaceli.pgv.StringValidation.lenBytes({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetLenBytes }});
{{- end -}}
{{- if $r.MinBytes }}
			cn.spaceli.pgv.StringValidation.minBytes({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMinBytes }});
{{- end -}}
{{- if $r.MaxBytes }}
			cn.spaceli.pgv.StringValidation.maxBytes({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMaxBytes }});
{{- end -}}
{{- if $r.Pattern }}
			cn.spaceli.pgv.StringValidation.pattern({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Pattern" }});
{{- end -}}
{{- if $r.Prefix }}
			cn.spaceli.pgv.StringValidation.prefix({{ errorName $ctx $index }}, {{ accessor $ctx }}, "{{ $r.GetPrefix }}");
{{- end -}}
{{- if $r.Contains }}
			cn.spaceli.pgv.StringValidation.contains({{ errorName $ctx $index }}, {{ accessor $ctx }}, "{{ $r.GetContains }}");
{{- end -}}
{{- if $r.NotContains }}
			cn.spaceli.pgv.StringValidation.notContains({{ errorName $ctx $index }}, {{ accessor $ctx }}, "{{ $r.GetNotContains }}");
{{- end -}}
{{- if $r.Suffix }}
			cn.spaceli.pgv.StringValidation.suffix({{ errorName $ctx $index }}, {{ accessor $ctx }}, "{{ $r.GetSuffix }}");
{{- end -}}
{{- if $r.GetEmail }}
			cn.spaceli.pgv.StringValidation.email({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetAddress }}
			cn.spaceli.pgv.StringValidation.address({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetHostname }}
			cn.spaceli.pgv.StringValidation.hostName({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIp }}
			cn.spaceli.pgv.StringValidation.ip({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIpv4 }}
			cn.spaceli.pgv.StringValidation.ipv4({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIpv6 }}
			cn.spaceli.pgv.StringValidation.ipv6({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetUri }}
			cn.spaceli.pgv.StringValidation.uri({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetUriRef }}
			cn.spaceli.pgv.StringValidation.uriRef({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetUuid }}
			cn.spaceli.pgv.StringValidation.uuid({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIgnoreEmpty }}
		}
{{- end -}}
{{- end -}}
`
