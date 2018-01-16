package generate

const CACHE_SERVICE_TEMPLATE = `package no.fint.consumer.models.{{ ToLower .Name  }};

import com.fasterxml.jackson.core.type.TypeReference;
import lombok.extern.slf4j.Slf4j;
import no.fint.cache.CacheService;
import no.fint.consumer.config.Constants;
import no.fint.consumer.config.ConsumerProps;
import no.fint.consumer.event.ConsumerEventUtil;
import no.fint.event.model.Event;
import no.fint.model.relation.FintResource;
import no.fint.model.felles.kompleksedatatyper.Identifikator;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;
import java.util.Arrays;
import java.util.List;
import java.util.Optional;

import {{ .Package }}.{{ .Name }};
import {{ GetActionPackage .Package }};

@Slf4j
@Service
public class {{ .Name }}CacheService extends CacheService<FintResource<{{ .Name }}>> {

    public static final String MODEL = {{ .Name }}.class.getSimpleName().toLowerCase();

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
    public Optional<FintResource<{{ $.Name }}>> get{{ $.Name }}By{{ ToTitle $ident.Name }}(String orgId, String {{ $ident.Name }}) {
        Identifikator needle = new Identifikator();
        needle.setIdentifikatorverdi({{ $ident.Name }});
        return getOne(orgId, (fintResource) -> needle.equals(fintResource.getResource().get{{ ToTitle $ident.Name }}()));
    }
{{ end }}

	@Override
    public void onAction(Event event) {
        update(event, new TypeReference<List<FintResource<{{ .Name }}>>>() {
        });
    }
}
`