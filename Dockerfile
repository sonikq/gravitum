# BUILD STAGE
FROM golang:1.23.1-alpine AS builder

WORKDIR /gravitum_test_task

COPY ../go.mod go.sum ./
RUN go mod download && go mod verify

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /gravitum_test_task/bin/main /gravitum_test_task/cmd/user_management/

# RUN STAGE
FROM ubuntu:22.04

ENV RUN_ADDRESS=localhost:3000
ENV DATABASE_DSN=postgres://user:password@host:port/DB?
ENV DB_POOL_WORKERS=150

ENV CTX_TIMEOUT=5000
ENV LOG_LEVEL=release
ENV SERVICE_NAME=user-management


WORKDIR /gravitum_test_task

COPY --from=builder /gravitum_test_task/bin/main .

EXPOSE 3000
CMD ["/gravitum_test_task/main", "-mode=release"]