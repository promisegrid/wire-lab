FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY . .

RUN go build -o /out/poc4-kernel ./kernel/cmd
RUN go build -o /out/poc4-relay ./relay/cmd
RUN go build -o /out/poc4-hello ./hello/cmd
RUN go build -o /out/poc4-echo ./echo/cmd
RUN go build -o /out/poc4-signed ./signed/cmd
RUN go build -o /out/poc4-fibonacci ./fibonacci/cmd
RUN go build -o /out/poc4-storage ./storage/cmd

FROM alpine:3.20

WORKDIR /poc4
COPY --from=build /out/ /usr/local/bin/
COPY scripts/ /usr/local/bin/
RUN chmod +x /usr/local/bin/run-alice.sh /usr/local/bin/run-bob.sh /usr/local/bin/run-carol.sh /usr/local/bin/run-dave.sh /usr/local/bin/run-ellen.sh

