# go-ecommerce-api

# Knowledge in project

- Migration
- Seeder
- Command
- MVC model
- Gorm, Gin, PostgreSql
- Docker, docker-compose

# Functionality

- Authentication (login, register, logout)
- Authorization (admin, user)
- Cart basic + checkout
- Manage resource (product, category) for admin
- Manage profile for user

# How to use ?

1. Install Golang with version > 1.2, Install gow to run server local: https://github.com/mitranim/gow
2. **git clone https://github.com/donghuynh99/go-ecommerce-api.git**
3. Go to project
4. **go mod tidy**
5. **docker-compose up -d** (Create postgreSql and get credential to connect)
6. **cp .env.example .env**
7. Update .env file like this

```
APP_NAME=Ecommerce API
APP_PORT=8080
APP_URL=http://localhost:8080
POSTGRESQL_DB=ecommerce_api
POSTGRESQL_HOST=localhost
POSTGRESQL_PORT=5432
POSTGRESQL_USER=root
POSTGRESQL_PASSWORD=pass
LANGUAGE=en
```

8. **go run main.go db:migrate**
9. **go run main.go generate:admin admin@gmail.com pass** (You also change email and password that you want)
10. Import file collection **ecommerce_api.json** to your postman application. And Import file enviroment **ercommerce_api_environment.json**.
11. **~/go/bin/gow run .** To serve your application.
