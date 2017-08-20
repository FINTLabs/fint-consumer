package cs

const CLASS_TEMPLATE = `using System;
{{ if .Using }}
{{ range $u := .Using }}
using {{ $u }};
{{- end -}}
{{ end }}

namespace {{ .Namespace }}
{
	public {{- if .Abstract }} abstract {{- end }} class {{ .Name }} {{ if .Extends -}} : {{ .Extends }} {{ end -}}
	{
		{{- if .Relations }}
		{{ $c := sub (len .Relations) 1 -}}
		public enum Relasjonsnavn
        {
		{{- range $i, $rel := .Relations }}
			{{ $rel }}{{if ne $i $c }},{{ end -}}
		{{ end }}
        }
        {{ end }}
	{{ if .Attributes }}
		{{ range $att := .Attributes -}}
			public {{ csType $att.Type}} {{ $att.Name }} { get; set; }
		{{ end -}}
	{{ end }}
	}
}
`
