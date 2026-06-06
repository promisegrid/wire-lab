FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc12-kernel ./cmd/poc12-kernel
RUN go build -o /out/poc12-relationship-agent ./cmd/poc12-relationship-agent
RUN go build -o /out/poc12-fulfillment ./cmd/poc12-fulfillment
RUN go build -o /out/poc12-postal-scale ./cmd/poc12-postal-scale
RUN go build -o /out/poc12-ups-label-printer ./cmd/poc12-ups-label-printer
RUN go build -o /out/poc12-printer-port ./cmd/poc12-printer-port
RUN go build -o /out/poc12-accounting ./cmd/poc12-accounting
RUN go build -o /out/poc12-supervisor ./cmd/poc12-supervisor
RUN go build -o /out/poc12-analyze ./cmd/poc12-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc12-kernel /usr/local/bin/poc12-kernel
COPY --from=build /out/poc12-relationship-agent /usr/local/bin/poc12-relationship-agent
COPY --from=build /out/poc12-fulfillment /usr/local/bin/poc12-fulfillment
COPY --from=build /out/poc12-postal-scale /usr/local/bin/poc12-postal-scale
COPY --from=build /out/poc12-ups-label-printer /usr/local/bin/poc12-ups-label-printer
COPY --from=build /out/poc12-printer-port /usr/local/bin/poc12-printer-port
COPY --from=build /out/poc12-accounting /usr/local/bin/poc12-accounting
COPY --from=build /out/poc12-supervisor /usr/local/bin/poc12-supervisor
COPY --from=build /out/poc12-analyze /usr/local/bin/poc12-analyze
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc12-supervisor"]
