FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/markitos-mdk ./cmd/app

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/markitos-mdk ./markitos-mdk
COPY --from=builder /src/index.html ./index.html
COPY --from=builder /src/faqs.html ./faqs.html
COPY --from=builder /src/articles.html ./articles.html
COPY --from=builder /src/videos.html ./videos.html
COPY --from=builder /src/git.html ./git.html
COPY --from=builder /src/css ./css
COPY --from=builder /src/js ./js

EXPOSE 8080

CMD ["./markitos-mdk"]
