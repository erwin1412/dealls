# Payslip System Developer Guide

## Introduction

This guide provides a comprehensive tutorial for developers to use the Payslip System, a Go-based web application for managing employee attendance, overtime, reimbursements, payroll processing, and payslip generation. The system uses a RESTful API built with Echo, GORM, and PostgreSQL, following SOLID principles for maintainability and scalability.

The API supports two user roles: **Admin** and **Employee**, with JWT-based authentication. Admins can register users, create attendance periods, run payroll, and view payroll summaries. Employees can submit attendance, overtime, reimbursements, and view their payslips. All actions are logged in an audit trail for accountability.

This document includes:
- A list of API endpoints with their roles and requirements.
- Setup instructions for running the system.
- Detailed usage examples for each endpoint.
- Error handling and audit logging details.
- Testing instructions using `curl`.

---

## System Overview

### Architecture
- **Framework**: Echo v4 for HTTP routing.
- **Database**: PostgreSQL with GORM ORM, using UUIDs for primary keys.
- **Authentication**: JWT with role-based access control (Admin, Employee).
- **Audit Logging**: Tracks create/update actions with user ID, IP address, and request ID.
- **Module Path**: `https://github.com/erwin1412/dealls`.
- **Structure**: Follows SOLID principles with separated handlers, services, repositories, and models.

### Models
- **User**: Stores username, password hash, role (admin/employee), and salary.
- **AttendancePeriod**: Defines payroll periods with start and end dates.
- **Attendance**: Records employee attendance for specific dates.
- **Overtime**: Tracks overtime hours (max 3 hours/day).
- **Reimbursement**: Stores employee expense claims.
- **Payroll**: Calculates base salary, overtime pay, reimbursement, and total pay.
- **AuditLog**: Logs actions with details (action, table, record ID, user, IP, etc.).

---

## Endpoint Summary

The following table summarizes the API endpoints, their URLs, methods, required roles, whether they need a `period_id`, and authentication requirements.

| Endpoint                | URL                                  | Method | Role            | Requires period_id? | Authentication     |
|-------------------------|--------------------------------------|--------|-----------------|---------------------|--------------------|
| Login                   | `{{baseUrl}}/login`                  | POST   | Admin, Employee | No                  | None               |
| Register                | `{{baseUrl}}/register`               | POST   | Admin Only      | No                  | Admin JWT          |
| Create Attendance Period| `{{baseUrl}}/attendance-period`      | POST   | Admin Only      | No (Generates it)   | Admin JWT          |
| Run Payroll             | `{{baseUrl}}/payroll/{{period_id}}`  | POST   | Admin Only      | Yes                 | Admin JWT          |
| Generate Payroll Summary| `{{baseUrl}}/payroll-summary/{{period_id}}` | GET | Admin Only      | Yes                 | Admin JWT          |
| Submit Attendance       | `{{baseUrl}}/attendance`             | POST   | Employee Only   | Yes                 | Employee JWT       |
| Submit Overtime         | `{{baseUrl}}/overtime`               | POST   | Employee Only   | Yes                 | Employee JWT       |
| Submit Reimbursement    | `{{baseUrl}}/reimbursement`          | POST   | Employee Only   | Yes                 | Employee JWT       |
| Generate Payslip        | `{{baseUrl}}/payslip/{{period_id}}`  | GET    | Employee Only   | Yes                 | Employee JWT       |

- **baseUrl**: Typically `http://localhost:8084` for local development.
- **period_id**: A UUID generated when creating an attendance period.
- **Authentication**: JWT tokens are required for all endpoints except `/login`. Tokens are included in the `Authorization` header as `Bearer <token>`.

---

## Setup Instructions

### Prerequisites
- **Go**: Version 1.23.0 or later.
- **PostgreSQL**: Running locally or via Docker.
- **Git**: To clone the repository.
- **curl**: For testing API endpoints.

### Environment Variables
Set the following environment variables:

```bash
export DATABASE_URL="host=localhost user=postgres password=1234 dbname=payslip port=5432 sslmode=disable"
export JWT_SECRET="your-secret-key"
export PORT="8084"
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/erwin1412/dealls.git
   cd go-deals
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Create the PostgreSQL database:
   ```bash
   psql -U postgres -c "CREATE DATABASE payslip;"
   ```
4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```
   The server starts at `http://localhost:8084` and seeds the database with an admin user (`username: admin`, `password: admin123`) and 100 employee users with random usernames and salaries.

