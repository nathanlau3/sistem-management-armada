FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/ ./cmd/...

FROM alpine:3.20

WORKDIR /app
COPY --from=build /out/ /usr/local/bin/
COPY movement.json .

EXPOSE 8080
CMD ["server"]
