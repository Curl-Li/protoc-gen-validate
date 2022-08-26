package java

const messageConstTpl = `{{- if .Rules }}{{ $r := .Rules -}}
{{- if and $r.Error .DefineErr }}
	private final RuntimeException {{ errorName . 0 }} = {{ error . $r.Error }};
{{- end -}}
{{- end -}}
`

const messageTpl = `{{ $f := .Field }}{{ $r := .Rules }}
	{{- if .MessageRules.GetSkip }}
			// skipping validation for {{ $f.Name }}
	{{- else -}}
		{{- if $r.GetRequired }}
			if ({{ hasAccessor . }}) {
				cn.spaceli.pgv.RequiredValidation.required({{ errorName . 0 }}, {{ accessor . }});
			} else {
				cn.spaceli.pgv.RequiredValidation.required({{ errorName . 0 }}, null);
			};
		{{- end -}}
		{{- if (isOfMessageType $f) }}
			// Validate {{ $f.Name }}
			if ({{ hasAccessor . }}) index.validatorFor({{ accessor . }}).assertValid({{ accessor . }});
		{{- end -}}
	{{- end -}}
`
