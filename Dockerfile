FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/
COPY config/ config/

RUN CGO_ENABLED=0 go build -o /service cmd/main.go

FROM alpine:3.19 AS run-stage

WORKDIR /

COPY --from=build-stage /service /service
COPY --from=build-stage /app/config/docker_config.yaml /config/docker_config.yaml

EXPOSE 8080

ENTRYPOINT ["/service"]
