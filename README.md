# Baby-Fans 宝贝积分管理系统

家长积分管理系统。Parents can manage their children's points, create point templates, and set up a points shop where children can exchange points for items.

## Tech Stack

- **Backend**: Go 1.25 + Gin
- **Frontend**: UniApp (Vue 3)
- **Database**: SQLite
- **Nginx**: Reverse proxy for SSL termination

## Quick Start

```bash
# Start both backend and frontend (H5 mode)
./start.sh

# Backend only
cd backend && go run ./cmd/server

# Frontend H5 dev
cd frontend && npm run dev:h5
```

## Project Structure

```
backend/          # Go backend
  cmd/server/     # Entry point
  internal/       # Handler, Service, Repository, Model layers
  config/         # Configuration
  nginx/          # Nginx configuration (deprecated, moved to root)
frontend/         # UniApp frontend
nginx/            # Nginx reverse proxy config
```

## Nginx

Nginx is used as a reverse proxy for SSL termination. See `nginx/nginx.conf`.

Configure your domain and SSL certificates in the nginx config.

## API Base URL

- Production: `https://occont.asia`
- Development: `https://occont.asia` (configured in frontend/.env.development)
