package java

const bytesConstTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.Const }}
	private final com.google.protobuf.ByteString {{ constantName $ctx "Const" }} = com.google.protobuf.ByteString.copyFrom({{ byteArrayLit $r.GetConst }});
{{- end -}}
{{- if $r.In }}
	private final com.google.protobuf.ByteString[] {{ constantName $ctx "In" }} = new com.google.protobuf.ByteString[]{
		{{- range $r.In }}
		com.google.protobuf.ByteString.copyFrom({{ byteArrayLit . }}),
		{{- end }}
	};
{{- end -}}
{{- if $r.NotIn }}
	private final com.google.protobuf.ByteString[] {{ constantName $ctx "NotIn" }} = new com.google.protobuf.ByteString[]{
		{{- range $r.NotIn }}
		com.google.protobuf.ByteString.copyFrom({{ byteArrayLit . }}),
		{{- end }}
	};
{{- end -}}
{{- if $r.Pattern }}
	private final com.google.re2j.Pattern {{ constantName $ctx "Pattern" }} = com.google.re2j.Pattern.compile({{ javaStringEscape $r.GetPattern }});
{{- end -}}
{{- if $r.Prefix }}
	private final byte[] {{ constantName $ctx "Prefix" }} = {{ byteArrayLit $r.GetPrefix }};
{{- end -}}
{{- if $r.Contains }}
	private final byte[] {{ constantName $ctx "Contains" }} = {{ byteArrayLit $r.GetContains }};
{{- end -}}
{{- if $r.Suffix }}
	private final byte[] {{ constantName $ctx "Suffix" }} = {{ byteArrayLit $r.GetSuffix }};
{{- end -}}
{{- if and $r.Error $ctx.DefineErr }}
	private final RuntimeException {{ errorName $ctx $index }} = {{ error $ctx $r.GetError }};
{{- end -}}
{{- end -}}
`

const bytesTpl = `{{ $ctx := . }}{{ $f := .Field }}{{ range $index, $r := .Rules.Rules -}}
{{- if $r.GetIgnoreEmpty }}
			if ( !{{ accessor $ctx }}.isEmpty() ) {
{{- end -}}
{{- if $r.Const }}
			cn.spaceli.pgv.ConstantValidation.constant({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Const" }});
{{- end -}}
{{- if $r.Len }}
			cn.spaceli.pgv.BytesValidation.length({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetLen }});
{{- end -}}
{{- if $r.MinLen }}
			cn.spaceli.pgv.BytesValidation.minLength({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMinLen }});
{{- end -}}
{{- if $r.MaxLen }}
			cn.spaceli.pgv.BytesValidation.maxLength({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ $r.GetMaxLen }});
{{- end -}}
{{- if $r.Pattern }}
			cn.spaceli.pgv.BytesValidation.pattern({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Pattern" }});
{{- end -}}
{{- if $r.Prefix }}
			cn.spaceli.pgv.BytesValidation.prefix({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Prefix" }});
{{- end -}}
{{- if $r.Contains }}
			cn.spaceli.pgv.BytesValidation.contains({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Contains" }});
{{- end -}}
{{- if $r.Suffix }}
			cn.spaceli.pgv.BytesValidation.suffix({{ errorName $ctx $index }}, {{ accessor $ctx }}, {{ constantName $ctx "Suffix" }});
{{- end -}}
{{- if $r.GetIp }}
			cn.spaceli.pgv.BytesValidation.ip({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIpv4 }}
			cn.spaceli.pgv.BytesValidation.ipv4({{ errorName $ctx $index }}, {{ accessor $ctx }});
{{- end -}}
{{- if $r.GetIpv6 }}
			cn.spaceli.pgv.BytesValidation.ipv6({{ errorName $ctx $index }}, {{ accessor $ctx }});
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