### Database Migration
The application automatically migrates the database schema on startup, creating tables for `User`, `AttendancePeriod`, `Attendance`, `Overtime`, `Reimbursement`, `Payroll`, and `AuditLog`. It also enables the `uuid-ossp` extension for UUID generation.

---

## API Usage Guide

Below are detailed instructions for each endpoint, including request formats, example responses, and error cases. Use `curl` for testing, and replace `<admin_token>`, `<employee_token>`, and `<period_id>` with actual values obtained during testing.

### 1. Login
- **Endpoint**: `POST {{baseUrl}}/login`
- **Role**: Admin, Employee
- **Authentication**: None
- **Description**: Authenticates a user and returns a JWT token.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/login -H "Content-Type: application/json" -d '{"username":"admin","password":"admin123"}'
  ```
- **Example Response**:
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "role": "admin"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid input"}`
  - 401: `{"error": "Invalid credentials"}`
- **Notes**:
  - Save the `token` for authenticated requests.
  - Try logging in as an employee (e.g., one of the seeded users) to get an employee token.

### 2. Register
- **Endpoint**: `POST {{baseUrl}}/register`
- **Role**: Admin Only
- **Authentication**: Admin JWT
- **Description**: Creates a new user (admin or employee).
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "role": "admin|employee"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/register -H "Content-Type: application/json" -H "Authorization: Bearer <admin_token>" -d '{"username":"newemployee","password":"password123","role":"employee"}'
  ```
- **Example Response**:
  ```json
  {
    "message": "User registered successfully",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "newemployee",
    "role": "employee"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Username already exists"}`, `{"error": "Password must be at least 6 characters"}`
  - 401: `{"error": "Unauthorized"}`
  - 403: `{"error": "Unauthorized"}` (if non-admin tries to register)
- **Notes**:
  - Username must be alphanumeric.
  - Employees are assigned a random salary between $2000 and $10000.
  - Audit log entry is created for each registration.

### 3. Create Attendance Period
- **Endpoint**: `POST {{baseUrl}}/attendance-period`
- **Role**: Admin Only
- **Authentication**: Admin JWT
- **Description**: Creates a new payroll period.
- **Request Body**:
  ```json
  {
    "start_date": "YYYY-MM-DD",
    "end_date": "YYYY-MM-DD"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/attendance-period -H "Content-Type: application/json" -H "Authorization: Bearer <admin_token>" -d '{"start_date":"2025-06-01","end_date":"2025-06-30"}'
  ```
- **Example Response**:
  ```json
  {
    "message": "Period created",
    "period_id": "789e1234-5678-90ab-cdef-123456789000"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid start date format"}`, `{"error": "End date must be after start date"}`
  - 401: `{"error": "Unauthorized"}`
- **Notes**:
  - Save the `period_id` for use in other endpoints.
  - Audit log entry is created.

### 4. Submit Attendance
- **Endpoint**: `POST {{baseUrl}}/attendance`
- **Role**: Employee Only
- **Authentication**: Employee JWT
- **Description**: Submits attendance for a specific date in a period.
- **Request Body**:
  ```json
  {
    "date": "YYYY-MM-DD",
    "period_id": "UUID"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/attendance -H "Content-Type: application/json" -H "Authorization: Bearer <employee_token>" -d '{"date":"2025-06-03","period_id":"789e1234-5678-90ab-cdef-123456789000"}'
  ```
- **Example Response**:
  ```json
  {
    "message": "Attendance submitted",
    "attendance_id": "456e7890-1234-56ab-cdef-678901234567"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid date format"}`, `{"error": "Cannot submit attendance on weekends"}`, `{"error": "Attendance already submitted for this date"}`
  - 401: `{"error": "Unauthorized"}`
  - 403: `{"error": "Payroll already processed for this period"}`
- **Notes**:
  - Attendance cannot be submitted for weekends (Saturday/Sunday).
  - Audit log entry is created.

### 5. Submit Overtime
- **Endpoint**: `POST {{baseUrl}}/overtime`
- **Role**: Employee Only
- **Authentication**: Employee JWT
- **Description**: Submits overtime hours for a specific date.
- **Request Body**:
  ```json
  {
    "date": "YYYY-MM-DD",
    "hours": number,
    "period_id": "UUID"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/overtime -H "Content-Type: application/json" -H "Authorization: Bearer <employee_token>" -d '{"date":"2025-06-03","hours":2,"period_id":"789e1234-5678-90ab-cdef-123456789000"}'
  ```
- **Example Response**:
  ```json
  {
    "message": "Overtime submitted",
    "overtime_id": "123e4567-89ab-cdef-1234-567890123456"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid date format"}`, `{"error": "Overtime cannot exceed 3 hours per day"}`
  - 401: `{"error": "Unauthorized"}`
  - 403: `{"error": "Payroll already processed for this period"}`
- **Notes**:
  - Maximum 3 hours of overtime per day.
  - Audit log entry is created.

### 6. Submit Reimbursement
- **Endpoint**: `POST {{baseUrl}}/reimbursement`
- **Role**: Employee Only
- **Authentication**: Employee JWT
- **Description**: Submits a reimbursement claim.
- **Request Body**:
  ```json
  {
    "amount": number,
    "description": "string",
    "period_id": "UUID"
  }
  ```
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/reimbursement -H "Content-Type: application/json" -H "Authorization: Bearer <employee_token>" -d '{"amount":100,"description":"Travel expenses","period_id":"789e1234-5678-90ab-cdef-123456789000"}'
  ```
- **Example Response**:
  ```json
  {
    "message": "Reimbursement submitted",
    "reimbursement_id": "789e1234-5678-90ab-cdef-678901234567"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid input"}`, `{"error": "Amount must be positive"}`, `{"error": "Description is required"}`
  - 401: `{"error": "Unauthorized"}`
  - 403: `{"error": "Payroll already processed for this period"}`
- **Notes**:
  - Audit log entry is created.

### 7. Run Payroll
- **Endpoint**: `POST {{baseUrl}}/payroll/{{period_id}}`
- **Role**: Admin Only
- **Authentication**: Admin JWT
- **Description**: Processes payroll for all employees in a period.
- **Request Body**: None
- **Example Request**:
  ```bash
  curl -X POST http://localhost:8084/payroll/789e1234-5678-90ab-cdef-123456789000 -H "Content-Type: application/json" -H "Authorization: Bearer <admin_token>"
  ```
- **Example Response**:
  ```json
  {
    "message": "Payroll processed"
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid period ID"}`, `{"error": "Payroll already processed for this period"}`
  - 401: `{"error": "Unauthorized"}`
  - 404: `{"error": "Period not found"}`
- **Notes**:
  - Calculates base salary (based on attendance), overtime pay (2x hourly rate), and reimbursement.
  - Audit log entries are created for each payroll record.
  - Cannot run payroll twice for the same period.

### 8. Generate Payroll Summary
- **Endpoint**: `GET {{baseUrl}}/payroll-summary/{{period_id}}`
- **Role**: Admin Only
- **Authentication**: Admin JWT
- **Description**: Retrieves a summary of payroll for a period.
- **Request Body**: None
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8084/payroll-summary/789e1234-5678-90ab-cdef-123456789000 -H "Authorization: Bearer <admin_token>"
  ```
- **Example Response**:
  ```json
  {
    "summary": [
      {"username": "employee1", "total_pay": 1234.56},
      {"username": "employee2", "total_pay": 2345.67}
    ],
    "total_payroll": 3580.23
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid period ID"}`
  - 401: `{"error": "Unauthorized"}`
  - 404: `{"error": "Payroll not found"}`
- **Notes**:
  - Requires payroll to be processed for the period.

### 9. Generate Payslip
- **Endpoint**: `GET {{baseUrl}}/payslip/{{period_id}}`
- **Role**: Employee Only
- **Authentication**: Employee JWT
- **Description**: Retrieves an employee’s payslip for a period.
- **Request Body**: None
- **Example Request**:
  ```bash
  curl -X GET http://localhost:8084/payslip/789e1234-5678-90ab-cdef-123456789000 -H "Authorization: Bearer <employee_token>"
  ```
- **Example Response**:
  ```json
  {
    "period": {
      "ID": "789e1234-5678-90ab-cdef-123456789000",
      "StartDate": "2025-06-01T00:00:00Z",
      "EndDate": "2025-06-30T00:00:00Z",
      ...
    },
    "attendance": [
      {"Date": "2025-06-03T00:00:00Z", ...}
    ],
    "overtime": [
      {"Date": "2025-06-03T00:00:00Z", "Hours": 2, ...}
    ],
    "reimbursements": [
      {"Amount": 100, "Description": "Travel expenses", ...}
    ],
    "base_salary": 1200.00,
    "overtime_pay": 200.00,
    "reimbursement_amount": 100.00,
    "total_pay": 1500.00
  }
  ```
- **Error Responses**:
  - 400: `{"error": "Invalid period ID"}`
  - 401: `{"error": "Unauthorized"}`
  - 404: `{"error": "Payroll not found"}`
- **Notes**:
  - Shows detailed attendance, overtime, and reimbursement records.
  - Requires payroll to be processed.

---

## Audit Logging

All create operations (e.g., register, create attendance period, submit attendance/overtime/reimbursement, run payroll) generate audit logs in the `audit_logs` table. Logs include:
- **Action**: e.g., `create`.
- **TableName**: e.g., `user`, `payroll`.
- **RecordID**: UUID of the affected record.
- **UserID**: UUID of the user performing the action.
- **IPAddress**: Client’s IP address.
- **RequestID**: Unique request ID from Echo middleware.
- **Details**: Descriptive message, e.g., `Created attendance period 789e1234... from 2025-06-01 to 2025-06-30`.
- **CreatedAt**: Timestamp.

To view audit logs:
```bash
psql -U postgres -d payslip -c "SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;"
```

---

## Error Handling

Common HTTP status codes:
- **200 OK**: Successful GET/POST.
- **201 Created**: Successful creation (e.g., register).
- **400 Bad Request**: Invalid input (e.g., wrong date format, missing fields).
- **401 Unauthorized**: Missing or invalid JWT.
- **403 Forbidden**: Insufficient role (e.g., employee accessing admin endpoint).
- **404 Not Found**: Resource not found (e.g., invalid period_id).
- **500 Internal Server Error**: Server-side error (rare, logged in console).

Error responses follow the format:
```json
{"error": "Error message"}
```

Check server logs for detailed errors:
```bash
tail -f /path/to/server.log
```

---

## Testing Workflow

Follow this sequence to test the full system:
1. **Start the Server**:
   ```bash
   go run cmd/api/main.go
   ```
2. **Login as Admin**:
   Use the default admin credentials (`admin`, `admin123`) to get an admin token.
3. **Create Attendance Period**:
   Generate a `period_id`.
4. **Login as Employee**:
   Use a seeded employee username (query `users` table or check server logs) to get an employee token.
5. **Submit Attendance/Overtime/Reimbursement**:
   Use the `period_id` and employee token.
6. **Run Payroll**:
   Use the admin token and `period_id`.
7. **Generate Payroll Summary**:
   Verify totals with the admin token.
8. **Generate Payslip**:
   View details with the employee token.
9. **Check Audit Logs**:
   Confirm all actions are logged.

---

## Development Notes

### Extending the System
- **Add Endpoints**: Create new handlers in `internal/api/handlers`, services in `internal/domain/services`, and repositories in `internal/infrastructure/repository`. Update interfaces in `internal/domain/interfaces`.
- **Validation**: Consider adding `github.com/go-playground/validator/v10` for stricter input validation.
- **Testing**: Write unit tests using `testing` and mocks (e.g., `github.com/stretchr/testify`).
  Example:
  ```go
  func TestPayrollService_GeneratePayrollSummary(t *testing.T) {
      // Mock setup and test cases
  }
  ```

### Common Issues
- **Database Connection**: Ensure PostgreSQL is running and `DATABASE_URL` is correct.
- **JWT Errors**: Verify `JWT_SECRET` is set and matches across requests.
- **UUID Parsing**: Ensure `period_id` is a valid UUID.
- **Role Restrictions**: Check user roles in `users` table if access is denied.

### Debugging
- Enable GORM logging in `internal/infrastructure/database/gorm.go` (already set to `logger.Info`).
- Check Echo request logs for request IDs and response times.
- Query the database directly to inspect data:
  ```bash
  psql -U postgres -d payslip -c "SELECT * FROM users;"
  ```

---

## Conclusion

This guide provides everything needed to use and extend the Payslip System API. By following the setup and testing instructions, developers can quickly integrate with the system. For further enhancements or issues, refer to the codebase at `https://github.com/erwin1412/dealls` or contact the development team.

**Version**: 1.0  
**Date**: June 25, 2025  
**Author**: Erwin