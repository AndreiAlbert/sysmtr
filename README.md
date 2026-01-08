# Sysmtr - System Monitoring Tool

Sysmtr is a full-stack system monitoring solution that provides real-time and historical tracking of CPU and RAM usage for distributed instances.

## Architecture

* **Agent (Go):** Collects metrics using `gopsutil` and sends them via gRPC streams.
* **Server (Go):** Manages gRPC connections, handles database persistence, and serves WebSocket/REST endpoints.
* **Web Dashboard (Angular):** Provides a clean UI for viewing live stats and history.
* **Database (PostgreSQL):** Persistent storage for all collected system stats.

## Tech Stack

* **Backend:** Go, gRPC, Protocol Buffers (proto3)
* **Frontend:** Angular 17, RxJS (WebSockets)
* **Database:** PostgreSQL 15
* **Infrastructure:** Docker, Docker Compose

## Getting Started

### Prerequisites

* Docker and Docker Compose installed on your machine.

### Installation & Deployment

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/AndreiAlbert/sysmtr.git](https://github.com/AndreiAlbert/sysmtr.git)
    cd sysmtr
    ```

2.  **Environment Configuration:**
    Create a `.env` file in the root directory (or use the defaults in `docker-compose.yml`) with the following variables:
    ```env
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=sysmonitor
    GRPC_PORT=50051
    HTTP_PORT=8080
    ```

3.  **Run with Docker Compose:**
    ```bash
    docker-compose up --build
    ```

### Accessing the Application

* **Web Dashboard:** `http://localhost:4200`
* **Backend API (HTTP):** `http://localhost:8080/history`
* **gRPC Server:** `localhost:50051`

## API Endpoints

| Endpoint | Method | Description |
| :--- | :--- | :--- |
| `/ws` | WS | WebSocket endpoint for real-time system stats. |
| `/history` | GET | Returns the last 50 recorded system stats in JSON format. |

## Development

### Local Agent Setup
If you want to run an agent locally without Docker:
```bash
cd backend/cmd/agent
go run main.go
