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

COPY certs/salt.pem /usr/share/ca-certificates/
RUN sed -i -e '$asalt.pem' /etc/ca-certificates.conf
RUN update-ca-certificates

FROM deps as build
RUN go build -o luigi-transformer github.com/saltpay/transaction-api-operational/src-luigi-transformer/cmd/web

FROM scratch
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /app/ ./

CMD ["/app/luigi-transformer"]
EXPOSE 8081
