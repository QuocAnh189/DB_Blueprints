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
    git clone https://github.com/QuocAnh189/DB_Blueprints
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
  make db_sql
  ```

- **Run the `GORM` example:**

  ```bash
  make gorm
  ```

- **Run the `sqlx` example:**
  ```bash
  make sqlx
  ```

## License

This project is distributed under the MIT License. See the `LICENSE` file for more information.
