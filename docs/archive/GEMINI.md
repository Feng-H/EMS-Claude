# GEMINI.md - EMS Equipment Management System

This file serves as the primary instructional context for the EMS project, detailing its architecture, business logic, and development workflows.

## 🚀 Project Overview

**EMS (Equipment Management System)** is a unified industrial equipment management platform. It is designed to handle ~50,000 devices across 10+ bases, 40+ factories, and 100+ workshops, serving 10,000+ users.

### Core Objectives
- **Anti-Fraud**: Prevents "cloud inspection" (fake reports) via mandatory on-site QR code scanning.
- **High Concurrency**: Optimized for peak loads (e.g., thousands of concurrent submissions during shift changes).
- **Process Closure**: Ensures a full lifecycle for equipment issues, from NG (Not Good) detection to repair and audit.
- **Knowledge Sharing**: Promotes successful repair cases into a searchable knowledge base.

## 🏗️ Technical Architecture

### Backend (Go 1.23+)
- **Framework**: Gin (HTTP) + GORM (ORM).
- **Authentication**: JWT-based stateless authentication.
- **Modes**: Supports both **Memory Mode** (volatile, for rapid prototyping) and **Database Mode** (PostgreSQL + Redis).
- **Key Libraries**: `viper` (config), `swaggo` (API docs), `go-qrcode`.

### Frontend (Vue 3 + TypeScript)
- **Tools**: Vite, Pinia (state management).
- **PC UI**: Element Plus (Management & Analytics).
- **Mobile UI**: Vant 4 (On-site operations: Inspection, Repair, Maintenance).
- **Design**: Responsive layout with global scrollbar hiding and optimized sidebar folding.

### Infrastructure
- **Database**: PostgreSQL 15+ (Main storage).
- **Cache**: Redis 7+ (Hot data and performance optimization).
- **Containerization**: Docker & Docker Compose for full-stack orchestration.

## 🛠️ Building and Running

### Development Environment & Scripts
The project uses a structured workflow to switch between development and production-like Docker environments.

- **Start Dev**: `./start-dev.sh` (Vite dev server + Backend + DB/Redis in Docker).
- **Sync to Prod**: `./sync-to-prod.sh` (Builds frontend and restarts all Docker containers).
- **Switch back to Dev**: `./switch-to-dev.sh` (Stops production containers and starts dev server).
- **Stop All**: `./stop-dev.sh`.

### Manual Execution

#### Backend
```bash
cd backend
swag init -g main.go  # Update API docs
go run main.go        # Starts on port 8080 (default)
```

#### Frontend
```bash
cd frontend
npm install
npm run dev           # Starts on port 5173
```

## 📂 Project Structure

- `backend/`: Go source code following a service-repository pattern.
    - `api/v1/`: HTTP handlers (Memory and Database versions).
    - `internal/`: Business logic (`service`), data access (`repository`), and `model` definitions.
    - `pkg/`: Common utilities (config, database, redis, jwt, qrcode).
- `frontend/`: Vue 3 application.
    - `src/api/`: Request wrappers for all modules.
    - `src/views/`: Page components divided by business module.
    - `src/styles/`: Design system and platform-specific styles (PC vs H5).
- `db/`: SQL migrations and seed data.
- `docker/`: Dockerfiles and Nginx configurations for deployment.

## 📋 Business Logic Highlights

1.  **Three-Tier Org Model**: Base (基地) -> Factory (工厂) -> Workshop (车间).
2.  **Inspection NG Workflow**: If an inspection item is marked as "NG", the system automatically triggers a repair work order.
3.  **Dedicated Maintenance**: Equipment can be bound to a specific worker for automatic repair dispatch.
4.  **QR Code Integration**: Every device has a unique QR code used for on-site verification and quick access to technical specs/history.
5.  **Knowledge Promotion**: Engineers can audit repair logs and convert high-quality solutions into Knowledge Base articles.

## ⚖️ License & Permissions

- **License**: GNU General Public License v3.0 (GPL-3.0).
- **Commercial Use**: Requires written authorization from the author if modifications are not to be open-sourced.

## 📝 Development Conventions

- **API Documentation**: Always run `swag init` in `backend/` after modifying API handlers.
- **Mode Switching**: Toggle between `memory` and `database` storage in `backend/config/config.yaml`.
- **H5 Access**: Mobile-specific pages are served via the `/h5` route (e.g., `http://localhost:5173/h5`).
- **Permissions**: Roles (Admin, Supervisor, Engineer, Worker, Operator) strictly control access to modules like user approval and plan auditing.
