FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .


FROM alpine:latest as production

RUN apk --no-cache add ca-certificates tzdata && \
    apk --no-cache upgrade && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

ENV TZ=Asia/Yerevan
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

WORKDIR /home/appuser/

COPY --from=builder /app/main .

RUN chown -R appuser:appgroup /home/appuser/

USER appuser

CMD ["./main"]
