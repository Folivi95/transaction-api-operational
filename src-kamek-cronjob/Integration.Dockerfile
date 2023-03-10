FROM public.ecr.aws/docker/library/golang:1.18.2 as deps

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY vendor ./vendor
COPY go.mod go.sum ./

COPY internal/ ./internal/

CMD [ "go", "test", "-count=1", "--tags=integration", "./..." ]
