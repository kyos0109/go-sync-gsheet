FROM golang:1.14-stretch AS base

WORKDIR /go/src/app

COPY / .

RUN go get -d -v ./...

RUN go install -v ./...

FROM gcr.io/distroless/base

COPY --from=base /go/bin/go-sync-gsheet /usr/bin/go-sync-gsheet

WORKDIR /app

ENTRYPOINT ["go-sync-gsheet"]