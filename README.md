# Dealls
Dealls backend engineer technical test

# HR Payroll System

A modular, scalable HR Payroll System built with Go and PostgreSQL.  
It handles employee attendance, overtime, reimbursements, and generates payroll for a given period.

---

## Features

- REST API
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

## Architecture

The architecture documentation can be read [here!](doc/ARCHITECTURE.md) (doc/ARCHITECTURE.md)

---

## Requirements

To run this project, make sure you have the following installed:

- **[Go 1.24+](https://go.dev/dl/)**  
- **[golang-migrate (postgresql driver)](https://github.com/golang-migrate/migrate)** – for database migrations  
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

## Testing
### Unit Test
Run the command below to run unit tests for internal components
```bash
go test ./internal/ -coverprofile="cover.out"
go tool cover -html="cover.out"
```
### Manual API Test
To test the API manually, you can follow these instructions and use your preferred way to send HTTP request (e.g. Postman, Insomnia, curl).
#### 1. Login as Admin
To login as admin simply send this POST HTTP request to `localhost:8080/login`
```json
{
    "username": "admin",
    "password": "admin"
}
```
> **_NOTE:_**  Change the password value to what you set in `.env` file.
and you will receive the access token in the response like below.
```json
{
	"message": "login success",
	"data": {
		"token": "token"
	}
}
```
Put the token in Authorization header as Bearer Type for every subsequent requests. For Example:
```bash
curl --request POST \
  --url http://localhost:8080/payroll/period \
  --header 'Authorization: Bearer <PUT YOUR TOKEN HERE>' \
  --header 'Content-Type: application/json' \
  --data '{
	"start_date": "2025-05-25T00:00:00+07:00",
	"end_date": "2025-06-25T00:00:00+07:00"
}'
```
### 1.5 Login as User
To login as user, just send a similar HTTP request but with username value that can be found in `hr.users` table and `password` as their password.

I have also provide 3 static users that can be used for testing
| user_id | username | password |
|---|---|---|
|  81d1bcd4-d5b3-4495-92ce-ef2c9b0f5e54 |  ani | password  |
|  8f29acd8-c18a-4e1c-9662-f102562bc893 |  budi | password  |
|  cc3a57a3-79cf-438e-9dc3-3a18bd86480b |  coki | password  |

### 2. Set Payroll Period
This endpoint is used to set the active payroll period for calculation.
```bash
curl --request POST \
  --url http://localhost:8080/payroll/period \
  --header 'Authorization: Bearer <PUT YOUR TOKEN HERE>' \
  --header 'Content-Type: application/json' \
  --data '{
	"start_date": "2025-05-25T00:00:00+07:00",
	"end_date": "2025-06-25T00:00:00+07:00"
}'
```
> **_NOTE:_**  This operation can only be done by admin. So use the admin's token you got from step 1.

### 3. Submit Attendance
This endpoint is used to submit attendance for specific user ID.
```bash
curl --request POST \
  --url http://localhost:8080/attendance \
  --header 'Authorization: Bearer <TOKEN>' \
  --header 'Content-Type: application/json' \
  --data '{
	"user_id":"cc3a57a3-79cf-438e-9dc3-3a18bd86480b",
	"timestamp": "2025-06-15T03:03:13.886Z"
}'
```
- `user_id` value denotes the which user's attendance is submitted.
- `timestamp` value denotes when the attendance happened. This is to allow retroactive filling by admin or similar cases.
> **_NOTE:_**  There is a TODO list to make this operation can be done only by the user itself and admin, by comparing the user ID in the body and the payload of the access token. But for now, the security measure done is just whether the request has valid access token.

### 4. Submit Overtime
This endpoint is used to submit overtime for specific user ID
```bash
curl --request POST \
  --url http://localhost:8080/overtime \
  --header 'Authorization: Bearer <TOKEN>' \
  --header 'Content-Type: application/json' \
  --data '{
	"user_id":"cc3a57a3-79cf-438e-9dc3-3a18bd86480b",
	"hours": 1,
	"timestamp": "2025-05-22T20:03:13.886+07:00"
}'
```
- `user_id` value denotes the which user's overtime is submitted.
- `hours` value denotes how many overtime hours worked.
- `timestamp` value denotes when the overtime work finished. This is to allow retroactive filling by admin or similar cases.
> **_NOTE:_**  There is a TODO list to make this operation can be done only by the user itself and admin, by comparing the user ID in the body and the payload of the access token. But for now, the security measure done is just whether the request has valid access token.

### 5. Submit Reimbursement
This endpoint is used to submit reimbursement request for specific user ID
```bash
curl --request POST \
  --url http://localhost:8080/reimbursement \
  --header 'Authorization: Bearer <TOKEN>' \
  --header 'Content-Type: application/json' \
  --data '{
	"user_id":"cc3a57a3-79cf-438e-9dc3-3a18bd86480b",
	"amount": 300000,
	"description": "buat judol hehe"
}'
```
- `user_id` value denotes the which user's overtime is submitted.
- `amount` value denotes how much is the amount requested.
- `description` value denotes the description for the reimbursement request.
> **_NOTE 1:_**  There is a TODO list to make this operation can be done only by the user itself and admin, by comparing the user ID in the body and the payload of the access token. But for now, the security measure done is just whether the request has valid access token.

> **_NOTE 2:_**  I don't use timestamp here because usually reimbursement is processed by when the request is made, instead of when the payment that is needed to be reimbursed is done.

### 6. Calculate Payroll
This endpoint is used to trigger payroll calculation for the active payroll period set in step 2. When done, there'll be immutable payslips data in `hr.payslips` table for the related active payroll period.
```bash
curl --request POST \
  --url http://localhost:8080/payroll/calculate \
  --header 'Authorization: Bearer <TOKEN>' \
```
The calculation process is explained below.
1. Get the active payroll period data.
2. Check whether this active payroll period is already processed/calculated.
3. Populate users/employees activities.
    1. Get all users attendances for the period.
    2. Get all users overtimes for the period.
    3. Get all users reimbursements for the period.
4. Get the salaries of the active users (listed in 3.1, 3.2, 3.3).
5. Setup channel for async process.
6. Spawn worker pool using goroutine.
7. Feed the worker with all the data from step 3 & 4.
8. Calculate each user's take home pay.
9. Store the details as payslips data in payslips table.
10. Mark the payroll period as processed.
> **_NOTE:_**  This operation can only be done by admin. So use the admin's token you got from step 1.

### 7. Get Payroll Period Summary
This endpoint is used to check the summary of the active payroll period
```bash
curl --request GET \
  --url http://localhost:8080/payroll/summary \
  --header 'Authorization: Bearer <TOKEN>' \
```
The response will contain how much the sum of the take home pay paid to employees, and its breakdown.
```json
{
	"message": "payroll summary in active period generated",
	"data": {
		"total_take_home_pay": 27902272.73,
		"payslips": [
			{
				"user_id": "8f29acd8-c18a-4e1c-9662-f102562bc893",
				"take_home_pay": 25284090.91,
				"name": "budi"
			},
			{
				"user_id": "cc3a57a3-79cf-438e-9dc3-3a18bd86480b",
				"take_home_pay": 2618181.82,
				"name": "coki"
			}
		]
	}
}
```
> **_NOTE 1:_**  This operation can only be done by admin. So use the admin's token you got from step 1.

> **_NOTE 2:_**  There is a TODO list to give this endpoint parameter to choose which payroll period to get the summary from. But for now, it can only be used to get the summary of the active payroll period.

### 8. Get User Payslips
This endpoint is used to get the payslips details of specific user ID in the active/latest payroll period.
```bash
curl --request GET \
  --url http://localhost:8080/payslip \
  --header 'Authorization: Bearer <TOKEN>' \
  --header 'Content-Type: application/json' \
  --data '{
	"user_id": "cc3a57a3-79cf-438e-9dc3-3a18bd86480b"
}'
```
The response will contain how the breakdown of the specified user ID's payslips.
```json
{
	"message": "payslip summary in active period fetched",
	"data": {
		"id": "ba63c186-b098-47ef-96e2-65db0f0353a1",
		"name": "coki",
		"user_id": "cc3a57a3-79cf-438e-9dc3-3a18bd86480b",
		"payroll_id": "af53a5f4-d489-4fa4-a29e-7bfe1b51006f",
		"base_salary": 17000000,
		"total_attendance": 3,
		"total_work_day": 22,
		"total_overtime_hour": 0,
		"overtime_bonus": 0,
		"reimbursement_list": [
			{
				"id": "9536ebba-cf42-48b4-9c45-830080d4bac2",
				"amount": 300000,
				"description": "buat judol juga hehe"
			}
		],
		"total_reimbursement_amount": 300000,
		"take_home_pay": 2618181.82
	}
}
```
> **_NOTE:_**  There is a TODO list to give this endpoint parameter to choose which payroll period to get the breakdown of the payslip from. But for now, it can only be used to get it from the active payroll period.