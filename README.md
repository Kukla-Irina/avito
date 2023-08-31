
<h1>Сервис динамического сегментирования пользователей</h1>

Сервис, хранящий пользователя и сегменты, в которых он состоит (создание, изменение, удаление сегментов, а также добавление и удаление пользователей в сегмент).

Инструменты, использованные при создании:
1. Реляционная БД PostgreSQL
2. Table Plus в качестве клиента БД
3. Язык Golang
4. Docker и Docker Compose для запуска сервиса

<h2>Подготовка</h2>

Перед запуском сервиса необходимо ввести данные БД, которая будет подключена. В файле docker-compose.yml необходимо ввести следующие данные:
1. DATABASE_URL: "host=go_db user={user} password={password} dbname={dbname} sslmode=disable" (пример DATABASE_URL: "host=go_db user=postgres password=postgres dbname=postgres sslmode=disable")
2. POSTGRES_PASSWORD: {password}
   POSTGRES_USER: {user}
   POSTGRES_DB: {dbname} 
   Пароль, пользователь и название БД должны совпадать с данными из п. 1.
   Пример:
   POSTGRES_PASSWORD: postgres
   POSTGRES_USER: postgres
   POSTGRES_DB: postgres
3. Номер порта, который использует БД
    ports:
      - "{port}:{port}"
    
    Пример:
    
    ports:
      - "5432:5432"

<h2>Запуск</h2>

Для запуска использовать команду § docker compose up. При запуске сервиса в БД создается пустая таблица с названием "segments".

<h2>Примеры запросов и ответов</h2>

<h3>Создание сегмента и добавление пользователю</h3>

Запрос

<code>curl --location 'localhost:8000/segments/add' \
--header 'Content-Type: application/json' \
--data '{
    "name": "AVITO_VOICE_MESSAGES",
    "userid": 1000
}'</code>

Ответ 

<code>{
    "id": 1,
    "name": "AVITO_VOICE_MESSAGES",
    "userid": 1000
}</code>

<h3>Удаление сегмента полностью из БД у всех пользователей</h3>

Запрос

<code>curl --location --request DELETE 'localhost:8000/segments/deleteall/AVITO_VOICE_MESSAGES'</code>

Ответ 

<code>"Segment deleted from all users"</code>

<h3>Удаление конкретного сегмента у конкретного пользователя</h3>

Запрос

<code>curl --location --request DELETE 'localhost:8000/segments/delete/AVITO_VOICE_MESSAGES/1000'</code>

Ответ 

<code>"Segment deleted"</code>

<h3>Поиск всех сегментов по id пользователя</h3>

Запрос

<code>curl --location 'localhost:8000/segments/1000'</code>

Ответ 

<code>[
    {
        "id": 5,
        "name": "AVITO_PERFORMANCE_VAS",
        "userid": 1000
    },
    {
        "id": 6,
        "name": "AVITO_DISCOUNT_30",
        "userid": 1000
    },
    {
        "id": 7,
        "name": "AVITO_VOICE_MESSAGES",
        "userid": 1000
    }
]</code>

<h3>Дополнительные комментарии к заданию</h3>

Было принято решение создать таблицу с привязкой названия сегмента к id пользователя (CONSTRAINT add UNIQUE (name, userid)). Это заставляет каждому добавлению сегмента пользователю присваивать отдельный id, что запрещает полное дублирование связки и позволяет сегментам не перетираться друг с другом и легко быть отсортированными по id пользователя. Так как оба значения не могут быть NULL, то считаем по умолчанию, что пользователь не включен ни в один сегмент, если по его id ничего не найдено. Пример заполненной таблицы находится в файле segments.sql