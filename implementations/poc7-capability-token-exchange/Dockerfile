FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc7-node ./cmd/poc7-node

FROM alpine:3.20
COPY --from=build /out/poc7-node /usr/local/bin/poc7-node
ENTRYPOINT ["/usr/local/bin/poc7-node"]
