## Database Blueprints in Go

Follow these steps to set up and run the project on your local machine.

### Prerequisites

- [Go](https://go.dev/dl/) (version 1.22+).
- [Docker](https://www.docker.com/products/docker-desktop/) and Docker Compose.
- [golang-migrate](https://github.com/golang-migrate/migrate) to run migrations.
- Install on macOS: `brew install golang-migrate`
- Or install with Go: `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### Installation Steps

1.  **Clone the repository**

    ```bash
    git clone [https://github.com/your-username/go-db-blueprints.git](https://github.com/your-username/go-db-blueprints.git)
    cd go-db-blueprints
    ```

2.  **Configure the Environment**
    Copy the `.env.example` file to `.env` and fill in your database credentials. The default values are already configured to work with the `docker-compose.yml` file provided.

    ```bash
    cp .env.example .env
    ```

    The `.env` file content:

    ```env
    DB_DRIVER=mysql
    DB_USER=user
    DB_PASSWORD=password
    DB_HOST=localhost
    DB_PORT=3306
    DB_NAME=blueprints_db
    ```

## How to Run the Examples

After completing the setup, you can run each example with the following commands:

- **Run the `database/sql` example:**

  ```bash
  go run ./cmd/database_sql/main.go
  ```

- **Run the `GORM` example:**

  ```bash
  go run ./cmd/gorm/main.go
  ```

- **Run the `sqlx` example:**
  ```bash
  go run ./cmd/sqlx/main.go
  ```

## How to Add a New Pattern

1.  Create a new directory in `blueprints/` (e.g., `blueprints/ent/`).
2.  Implement the repository logic in that directory.
3.  Create a new example directory in `cmd/` (e.g., `cmd/example-ent/`).
4.  Write the `main.go` file inside it to call the logic from the `blueprints` package you just created.
5.  Update this `README.md` to add your new pattern to the list!

## License

This project is distributed under the MIT License. See the `LICENSE` file for more information.
