# Asynchronous Notification Service with Go & RabbitMQ

This repository contains a Proof of Concept (PoC) for a decoupled, resilient, and scalable notification microservice, built with Go and RabbitMQ.

## 🎯 Project Goal

### The Business Problem
In modern software architectures, especially those based on microservices, common tasks like sending emails, SMS, or push notifications can become performance bottlenecks and single points of failure. If a critical service (e.g., Billing of whatever) depends on a synchronous call to an external email API, the slowness or unavailability of that provider can directly impact core business functionality.

This application solves that problem by implementing a **centralized and asynchronous notification service**. It allows any other service on the platform to "fire and forget" a notification request, with the guarantee that it will be processed safely and resiliently in the background, without impacting the performance of the original operation.

### Learning Objectives
This project serves as a practical lab to study and put in practice the following backend engineering concepts:
* **Event-Driven Architecture:** Understanding how to decouple services using a Message Broker.
* **Messaging Patterns:** Implementing the Producer/Consumer and Work Queue patterns.
* **System Resilience:** Building services that can withstand temporary component failures.
* **Concurrent Go:** Using Go to build high-performance network services.
* **Containerization:** Orchestrating a multi-service environment with Docker and Docker Compose.

## 📖 Hypothetical Use Case: The "Nexus" Platform

To provide context, imagine this service is part of a larger SaaS platform called "Nexus," which has several microservices:
* **Accounts Service:** Manages user registration and profiles.
* **Billing Service:** Handles subscriptions and processes monthly payments.
* **Marketing Service:** Sends newsletters.

All of these services are clients (Producers) of our Notification API. When the `Accounts Service` needs to send a welcome email, it simply makes a fast, lightweight `POST` call to our API and moves on. Our API enqueues the message, and our `Worker` processes it in the background, ensuring the user registration flow is always fast and never fails because of the email system.

## 🛠️ Tech Stack & Concepts

### Core Technologies
* **Go (v1.25+):** The programming language used for the API and the Worker.
* **Echo:** A high-performance, minimalist Go web framework for the API.
* **RabbitMQ:** The Message Broker, responsible for managing queues and guaranteeing message delivery.
* **Docker & Docker Compose:** For containerizing and orchestrating the entire development environment.

### Messaging Patterns (Core Focus)
This project implements the following fundamental messaging concepts:

* **Producer/Consumer Architecture:**
    * **Producer:** Our **REST API** acts as the producer. Its sole responsibility is to receive HTTP requests, validate the data, and publish a message to the queue.
    * **Consumer:** Our **Worker** acts as the consumer. It is a background process that connects to the queue, receives messages, and performs the task (in this case, "sending" the email).

* **Work Queue (Point-to-Point):**
    * We use a simple queue (`email_notifications`) where each message is processed by a single worker. This pattern allows for easy scaling by adding more worker instances to consume from the same queue.

