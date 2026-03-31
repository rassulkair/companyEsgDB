FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=build /app/app /app/app
COPY --from=build /app/web /app/web
COPY --from=build /app/internal/db /app/internal/db

EXPOSE 6061

CMD ["/app/app"]