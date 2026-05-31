FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc8-node ./cmd/poc8-node

FROM alpine:3.20
COPY --from=build /out/poc8-node /usr/local/bin/poc8-node
ENTRYPOINT ["/usr/local/bin/poc8-node"]
