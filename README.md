# companyEsgDB

This project is designed to help manage company information with a focus on **ESG** (Environmental, Social, and Governance) metrics. The application allows you to store and manage company data, including procurement methods, contacts, and ESG reports.

## Features
- Add, update, and delete companies
- Track procurement methods
- Extract data from websites using web scraping
- Export data to CSV and Excel formats

## Technologies Used
- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Frontend**: HTML, Bootstrap 5
- **Web Scraping**: Goquery, regex, custom parsers
- **Containerization**: Docker
- **Job Scheduling**: Cron for automatic updates

## How to Run

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/companyEsgDb.git
    cd companyEsgDb
    ```

2. Install dependencies and run using Docker:
    - Make sure you have [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
    - Create a `.env` file from the example:
    ```bash
    cp .env.example .env
    ```
    - Configure the `.env` file with the appropriate values (explained below).
    - Run the application with Docker Compose:
    ```bash
    docker-compose up --build
    ```

    The application will be accessible at `http://localhost:6061`.

## Environment Variables

The `.env` file is required for setting the environment variables. Here is an example:


DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=qwerty
DB_NAME=goDB
APP_PORT=6061


- **DB_HOST**: Database hostname (use `db` if using Docker Compose).
- **DB_PORT**: Database port (default `5432`).
- **DB_USER**: PostgreSQL user (default `postgres`).
- **DB_PASSWORD**: PostgreSQL password (default `qwerty`).
- **DB_NAME**: Database name (default `goDB`).
- **APP_PORT**: Application port (default `6061`).

## Docker Usage

The project uses Docker for containerization. Here's a breakdown of the services defined in the `docker-compose.yml` file:

1. **db (PostgreSQL)**:
    - **Image**: `postgres:16`
    - **Ports**: Exposes port `5433` for accessing the database.
    - **Healthcheck**: Ensures the database is ready before other services start.
    - **Volumes**: Persist data in a Docker volume.

2. **app (Go application)**:
    - **Ports**: Exposes port `6061` for the application.
    - **Depends on**: Waits until the database is healthy before starting.
    - **Build**: Builds the app from the Dockerfile.
    - **Environment variables**: The app uses the environment variables from the `.env` file to connect to the database.


Logs

You can view logs for both the app and the PostgreSQL container using:

docker-compose logs -f

To view logs for a specific service (e.g., the app):

docker-compose logs -f app
License

This project is licensed under the MIT License — see the LICENSE
 file for details.
