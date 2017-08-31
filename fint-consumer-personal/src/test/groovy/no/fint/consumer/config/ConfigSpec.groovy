package no.fint.consumer.config

import spock.lang.Specification

class ConfigSpec extends Specification {

    def "Get full path including configured context path"() {
        given:
        def config = new Config(contextPath: '/test1')

        when:
        def fullPath = config.fullPath('/test2')

        then:
        fullPath == '/test1/test2'
    }
}
