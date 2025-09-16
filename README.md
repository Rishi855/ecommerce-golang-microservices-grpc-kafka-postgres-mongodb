# 🛒 E-commerce Microservices with Go, gRPC, Kafka, Postgres & MongoDB

This project is a **microservices-based e-commerce system** built with **Golang**, using **gRPC** for communication, **Kafka** for event-driven messaging, and **Postgres/MongoDB** for data storage.  
It demonstrates how to design and implement a **scalable, event-driven, cloud-native system** using modern technologies.

---

## 📌 Architecture Overview

### Components
- **Client (Users)** → Interacts with the system via **HTTP (port 8000)** through the API Gateway.
- **API Gateway** → Entry point for all client requests, routes traffic to microservices over **gRPC**.
- **Auth Service** → Handles authentication, user management, and JWT.  
  - DB: **Postgres** (port `50051`)
- **Order Service** → Manages orders, order states, and inventory tracking.  
  - DB: **Postgres** (port `50052`)
- **Product Service** → Handles product catalog, availability, and pricing.  
  - DB: **Postgres** (port `50053`)
- **Notification Service** → Sends notifications (e.g., email, SMS) on order status changes.  
  - DB: **Postgres**
- **Log Service** → Centralized logging of service activity and events.  
  - DB: **MongoDB**

### Event-Driven with Kafka
**Kafka** is used to handle events such as:
- Order Created  
- Product Stock Updated  
- Notification Events  
- Logging & Monitoring  

<img width="2183" height="1143" alt="image" src="https://github.com/user-attachments/assets/2d4b1b0b-fab4-43f0-9453-5dffb9e5c050" />

---

## 🚀 Features

- ✅ **Microservices architecture** with independent scaling  
- ✅ **gRPC communication** for high-performance service-to-service communication  
- ✅ **API Gateway** for unified client access  
- ✅ **Event-driven messaging** with Kafka  
- ✅ **Postgres** for transactional data  
- ✅ **MongoDB** for logs & analytics  
- ✅ **Dockerized deployment** for easy local setup  

---

## 🛠️ Tech Stack

- **Language**: Golang  
- **API Gateway**: gRPC + HTTP  
- **Message Broker**: Apache Kafka  
- **Databases**:  
  - Postgres (Product, Order, Auth, Notification)  
  - MongoDB (Logs)  
- **Containerization**: Docker & Docker Compose  

---

## 📂 Project Structure
<img width="587" height="500" alt="image" src="https://github.com/user-attachments/assets/7f9d267a-1135-40a6-a234-25be928cc88c" />


# ⚡ Getting Started

### 1️⃣ Clone Repository
```bash
git clone https://github.com/Rishi855/ecommerce-golang-microservices-grpc-kafka-postgres-mongodb.git
cd ecommerce-golang-microservices-grpc-kafka-postgres-mongodb
```
### 2️⃣ Start Services with Docker
```bash
docker-compose up --build
```
### 3️⃣ Access API Gateway
```bash
http://localhost:8000
```
# 🧪 Example Workflow

- User registers/login → handled by Auth Service

- Browse products → from Product Service

- Place order → triggers Order Service, publishes Kafka event

- Kafka events → update Product stock, send Notification, log activity in MongoDB

- Logs & Notifications → visible via Log Service and Notification Service

# 📬 Future Improvements

- Add payment service

- Implement API rate limiting at Gateway

- Add monitoring with Prometheus + Grafana

- Enhance CI/CD with GitHub Actions
