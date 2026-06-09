FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc13-kernel ./cmd/poc13-kernel
RUN go build -o /out/poc13-relationship-agent ./cmd/poc13-relationship-agent
RUN go build -o /out/poc13-fulfillment ./cmd/poc13-fulfillment
RUN go build -o /out/poc13-postal-scale ./cmd/poc13-postal-scale
RUN go build -o /out/poc13-ups-label-printer ./cmd/poc13-ups-label-printer
RUN go build -o /out/poc13-printer-port ./cmd/poc13-printer-port
RUN go build -o /out/poc13-accounting ./cmd/poc13-accounting
RUN go build -o /out/poc13-supervisor ./cmd/poc13-supervisor
RUN go build -o /out/poc13-analyze ./cmd/poc13-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc13-kernel /usr/local/bin/poc13-kernel
COPY --from=build /out/poc13-relationship-agent /usr/local/bin/poc13-relationship-agent
COPY --from=build /out/poc13-fulfillment /usr/local/bin/poc13-fulfillment
COPY --from=build /out/poc13-postal-scale /usr/local/bin/poc13-postal-scale
COPY --from=build /out/poc13-ups-label-printer /usr/local/bin/poc13-ups-label-printer
COPY --from=build /out/poc13-printer-port /usr/local/bin/poc13-printer-port
COPY --from=build /out/poc13-accounting /usr/local/bin/poc13-accounting
COPY --from=build /out/poc13-supervisor /usr/local/bin/poc13-supervisor
COPY --from=build /out/poc13-analyze /usr/local/bin/poc13-analyze
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc13-supervisor"]
