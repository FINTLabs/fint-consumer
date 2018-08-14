FROM golang:alpine AS builder
WORKDIR /go/src/app/vendor/github.com/FINTprosjektet/fint-consumer
ARG VERSION=0.0.0
COPY . .
RUN go install -v -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-consumer --version

FROM alpine
COPY --from=builder /go/bin/fint-consumer /usr/bin/fint-consumer
WORKDIR /src
VOLUME [ "/src" ]
ENTRYPOINT [ "/usr/bin/fint-consumer" ]
