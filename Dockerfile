# Stage 1: Build the static files
FROM node:19.9.0 AS frontend-builder
WORKDIR /ui
COPY /ui/package.json /ui/package-lock.json ./
RUN npm ci
COPY /ui .
RUN npm run export

# Stage 2: Build the binary
FROM golang:1.20 AS binary-builder
WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /ui/dist ./ui/dist
RUN CGO_ENABLED=0 go build -ldflags "-w" -a -buildvcs=false -o main ./cmd/web

# Stage 3: Run the binary
FROM gcr.io/distroless/static
EXPOSE 8080
EXPOSE 8081
EXPOSE 3000
WORKDIR /app
COPY --from=binary-builder --chown=nonroot:nonroot /builder/main .
COPY client_secret.json client_secret.json
COPY .env .env
ENTRYPOINT ["./main"]