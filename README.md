# Note Taking App

## Goals
- Build a backend in Go with a clean, layered architecture (handlers, repositories, routes).
- Gain familiarity with Go’s standard library for HTTP, JSON, and database interactions.
- Use Postgres as the primary database.
- Implement Dockerfiles for each component and run them together with docker-compose.
- Deploy a simple frontend (React) to interact with the backend API.

## Architecture

This project follows a **Layered Architecture** pattern, separating concerns into distinct layers:
- Frontend (React.js): 
  - Provides a simple UI for interacting with lists and items.
- Backend (Go):
  - Exposes REST endpoints
  - Handles request/response logic
  - Uses repository layer for database access
- Database (Postgres):
  - Stores lists and items.

## Containerization

Components:
- Backend: Builds the Go binary in a multi-stage build and runs it in a lightweight Alpine container.
- Frontend: Built and served with Node.
- Database: Uses the official Postgres image with mounted volumes for persistence.

Docker Compose is used to orchestrate the system so everything can run with a single command:
```
docker compose up --build
```
_Note: a lot of the code, especially the frontend, is not complete. Creating the current images and running the app as it is with Docker was just a way to make sure I understood how to get all containers communicating._

## Next Steps
- Add authentication/authorization.
- Expand frontend functionality and styling.
- Write automated tests for backend services.

