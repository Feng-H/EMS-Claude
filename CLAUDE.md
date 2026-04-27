# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

EMS (Equipment Management System) is a unified equipment management platform for group factories. The system manages approximately 50,000 devices across 10+ bases, 40+ factories, and 100+ workshops, serving 10,000+ users. Current version: v0.2.0-beta (~65% complete).

### Core Problem Being Solved

- Large-scale equipment distributed across multiple locations
- Fake inspection reports ("cloud inspection" fraud)
- Slow maintenance response times
- Data silos between different factory units
- High concurrent access during shift changes

## User Roles and Permissions

| Role | Level | Responsibilities | Access Scope |
|------|-------|------------------|--------------|
| **admin** | 5 | Basic data maintenance, user/permission configuration, approve accounts | Global access |
| **supervisor** | 4 | View reports, monitor equipment status, analyze OEE | Data view for assigned factory/base |
| **engineer** | 3 | Create inspection/maintenance plans, audit repair records, manage knowledge base | Planning and audit for assigned factory |
| **maintenance** | 2 | Execute repair tasks, Level 2 maintenance, consume spare parts | Receive work orders, fill repair records |
| **operator** | 1 | Execute daily inspections, Level 1 maintenance, report faults | Access only to equipment they operate |

Users can apply for accounts (requires admin approval). First login requires password change.

## Tech Stack

### Backend (Go 1.23)
- **Framework**: Gin
- **ORM**: GORM + PostgreSQL driver
- **Cache**: Redis (go-redis/v9)
- **Auth**: JWT (golang-jwt/v5) + bcrypt
- **Config**: Viper (YAML + environment variable overrides via `EMS_` prefix)
- **LLM**: OpenAI-compatible client (default: SiliconFlow/DeepSeek-V3)
- **QR Code**: go-qrcode
- **API Docs**: Swagger (swaggo)

### Frontend (Vue 3 + TypeScript)
- **Build**: Vite 7
- **State**: Pinia 3
- **Router**: Vue Router 4 with auth guards and role-based access
- **PC UI**: Element Plus 2.13 + @element-plus/icons-vue
- **Mobile UI**: Vant 4.8
- **Charts**: ECharts 5.6
- **HTTP**: Axios

## Architecture

### Backend Architecture

```
main.go (routing, dual-mode setup)
  |
  +-- api/v1/              HTTP handlers (Gin controllers)
  |     auth.go, equipment.go, inspection.go, repair.go,
  |     maintenance.go, sparepart.go, analytics.go, knowledge.go,
  |     lark.go, memory.go
  |
  +-- internal/
  |     model/model.go        All GORM models (30+ entities)
  |     dto/                  Request/response DTOs (per domain)
  |     middleware/auth.go    JWT auth + role-based access control
  |     service/              Business logic layer
  |     repository/           Data access layer (GORM for DB, memory.Store for dev)
  |
  +-- internal/agent/         AI Agent subsystem (self-contained domain module)
  |     controller/           Agent HTTP handlers
  |     service/agent.go      Core orchestration (chat, skills, learning)
  |     repository/           Agent data access (DB + memory modes)
  |     analyzer/             Domain analyzers (predictive, repair audit, maintenance)
  |     tool/                 Agent tools (retrieval, repair data, maintenance data)
  |     policy/policy.go      Authorization/scoping per factory
  |     prompt/prompt.go      LLM prompt templates
  |     dto/                  Agent-specific DTOs
  |
  +-- pkg/                    Shared/reusable packages
        config/               Viper config loader
        database/             PostgreSQL connection (GORM)
        redis/                Redis client + CacheService
        jwt/                  JWT token generation/parsing
        llm/                  OpenAI-compatible LLM HTTP client
        lark/                 Lark/Feishu API client
        memory/               In-memory data store (full map-based mock)
        qrcode/               QR code generation
        trace/                Trace ID generator
```

### Dual Storage Mode

The backend supports two runtime modes via `config.yaml` `storage.mode`:

| Mode | Description | Use Case |
|------|-------------|----------|
| `memory` | All data in-memory maps, no DB needed. Seeds mock data. | Local development, demos |
| `database` | Full PostgreSQL + Redis with GORM AutoMigrate. | Production, staging |

Both modes share the same API routes and handlers. Each service layer checks the mode to branch between `memory.Store` and GORM queries.

### Frontend Architecture

