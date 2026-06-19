FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc16-kernel ./cmd/poc16-kernel
RUN go build -o /out/poc16-relationship-agent ./cmd/poc16-relationship-agent
RUN go build -o /out/poc16-fulfillment ./cmd/poc16-fulfillment
RUN go build -o /out/poc16-postal-scale ./cmd/poc16-postal-scale
RUN go build -o /out/poc16-ups-label-printer ./cmd/poc16-ups-label-printer
RUN go build -o /out/poc16-printer-port ./cmd/poc16-printer-port
RUN go build -o /out/poc16-accounting ./cmd/poc16-accounting
RUN go build -o /out/poc16-wasm-agent ./cmd/poc16-wasm-agent
RUN go build -o /out/poc16-stdio-adapter ./cmd/poc16-stdio-adapter
RUN go build -o /out/poc16-stdio-worker ./cmd/poc16-stdio-worker
RUN go build -o /out/poc16-supervisor ./cmd/poc16-supervisor
RUN go build -o /out/poc16-event-collector ./cmd/poc16-event-collector
RUN go build -o /out/poc16-analyze ./cmd/poc16-analyze
RUN go build -o /out/poc16-cbor-diag ./cmd/poc16-cbor-diag

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc16-kernel /usr/local/bin/poc16-kernel
COPY --from=build /out/poc16-relationship-agent /usr/local/bin/poc16-relationship-agent
COPY --from=build /out/poc16-fulfillment /usr/local/bin/poc16-fulfillment
COPY --from=build /out/poc16-postal-scale /usr/local/bin/poc16-postal-scale
COPY --from=build /out/poc16-ups-label-printer /usr/local/bin/poc16-ups-label-printer
COPY --from=build /out/poc16-printer-port /usr/local/bin/poc16-printer-port
COPY --from=build /out/poc16-accounting /usr/local/bin/poc16-accounting
COPY --from=build /out/poc16-wasm-agent /usr/local/bin/poc16-wasm-agent
COPY --from=build /out/poc16-stdio-adapter /usr/local/bin/poc16-stdio-adapter
COPY --from=build /out/poc16-stdio-worker /usr/local/bin/poc16-stdio-worker
COPY --from=build /out/poc16-supervisor /usr/local/bin/poc16-supervisor
COPY --from=build /out/poc16-event-collector /usr/local/bin/poc16-event-collector
COPY --from=build /out/poc16-analyze /usr/local/bin/poc16-analyze
COPY --from=build /out/poc16-cbor-diag /usr/local/bin/poc16-cbor-diag
COPY config.json /app/config.json
ENTRYPOINT ["/usr/local/bin/poc16-supervisor"]
