node('master') {
    stage('Prepare') {
        checkout scm
        sh 'git log --oneline | nl -nln | perl -lne \'if (/^(\\d+).*Version (\\d+\\.\\d+\\.\\d+)/) { print "$2-$1"; exit; }\' > version.txt'
        stash includes: 'version.txt', name: 'version'
    }
}
node('docker') {
    stage('Build') {
        checkout scm
        String goPath = "/go/src/app/vendor/github.com/FINTprosjektet/fint-consumer"
        docker.image('golang').inside("-v /tmp:/tmp -v ${pwd()}:${goPath}") {
            sh "go-wrapper download github.com/mitchellh/gox && go-wrapper install github.com/mitchellh/gox"
            unstash 'version'
            VERSION=readFile('version.txt').trim()
            sh "cd ${goPath}; gox -output=\"./{{.Dir}}-${VERSION}-{{.OS}}\" -rebuild -osarch=\"darwin/amd64 windows/amd64\" -ldflags='-X main.Version=${VERSION}'"
            stash name: 'artifacts', includes: 'fint-consumer-*'
        }
    }
    stage('Publish') {
        unstash 'version'
        unstash 'artifacts'
        VERSION=readFile('version.txt').trim()
        archiveArtifacts 'fint-consumer-*'
    }
    stage('Cleanup') {
        cleanWs
    }
}
