# golang:1.20.3-alpine - Название образа из DockerHub, включает минималистичный линукс + golang
# с FROM начинается этап сборки докера
# builder - присваивается название этапа сборки
FROM golang:1.23.2-alpine3.20 AS builder

# Копируем все из текущей директории в указанный вторым аргументом путь в докер образе
COPY . /github.com/algol/auth/source/
# Устанавливаем папку как рабочую директорию
WORKDIR /github.com/algol/auth/source/   
# Скачиваем зависимости
RUN go mod download
# Собираем приложение, первым параметром название бинарника, вторым - путь до мейна
RUN go build -o ./bin/auth_server cmd/main.go cmd/auth_db.go       

# второй этап сборки
# загружаем минимальный линукс без golang (он нужен был только для сборки приложения в п.1)
FROM alpine:latest
# Устанавливаем рабочую директорию
WORKDIR /root/
# Копируем из п.1 собранный бинарник в рабочую директорию 
COPY --from=builder /github.com/algol/auth/source/bin/auth_server .

CMD ["./auth_server"]      