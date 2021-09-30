package generate

const CONTROLLER_TEMPLATE = `package no.fint.consumer.models.{{ modelPkg .Package  }}{{ ToLower .Name }};

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.google.common.collect.ImmutableMap;
import io.swagger.annotations.Api;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;

import no.fint.audit.FintAuditService;

import no.fint.cache.exceptions.*;
import no.fint.consumer.config.Constants;
import no.fint.consumer.config.ConsumerProps;
import no.fint.consumer.event.ConsumerEventUtil;
import no.fint.consumer.event.SynchronousEvents;
import no.fint.consumer.exceptions.*;
import no.fint.consumer.status.StatusCache;
import no.fint.consumer.utils.EventResponses;
import no.fint.consumer.utils.RestEndpoints;

import no.fint.event.model.*;

import no.fint.relations.FintRelationsMediaType;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.util.UriComponentsBuilder;

import no.fint.security.access.policy.FintUserDetails;
import org.springframework.security.core.annotation.AuthenticationPrincipal;

import javax.servlet.http.HttpServletRequest;
import java.net.UnknownHostException;
import java.net.URI;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.TimeUnit;
import java.util.stream.Stream;

import {{ resourcePkg .Package }}.{{ .Name }}Resource;
import {{ resourcePkg .Package }}.{{ .Name }}Resources;
import {{ GetActionPackage .Package }};

@Slf4j
@Api(tags = {"{{ .Name }}"})
@CrossOrigin
@RestController
@RequestMapping(name = "{{ .Name }}", value = RestEndpoints.{{ ToUpper .Name }}, produces = {FintRelationsMediaType.APPLICATION_HAL_JSON_VALUE, MediaType.APPLICATION_JSON_UTF8_VALUE})
public class {{ .Name }}Controller {

    @Autowired(required = false)
    private {{ .Name }}CacheService cacheService;

    @Autowired
    private FintAuditService fintAuditService;

    @Autowired
    private {{ .Name }}Linker linker;

    @Autowired
    private ConsumerProps props;

    @Autowired
    private StatusCache statusCache;

    @Autowired
    private ConsumerEventUtil consumerEventUtil;

    @Autowired
    private ObjectMapper objectMapper;

    @Autowired
    private SynchronousEvents synchronousEvents;

    @GetMapping("/last-updated")
    public Map<String, String> getLastUpdated(@AuthenticationPrincipal FintUserDetails userDetails) {
        String orgId = userDetails.getOrgId();
        if (cacheService == null) {
            throw new CacheDisabledException("{{ .Name }} cache is disabled.");
        }
        if (props.isOverrideOrgId() || orgId == null) {
            orgId = props.getDefaultOrgId();
        }
        String lastUpdated = Long.toString(cacheService.getLastUpdated(orgId));
        return ImmutableMap.of("lastUpdated", lastUpdated);
    }

    @GetMapping("/cache/size")
    public ImmutableMap<String, Integer> getCacheSize(@AuthenticationPrincipal FintUserDetails userDetails) {
        String orgId = userDetails.getOrgId();
        if (cacheService == null) {
            throw new CacheDisabledException("{{ .Name }} cache is disabled.");
        }
        if (props.isOverrideOrgId() || orgId == null) {
            orgId = props.getDefaultOrgId();
        }
        return ImmutableMap.of("size", cacheService.getCacheSize(orgId));
    }

    @GetMapping
    public {{ .Name }}Resources get{{ .Name }}(
            @RequestParam(defaultValue = "0") long sinceTimeStamp,
            @RequestParam(defaultValue = "0") int size,
            @RequestParam(defaultValue = "0") int offset,
            @AuthenticationPrincipal FintUserDetails userDetails,
            HttpServletRequest request) {
        String client = userDetails.getUsername();
        String orgId = userDetails.getOrgId();
        if (cacheService == null) {
            throw new CacheDisabledException("{{ .Name }} cache is disabled.");
        }
        if (props.isOverrideOrgId() || orgId == null) {
            orgId = props.getDefaultOrgId();
        }
        if (client == null) {
            client = props.getDefaultClient();
        }
        log.debug("OrgId: {}, Client: {}", orgId, client);

        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction .Package }}.GET_ALL_{{ ToUpper .Name }}, client);
        event.setOperation(Operation.READ);
        if (StringUtils.isNotBlank(request.getQueryString())) {
            event.setQuery("?" + request.getQueryString());
        }
        fintAuditService.audit(event);
        fintAuditService.audit(event, Status.CACHE);

        Stream<{{ .Name }}Resource> resources;
        if (size > 0 && offset >= 0 && sinceTimeStamp > 0) {
            resources = cacheService.streamSliceSince(orgId, sinceTimeStamp, offset, size);
        } else if (size > 0 && offset >= 0) {
            resources = cacheService.streamSlice(orgId, offset, size);
        } else if (sinceTimeStamp > 0) {
            resources = cacheService.streamSince(orgId, sinceTimeStamp);
        } else {
            resources = cacheService.streamAll(orgId);
        }

        fintAuditService.audit(event, Status.CACHE_RESPONSE, Status.SENT_TO_CLIENT);

        return linker.toResources(resources, offset, size, cacheService.getCacheSize(orgId));
    }

{{ range $i, $ident := .Identifiers }}
    @GetMapping("/{{ ToLower $ident.Name }}/{id:.+}")
    public {{$.Name}}Resource get{{ $.Name }}By{{ ToTitle $ident.Name }}(
            @PathVariable String id,
            @AuthenticationPrincipal FintUserDetails userDetails) throws InterruptedException {
        String client = userDetails.getUsername();
        String orgId = userDetails.getOrgId();
        if (props.isOverrideOrgId() || orgId == null) {
            orgId = props.getDefaultOrgId();
        }
        if (client == null) {
            client = props.getDefaultClient();
        }
        log.debug("{{ $ident.Name }}: {}, OrgId: {}, Client: {}", id, orgId, client);

        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction $.Package }}.GET_{{ ToUpper $.Name }}, client);
        event.setOperation(Operation.READ);
        event.setQuery("{{ $ident.Name }}/" + id);

        if (cacheService != null) {
            fintAuditService.audit(event);
            fintAuditService.audit(event, Status.CACHE);

            Optional<{{ $.Name }}Resource> {{ ToLower $.Name }} = cacheService.get{{ $.Name }}By{{ ToTitle $ident.Name }}(orgId, id);

            fintAuditService.audit(event, Status.CACHE_RESPONSE, Status.SENT_TO_CLIENT);

            return {{ ToLower $.Name }}.map(linker::toResource).orElseThrow(() -> new EntityNotFoundException(id));

        } else {
            BlockingQueue<Event> queue = synchronousEvents.register(event);
            consumerEventUtil.send(event);

            Event response = EventResponses.handle(queue.poll(5, TimeUnit.MINUTES));

            if (response.getData() == null ||
                    response.getData().isEmpty()) throw new EntityNotFoundException(id);

            {{ $.Name }}Resource {{ ToLower $.Name }} = objectMapper.convertValue(response.getData().get(0), {{ $.Name }}Resource.class);

            fintAuditService.audit(response, Status.SENT_TO_CLIENT);

            return linker.toResource({{ ToLower $.Name }});
        }    
    }
{{ end }}

{{ if .Writable }}
    // Writable class
    @GetMapping("/status/{id}")
    public ResponseEntity getStatus(
            @PathVariable String id,
            @AuthenticationPrincipal FintUserDetails userDetails) {
        String client = userDetails.getUsername();
        String orgId = userDetails.getOrgId();
        log.debug("/status/{} for {} from {}", id, orgId, client);
        return statusCache.handleStatusRequest(id, orgId, linker, {{.Name}}Resource.class);
    }

    @PostMapping
    public ResponseEntity post{{.Name}}(
            @AuthenticationPrincipal FintUserDetails userDetails,
            @RequestBody {{.Name}}Resource body,
            @RequestParam(name = "validate", required = false) boolean validate
    ) {
        String client = userDetails.getUsername();
        String orgId = userDetails.getOrgId();
        log.debug("post{{.Name}}, Validate: {}, OrgId: {}, Client: {}", validate, orgId, client);
        log.trace("Body: {}", body);
        linker.mapLinks(body);
        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction .Package}}.UPDATE_{{ ToUpper .Name }}, client);
        event.addObject(objectMapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS).convertValue(body, Map.class));
        event.setOperation(validate ? Operation.VALIDATE : Operation.CREATE);
        consumerEventUtil.send(event);

        statusCache.put(event.getCorrId(), event);

        URI location = UriComponentsBuilder.fromUriString(linker.self()).path("status/{id}").buildAndExpand(event.getCorrId()).toUri();
        return ResponseEntity.status(HttpStatus.ACCEPTED).location(location).build();
    }

  {{ range $i, $ident := .Identifiers }}
    @PutMapping("/{{ ToLower $ident.Name }}/{id:.+}")
    public ResponseEntity put{{ $.Name }}By{{ ToTitle $ident.Name }}(
            @PathVariable String id,
            @AuthenticationPrincipal FintUserDetails userDetails,
            @RequestBody {{$.Name}}Resource body
    ) {
        String client = userDetails.getUsername();
        String orgId = userDetails.getOrgId();
        log.debug("put{{$.Name}}By{{ ToTitle $ident.Name}} {}, OrgId: {}, Client: {}", id, orgId, client);
        log.trace("Body: {}", body);
        linker.mapLinks(body);
        Event event = new Event(orgId, Constants.COMPONENT, {{ GetAction $.Package}}.UPDATE_{{ ToUpper $.Name }}, client);
        event.setQuery("{{ ToLower $ident.Name }}/" + id);
        event.addObject(objectMapper.disable(SerializationFeature.FAIL_ON_EMPTY_BEANS).convertValue(body, Map.class));
        event.setOperation(Operation.UPDATE);
        fintAuditService.audit(event);

        consumerEventUtil.send(event);

        statusCache.put(event.getCorrId(), event);

        URI location = UriComponentsBuilder.fromUriString(linker.self()).path("status/{id}").buildAndExpand(event.getCorrId()).toUri();
        return ResponseEntity.status(HttpStatus.ACCEPTED).location(location).build();
    }
  {{ end }}
                
{{- end }}

    //
    // Exception handlers
    //
    @ExceptionHandler(EventResponseException.class)
    public ResponseEntity handleEventResponseException(EventResponseException e) {
        return ResponseEntity.status(e.getStatus()).body(e.getResponse());
    }

    @ExceptionHandler(UpdateEntityMismatchException.class)
    public ResponseEntity handleUpdateEntityMismatch(Exception e) {
        return ResponseEntity.badRequest().body(ErrorResponse.of(e));
    }

    @ExceptionHandler(EntityNotFoundException.class)
    public ResponseEntity handleEntityNotFound(Exception e) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(ErrorResponse.of(e));
    }

    @ExceptionHandler(CreateEntityMismatchException.class)
    public ResponseEntity handleCreateEntityMismatch(Exception e) {
        return ResponseEntity.badRequest().body(ErrorResponse.of(e));
    }

    @ExceptionHandler(EntityFoundException.class)
    public ResponseEntity handleEntityFound(Exception e) {
        return ResponseEntity.status(HttpStatus.FOUND).body(ErrorResponse.of(e));
    }

    @ExceptionHandler(CacheDisabledException.class)
    public ResponseEntity handleBadRequest(Exception e) {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(ErrorResponse.of(e));
    }

    @ExceptionHandler(UnknownHostException.class)
    public ResponseEntity handleUnkownHost(Exception e) {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(ErrorResponse.of(e));
    }

    @ExceptionHandler(CacheNotFoundException.class)
    public ResponseEntity handleCacheNotFound(Exception e) {
        return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(ErrorResponse.of(e));
    }

}

`
