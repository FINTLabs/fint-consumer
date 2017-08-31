package setup




const CONSUMER_PROPS_TEMPLATE = `package no.fint.consumer.config;

import lombok.Getter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

@Getter
@Component
public class ConsumerProps {

{{ range $i, $model := . }}
	public static final String CACHE_INITIALDELAY_{{ ToUpper .Name }} = "${fint.consumer.cache.initialDelay.{{ .Name }}:{{ GetInitialRate $i }}";
	public static final String CACHE_FIXEDRATE_{{ ToUpper .Name }} = "${fint.consumer.cache.fixedRate.{{ .Name }}:900000}";
{{end }}


        @Value("${fint.events.orgIds:fint.no}")
        private String[] orgs;

}
`
