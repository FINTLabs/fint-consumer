package generate

const LINKER_TEMPLATE = `package no.fint.consumer.models.{{ modelPkg .Package  }}{{ ToLower .Name }};

import no.fint.model.resource.Link;
import {{ resourcePkg .Package }}.{{ .Name }}Resource;
import {{ resourcePkg .Package }}.{{ .Name }}Resources;
import no.fint.relations.FintLinker;
import org.springframework.stereotype.Component;

import java.util.Collection;
import java.util.stream.IntStream;
import java.util.stream.Stream;

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
        return toResources(collection.stream(), 0, 0, collection.size());
    }

    @Override
    public {{ .Name }}Resources toResources(Stream<{{ .Name }}Resource> stream, int offset, int size, int totalItems) {
        {{ .Name }}Resources resources = new {{ .Name }}Resources();
        stream.map(this::toResource).forEach(resources::addResource);
        addPagination(resources, offset, size, totalItems);
        return resources;
    }

    @Override
    public String getSelfHref({{ $.Name }}Resource {{ ToLower $.Name  }}) {
        return getAllSelfHrefs({{ ToLower $.Name  }}).findFirst().orElse(null);
    }

    @Override
    public Stream<String> getAllSelfHrefs({{ $.Name }}Resource {{ ToLower $.Name  }}) {
        Stream.Builder<String> builder = Stream.builder();
        {{ range $i, $ident := .Identifiers -}}
        if (!isNull({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}()) && !isEmpty({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi())) {
            builder.add(createHrefWithId({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi(), "{{ ToLower $ident.Name }}"));
        }
        {{ end }}
        return builder.build();
    }

    int[] hashCodes({{ $.Name }}Resource {{ ToLower $.Name }}) {
        IntStream.Builder builder = IntStream.builder();
        {{ range $i, $ident := .Identifiers -}}
        if (!isNull({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}()) && !isEmpty({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi())) {
            builder.add({{ ToLower $.Name  }}.get{{ ToTitle $ident.Name }}().getIdentifikatorverdi().hashCode());
        }
        {{ end }}
        return builder.build().toArray();
    }

}

`
