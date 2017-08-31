package no.fint.consumer.event

import no.fint.audit.FintAuditService
import no.fint.event.model.Event
import no.fint.event.model.Status
import no.fint.events.FintEvents
import no.fint.events.FintEventsHealth
import spock.lang.Specification

class ConsumerEventUtilSpec extends Specification {
    private ConsumerEventUtil consumerEventUtil
    private FintEvents fintEvents
    private FintEventsHealth fintEventsHealth
    private FintAuditService fintAuditService

    void setup() {
        fintEvents = Mock(FintEvents)
        fintAuditService = Mock(FintAuditService)
        fintEventsHealth = Mock(FintEventsHealth)
        consumerEventUtil = new ConsumerEventUtil(fintEvents: fintEvents, fintEventsHealth: fintEventsHealth, fintAuditService: fintAuditService)
    }

    def "Send and receive Event"() {
        given:
        def event = new Event(orgId: 'rogfk.no', corrId: '123')

        when:
        def response = consumerEventUtil.healthCheck(event)

        then:
        2 * fintAuditService.audit(_ as Event, _ as Status)
        1 * fintEventsHealth.sendHealthCheck('rogfk.no', '123', event) >> event
        response.isPresent()
    }
}
