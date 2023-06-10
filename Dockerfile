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
ENV DB_HOST=db
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=mysecretpassword
ENV DB_NAME=mydb
WORKDIR /app
COPY --from=binary-builder --chown=nonroot:nonroot /builder/main .
ENTRYPOINT ["./main"]



# old
# FROM golang:1.20

# WORKDIR /app

# COPY . .

# RUN go build -o main ./cmd/web

# EXPOSE 4000

# ENV DB_HOST=db
# ENV DB_PORT=5432
# ENV DB_USER=postgres
# ENV DB_PASSWORD=mysecretpassword
# ENV DB_NAME=mydb

# CMD ["./main"]