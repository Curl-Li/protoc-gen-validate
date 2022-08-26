package java

const requiredTpl = `{{ $f := .Field }}
	{{- if .Rules.GetRequired }}
		if ({{ hasAccessor . }}) {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName . 0 }}, {{ accessor . }});
		} else {
			cn.spaceli.pgv.RequiredValidation.required({{ errorName . 0 }}, null);
		};
	{{- end -}}
`
