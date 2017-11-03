$goPath = "/go/src/github.com/FINTprosjektet/fint-consumer"
#docker run -i -v ${PWD}:${goPath} -w $goPath golang bash -c "go get github.com/rancher/trash && go install github.com/rancher/trash && trash --update && pwd"
docker run -i -v ${PWD}:${goPath} -w $goPath -e GOOS=windows golang go build
