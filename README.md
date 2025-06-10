# Dealls
Dealls backend engineer technical test

# HR Payroll System

A modular, scalable HR Payroll System built with Go and PostgreSQL.  
It handles employee attendance, overtime, reimbursements, and generates payroll for a given period.

---

## Features

- REST API with OpenAPI 3 support
- Modular architecture (handler → logic → repository)
- Payroll generation based on:
  - Attendance records
  - Overtime records
  - Reimbursements
  - Salary configuration
- Concurrent payslip generation with limited worker pool
- Clean separation of logic and infrastructure
- Database migration support
- Dockerized environment

---

## Requirements

To run this project, make sure you have the following installed:

- **[Go 1.24+](https://go.dev/dl/)**  
- **[golang-migrate](https://github.com/golang-migrate/migrate)** – for database migrations  
- **[Docker & Docker Compose](https://docs.docker.com/compose/)** – for local development setup

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/rahadianir/dealls.git
cd dealls
```

### 2. Setup environment file
Setup the values in `.env` and `docker-compose.yml` file as preferred. The default value should suffice to run the service.

### 3. Run infrastructure and migration setup
```bash
chmod +x setup.sh
./setup.sh
```
This will store `.env` values into the environment variables, spin up a dockerized postgresql instance, download all needed go packages for the app, and run the database migration schema into the database.

### 4. Run the server
```bash
go run main.go http
```
This will run the app that listens to port `:8080`.