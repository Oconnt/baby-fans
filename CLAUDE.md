# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Baby-Fans is a full-stack parent-child points management application (家长积分管理系统). Parents can manage their children's points, create point templates, and set up a points shop where children can exchange points for items.

## Technology Stack

- **Backend**: Go 1.25 with Gin framework
- **Frontend**: UniApp (Vue 3) - cross-platform mini-program framework
- **Database**: SQLite (gorm.io/glebarez/sqlite)
- **Authentication**: JWT tokens

## Common Commands

```bash
# Start both backend and frontend (H5 mode)
./start.sh

# Start backend only
cd backend && go run cmd/server/main.go

# Start frontend only (H5 mode)
cd frontend && npm run dev:h5

# Build frontend for WeChat Mini Program
cd frontend && npm run dev:mp-weixin

# Run backend tests
cd backend && go test ./...

# Run specific test
cd backend && go test -v ./internal/service/...
```

## Architecture

### Backend (Go)

Uses a layered architecture:
- **Handler** (`internal/api/handler/`): HTTP request handling, validation
- **Service** (`internal/service/`): Business logic
- **Repository** (`internal/repository/`): Database operations
- **Model** (`internal/model/`): Data structures

Key files:
- `backend/cmd/server/main.go`: Entry point
- `backend/internal/api/router.go`: Route definitions, middleware setup
- `backend/internal/api/handler/handlers.go`: API handlers
- `backend/internal/service/services.go`: Business logic
- `backend/internal/model/models.go`: Database models

### Frontend (UniApp/Vue)

Pages are configured in `pages.json` with tab bar navigation.

Key pages:
- `pages/login/login.vue`: Login page
- `pages/register/register.vue`: Registration page
- `pages/home-parent/home-parent.vue`: Child management (tab)
- `pages/tags/tags.vue`: Points templates (tab)
- `pages/shop/shop.vue`: Points shop (tab)
- `pages/mine/mine.vue`: Profile (tab)
- `pages/records/records.vue`: Exchange records
- `pages/points/points.vue`: Points modification records

API utility: `frontend/src/utils/request.ts`

### API Structure

Authentication uses JWT tokens passed in Authorization header:
- Public: `/login/face`, `/login/code`, `/register`, `/api/v1/auth/wechat/login`, `/parent/items`
- Parent-only: `/parent/*` (requires parent role)
- Child-only: `/child/*` (requires child role)

Database file: `backend/baby-fans.db`

## Data Models

- **User**: id, name, role (parent/child), password, login_code, points, openid, unionid, nickname, avatar_url
- **UserBinding**: parent-child binding with bind_code
- **PointsRecord**: user points history with operator tracking
- **ShopItem**: items available for points exchange
- **Redemption**: exchange requests with status (pending/completed/cancelled)
- **PointsTemplate**: quick point templates for parents
