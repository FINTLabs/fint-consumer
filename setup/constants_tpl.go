package setup


const CONSTANTS_TEMPLATE = `
package no.fint.consumer.config;

public enum Constants {
;

    public static final String COMPONENT = "{{ .Name }}";
    public static final String COMPONENT_CONSUMER = COMPONENT + " consumer";
    public static final String CACHE_SERVICE = "CACHE_SERVICE";

}
`