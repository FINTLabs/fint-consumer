package generate

const LINKER_TEMPLATE = `package no.fint.consumer.models.{{ modelPkg .Package  }}{{ ToLower .Name }};

import no.fint.model.resource.Link;
import {{ resourcePkg .Package }}.{{ .Name }}Resource;
import {{ resourcePkg .Package }}.{{ .Name }}Resources;
import no.fint.relations.FintLinker;
import org.springframework.stereotype.Component;

import java.util.Collection;

import static java.util.Objects.isNull;
import static org.springframework.util.StringUtils.isEmpty;


@Component
public class {{ .Name }}Linker extends FintLinker<{{ .Name }}Resource> {

    public {{ .Name }}Linker() {
        super({{ .Name }}Resource.class);
    }

    public void mapLinks({{.Name}}Resource resource) {
        super.mapLinks(resource);
    }

    @Override
    public {{ .Name }}Resources toResources(Collection<{{ .Name }}Resource> collection) {
        {{ .Name }}Resources resources = new {{ .Name }}Resources();
        collection.stream().map(this::toResource).forEach(resources::addResource);
        resources.addSelf(Link.with(self()));
        return resources;
    }

    @Override
    public String getSelfHref({{ $.Name }}Resource {{ ToLower $.Name  }}) {
        {{ range $i, $ident := .Identifiers -}}
        if (!isNull({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}()) && !isEmpty({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi())) {
            return createHrefWithId({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi(), "{{ ToLower $ident.Name }}");
        }
        {{ end }}
        return null;
    }
    
}

`
