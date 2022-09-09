FROM golang:alpine AS build

RUN apk add --update git
WORKDIR /go/src/github.com/jhonnyesquivel/quasar-op
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/quasar-op-api main.go

# Building image with the binary
FROM scratch
COPY --from=build /go/bin/quasar-op-api /go/bin/quasar-op-api
ENTRYPOINT ["/go/bin/quasar-op-api"]
