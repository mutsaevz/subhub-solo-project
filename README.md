
## Функционал сайта

- Авторизация
- Просмотр категорий рецептов
- Поиск рецепта по названию
- Добавление рецепта в избранное, а также удаление оттуда
- Лайк на понравившийся рецепт
- Комментирование рецепта
- Просмотр раздела "Избранное"

![gif](https://github.com/Kharsaeva/project-recipes/blob/Aisha/IMG_1338.gif)

- При авторизации есть возможность добавлять новые рецепты, а также удалять их.

![gif](https://github.com/Kharsaeva/project-recipes/blob/Aisha/IMG_1337.gif)

--


## 🏗 Архитектурные решения

Проект построен с использованием **Clean Architecture**:

#### HTTP (Gin handlers)

↓

#### Service (business logic)

↓

#### Repository (data access)

↓

PostgreSQL / Redis

---

## 💻 Технологический стек

<p>
  <img alt="Golang" src="https://img.shields.io/badge/-Golang-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="Gin" src="https://img.shields.io/badge/-Gin-008ECF?style=for-the-badge&logo=gin&logoColor=white" />
  <img alt="GORM" src="https://img.shields.io/badge/-GORM-2D3748?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/-PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img alt="Redis" src="https://img.shields.io/badge/-Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" />
  <img alt="REST API" src="https://img.shields.io/badge/-REST_API-005571?style=for-the-badge&logo=fastapi&logoColor=white" />
  <img alt="JWT" src="https://img.shields.io/badge/-JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/-Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img alt="Docker Compose" src="https://img.shields.io/badge/-Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img alt="Swagger" src="https://img.shields.io/badge/-Swagger_(OpenAPI)-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" />
  <img alt="Unit Tests" src="https://img.shields.io/badge/-Unit_Tests-25A162?style=for-the-badge&logo=go&logoColor=white" />
</p>

---

Такой подход позволяет:
- тестировать бизнес-логику изолированно  
- изменять способ хранения данных без переписывания логики  
- масштабировать функциональность без архитектурных изменений  

---

## ✅ Индикаторы качества реализации

- [x] Clean Architecture  
- [x] Разделение `handler / service / repository`  
- [x] Интерфейсы и слабая связанность  
- [x] Redis-кеширование с инвалидацией  
- [x] JWT-аутентификация и роли  
- [x] Unit-тесты сервисного слоя  
- [x] Docker и docker-compose  
- [x] Swagger-документация  

---

<p align="center">
  <strong>
    Проект демонстрирует не только знание Go,<br>
    но и понимание принципов построения поддерживаемых backend-сервисов.
  </strong>
</p>

---
