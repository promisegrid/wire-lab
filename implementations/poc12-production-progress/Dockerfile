FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc12-node ./cmd/poc12-node
RUN go build -o /out/poc12-supervisor ./cmd/poc12-supervisor
RUN go build -o /out/poc12-analyze ./cmd/poc12-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc12-node /usr/local/bin/poc12-node
COPY --from=build /out/poc12-supervisor /usr/local/bin/poc12-supervisor
COPY --from=build /out/poc12-analyze /usr/local/bin/poc12-analyze
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc12-supervisor"]
