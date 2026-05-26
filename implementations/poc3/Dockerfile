FROM golang:1.22-alpine AS build

WORKDIR /src
COPY . .

RUN go test ./...
RUN go build -o /out/poc3-kernel ./kernel/cmd
RUN go build -o /out/poc3-hello ./hello/cmd
RUN go build -o /out/poc3-echo ./echo/cmd
RUN go build -o /out/poc3-signed ./signed/cmd

FROM alpine:3.20

COPY --from=build /out/poc3-kernel /usr/local/bin/poc3-kernel
COPY --from=build /out/poc3-hello /usr/local/bin/poc3-hello
COPY --from=build /out/poc3-echo /usr/local/bin/poc3-echo
COPY --from=build /out/poc3-signed /usr/local/bin/poc3-signed
COPY scripts/run-alice.sh /usr/local/bin/run-alice.sh
COPY scripts/run-bob.sh /usr/local/bin/run-bob.sh

RUN chmod +x /usr/local/bin/run-alice.sh /usr/local/bin/run-bob.sh

WORKDIR /work
