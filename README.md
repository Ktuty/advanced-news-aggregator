<div align="center">
<a href="https://git.io/typing-svg"><img src="https://readme-typing-svg.herokuapp.com?font=Tektur&weight=600&size=35&duration=4000&color=53E1FF&center=true&vCenter=true&width=435&height=100&lines=news+aggregator" alt="Typing SVG" /></a>
</div>

# Технологии

- Golang
- PostgreSQL
- REST API
- microservices

# Подготовка микросервисов

1. **Клонирование репозитория**:
   ```sh
   git clone https://github.com/Ktuty/advanced-news-aggregator

# APIGateway

1. **Установка зависимостей в директории APIGateway**:
   ```go
   go mod download
   go mod tidy

# censorship-service

1. **Установка зависимостей в директории censorship-service**:
   ```go
   go mod download
   go mod tidy

# comment-service

1. **Установка зависимостей в директории comment-service**:
   ```go
   go mod download
   go mod tidy

2. **Проверка и изменение файлов конфигурации .env:**
   ```env
   DATABASE_URL: "yuor_DB_URL"

# news-aggregator

1. **Установка зависимостей в директории news-aggregator**:
   ```go
   go mod download
   go mod tidy
   
2. **Проверка и изменение файлов конфигурации .env:**
   ```env
   DATABASE_URL: "yuor_DB_URL"

3. **Проверка и изменение пути для JSON файла sites.json :**
   news-aggregator/internal/rss/rss.go:
   ```sh
   C:/Users/User/.../news-aggregator/internal/rss/sites.json

4. **Настройка файла sites.json:**
   Example:
   ```json
     {
      "rss":[
        "https://habr.com/ru/rss/hub/go/all/?fl=ru",
        "https://habr.com/ru/rss/best/daily/?fl=ru",
        "https://cprss.s3.amazonaws.com/golangweekly.com.xml",
        "https://go.dev/blog/feed.atom?format=xml",
        "https://blog.jetbrains.com/go/feed/"
      ],
      "request_period": 1
    }

Запустите все микросервисы и приложение готово к локальному использованию
