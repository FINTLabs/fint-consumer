package generate

const CACHE_SERVICE_TEMPLATE = `package no.fint.consumer.models.{{ ToLower .Name  }};

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;

import lombok.extern.slf4j.Slf4j;

import no.fint.cache.CacheService;
import no.fint.consumer.config.Constants;
import no.fint.consumer.config.ConsumerProps;
import no.fint.consumer.event.ConsumerEventUtil;
import no.fint.event.model.Event;
import no.fint.model.felles.kompleksedatatyper.Identifikator;
import no.fint.model.relation.FintResource;
import no.fint.model.resource.Link;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

import {{ .Package }}.{{ .Name }};
import {{ resourcePkg .Package }}.{{ .Name }}Resource;
import {{ GetActionPackage .Package }};

@Slf4j
@Service
public class {{ .Name }}CacheService extends CacheService<{{ .Name }}Resource> {

    public static final String MODEL = {{ .Name }}.class.getSimpleName().toLowerCase();

    @Value("${fint.consumer.compatibility.fintresource:true}")
    private boolean fintResourceCompatibility;

    @Autowired
    private ConsumerEventUtil consumerEventUtil;

    @Autowired
    private ConsumerProps props;

    public {{ .Name }}CacheService() {
        super(MODEL, {{ GetAction .Package }}.GET_ALL_{{ ToUpper .Name }});
    }

    @PostConstruct
    public void init() {
        Arrays.stream(props.getOrgs()).forEach(this::createCache);
    }

    @Scheduled(initialDelayString = ConsumerProps.CACHE_INITIALDELAY_{{ ToUpper .Name }}, fixedRateString = ConsumerProps.CACHE_FIXEDRATE_{{ ToUpper .Name }})
    public void populateCacheAll() {
        Arrays.stream(props.getOrgs()).forEach(this::populateCache);
    }

    public void rebuildCache(String orgId) {
		flush(orgId);
		populateCache(orgId);
	}

    private void populateCache(String orgId) {
		log.info("Populating {{ .Name }} cache for {}", orgId);
        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction .Package}}.GET_ALL_{{ ToUpper .Name }}, Constants.CACHE_SERVICE);
        consumerEventUtil.send(event);
    }

{{ range $i, $ident := .Identifiers }}
    public Optional<{{ $.Name }}Resource> get{{ $.Name }}By{{ ToTitle $ident.Name }}(String orgId, String {{ $ident.Name }}) {
        return getOne(orgId, (resource) -> Optional
                .ofNullable(resource)
                .map({{ $.Name }}Resource::get{{ ToTitle $ident.Name }})
                .map(Identifikator::getIdentifikatorverdi)
                .map(_id -> _id.equals({{ $ident.Name }}))
                .orElse(false));
    }
{{ end }}

	@Override
    public void onAction(Event event) {
        if (fintResourceCompatibility && !event.getData().isEmpty() && event.getData().get(0) instanceof FintResource) {
            log.info("Compatibility: Converting FintResource<{{.Name}}Resource> to {{.Name}}Resource ...");
            ObjectMapper objectMapper = new ObjectMapper();
            objectMapper.configure(SerializationFeature.FAIL_ON_EMPTY_BEANS, false);
            List<FintResource<{{.Name}}Resource>> original = objectMapper.convertValue(event.getData(), new TypeReference<List<FintResource<{{.Name}}Resource>>>() {
            });
            List<{{.Name}}Resource> replacement = original.stream().map(fintResource -> {
                {{.Name}}Resource resource = fintResource.getResource();
                fintResource.getRelations().forEach(relation -> resource.addLink(relation.getRelationName(), Link.with(relation.getLink())));
                return resource;
            }).collect(Collectors.toList());
            event.setData(replacement);
        }
        update(event, new TypeReference<List<{{ .Name }}Resource>>() {
        });
    }
}
`
