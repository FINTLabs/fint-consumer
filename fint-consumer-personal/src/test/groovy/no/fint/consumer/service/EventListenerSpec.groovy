package no.fint.consumer.service

import no.fint.consumer.event.EventListener
import no.fint.event.model.Event
import spock.lang.Specification

class EventListenerSpec extends Specification {
    private EventListener eventListener

    void setup() {
        eventListener = new EventListener()
    }

    def "No exception is thrown when receiving event"() {
        when:
        eventListener.recieve(new Event(corrId: '123'))

        then:
        noExceptionThrown()
    }
}
