package setup

const LINKMAPPER_TEMPLATE = `package no.fint.consumer.config;

import no.fint.consumer.utils.RestEndpoints;
import java.util.Map;
import com.google.common.collect.ImmutableMap;

public class LinkMapper {

	public static Map<String, String> linkMapper(String contextPath) {
		return ImmutableMap.<String,String>builder()
		{{- range $i, $model := . }}
			.put({{ ToTitle .Name }}.class.getName(), contextPath + RestEndpoints.{{ ToUpper .Name }})
		{{- end }}
			/* .put(TODO,TODO) */
			.build();
	}

}
`
