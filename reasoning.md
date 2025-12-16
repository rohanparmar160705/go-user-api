# 🧠 Design Reasoning & Architectural Decisions

This document outlines the approach, architectural choices, and key decisions made during the development of the Go User API.

## 🏗️ Architectural Approach

I adopted a **Layered Architecture (Clean Architecture)** pattern to ensure separation of concerns, testability, and maintainability.

### The 3-Layer Flow:
1.  **Handler Layer (`/internal/handler`)**: 
    *   **Responsibility**: Manages HTTP transport (parsing JSON, reading params, validation).
    *   **Reasoning**: Decouples the API framework (Fiber) from the business logic. If we switched to `Gin` or standard `net/http` later, we would only rewrite this layer.
2.  **Service Layer (`/internal/service`)**: 
    *   **Responsibility**: Contains core business logic (e.g., Age Calculation).
    *   **Reasoning**: Ensures that business rules are centralized. The API handler doesn't know *how* usage is calculated, only that it needs it. This allowed me to easily write **Unit Tests** for the age logic without needing a running database/server.
3.  **Repository Layer (`/internal/repository`)**: 
    *   **Responsibility**: Pure data access using SQLC.
    *   **Reasoning**: abstracts the SQL details. The Service layer just asks for "CreateUser" and doesn't care if it's Postgres or MySQL.

---

## 🔑 Key Technical Decisions

### 1. SQL Generator (SQLC) vs. ORM (GORM)
**Decision**: strictly adhered to using **SQLC** as requested, but I would have chosen it regardless.
*   **Why**: standard ORMs (like GORM) rely heavily on runtime reflection, which impacts performance and can hide errors until the code actually runs.
*   **Benefit**: SQLC generates **type-safe Go code** from raw SQL queries. If I make a typo in my SQL schema or query, the Go compiler catches it immediately (Compile-time safety). It provides the performance of raw SQL with the developer experience of an ORM.

### 2. Dynamic Age Calculation
**Decision**: Store `DOB` (Date of Birth) but return `Age`.
*   **Why**: Storing `Age` in the database is an anti-pattern because it becomes stale instantly (every day at midnight).
*   **Implementation**: The `Age` is calculated strictly in the **Service Layer** at the moment of retrieval using `time.Now()`. This ensures 100% data accuracy without needing background jobs to update user rows.

### 3. Middleware Strategy
**Decision**: Implemented global middleware for **Request ID** and **Logging**.
*   **Why**: In a production environment (like Railway/AWS), tracing a specific failed request is impossible without a unique ID.
*   **Implementation**: I used Fiber's middleware to inject a UUID into the `X-Request-ID` header and log the processing duration. This satisfies the "Observability" bonus requirement.

### 4. Dependency Injection
**Decision**: `main.go` initializes dependencies (Database, Config) and passes them down.
*   **Why**: Global state prevents parallel testing. By passing the `Repository` into the `Service`, and `Service` into the `Handler`, the application is modular and easy to mock during testing.

---

## 🐳 Deployment & Ops

*   **Multi-Stage Dockerfile**: I used a multi-stage build process.
    *   *Stage 1 (Build)*: Uses the full `golang` image to compile.
    *   *Stage 2 (Run)*: Copies only the binary to a lightweight `Alpine` image.
    *   *Result*: The final container is extremely small (~25MB) and secure (no source code included).
*   **CI/CD**: Implemented GitHub Actions to automatically run tests and deploy to Railway on every push, ensuring production readiness.

## ✅ Requirements Coverage

| Feature | Implementation Details |
| :--- | :--- |
| **GoFiber** | Used for high-performance routing. |
| **SQLC + Postgres** | Type-safe DB layer generated from `db/queries`. |
| **Zap Logging** | Structured JSON logging for cloud compatibility. |
| **Validator** | Input validation struct tags (`validate:"required"`). |
| **Project Structure** | Strictly followed standard Go layout (`cmd`, `internal`, etc). |
| **Dynamic Age** | Calculated in `user_service.go`. |
| **Bonuses** | Docker, Pagination, Unit Tests, and Middleware all implemented. |
