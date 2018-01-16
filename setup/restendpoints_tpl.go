package setup

const RESTENDPOINTS_TEMPLATE = `package no.fint.consumer.utils;

public enum RestEndpoints {
    ;

    public static final String ADMIN = "/admin";
{{- range $i, $model := . }}
	public static final String {{ ToUpper .Name }} = "/{{ ToLower .Name }}";
{{- end }}

}
`
