package generate


const RESOURCE_ASSEMBLER_TEMPLATE = `package no.fint.consumer.{{ ToLower .Name  }};

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

    @Override
    public FintResourceSupport assemble({{ .Name }} {{ ToLower .Name  }} , FintResource<{{ .Name }}> fintResource) {
        return createResourceWithId({{ ToLower .Name  }}.get***fixme***().getIdentifikatorverdi(), fintResource, "***fixme***");
    }
}

`