# SoAt - Go Web Service

A robust and scalable web service built with Go, containerized with Docker, and backed by PostgreSQL and Redis. This project uses a clean, layered architecture to separate concerns and improve maintainability.

---

## ✨ Features

- **Layered Go Backend**: High-performance backend organized by features (Users, Posts, Friendships) across controllers, services, and models.
- **Containerized Environment**: Uses Docker and Docker Compose for consistent development, testing, and production environments.
- **Relational Database**: Integrates with PostgreSQL for reliable data storage, managed by GORM.
- **In-Memory Caching**: Utilizes Redis for fast data access and caching.
- **RESTful API Structure**: Routes are clearly defined and grouped by feature.

---

## 📋 Prerequisites

Before you begin, ensure you have the following installed on your local machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/) (Typically included with Docker Desktop)

---

## 🚀 Getting Started

Follow these steps to get your development environment up and running.

### 1. Clone the Repository

```bash
git clone <your-repository-url>
cd SoAt
```

### 2. Configure Environment Variables

This project uses a `docker-compose.yaml` file to manage services and their configurations. The environment variables for the Go application are set directly in this file for simplicity in a containerized setup.

All necessary variables for the database and Redis connections are already configured in `docker-compose.yaml` to work within the Docker network.

### 3. Build and Run the Application

Use Docker Compose to build the application's Docker image and start all the services (Go app, PostgreSQL, Redis).

```bash
docker compose up --build
```

- `--build`: This flag forces Docker to rebuild your application's image. This is necessary whenever you make changes to the source code (e.g., `.go` files) or the `Dockerfile`.

Your application will be accessible at `http://localhost:3000`.

### 4. Stopping the Application

To stop all running containers, press `Ctrl + C` in the terminal where the application is running. To remove the containers and the network, run:

```bash
docker compose down
```

To remove the persistent database volume as well (deleting all data), run:
```bash
docker compose down -v
```
---

## 🚀 API Usage Examples

You can interact with the running server using any API client. The following examples use `curl`.

### 1. Create a New User

Send a `POST` request to `/users` with the user's details in the request body.

**Request:**
```bash
curl -X POST http://localhost:3000/users \
-H "Content-Type: application/json" \
-d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword123"
}'
```

**Expected Response:**
The server will respond with the newly created user object, including its unique ID and timestamps.

### 2. Get All Users

Send a `GET` request to `/users` to retrieve a list of all users.

**Request:**
```bash
curl -X GET http://localhost:3000/users
```

**Expected Response:**
A JSON array containing all user objects in the database.

### 3. Get a Single User by ID

Send a `GET` request to `/users/{id}`, replacing `{id}` with the actual user ID you want to retrieve.

**Request:**
```bash
# Replace {user-id} with an actual ID from the GET all users response
curl -X GET http://localhost:3000/users/{user-id}
```

### 4. Delete a User by ID

Send a `DELETE` request to `/users/{id}` to remove a specific user.

**Request:**
```bash
# Replace {user-id} with the ID of the user you want to delete
curl -X DELETE http://localhost:3000/users/{user-id}
```

**Expected Response:**
A success message indicating the user has been deleted.

---

## Screenshots

![Used Insomnia for sending POST,GET,DELETE,PUT](assets/images/S0.png)
### <b>Used Insomnia for sending REST API Calls<b>

![Building Docker app](assets/images/S1.png)
### <b>Building the Docker App<b>
![Docker Repository](assets/images/S2.png)
### <b>Docker Repository<b>
![After using Docker Compose on your system](assets/images/S3.png)
### <b>After using Docker Compose on your system<b>
![Docker Desktop for checking logs](assets/images/S4.png)
### <b>Docker Desktop Interface<b>

## 📂 Project Structure

The project follows a layered architecture to separate responsibilities:

```
.
├── cmd/                # Main application entry point
│   ├── app/
│   └── main.go
├── controllers/        # Handles HTTP request/response logic
│   ├── friendships/
│   ├── posts/
│   └── users/
├── internals/          # Core shared packages
│   ├── cache/          # Redis connection logic
│   ├── constants/      # Application-wide constants
│   ├── database/       # PostgreSQL connection logic
│   ├── dto/            # Data Transfer Objects
│   ├── notifications/
│   └── server/         # HTTP server, handlers, and middleware
├── models/             # GORM database models
│   ├── friendships/
│   ├── posts/
│   └── users/
├── routes/             # API route definitions
│   ├── friendships/
│   ├── posts/
│   └── users/
├── services/           # Business logic for each feature
│   ├── friendships/
│   ├── posts/
│   └── users/
├── .dockerignore       # Specifies files to ignore in Docker build context
├── .env                # Local environment variables (ignored by Git)
├── .gitignore          # Specifies files to ignore for Git
├── Dockerfile          # Instructions for building the Go application image
├── docker-compose.yaml # Defines and configures all services
├── go.mod              # Go module dependencies
├── go.sum              # Go module checksums
└── README.md           # You are here!
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue or email me ( dibyendusahoo03@gmail.com ).