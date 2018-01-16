package generate


const RESOURCE_ASSEMBLER_TEMPLATE = `package no.fint.consumer.models.{{ ToLower .Name  }};

import {{ .Package }}.{{ .Name }};
import no.fint.model.relation.FintResource;
import no.fint.relations.FintResourceAssembler;
import no.fint.relations.FintResourceSupport;
import org.springframework.stereotype.Component;

@Component
public class {{ .Name }}Assembler extends FintResourceAssembler<{{ .Name }}> {

    public {{ .Name }}Assembler() {
        super({{ .Name }}Controller.class);
    }

{{ range $i, $ident := .Identifiers -}}
  {{ if not $ident.Optional }}
    @Override
    public FintResourceSupport assemble({{ $.Name }} {{ ToLower $.Name  }} , FintResource<{{ $.Name }}> fintResource) {
        return createResourceWithId({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi(), fintResource, "{{ $ident.Name }}");
    }
    {{ break }}
  {{ end }}
{{ end }}
    
}

`