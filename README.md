<div align="center">

# 🚀 deployes

### GitHub Deployment Automation Platform

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Angular](https://img.shields.io/badge/Angular-21+-DD0031?style=for-the-badge&logo=angular&logoColor=white)](https://angular.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

[![CI](https://github.com/aliicolak/deployes/workflows/CI/badge.svg)](https://github.com/aliicolak/deployes/actions)
[![codecov](https://codecov.io/gh/aliicolak/deployes/branch/main/graph/badge.svg)](https://codecov.io/gh/aliicolak/deployes)
[![Go Report Card](https://goreportcard.com/badge/github.com/aliicolak/deployes)](https://goreportcard.com/report/github.com/aliicolak/deployes)
[![Maintainability](https://api.codeclimate.com/v1/badges/YOUR_TOKEN/maintainability)](https://codeclimate.com/github/aliicolak/deployes/maintainability)

**deployes is a self-hosted, modern deployment automation platform that enables you to deploy your GitHub projects to remote servers with a single click or automatically via webhooks.**

[🇹🇷 Türkçe Dokümantasyon](README.tr.md)

---

<img src="https://raw.githubusercontent.com/aliicolak/deployes/main/docs/screenshot-dashboard.png" alt="deployes Dashboard" width="800"/>

</div>

## ✨ Features

### 🎯 Core Features
- **One-Click Deployments** - Deploy any project to any server with a single click
- **GitHub Webhook Integration** - Automatic deployments on push events
- **Branch-Based Triggers** - Configure specific branches to trigger deployments
- **Multi-Server Management** - Manage unlimited servers from a single dashboard
- **Encrypted Credentials** - SSH keys and secrets are encrypted at rest (AES-256)

### 🔐 Security
- **JWT Authentication** - Secure token-based authentication
- **SSH Key Management** - Auto-generated deploy keys for private repositories
- **Environment Secrets** - Encrypted environment variables for deployments
- **No Password Storage** - SSH key-based server authentication only

### 📊 Monitoring & Logs
- **Real-time Deployment Logs** - WebSocket-based live streaming
- **Deployment History** - Complete history with status tracking
- **Dashboard Analytics** - Overview of deployment statistics

### 🎨 Modern UI/UX
- **Dark/Light Mode** - Beautiful themes with smooth transitions
- **Responsive Design** - Works on desktop, tablet, and mobile
- **Deploy Script Templates** - Pre-built templates for Node.js, Python, Go, .NET, Docker

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        deployes Architecture                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐         ┌──────────────┐                     │
│   │   Angular    │◄───────►│   Go API     │                     │
│   │   Frontend   │  REST   │   Backend    │                     │
│   │   (Port 4200)│  + WS   │  (Port 8080) │                     │
│   └──────────────┘         └──────┬───────┘                     │
│                                   │                              │
│                    ┌──────────────┼──────────────┐              │
│                    ▼              ▼              ▼              │
│           ┌────────────┐  ┌────────────┐  ┌────────────┐       │
│           │ PostgreSQL │  │   GitHub   │  │  Remote    │       │
│           │  Database  │  │  Webhooks  │  │  Servers   │       │
│           └────────────┘  └────────────┘  └────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+
- **Node.js** 18+
- **Docker** & Docker Compose
- **Git**

### 1. Clone the Repository

```bash
git clone https://github.com/aliicolak/deployes.git
cd deployes
```

### 2. Configure Environment Variables

```bash
cp .env.example .env
# Edit .env and set your own secure values
```

> ⚠️ **Never commit `.env` to git.** See [`.env.example`](.env.example) for required variables.

### 3. Start the Database

```bash
docker compose up -d

### 4. Run the Backend

```bash
go run ./cmd/api
```

### 5. Run the Frontend

```bash
cd web
npm install
npm start
```

### 6. Open in Browser

Navigate to [http://localhost:4200](http://localhost:4200)

---

## 📁 Project Structure

```
deployes/
├── cmd/
│   └── api/                 # Application entry point
├── internal/
│   ├── application/         # Business logic services
│   │   ├── deployment/      # Deployment service
│   │   ├── project/         # Project management
│   │   ├── server/          # Server management
│   │   └── webhook/         # Webhook handling
│   ├── domain/              # Domain entities & interfaces
│   ├── handlers/            # HTTP handlers
│   └── infrastructure/      # Database repositories
├── pkg/
│   └── utils/               # Utility functions
├── web/                     # Angular frontend
│   └── src/
│       ├── app/
│       │   ├── core/        # Services, guards, interceptors
│       │   ├── features/    # Feature components
│       │   └── shared/      # Shared components
│       └── assets/
├── docker-compose.yml
├── go.mod
└── README.md
```

---

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | ✅ | See `.env.example` |
| `JWT_SECRET` | JWT signing key (min 32 chars) | ✅ | Generate a random key |
| `ENCRYPTION_KEY` | AES encryption key (exactly 32 chars) | ✅ | Generate a random key |
| `APP_PORT` | API server port | ❌ | `8080` (default) |
| `ALLOWED_ORIGINS` | CORS allowed origins | ❌ | `http://localhost:4200` |

---

## 📖 Usage Guide

### Adding a Server

1. Navigate to **Servers** page
2. Click **"+ Yeni Sunucu"** (New Server)
3. Enter server details:
   - **Name**: Display name for the server
   - **Host**: IP address or hostname
   - **Port**: SSH port (default: 22)
   - **Username**: SSH user
   - **SSH Key**: Private key for authentication
4. Click **Save**

### Adding a Project

1. Navigate to **Projects** page
2. Click **"+ Yeni Proje"** (New Project)
3. Enter project details:
   - **Project Name**: Display name
   - **GitHub Repo URL**: Full repository URL
   - **Branch**: Target branch (e.g., `main`, `master`)
   - **Deploy Script**: Shell commands to execute
4. For private repos, copy the generated **Deploy Key** to GitHub

### Creating a Deployment

1. Navigate to **Deployments** page
2. Click **"+ Yeni Deployment"** (New Deployment)
3. Select the **Project** and **Server**
4. Click **Deploy**
5. Watch real-time logs in the terminal

### Setting Up Webhooks

1. Navigate to **Webhooks** page
2. Create a new webhook for your project/server pair
3. Copy the generated **Webhook URL**
4. Add it to GitHub: `Settings → Webhooks → Add webhook`
5. Content type: `application/json`
6. Select: **Just the push event**

---

## 🛠️ Deploy Script Templates

### Node.js + PM2
```bash
#!/bin/bash
set -e
npm install
npm run build
pm2 reload ecosystem.config.js || pm2 restart all
```

### Docker Compose
```bash
#!/bin/bash
set -e
docker-compose pull
docker-compose up -d --build
docker system prune -f
```

### .NET / ASP.NET Core
```bash
#!/bin/bash
set -e
dotnet restore
dotnet publish -c Release -o ./publish
systemctl restart myapp.service
```

---

## 🔒 Security Best Practices

1. **Use strong secrets** - Generate random, long JWT and encryption keys
2. **SSH keys only** - Never use password-based SSH authentication
3. **Firewall rules** - Restrict access to ports 8080 and 4200
4. **HTTPS** - Use a reverse proxy (nginx) with SSL certificates
5. **Regular updates** - Keep dependencies and Docker images updated

---

## 🤝 Contributing

We love contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Quick Start

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/deployes.git`
3. Create a feature branch: `git checkout -b feature/amazing-feature`
4. Make your changes and test: `make test`
5. Commit your changes: `git commit -m 'feat: add amazing feature'`
6. Push to the branch: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run linter
make lint

# Build
make build
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for full guidelines.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👨‍💻 Authors

**Ali Çolak**

- GitHub: [@aliicolak](https://github.com/aliicolak)

**Alper Şahin**

- GitHub: [@Alpersahin11](https://github.com/Alpersahin11)
---

<div align="center">

**⭐ Star this repository if you find it helpful!**

Made with ❤️ using Go and Angular

</div>
