package generate

const CACHE_SERVICE_TEMPLATE = `package no.fint.consumer.models.{{ modelPkg .Package  }}{{ ToLower .Name }};

import com.fasterxml.jackson.databind.JavaType;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;

import lombok.extern.slf4j.Slf4j;

import no.fint.cache.CacheService;
import no.fint.cache.model.CacheObject;
import no.fint.consumer.config.Constants;
import no.fint.consumer.config.ConsumerProps;
import no.fint.consumer.event.ConsumerEventUtil;
import no.fint.event.model.Event;
import no.fint.event.model.ResponseStatus;
import no.fint.model.felles.kompleksedatatyper.Identifikator;
import no.fint.relations.FintResourceCompatibility;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
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
    private boolean checkFintResourceCompatibility;

    @Autowired
    private FintResourceCompatibility fintResourceCompatibility;

    @Autowired
    private ConsumerEventUtil consumerEventUtil;

    @Autowired
    private ConsumerProps props;

    @Autowired
    private {{ .Name }}Linker linker;

    private JavaType javaType;

    private ObjectMapper objectMapper;

    public {{ .Name }}CacheService() {
        super(MODEL, {{ GetAction .Package }}.GET_ALL_{{ ToUpper .Name }}, {{ GetAction .Package }}.UPDATE_{{ ToUpper .Name }});
        objectMapper = new ObjectMapper();
        javaType = objectMapper.getTypeFactory().constructCollectionType(List.class, {{ .Name }}Resource.class);
        objectMapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);
    }

    @PostConstruct
    public void init() {
        props.getAssets().forEach(this::createCache);
    }

    @Scheduled(initialDelayString = Constants.CACHE_INITIALDELAY_{{ ToUpper .Name }}, fixedRateString = Constants.CACHE_FIXEDRATE_{{ ToUpper .Name }})
    public void populateCacheAll() {
        props.getAssets().forEach(this::populateCache);
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
        return getOne(orgId, {{ $ident.Name }}.hashCode(),
            (resource) -> Optional
                .ofNullable(resource)
                .map({{ $.Name }}Resource::get{{ ToTitle $ident.Name }})
                .map(Identifikator::getIdentifikatorverdi)
                .map(_id -> _id.equals({{ $ident.Name }}))
                .orElse(false));
    }
{{ end }}

	@Override
    public void onAction(Event event) {
        List<{{ .Name }}Resource> data;
        if (checkFintResourceCompatibility && fintResourceCompatibility.isFintResourceData(event.getData())) {
            log.info("Compatibility: Converting FintResource<{{ .Name }}Resource> to {{ .Name }}Resource ...");
            data = fintResourceCompatibility.convertResourceData(event.getData(), {{ .Name }}Resource.class);
        } else {
            data = objectMapper.convertValue(event.getData(), javaType);
        }
        data.forEach(linker::mapLinks);
        if ({{ GetAction .Package }}.valueOf(event.getAction()) == {{ GetAction .Package }}.UPDATE_{{ ToUpper .Name }}) {
            if (event.getResponseStatus() == ResponseStatus.ACCEPTED || event.getResponseStatus() == ResponseStatus.CONFLICT) {
                List<CacheObject<{{ .Name }}Resource>> cacheObjects = data
                    .stream()
                    .map(i -> new CacheObject<>(i, linker.hashCodes(i)))
                    .collect(Collectors.toList());
                addCache(event.getOrgId(), cacheObjects);
                log.info("Added {} cache objects to cache for {}", cacheObjects.size(), event.getOrgId());
            } else {
                log.debug("Ignoring payload for {} with response status {}", event.getOrgId(), event.getResponseStatus());
            }
        } else {
            List<CacheObject<{{ .Name }}Resource>> cacheObjects = data
                    .stream()
                    .map(i -> new CacheObject<>(i, linker.hashCodes(i)))
                    .collect(Collectors.toList());
            updateCache(event.getOrgId(), cacheObjects);
            log.info("Updated cache for {} with {} cache objects", event.getOrgId(), cacheObjects.size());
        }
    }
}
`
