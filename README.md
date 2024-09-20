---
# Workflow Management System (Go-Gin)

## Table of Contents

1. [Introduction](#introduction)
2. [Features](#features)
3. [Technologies](#technologies)
4. [Installation](#installation)
5. [Environment Variables](#environment-variables)
6. [Database Setup](#database-setup)
7. [API Endpoints](#api-endpoints)
8. [Authentication](#authentication)
9. [Running the Project](#running-the-project)
10. [Graceful Shutdown](#graceful-shutdown)
11. [Contributing](#contributing)
12. [License](#license)
---

## Introduction

This project is a **Workflow Management System** built using **Go** and the **Gin** web framework. It offers a robust API for managing workflows with features such as JWT-based authentication, item management (CRUD operations), and role-based access control. The system is designed to be highly modular and extendable.

---

## Features

- **JWT Authentication**: Secure user authentication and authorization using JWT tokens.
- **Item Management**: Full support for creating, reading, updating, and deleting (CRUD) items.
- **Role-Based Access Control**: Admin-only routes for critical operations such as bulk status updates.
- **CORS Support**: Allows cross-origin resource sharing with multiple front-end clients.
- **Database Versioning**: Keep track of applied migrations with Goose and versioning.

---

## Technologies

- **Go (Golang)**: Backend programming language.
- **Gin Framework**: Lightweight HTTP web framework for Go.
- **GORM**: ORM for Go, with PostgreSQL as the database driver.
- **JWT (JSON Web Token)**: For stateless authentication.
- **Cors**: Middleware for handling CORS (Cross-Origin Resource Sharing).
- **Goose**: Database migration tool to track schema changes.
- **PostgreSQL**: Database for storing items and user data.

---

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Kiratopat-s/midterm-test
   cd midterm-test
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Docker up PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

---

## Database Setup

**Migreate Goose**:
    
    ```bash
    cd mirgrations

    goose postgres "postgres://postgres:postgres@localhost:5432/iws" up
    ```
---

## API Endpoints

The following endpoints are available:

| Method | Endpoint                    | Description                                 | Auth Required |
| ------ | --------------------------- | ------------------------------------------- | ------------- |
| GET    | `/version`                  | Get the current database version            | No            |
| GET    | `/hello`                    | Simple Hello World response                 | No            |
| GET    | `/hello-verifytoken`        | Hello World with JWT verification           | Yes           |
| POST   | `/items`                    | Create a new item                           | Yes           |
| GET    | `/items`                    | Fetch all items                             | Yes           |
| GET    | `/items/:id`                | Fetch an item by ID                         | Yes           |
| PUT    | `/items/:id`                | Update an item by ID                        | Yes           |
| PATCH  | `/items/:id`                | Update the status of an item                | Yes           |
| PATCH  | `/items/update/status/many` | Update the status of multiple items (Admin) | Yes (Admin)   |
| DELETE | `/items/:id`                | Delete an item                              | Yes           |
| DELETE | `/items/delete/many`        | Delete multiple items                       | Yes           |
| GET    | `/items/status/count/user`  | Count items by user and status              | Yes           |
| POST   | `/login`                    | User login                                  | No            |
| POST   | `/register`                 | User registration                           | No            |

---

## Authentication

This system uses **JWT** for user authentication. The JWT token must be sent with each request that requires authentication.

- **Token Verification Middleware (`verifyToken`)**: Protects routes by ensuring that users provide valid JWT tokens.
- **Admin Guard (`verifyAdmin`)**: Restricts access to admin-only routes.

---

## Running the Project

1. **Start the server**:

   ```bash
   DATABASE_URL=postgres://postgres:postgres@localhost:5432/iws PORT=2024 JWT_SECRET=kirato go run cmd/main.go
   ```

   The server will start on the port defined in the environment variables or default to `8080`.

---

## Graceful Shutdown

This project includes a graceful shutdown process that listens for system signals (SIGINT, SIGTERM) and allows the server to complete existing requests before shutting down. It waits for up to **60 seconds** to ensure all active connections are closed.

---

<!-- ## Contributing

Contributions are welcome! If you would like to contribute:

1. Fork this repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add a new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

--- -->