```
src/
  api/              Axios-based API modules (per domain)
    request.ts        Base Axios instance with JWT interceptor
    auth.ts, equipment.ts, inspection.ts, repair.ts,
    maintenance.ts, sparepart.ts, analytics.ts,
    knowledge.ts, agent.ts, user.ts
  router/index.ts    Route definitions with auth guards + role checks
  stores/
    auth.ts          Token, user info, login/logout, role hierarchy
    theme.ts         Dark/light theme toggle (persisted)
  views/
    layout/          MainLayout.vue (PC), MobileLayout.vue (H5)
    auth/            Login, ChangePassword
    equipment/       Equipment list, detail, Organization
    inspection/      Templates, Tasks, Execute
    repair/          Orders, Execute, Report
    maintenance/     Plans, Tasks, Execute
    sparepart/       SparePart list
    analytics/       Charts & statistics
    knowledge/       Knowledge base
    agent/           AI Management Assistant cockpit
    user/            User management
    h5/              Mobile views (Vant-based)
  components/        Shared: StatItem, MobileHeader, MobileActionBar, SparePartSelector
  composables/       useDeviceDetection
  utils/             device.ts (isMobile/isTablet/isDesktop)
  styles/            design-system.css, pc.css, h5.css, utilities.css
```

PC and mobile use **separate route trees** and **separate UI libraries**:
- PC: `/` routes + `MainLayout` + Element Plus
- Mobile: `/h5` routes + `MobileLayout` + Vant

## Core Functionality (Implemented)

### 1. Equipment Ledger Management
- Three-tier organization: Base (基地) -> Factory (工厂) -> Workshop (车间)
- Equipment profiles with technical parameters, financial fields (purchase_price, scrap_value, hourly_loss, service_life), and dedicated maintenance engineer binding
- QR code for each equipment (unique identifier for on-site scanning)
- Equipment types management

### 2. Inspection Management (点检)
- Configurable templates by equipment type with check items
- Auto-generated daily/periodic inspection tasks
- GPS coordinate recording on inspection start
- My-tasks and my-stats endpoints for mobile workers

### 3. Repair Management (维修)
- Full 7-state closed-loop workflow: pending -> assigned -> in_progress -> testing -> confirmed -> audited -> closed
- Auto-dispatch: priority to equipment-bound maintenance worker, else public pool or engineer assignment
- Spare parts consumption linked to repair orders with cost tracking
- Priority levels, repair logs, and cost detail breakdowns

### 4. Maintenance Management (保养)
- Tiered maintenance levels (Level 1 by operators, Level 2 by maintenance workers)
- Plan items as checklists, task generation from plans
- Flexible scheduling with execution window
- GPS coordinate recording on maintenance start

### 5. Spare Parts Management (备件)
- Spare part catalog with safety stock thresholds
- Per-factory inventory management
- Stock in/out with consumption tracking linked to work orders
- Low stock alert system

### 6. Analytics (统计分析)
- Dashboard overview with key metrics
- MTTR (Mean Time To Repair), MTBF (Mean Time Between Failures)
- Trend data, failure analysis, top failure equipment

### 7. Knowledge Base (知识库)
- Articles with fault phenomenon, cause, solution, tags
- Search functionality
- Convert excellent repair records to knowledge base entries
- Equipment manual documents with text chunking for retrieval

### 8. AI Agent / Smart O&M Assistant
- Multi-turn conversational AI chat with LLM (OpenAI-compatible API)
- Skill system: pre-defined analysis workflows (maintenance recommendation, repair audit, predictive analysis)
- Self-learning: async extraction of knowledge and skills from chat history
- Predictive analytics: RUL (Remaining Useful Life), TCO (Total Cost of Ownership), sub-health symptom detection, retirement evaluation
- Evidence-based analysis with traceable data chains
- Factory-level data isolation via policy service
- Frontend cockpit: chat interface, audit mode, knowledge review, real-time equipment health panel

### 9. Lark (Feishu) Integration
- Webhook handler for Lark events (URL verification + message events)
- Signature verification, auto-cached tenant access token
- Text and interactive card message sending
- User account binding (Lark OpenID <-> EMS user) via `/h5/bind-lark`
- Agent bridge: Lark messages forwarded to AI agent, responses sent back via Lark

## API Structure

All API routes prefixed with `/api/v1`. Base URL configured via `VITE_API_BASE_URL` (dev: `http://localhost:8080/api/v1`).

**Public**: `POST /auth/login`, `POST /auth/refresh`, `POST /auth/apply`, `POST /lark/webhook`

**Protected** (JWT required): `/auth/me`, `/auth/change-password`, `/auth/bind-lark`, `/users/*`, `/organization/*`, `/equipment/*`, `/inspection/*`, `/repair/*`, `/maintenance/*`, `/spareparts/*`, `/analytics/*`, `/knowledge/*`, `/agent/*`

## Common Development Commands

### Backend (Go)
```bash
cd backend

# Run in memory mode (no DB required, mock data)
go run main.go

# Run in database mode (requires PostgreSQL + Redis)
# Set EMS_STORAGE_MODE=database in .env or use config.docker.yaml
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
# Schema and data are auto-created by GORM AutoMigrate + Go seeder on startup
# db/schema.sql and db/seeds/seed.sql are reference documentation only

# Start PostgreSQL and Redis via Docker
docker-compose up -d postgres redis

# Connect to database
psql -h localhost -p 5432 -U ems -d ems_db
# Password: ems123
```

