FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc10-node ./cmd/poc10-node

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/poc10-node /usr/local/bin/poc10-node
COPY config.example.json /app/config.example.json
ENTRYPOINT ["/usr/local/bin/poc10-node"]
