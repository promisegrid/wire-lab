FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc13-supervisor ./cmd/poc13-supervisor
RUN go build -o /out/poc13-agent ./cmd/poc13-agent
RUN go build -o /out/poc13-analyze ./cmd/poc13-analyze

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc13-supervisor /usr/local/bin/poc13-supervisor
COPY --from=build /out/poc13-agent /usr/local/bin/poc13-agent
COPY --from=build /out/poc13-analyze /usr/local/bin/poc13-analyze
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc13-supervisor"]
