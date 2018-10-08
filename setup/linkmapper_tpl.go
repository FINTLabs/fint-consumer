package setup

const LINKMAPPER_TEMPLATE = `package no.fint.consumer.config;

import no.fint.consumer.utils.RestEndpoints;
import java.util.Map;
import com.google.common.collect.ImmutableMap;

import no.fint.model.{{.Component}}.{{.Package}}.*;
import no.fint.model.felles.*;

public class LinkMapper {

	public static Map<String, String> linkMapper(String contextPath) {
		return ImmutableMap.<String,String>builder()
		{{- range $i, $model := .Models }}
			.put({{ ToTitle .Name }}.class.getName(), contextPath + RestEndpoints.{{ ToUpper .Name }})
		{{- end }}
			/* .put(TODO,TODO) */
			.build();
	}

}
`
