# ================================
# STAGE 1 : Build
# ================================
FROM golang:1.25.5-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/microservice ./cmd/api

# ===============================
# STAGE 2: RUNTIME
# ===============================
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/microservice /microservice

EXPOSE 3200
ENTRYPOINT [ "/microservice" ]