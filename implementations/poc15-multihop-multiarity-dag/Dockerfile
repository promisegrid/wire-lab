FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc15-kernel ./cmd/poc15-kernel
RUN go build -o /out/poc15-relationship-agent ./cmd/poc15-relationship-agent
RUN go build -o /out/poc15-fulfillment ./cmd/poc15-fulfillment
RUN go build -o /out/poc15-postal-scale ./cmd/poc15-postal-scale
RUN go build -o /out/poc15-ups-label-printer ./cmd/poc15-ups-label-printer
RUN go build -o /out/poc15-printer-port ./cmd/poc15-printer-port
RUN go build -o /out/poc15-accounting ./cmd/poc15-accounting
RUN go build -o /out/poc15-wasm-agent ./cmd/poc15-wasm-agent
RUN go build -o /out/poc15-stdio-adapter ./cmd/poc15-stdio-adapter
RUN go build -o /out/poc15-stdio-worker ./cmd/poc15-stdio-worker
RUN go build -o /out/poc15-supervisor ./cmd/poc15-supervisor
RUN go build -o /out/poc15-event-collector ./cmd/poc15-event-collector
RUN go build -o /out/poc15-analyze ./cmd/poc15-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc15-kernel /usr/local/bin/poc15-kernel
COPY --from=build /out/poc15-relationship-agent /usr/local/bin/poc15-relationship-agent
COPY --from=build /out/poc15-fulfillment /usr/local/bin/poc15-fulfillment
COPY --from=build /out/poc15-postal-scale /usr/local/bin/poc15-postal-scale
COPY --from=build /out/poc15-ups-label-printer /usr/local/bin/poc15-ups-label-printer
COPY --from=build /out/poc15-printer-port /usr/local/bin/poc15-printer-port
COPY --from=build /out/poc15-accounting /usr/local/bin/poc15-accounting
COPY --from=build /out/poc15-wasm-agent /usr/local/bin/poc15-wasm-agent
COPY --from=build /out/poc15-stdio-adapter /usr/local/bin/poc15-stdio-adapter
COPY --from=build /out/poc15-stdio-worker /usr/local/bin/poc15-stdio-worker
COPY --from=build /out/poc15-supervisor /usr/local/bin/poc15-supervisor
COPY --from=build /out/poc15-event-collector /usr/local/bin/poc15-event-collector
COPY --from=build /out/poc15-analyze /usr/local/bin/poc15-analyze
COPY config.json /app/config.json
ENTRYPOINT ["/usr/local/bin/poc15-supervisor"]
