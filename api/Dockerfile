FROM golang:1.26.4-alpine3.22 AS build
WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /out/marketplace-api ./cmd/server

FROM alpine:3.22
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/marketplace-api /app/marketplace-api
EXPOSE 8098
ENV BIND_ADDRESS=:8098
ENTRYPOINT ["/app/marketplace-api"]
