package setup

const LINKMAPPER_TEMPLATE = `package no.fint.consumer.config;

import no.fint.consumer.utils.RestEndpoints;
import java.util.Map;
import com.google.common.collect.ImmutableMap;

{{ if .Package -}}
import no.fint.model.{{.Component}}.{{.Package}}.*;
{{- else -}}
import no.fint.model.{{.Component}}.*;
{{- end }}

public class LinkMapper {

    public static Map<String, String> linkMapper(String contextPath) {
        return ImmutableMap.<String,String>builder()
        {{- range $i, $model := .Models }}
            .put({{ ToTitle .Name }}.class.getName(), contextPath + RestEndpoints.{{ ToUpper .Name }})
        {{- end }}
        {{- range $i, $assoc := .Assocs }}
            .put("{{ .TargetPackage }}.{{ .Target }}", "/{{ ToUri .TargetPackage }}/{{ ToLower .Target }}")
        {{- end }}
            /* .put(TODO,TODO) */
            .build();
    }

}
`
