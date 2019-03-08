pipeline {
    agent {
        label 'docker'
    }
    stages {
        stage('Build') {
            when {
                not { buildingTag() }
            }
            steps {
                sh "docker build --build-arg VERSION=${BUILD_TAG} ."
            }
        }
        stage('Deploy') {
            when {
                tag pattern: "v\\d+\\.\\d+\\.\\d+(-\\w+-\\d+)?", comparator: "REGEXP"
            }
            steps {
                script {
                    VERSION = TAG_NAME[1..-1]
                }
                sh "echo Version is ${VERSION}"
                sh "docker build --tag fint/fint-consumer:${VERSION} --build-arg VERSION=${VERSION} ."
                withDockerRegistry([credentialsId: 'asgeir-docker', url: '']) {
                    sh "docker push fint/fint-consumer:${VERSION}"
                }
            }
        }
    }
}
