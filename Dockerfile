# --- build ---
FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app ./cmd/app

# --- run ---
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/bin/app /app/app
ENV HTTP_ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
