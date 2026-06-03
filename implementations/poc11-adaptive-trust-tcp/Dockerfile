FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc11-node ./cmd/poc11-node
RUN go build -o /out/poc11-supervisor ./cmd/poc11-supervisor

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc11-node /usr/local/bin/poc11-node
COPY --from=build /out/poc11-supervisor /usr/local/bin/poc11-supervisor
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc11-supervisor"]
