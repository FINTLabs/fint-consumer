package generate

const LINKER_TEMPLATE = `package no.fint.consumer.models.{{ ToLower .Name  }};

import {{ resourcePkg .Package }}.{{ .Name }}Resource;
import no.fint.relations.FintLinker;
import org.springframework.stereotype.Component;

@Component
public class {{ .Name }}Linker extends FintLinker<{{ .Name }}Resource> {

    public {{ .Name }}Linker() {
        super({{ .Name }}Resource.class);
    }

{{ range $i, $ident := .Identifiers -}}
  {{ if not $ident.Optional }}
    @Override
    public String getSelfHref({{ $.Name }}Resource {{ ToLower $.Name  }}) {
        return createHrefWithId({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi(), "{{ ToLower $ident.Name }}");
    }
    {{/* This only works in go1.10rc1 -- sorry :( */}}
    {{ break }}
  {{ end }}
{{ end }}
    
}

`
