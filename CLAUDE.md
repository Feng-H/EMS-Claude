# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EMS (Equipment Management System) is a unified equipment management platform for group factories. The system manages approximately 50,000 devices across 10+ bases, 40+ factories, and 100+ workshops, serving 10,000+ users.

### Core Problem Being Solved

- Large-scale equipment distributed across multiple locations
- Fake inspection reports ("cloud inspection" fraud)
- Slow maintenance response times
- Data silos between different factory units
- High concurrent access during shift changes

## User Roles and Permissions

| Role | Responsibilities | Access Scope |
|------|------------------|--------------|
| **System Administrator** | Basic data maintenance, user/permission configuration, data imports | Global configuration access |
| **Equipment Supervisor** | View reports, monitor equipment status, analyze OEE | Data view access for assigned factory/base |
| **Equipment Engineer** | Create inspection/maintenance plans, audit repair records, manage knowledge base | Planning and audit access for assigned factory |
| **Maintenance Worker** | Execute repair tasks, Level 2 maintenance, consume spare parts | Receive work orders, fill repair records |
| **Operator** | Execute daily inspections, Level 1 maintenance, report faults | Access only to equipment they operate |

## Core Functionality

### 1. Equipment Ledger Management
- Three-tier organization: Base (基地) -> Factory (工厂) -> Workshop (车间)
- Equipment profiles with technical parameters and location
- QR code for each equipment (unique identifier for on-site scanning)
- Batch import via Excel/JSON

### 2. Inspection Management (点检)
- **Anti-fraud**: Must scan equipment QR code on-site to start inspection
- Configurable templates by equipment type (CNC, welding machine, etc.)
- Auto-generated daily/periodic inspection tasks
- NG (Not Good) items automatically trigger repair workflow

### 3. Repair Management (维修)
- Mobile fault reporting with photo upload via QR code scan
- Auto-dispatch:优先派给设备绑定的专属维修工，无绑定则推送到公共池或由工程师指派
- Closed-loop workflow: Report -> Accept -> Repair -> Test -> Requester Confirm -> Engineer Audit -> Archive
- Spare parts consumption recorded during repair, auto-deducts inventory

### 4. Maintenance Management (保养)
- Tiered maintenance levels:
  - Level 1 (一级保养): By operators
  - Level 2 (二级保养): By maintenance workers
  - Precision maintenance
- Flexible scheduling with execution window (e.g., ±3 days)
- Visual calendar for maintenance plans

### 5. Spare Parts Management (备件)
- Factory-level inventory management
- Consumption must be linked to work orders for cost traceability

### 6. Analytics (统计分析)
- Core metrics: MTTR (Mean Time To Repair), MTBF (Mean Time Between Failures)
- Execution rates: Inspection completion rate, on-time maintenance completion rate
- OEE (Overall Equipment Effectiveness) display

### 7. Knowledge Base (知识库)
- Excellent repair records can be "promoted" to knowledge base entries
- On-site search by fault symptoms

## Technical Requirements

### High Concurrency
- Must support thousands of users concurrently submitting inspection data during morning shift start

### Multi-Platform Support
- **PC Client**: Management configuration, data analytics (Vue3 + Element Plus)
- **Mobile Client**: On-site operations (Vue3 + Vant H5), integrated into enterprise App or standalone

### Deployment
- Docker containerization
- Cloud-based elastic scaling

### Security
- Username/password authentication
- SSO interface reserved for future integration

## Architecture Notes

This is a greenfield project - the codebase is currently being set up. When implementing:
- Follow the three-tier organization model (Base -> Factory -> Workshop)
- Design for high concurrency from the ground up
- QR code scanning is a critical anti-fraud mechanism for inspections
- All workflows should support mobile-first experience for on-site workers

## Common Development Commands

### Backend (Go)
```bash
cd backend

# Run development server
go run main.go

# Run tests
go test ./...

# Build binary
go build -o bin/ems main.go

# Generate Swagger docs
swag init -g main.go
```

### Frontend (Vue3)
```bash
cd frontend

# Install dependencies
npm install

# Run dev server (http://localhost:5173)
npm run dev

# Build for production
npm run build

# Run linting
npm run lint
```

### Database
```bash
# Start PostgreSQL and Redis via Docker
cd docker
docker-compose up -d postgres redis

# Connect to database
psql -h localhost -U ems -d ems_db

# Run schema
psql -h localhost -U ems -d ems_db -f db/schema.sql

# Run seed data
psql -h localhost -U ems -d ems_db -f db/seeds/seed.sql
```

### Docker (Full Stack)
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## Key File Locations

- **Backend entry**: `backend/main.go`
- **Database schema**: `db/schema.sql`
- **Seed data**: `db/seeds/seed.sql`
- **Frontend config**: `frontend/vite.config.ts`
- **Docker compose**: `docker/docker-compose.yml`
- **API config**: `backend/config/config.yaml`
