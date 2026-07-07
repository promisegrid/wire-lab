FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc-agent ./cmd/poc-agent
RUN go build -o /out/poc-event-collector ./cmd/poc-event-collector
RUN go build -o /out/poc-analyze ./cmd/poc-analyze
RUN go build -o /out/poc-cbor-diag ./cmd/poc-cbor-diag
RUN go build -o /out/grid ./cmd/grid

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc-agent /usr/local/bin/poc-agent
COPY --from=build /out/poc-event-collector /usr/local/bin/poc-event-collector
COPY --from=build /out/poc-analyze /usr/local/bin/poc-analyze
COPY --from=build /out/poc-cbor-diag /usr/local/bin/poc-cbor-diag
COPY --from=build /out/grid /usr/local/bin/grid
ENTRYPOINT ["/usr/local/bin/poc-agent"]
