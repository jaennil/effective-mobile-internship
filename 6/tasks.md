### Задание 1. **Создание Docker контейнера**
- Создайте Docker контейнер с простым веб-сервером (например, на основе Nginx), который отображает "Hello, Docker!" при обращении к корневому URL.
    - Создайте файл **`Dockerfile`**
    - Создайте файл **`index.html`**
    - Постройте и запустите контейнер
    - Теперь, когда вы открываете браузер и переходите по адресу **`http://localhost:8080`**, вы должны увидеть "Hello, Docker!".

**`index.html`**
```html
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
    </head>
    <body>
        <h1>Hello, Docker!</h1>
    </body>
</html>
```

**`Dockerfile`**
```Dockerfile
FROM nginx:alpine

COPY index.html /usr/share/nginx/html/index.html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

**`Commands`**
```sh
docker build -t 6 .
docker run -p 8080:80 6
```

### Задание 2. **Передача аргументов в Docker контейнер**
- Измените предыдущий Dockerfile так, чтобы текст приветствия можно было передать в контейнер через аргументы

**`Dockerfile`**
```Dockerfile
FROM nginx:alpine

COPY index.html /usr/share/nginx/html/index.html

ARG GREETING="Hello, Docker!"

ENV GREETING_TEXT=$GREETING

RUN sed -i "s/<h1>.*<\/h1>/<h1>$GREETING_TEXT<\/h1>/" /usr/share/nginx/html/index.html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

**`Commands`**
```sh
docker build -t 6 . --build-arg GREETING="Hello, Custom Greeting!"
docker run -p 8080:80 6
```

### Задание 3. **Работа с многими контейнерами**
- Создайте файл **`docker-compose.yml`**, чтобы запустить одновременно два контейнера: один с веб-сервером Nginx, другой с базой данных PostgreSQL.

**`docker-compose.yml`**
```yml
version: '3.8'

services:
  web:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "8080:80"
    volumes:
      - ./html:/usr/share/nginx/html
    depends_on:
      - db

  db:
    image: postgres:alpine
    container_name: postgres
    environment:
      POSTGRES_USER: db_user
      POSTGRES_PASSWORD: db_pass
      POSTGRES_DB: db_name
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**`Commands`**
```sh
docker-compose up -d
```

### Задание 4. **Передача данных между контейнерами**
- Создайте Docker контейнер с приложением на Go, которое отправляет запрос к базе данных PostgreSQL (созданной в предыдущем задании) и выводит результат.
    - Создайте Dockerfile для приложения Go
    - Создайте простое приложение на Go (**`main.go`**)
    - Создайте **`docker-compose.yml`** для обоих контейнеров
    - Запустите контейнеры
    - Теперь ваше приложение Go отправляет запрос к базе данных PostgreSQL, запущенной в соседнем контейнере.
