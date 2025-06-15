# Habits Backend
This is a Go backend for an Atomic Habits app. 


## Resources
- Habits - represents an individual habit about which a user has to be reminded about
- User - represents an individual registered user's information


## How to run?

### Using Docker

1. Make sure you have Docker and Docker Compose installed on your system.
2. Run the application using Docker Compose
The API will be available at http://localhost:8080

### Running Locally (Recommender)

1. Install Go (version 1.16 or later)
2. Run the postgres container with `docker run -e POSTGRES_PASSWORD=habits -e POSTGRES_USER=habits postgres:latest `
3. Run database migrations:
   ```
   go run cmd\migrate\main.go up
   ```
5. Start the API server:
   ```
   go run cmd\api\main.go
   ```

The API will be available at http://localhost:8080
