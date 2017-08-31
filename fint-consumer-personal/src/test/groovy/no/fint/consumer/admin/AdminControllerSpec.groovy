package no.fint.consumer.admin

import no.fint.consumer.event.ConsumerEventUtil
import no.fint.event.model.DefaultActions
import no.fint.event.model.Event
import no.fint.event.model.HeaderConstants
import no.fint.test.utils.MockMvcSpecification
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.setup.MockMvcBuilders

class AdminControllerSpec extends MockMvcSpecification {
    private AdminController healthController
    private MockMvc mockMvc
    private ConsumerEventUtil consumerEventUtil

    void setup() {
        consumerEventUtil = Mock(ConsumerEventUtil)
        healthController = new AdminController(consumerEventUtil: consumerEventUtil)
        mockMvc = MockMvcBuilders.standaloneSetup(healthController).build()
    }

    def "Check response on healthcheck"() {
        when:
        def response = mockMvc.perform(get("/admin/health").header(HeaderConstants.ORG_ID, "mock.no").header(HeaderConstants.CLIENT, "mock"))

        then:
        1 * consumerEventUtil.healthCheck(_ as Event) >> Optional.of(new Event(action: DefaultActions.HEALTH.name()))
        response.andExpect(status().isOk())
                .andExpect(jsonPath('$.action')
                .value(equalTo('HEALTH')))
    }

    def "Check response on healthcheck is empty"() {
        when:
        def response = mockMvc.perform(get("/admin/health").header(HeaderConstants.ORG_ID, "mock.no").header(HeaderConstants.CLIENT, "mock"))

        then:
        1 * consumerEventUtil.healthCheck(_ as Event) >> Optional.empty()
        response.andExpect(status().is5xxServerError())
                .andExpect(jsonPath('$.action').value(equalTo('HEALTH')))
                .andExpect(jsonPath('$.message').value(equalTo('No response from adapter')))
    }
}
