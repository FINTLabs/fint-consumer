package generate

const CONTROLLER_TEMPLATE = `package no.fint.consumer.models.{{ ToLower .Name  }};

import com.google.common.collect.ImmutableMap;
import lombok.extern.slf4j.Slf4j;
import no.fint.audit.FintAuditService;
import no.fint.consumer.config.Constants;
import no.fint.consumer.utils.RestEndpoints;
import no.fint.event.model.Event;
import no.fint.event.model.HeaderConstants;
import no.fint.event.model.Status;

import no.fint.model.relation.FintResource;
import no.fint.relations.FintRelationsMediaType;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;
import java.util.Optional;

import {{ .Package }}.{{ .Name }};
import {{ GetActionPackage .Package }};

@Slf4j
@CrossOrigin
@RestController
@RequestMapping(value = "/{{ ToLower .Name }}", produces = {FintRelationsMediaType.APPLICATION_HAL_JSON_VALUE, MediaType.APPLICATION_JSON_UTF8_VALUE})
public class {{ .Name }}Controller {

    @Autowired
    private {{ .Name }}CacheService cacheService;

    @Autowired
    private FintAuditService fintAuditService;

    @Autowired
    private {{ .Name }}Assembler assembler;

    @GetMapping("/last-updated")
    public Map<String, String> getLastUpdated(@RequestHeader(HeaderConstants.ORG_ID) String orgId) {
        String lastUpdated = Long.toString(cacheService.getLastUpdated(orgId));
        return ImmutableMap.of("lastUpdated", lastUpdated);
    }

    @GetMapping("/cache/size")
     public ImmutableMap<String, Integer> getCacheSize(@RequestHeader(HeaderConstants.ORG_ID) String orgId) {
        return ImmutableMap.of("size", cacheService.getAll(orgId).size());
    }

    @PostMapping("/cache/rebuild")
    public void rebuildCache(@RequestHeader(HeaderConstants.ORG_ID) String orgId) {
        cacheService.rebuildCache(orgId);
    }

    @GetMapping
    public ResponseEntity get{{ .Name }}(@RequestHeader(HeaderConstants.ORG_ID) String orgId,
            @RequestHeader(HeaderConstants.CLIENT) String client,
            @RequestParam(required = false) Long sinceTimeStamp) {
        log.info("OrgId: {}, Client: {}", orgId, client);

        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction .Package }}.GET_ALL_{{ ToUpper .Name }}, client);
        fintAuditService.audit(event);

        fintAuditService.audit(event, Status.CACHE);

        List<FintResource<{{ .Name }}>> {{ ToLower .Name }};
        if (sinceTimeStamp == null) {
            {{ ToLower .Name }} = cacheService.getAll(orgId);
        } else {
            {{ ToLower .Name }} = cacheService.getAll(orgId, sinceTimeStamp);
        }

        fintAuditService.audit(event, Status.CACHE_RESPONSE, Status.SENT_TO_CLIENT);

        return assembler.resources({{ ToLower .Name }});
    }

{{ range $i, $ident := .Identifiers }}
    @GetMapping("/{{ ToLower $ident.Name }}/{id}")
    public ResponseEntity get{{ $.Name }}By{{ ToTitle $ident.Name }}(@PathVariable String id,
            @RequestHeader(HeaderConstants.ORG_ID) String orgId,
            @RequestHeader(HeaderConstants.CLIENT) String client) {
        log.info("{{ ToTitle $ident.Name }}: {}, OrgId: {}, Client: {}", id, orgId, client);

        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction $.Package }}.GET_{{ ToUpper $.Name }}, client);
        fintAuditService.audit(event);

        fintAuditService.audit(event, Status.CACHE);

        Optional<FintResource<{{ $.Name }}>> {{ ToLower $.Name }} = cacheService.get{{ $.Name }}By{{ ToTitle $ident.Name }}(orgId, id);

        fintAuditService.audit(event, Status.CACHE_RESPONSE, Status.SENT_TO_CLIENT);

        if ({{ ToLower $.Name }}.isPresent()) {
            return assembler.resource({{ ToLower $.Name }}.get());
        } else {
            return ResponseEntity.notFound().build();
        }
    }
{{ end }}
    
}

`