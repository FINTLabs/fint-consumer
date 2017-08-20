package generate

const CONTROLLER_TEMPLATE = `package no.fint.consumer.{{ ToLower .Name  }};

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
@RequestMapping(value = RestEndpoints.{{ ToUpper .Name }}, produces = {FintRelationsMediaType.APPLICATION_HAL_JSON_VALUE, MediaType.APPLICATION_JSON_UTF8_VALUE})
public class {{ .Name }}Controller {

    @Autowired
    private {{ .Name }}CacheService cacheService;

    @Autowired
    private FintAuditService fintAuditService;

    @Autowired
    private {{ .Name }}Assembler assembler;

    @RequestMapping(value = "/last-updated", method = RequestMethod.GET)
    public Map<String, String> getLastUpdated(@RequestHeader(value = HeaderConstants.ORG_ID, defaultValue = Constants.DEFAULT_HEADER_ORGID) String orgId) {
        String lastUpdated = Long.toString(cacheService.getLastUpdated(orgId));
        return ImmutableMap.of("lastUpdated", lastUpdated);
    }

    @RequestMapping(method = RequestMethod.GET)
    public ResponseEntity get{{ .Name }}(@RequestHeader(value = HeaderConstants.ORG_ID, defaultValue = Constants.DEFAULT_HEADER_ORGID) String orgId,
                                               @RequestHeader(value = HeaderConstants.CLIENT, defaultValue = Constants.DEFAULT_HEADER_CLIENT) String client,
                                               @RequestParam(required = false) Long sinceTimeStamp) {
        log.info("OrgId: {}", orgId);
        log.info("Client: {}", client);
        log.info("SinceTimeStamp: {}", sinceTimeStamp);

        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction .Package }}.GET_ALL_{{ ToUpper .Name }}, client);
        fintAuditService.audit(event);

        event.setStatus(Status.CACHE);
        fintAuditService.audit(event);

        List<FintResource<{{ .Name }}>> {{ ToLower .Name }};
        if (sinceTimeStamp == null) {
            {{ ToLower .Name }} = cacheService.getAll(orgId);
        } else {
            {{ ToLower .Name }} = cacheService.getAll(orgId, sinceTimeStamp);
        }

        event.setStatus(Status.CACHE_RESPONSE);
        fintAuditService.audit(event);

        event.setStatus(Status.SENT_TO_CLIENT);
        fintAuditService.audit(event);

        return assembler.resources({{ ToLower .Name }});
    }

    @RequestMapping(value = "/***fixme***/{id}", method = RequestMethod.GET)
    public ResponseEntity get{{ .Name }}(@PathVariable String id,
                                             @RequestHeader(value = HeaderConstants.ORG_ID, defaultValue = Constants.DEFAULT_HEADER_ORGID) String orgId,
                                             @RequestHeader(value = HeaderConstants.CLIENT, defaultValue = Constants.DEFAULT_HEADER_CLIENT) String client) {
        log.info("OrgId: {}", orgId);
        log.info("Client: {}", client);

        Event event = new Event(orgId, Constants.COMPONENT, PersonalActions.GET_{{ ToUpper .Name }}, client);
        fintAuditService.audit(event);

        event.setStatus(Status.CACHE);
        fintAuditService.audit(event);

        Optional<FintResource<{{ .Name }}>> {{ ToLower .Name }} = cacheService.get{{ .Name }}(orgId, id);

        event.setStatus(Status.CACHE_RESPONSE);
        fintAuditService.audit(event);

        event.setStatus(Status.SENT_TO_CLIENT);
        fintAuditService.audit(event);

        if ({{ ToLower .Name }}.isPresent()) {
            return assembler.resource({{ ToLower .Name }}.get());
        } else {
            return ResponseEntity.notFound().build();
        }
    }

    /*
     * TODO: Add endpoints for all Identifikatore.
     */

}

`
