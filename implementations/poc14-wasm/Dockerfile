FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc14-kernel ./cmd/poc14-kernel
RUN go build -o /out/poc14-relationship-agent ./cmd/poc14-relationship-agent
RUN go build -o /out/poc14-fulfillment ./cmd/poc14-fulfillment
RUN go build -o /out/poc14-postal-scale ./cmd/poc14-postal-scale
RUN go build -o /out/poc14-ups-label-printer ./cmd/poc14-ups-label-printer
RUN go build -o /out/poc14-printer-port ./cmd/poc14-printer-port
RUN go build -o /out/poc14-accounting ./cmd/poc14-accounting
RUN go build -o /out/poc14-wasm-agent ./cmd/poc14-wasm-agent
RUN go build -o /out/poc14-stdio-adapter ./cmd/poc14-stdio-adapter
RUN go build -o /out/poc14-stdio-worker ./cmd/poc14-stdio-worker
RUN go build -o /out/poc14-supervisor ./cmd/poc14-supervisor
RUN go build -o /out/poc14-event-collector ./cmd/poc14-event-collector
RUN go build -o /out/poc14-analyze ./cmd/poc14-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc14-kernel /usr/local/bin/poc14-kernel
COPY --from=build /out/poc14-relationship-agent /usr/local/bin/poc14-relationship-agent
COPY --from=build /out/poc14-fulfillment /usr/local/bin/poc14-fulfillment
COPY --from=build /out/poc14-postal-scale /usr/local/bin/poc14-postal-scale
COPY --from=build /out/poc14-ups-label-printer /usr/local/bin/poc14-ups-label-printer
COPY --from=build /out/poc14-printer-port /usr/local/bin/poc14-printer-port
COPY --from=build /out/poc14-accounting /usr/local/bin/poc14-accounting
COPY --from=build /out/poc14-wasm-agent /usr/local/bin/poc14-wasm-agent
COPY --from=build /out/poc14-stdio-adapter /usr/local/bin/poc14-stdio-adapter
COPY --from=build /out/poc14-stdio-worker /usr/local/bin/poc14-stdio-worker
COPY --from=build /out/poc14-supervisor /usr/local/bin/poc14-supervisor
COPY --from=build /out/poc14-event-collector /usr/local/bin/poc14-event-collector
COPY --from=build /out/poc14-analyze /usr/local/bin/poc14-analyze
COPY config.json /app/config.json
ENTRYPOINT ["/usr/local/bin/poc14-supervisor"]
