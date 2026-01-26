---

<p align="center">
  <strong>Effective Mobile — Test Assignment</strong><br>
  <sub>REST API for subscription aggregation</sub>
</p>

---

## 📍 Контекст проекта

Данный проект реализован в рамках тестового задания на позицию **Junior Golang Developer**.

При этом при разработке был сделан **осознанный акцент не на минимальное выполнение требований**,  
а на применение архитектурных и инженерных практик, используемых в production-разработке.

> Цель проекта — продемонстрировать подход к проектированию backend-сервисов,  
> качество архитектуры и понимание жизненного цикла приложения.

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

<p>
  <img alt="Go" src="https://img.shields.io/badge/-Golang-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="Gin" src="https://img.shields.io/badge/-Gin-008ECF?style=for-the-badge&logo=gin&logoColor=white" />
  <img alt="GORM" src="https://img.shields.io/badge/-GORM-2D3748?style=for-the-badge&logo=go&logoColor=white" />
  <img alt="PostgreSQL" src="https://img.shields.io/badge/-PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white" />
  <img alt="Redis" src="https://img.shields.io/badge/-Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white" />
  <img alt="Docker" src="https://img.shields.io/badge/-Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" />
  <img alt="Swagger" src="https://img.shields.io/badge/-Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" />
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


---
