FROM golang:1.24-alpine

WORKDIR /app

# Копируем весь проект
COPY .. .

# Переходим в папку product-service
WORKDIR /app/product-service

# Собираем
RUN go build -o main .

EXPOSE 8081

CMD ["./main"]
