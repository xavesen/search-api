FROM golang:1.22.6-alpine3.20
WORKDIR /work/search-api
COPY . .
RUN rm -f main
RUN go get -d -v ./... && go install -v ./...
RUN go build -tags=viper_bind_struct cmd/main.go

FROM alpine:3.20.2
WORKDIR /work/search-api
COPY --from=0 /work/search-api/main .
RUN apk add --no-cache curl
CMD ["/work/search-api/main"]