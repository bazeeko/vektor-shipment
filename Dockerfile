FROM golang:1.26.1-alpine as builder
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM scratch
COPY --from=builder /build/server /server
COPY /configs/values.yaml /configs/values.yaml
COPY migrations migrations

ENV CONFIG_FILE_PATH="/configs/values.yaml"

ENTRYPOINT ["/server"]