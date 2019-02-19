FROM golang AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/app/vendor/github.com/FINTLabs/fint-consumer
ARG VERSION=0.0.0
COPY . .
RUN go install -v -ldflags "-X main.Version=${VERSION}"
RUN /go/bin/fint-consumer --version

FROM gcr.io/distroless/static
VOLUME [ "/src" ]
WORKDIR /src
COPY --from=builder /go/bin/fint-consumer /usr/bin/fint-consumer
ENTRYPOINT [ "/usr/bin/fint-consumer" ]