### Docker
```bash
# Development (root docker-compose: all services)
docker-compose up -d --build

# Production (deploy/docker-compose: with nginx)
cd deploy && docker-compose up -d --build

# Full rebuild with clean database (when schema changes)
docker-compose down -v && docker-compose up -d --build

# View logs
docker-compose logs -f backend
```

### Scripts
```bash
# Start full dev environment
./start-dev.sh

# Stop dev environment
./stop-dev.sh

# Switch to dev mode
./switch-to-dev.sh

# Sync to production
./sync-to-prod.sh
```

## Key File Locations

### Backend
- **Entry point**: `backend/main.go`
- **Models**: `backend/internal/model/model.go`
- **Handlers**: `backend/api/v1/*.go`
- **Services**: `backend/internal/service/*.go`
- **Repositories**: `backend/internal/repository/*.go`
- **Agent subsystem**: `backend/internal/agent/`
- **Auth middleware**: `backend/internal/middleware/auth.go`
- **Dev config**: `backend/config/config.yaml` (memory mode)
- **Docker config**: `backend/config/config.docker.yaml` (database mode)

### Frontend
- **Vite config**: `frontend/vite.config.ts`
- **Router**: `frontend/src/router/index.ts`
- **API layer**: `frontend/src/api/*.ts`
- **Stores**: `frontend/src/stores/auth.ts`, `frontend/src/stores/theme.ts`
- **PC layout**: `frontend/src/views/layout/MainLayout.vue`
- **Mobile layout**: `frontend/src/views/layout/MobileLayout.vue`
- **AI cockpit**: `frontend/src/views/agent/ManagementAssistantView.vue`

### Infrastructure
- **DB schema**: `db/schema.sql` (reference only; runtime schema managed by GORM AutoMigrate)
- **Seed data**: `db/seeds/seed.sql` (reference only; runtime data from Go seeder `repository.SeedDatabase()`)
- **Migrations**: `backend/scripts/migrations/` (manual performance indexes, optional)
- **Root docker-compose**: `docker-compose.yml` (postgres + redis + backend + frontend)
- **Deploy docker-compose**: `deploy/docker-compose.yml` (+ nginx, for production)
- **Nginx config**: `frontend/nginx.conf` (copied into frontend image at build time)
- **Env template**: `.env.example`

### Documentation
- **Historical docs**: `docs/archive/` (completed Agent phase docs, old deployment guides, business requirements)
- **README.md**: Project overview + detailed Feishu bot setup guide

## Environment Variables

Key variables (set in `.env`, all prefixed with `EMS_`):

| Variable | Description | Default |
|----------|-------------|---------|
| `EMS_STORAGE_MODE` | `memory` or `database` | `memory` |
| `EMS_DATABASE_HOST` | PostgreSQL host | `localhost` |
| `EMS_DATABASE_PORT` | PostgreSQL port | `5432` |
| `EMS_DATABASE_USER` | PostgreSQL user | `ems` |
| `EMS_DATABASE_PASSWORD` | PostgreSQL password | `ems123` |
| `EMS_DATABASE_NAME` | Database name | `ems_db` |
| `EMS_REDIS_HOST` | Redis host | `localhost` |
| `EMS_REDIS_PORT` | Redis port | `6379` |
| `EMS_LLM_PROVIDER` | LLM provider name | `openai` |
| `EMS_LLM_BASE_URL` | LLM API base URL | (SiliconFlow) |
| `EMS_LLM_API_KEY` | LLM API key | - |
| `EMS_LLM_MODEL` | LLM model name | `deepseek-ai/DeepSeek-V3` |
| `EMS_LARK_APP_ID` | Lark app ID | - |
| `EMS_LARK_APP_SECRET` | Lark app secret | - |
| `EMS_DOMAIN` | Application domain | `ems.example.com` |

## Architecture Notes

- Follow the three-tier organization model (Base -> Factory -> Workshop)
- QR code scanning is a critical mechanism for on-site equipment identification
- Mobile-first experience for on-site workers (separate H5 views with Vant)
- Dual storage mode allows full development without external dependencies
- The Agent subsystem is self-contained with its own controller/service/repository layers
- Factory-level data isolation is enforced by the Agent policy service
- All agent conclusions should be backed by traceable evidence from database records
- When adding new features, implement both `memory` and `database` mode paths

### Important: GORM Column Naming

GORM converts Go field names to snake_case column names. Consecutive uppercase letters are treated as one word:
- `LarkOpenID` -> column `lark_open_id` (NOT `lark_openid`)
- `EquipmentTypeID` -> column `equipment_type_id`

When writing raw SQL queries in repository layer, always use the GORM-generated column name, not the Go field name. Test with `docker compose exec backend sh -c "psql -U ems -d ems_db -c '\\d users'"` to verify column names.
