FROM golang:1.25.7-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/poc9-node ./cmd/poc9-node

FROM alpine:3.20
COPY --from=build /out/poc9-node /usr/local/bin/poc9-node
ENTRYPOINT ["/usr/local/bin/poc9-node"]
