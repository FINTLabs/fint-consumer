package setup

const CONSTANTS_TEMPLATE = `
package no.fint.consumer.config;

public enum Constants {
;

    public static final String COMPONENT = "{{ .Name }}";
    public static final String COMPONENT_CONSUMER = COMPONENT + " consumer";
    public static final String CACHE_SERVICE = "CACHE_SERVICE";

{{ range $i, $model := .Models }}    
    public static final String CACHE_INITIALDELAY_{{ ToUpper .Name }} = "${fint.consumer.cache.initialDelay.{{ .Name }}:{{ GetInitialRate $i }}}";
    public static final String CACHE_FIXEDRATE_{{ ToUpper .Name }} = "${fint.consumer.cache.fixedRate.{{ .Name }}:900000}";
{{end }}    

}
`
