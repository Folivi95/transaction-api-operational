FROM public.ecr.aws/docker/library/golang:1.18 as deps

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPRIVATE=github.com/saltpay/*

RUN git config --global url."git@github.com".insteadOf "https://github.com"

WORKDIR /app
COPY vendor ./vendor
COPY go.mod go.sum ./

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY black-box-tests/ ./black-box-tests/

FROM deps as build
RUN go build -o mario-api github.com/saltpay/transaction-api-operational/src-mario-api/cmd/web

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/ ./

CMD ["/app/mario-api"]
EXPOSE 8080