* **Decoupling:**
    * The API does not know (and doesn't need to know) anything about the Worker. It just hands the message off to RabbitMQ. This means we can update, restart, or scale the Worker independently without ever taking the API offline.

* **Persistence and Durability:**
    * To ensure no notifications are lost even if the RabbitMQ server restarts, we implement durability at two layers:
        1.  **Durable Queues:** The `email_notifications` queue is declared as `durable`.
        2.  **Persistent Messages:** Each message is published with the `Persistent` delivery mode, instructing RabbitMQ to save it to disk.

* **Message Acknowledgments (`Ack`):**
    * The Worker only sends an `ack` (acknowledgment) to RabbitMQ **after** successfully processing a message. If the Worker fails mid-process, the message is not `acked`, and RabbitMQ will re-queue it to be delivered to another worker (or the same one when it restarts), guaranteeing at-least-once delivery.

## 🚀 How to Run

**Prerequisites:**
* Docker
* Docker Compose

**Steps:**
1.  Clone this repository.
2.  In the project root, create a `.env` file with the following content:
    ```env
    # Port the API will be exposed on your host machine
    API_PORT_DOCKER=9919

    # RabbitMQ connection string for the containers
    RABBITMQ_URL=amqp://user:password@rabbitmq:5672/
    ```
3.  Build and run the containers:
    ```bash
    docker-compose up --build
    ```
4.  The API will be available at `http://localhost:9919`.
5.  The RabbitMQ management UI will be available at `http://localhost:15672` (login: `user`, password: `password`).

## API Endpoints

### Send an Email

* **`POST /api/v1/notifications/email`**

**Request Body (JSON):**
```json
{
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "This is the body of your message."
}
```

**Success Response**

- **Code:** `202 Accepted`  
- **Body:**
```json
{
  "message": "request accepted and is being processed"
}
```

# Asynchronous Notification Service with Go & RabbitMQ

This repository contains a Proof of Concept (PoC) for a decoupled, resilient, and scalable notification microservice, built with Go and RabbitMQ.

## 🎯 Project Goal

### The Business Problem
In modern software architectures, especially those based on microservices, common tasks like sending emails, SMS, or push notifications can become performance bottlenecks and single points of failure. If a critical service (e.g., Billing of whatever) depends on a synchronous call to an external email API, the slowness or unavailability of that provider can directly impact core business functionality.

This application solves that problem by implementing a **centralized and asynchronous notification service**. It allows any other service on the platform to "fire and forget" a notification request, with the guarantee that it will be processed safely and resiliently in the background, without impacting the performance of the original operation.

### Learning Objectives
This project serves as a practical lab to study and put in practice the following  backend engineering concepts:
* **Event-Driven Architecture:** Understanding how to decouple services using a Message Broker.
* **Messaging Patterns:** Implementing the Producer/Consumer and Work Queue patterns.
* **System Resilience:** Building services that can withstand temporary component failures.
* **Concurrent Go:** Using Go to build high-performance network services.
* **Containerization:** Orchestrating a multi-service environment with Docker and Docker Compose.

## 📖 Hypothetical Use Case: The "Nexus" Platform

To provide context, imagine this service is part of a larger SaaS platform called "Nexus," which has several microservices:
* **Accounts Service:** Manages user registration and profiles.
* **Billing Service:** Handles subscriptions and processes monthly payments.
* **Marketing Service:** Sends newsletters.

All of these services are clients (Producers) of our Notification API. When the `Accounts Service` needs to send a welcome email, it simply makes a fast, lightweight `POST` call to our API and moves on. Our API enqueues the message, and our `Worker` processes it in the background, ensuring the user registration flow is always fast and never fails because of the email system.

## 🛠️ Tech Stack & Concepts

### Core Technologies
* **Go (v1.25+):** The programming language used for the API and the Worker.
* **Echo:** A high-performance, minimalist Go web framework for the API.
* **RabbitMQ:** The Message Broker, responsible for managing queues and guaranteeing message delivery.
* **Docker & Docker Compose:** For containerizing and orchestrating the entire development environment.

### Messaging Patterns (Core Focus)
This project implements the following fundamental messaging concepts:

* **Producer/Consumer Architecture:**
    * **Producer:** Our **REST API** acts as the producer. Its sole responsibility is to receive HTTP requests, validate the data, and publish a message to the queue.
    * **Consumer:** Our **Worker** acts as the consumer. It is a background process that connects to the queue, receives messages, and performs the task (in this case, "sending" the email).

* **Work Queue (Point-to-Point):**
    * We use a simple queue (`email_notifications`) where each message is processed by a single worker. This pattern allows for easy scaling by adding more worker instances to consume from the same queue.

* **Decoupling:**
    * The API does not know (and doesn't need to know) anything about the Worker. It just hands the message off to RabbitMQ. This means we can update, restart, or scale the Worker independently without ever taking the API offline.

* **Persistence and Durability:**
    * To ensure no notifications are lost even if the RabbitMQ server restarts, we implement durability at two layers:
        1.  **Durable Queues:** The `email_notifications` queue is declared as `durable`.
        2.  **Persistent Messages:** Each message is published with the `Persistent` delivery mode, instructing RabbitMQ to save it to disk.

* **Message Acknowledgments (`Ack`):**
    * The Worker only sends an `ack` (acknowledgment) to RabbitMQ **after** successfully processing a message. If the Worker fails mid-process, the message is not `acked`, and RabbitMQ will re-queue it to be delivered to another worker (or the same one when it restarts), guaranteeing at-least-once delivery.

## 🚀 How to Run

**Prerequisites:**
* Docker
* Docker Compose

**Steps:**
1.  Clone this repository.
2.  In the project root, create a `.env` file with the following content:
    ```env
    # Port the API will be exposed on your host machine
    API_PORT_DOCKER=9919

    # RabbitMQ connection string for the containers
    RABBITMQ_URL=amqp://user:password@rabbitmq:5672/
    ```
3.  Build and run the containers:
    ```bash
    docker-compose up --build
    ```
4.  The API will be available at `http://localhost:9919`.
5.  The RabbitMQ management UI will be available at `http://localhost:15672` (login: `user`, password: `password`).

## API Endpoints

### Send an Email

* **`POST /api/v1/notifications/email`**

**Request Body (JSON):**
```json
{
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "This is the body of your message."
}
```

**Success Response**

- **Code:** `202 Accepted`  
- **Body:**
```json
{
  "message": "request accepted and is being processed"
}
```