Данный репозиторий содержит сервис авторизации по курсу ["Микросервисы как в BigTech-компаниях"](https://olezhek28.courses/).

Автор курса: [Олег Козырев](https://www.linkedin.com/in/olezhek28/)

Локально миграция накатывается в docker-compose.yaml, командой в терминале
> docker-compose up -d
Локально остановить контейнер с пострес можно командой
> docker-compose down

Адрес локального сервиса localhost:50051
Адрес удаленного сервиса 45.130.9.109:50061

Адрес локальной базы localhost:54321
Адрес удаленной базы 45.130.9.109:54331

TODO
- [x] Переименовать DbWorker
- [x] Передавать контекст при вызове функций, а не хранить его в структуре
- [x] Вынести в константы поля БД
- [x] Задавать поле createdAt на уровне базы
- [x] Убрать Enum из Postgres
- [x] Update делать через Patch
- [x] Убрать log.Fatal, заменить на возврат ошибок
- [x] Доработать возврат ошибок на верхний уровень
- Убрать мертвый код