# Venturo Golang Core

A robust, slim, and well-organized backend codebase core using Golang. This project serves as a comprehensive template and training tool for building professional, modern web services. It is built with a "Fat Model" (Active Record) pattern for simplicity and speed, and demonstrates advanced features like asynchronous job processing, graceful shutdowns, and structured logging.

## Table of Contents

  - [✨ Key Features & Concepts](https://www.google.com/search?q=%23-key-features--concepts)
  - [🏗️ Architectural Patterns](https://www.google.com/search?q=%23%EF%B8%8F-architectural-patterns)
  - [📂 Codebase Structure](https://www.google.com/search?q=%23-codebase-structure)
  - [⚙️ Environment Variables](https://www.google.com/search?q=%23%EF%B8%8F-environment-variables)
  - [🚀 Getting Started](https://www.google.com/search?q=%23-getting-started)
  - [🗄️ Database Migrations](https://www.google.com/search?q=%23-database-migrations)
  - [📖 API Documentation](https://www.google.com/search?q=%23-api-documentation)
  - [🧩 Creating a New Feature](https://www.google.com/search?q=%23-creating-a-new-feature)
  - [📚 Core Libraries](https://www.google.com/search?q=%23-core-libraries)

-----

## Answer for Learning Modules
These answer are on the [`answers/`](answers/) folder.

- [Module 1 B: Core Go Concepts - Challenges (with Starters)](answers/module-1)
- [Module 2 A: Concurrency & Advanced Reporting](answers/module-2-a)


## ✨ Key Features & Concepts

This codebase serves as a training tool and demonstrates several important backend concepts:

  * **Clean Architecture:** A clear separation of concerns is enforced between different layers of the application. The request flows from a **Handler** (which deals with HTTP) to a **Service** (which contains business logic) to a **Model** (which handles database interaction). This makes the code modular and easy to maintain.

  * **Authentication & Authorization:** A complete JWT-based authentication flow allows users to register and log in. Protected endpoints use a custom middleware to validate tokens. Authorization logic is implemented in the service layer to ensure users can only modify their own data.

  * **Asynchronous Processing:** Long-running tasks, like file uploads, are handled in the background using **goroutines**. This provides an immediate response to the user, improving their experience. A `sync.WaitGroup` is used to track these background jobs, ensuring they can complete before the server shuts down.

  * **Graceful Shutdown:** The application listens for OS signals (like `Ctrl+C`) to shut down gracefully. It waits for all background processes to finish before exiting, preventing data loss or corruption.

  * **Structured Logging:** Uses Go's standard `slog` library to produce machine-readable JSON logs. This is crucial for production environments, as it allows logs to be easily searched, filtered, and analyzed by log management platforms.

  * **Robust Tooling:** The entire development environment is containerized with **Docker** and `docker-compose`. For rapid development, **Air** is configured for hot-reloading, automatically recompiling and restarting the server whenever a Go file is saved.

  * **Standardized API:** All API responses and validation errors follow a consistent, predictable JSON structure defined in the `pkg/response` package. This creates a professional and easy-to-use API for any client.

-----

## 🏗️ Architectural Patterns

This project intentionally uses several common software design patterns.

### \#\#\# Dependency Injection (DI)

Instead of creating dependencies inside functions, we "inject" them from a higher level. For example, when the server starts, it creates the database connection and the logger, and then passes them into the services that need them. This makes our components decoupled and much easier to test.

### \#\#\# Adapter Pattern

To interact with third-party services (like Google Cloud Storage), we use an adapter. We first define a generic `StorageAdapter` interface, which defines the methods we need (e.g., `Upload`). Then, we create a concrete struct (`GCSAdapter`) that implements this interface. This allows us to easily swap out GCS for another service (like AWS S3) in the future by simply creating a new adapter, without changing any of our business logic.

### \#\#\# Fat Model (Active Record) Pattern

For simplicity, this project places database logic directly into methods on the model structs (e.g., `user.Save(db)`). This pattern, similar to Laravel's Eloquent, makes CRUD operations very straightforward and keeps the number of layers to a minimum, which is great for learning and for smaller projects.

-----

## 📂 Codebase Structure

The project follows a clean architecture pattern designed for simplicity and scalability.

```
/
├── cmd/                  # Application entry points (main packages)
│   ├── migrate/          # The database migration tool.
│   └── server/           # The main API server.
├── configs/              # Configuration loading from the .env file.
├── database/             # SQL migration files managed by golang-migrate.
├── docs/                 # Auto-generated Swagger API documentation files.
├── internal/
│   ├── adapter/          # Adapters for 3rd party services (e.g., GCS).
│   ├── database/         # Database connection and migration logic.
│   ├── handler/http/     # HTTP Handlers (Controllers). They parse requests and call services.
│   ├── model/            # Data models and their database methods (Fat Model).
│   └── server/           # Server setup, dependency injection, and routing.
├── pkg/
│   ├── logger/           # Structured logger configuration.
│   ├── response/         # Standardized API response helpers.
│   ├── uploader/         # Generic file upload utility.
│   └── validator/        # Input validation utility.
├── .air.toml             # Configuration for hot-reloading with Air.
├── docker-compose.yml    # Docker service definitions for app and db.
└── Dockerfile            # Docker build instructions.
```

-----

## ⚙️ Environment Variables

Create a `.env` file in the project root and configure the following variables:

| Variable         | Description                                     | Example                      |
| ---------------- | ----------------------------------------------- | ---------------------------- |
| `DB_HOST`        | The hostname of the MySQL database.             | `db` (for Docker), `127.0.0.1` (for local) |
| `DB_PORT`        | The port your MySQL database is running on.     | `3306`                       |
| `DB_USER`        | The username for the MySQL database.            | `root`                       |
| `DB_PASSWORD`    | The password for the database user.             | `your_password`              |
| `DB_NAME`        | The name of the database to use.                | `venturo_db`                 |
| `JWT_SECRET_KEY` | A long, random, secret string for signing JWTs. | `super-secret-key`           |

-----

## 🚀 Getting Started

*Instructions for both Docker and Local setups are provided.*

### Prerequisites

  * Go 1.24+
  * Docker & Docker Compose
  * MySQL (only if running locally without Docker)

### Configuration

1.  Copy the example environment file: `cp .env.example .env`
2.  Edit the `.env` file with your specific configuration.

### Running with Docker (Recommended)

1.  **Update `.env` file:** Ensure `DB_HOST` is set to `db`.
2.  **Build and Run:**
    ```bash
    docker-compose up --build
    ```
3.  **Run Migrations:** Open a new terminal and run:
    ```bash
    docker-compose run --rm app go run ./cmd/migrate/main.go up
    ```

The API server will be available at `http://localhost:3000`.

### Local Development (Without Docker)

1.  **Install `migrate` CLI:**
    ```bash
    go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
2.  **Update `.env` file:** Ensure `DB_HOST` is set to `127.0.0.1`.
3.  **Install Go Dependencies:**
    ```bash
    go mod tidy
    ```
4.  **Run Database Migrations:**
    ```bash
    migrate -database 'mysql://your_user:your_password@tcp(127.0.0.1:41181)/venturo-golang-core' -path database/migrations up
    ```
5.  **Run the Server:**
    ```bash
    go run ./cmd/server/main.go
    ```

The API server will be available at `http://localhost:3000`.

-----

## 🗄️ Database Migrations

Use the appropriate command for your chosen development environment to manage the database schema. The `fresh` command is particularly useful for resetting the database during development.

-----

## 📖 API Documentation

This project uses `swaggo` to generate interactive API documentation from code comments.

1.  **Example Annotation:** Comments are added directly above handler functions like this:
    ```go
    // @Summary      Get User Profile
    // @Description  Retrieves the profile of the currently authenticated user.
    // @Tags         User
    // @Produce      json
    // @Security     ApiKeyAuth
    // @Success      200  {object}  response.ApiResponse{data=model.User}
    // @Router       /profile [get]
    func (h *UserHandler) GetProfile(c *fiber.Ctx) error { ... }
    ```
2.  **Generate or Update Docs:** After changing annotations, run this command locally:
    ```bash
    swag init -g ./cmd/server/main.go --parseDependency --parseInternal
    ```
3.  **View Docs:** Start the server and navigate to:
    **http://localhost:3000/swagger/index.html**

-----

## 🧩 Creating a New Feature

Follow this pattern to add a new feature (e.g., "Comments"):

1.  **Create the Migration:** Use `migrate create -ext sql ...` to generate files. Edit the SQL to define your new table.
2.  **Run the Migration:** Apply the migration using the command for your environment.
3.  **Create the Model:** Create `internal/model/comment_model.go`, defining the struct and its database methods.
4.  **Create the Service:** Create `internal/service/comment_service.go` for the business logic.
5.  **Create the Handler:** Create `internal/handler/http/comment_handler.go`. Add `validate` tags to your payload structs and use the helpers from `pkg/response` and `pkg/validator`.
6.  **Register Routes:** In `internal/server/routes.go`, initialize your new service/handler and register the new API endpoints.

-----

## 📚 Core Libraries

  * [**Fiber**](https://github.com/gofiber/fiber): Web framework.
  * [**GORM**](https://github.com/go-gorm/gorm): Database ORM.
  * [**golang-migrate**](https://github.com/golang-migrate/migrate): Database migrations.
  * [**go-playground/validator**](https://github.com/go-playground/validator): Struct validation.
  * [**swaggo**](https://github.com/swaggo/swag): API documentation generator.
  * [**Air**](https://github.com/air-verse/air): Live-reloading for development.
  * [**go-jwt**](https://github.com/golang-jwt/jwt): JSON Web Token implementation.
  * [**Bcrypt**](https://pkg.go.dev/golang.org/x/crypto/bcrypt): Password hashing.