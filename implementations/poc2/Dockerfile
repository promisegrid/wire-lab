FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . ./
RUN go test ./...
RUN go build -o /out/poc2 .

FROM alpine:3.20
COPY --from=build /out/poc2 /usr/local/bin/poc2-kernel
COPY --from=build /out/poc2 /usr/local/bin/poc2-hello
COPY scripts/run-alice.sh /usr/local/bin/run-alice.sh
COPY scripts/run-bob.sh /usr/local/bin/run-bob.sh
RUN chmod 755 /usr/local/bin/run-alice.sh /usr/local/bin/run-bob.sh
WORKDIR /work
