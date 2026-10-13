FROM golang:alpine AS build

RUN apk upgrade --no-cache zlib
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY app ./app
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /payments ./cmd/payments

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /payments /payments
ENV GIN_MODE=release PORT=8080
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/payments"]
