# Peekaping 🚀

**A modern, self-hosted uptime monitoring solution**

Peekaping is a powerful, feature-rich uptime monitoring system similar to Uptime Kuma, built with Go and React. Monitor your websites, APIs, and services with real-time notifications, beautiful status pages, and comprehensive analytics.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go](https://img.shields.io/badge/go-%23007d9c.svg?style=flat&logo=go&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=flat&logo=react&logoColor=%2361dafb)
![TypeScript](https://img.shields.io/badge/typescript-%23007acc.svg?style=flat&logo=typescript&logoColor=white)
![MongoDB](https://img.shields.io/badge/mongodb-4ea94b.svg?style=flat&logo=mongodb&logoColor=white)

![Peekaping Dashboard](./pictures/monitor.png)

## ✨ Features

### 🔍 **Monitoring Types**
- **HTTP/HTTPS Monitoring** - Monitor websites, APIs, and web services
- **Push Monitoring** - Monitor services that push heartbeats to Peekaping

### 📊 **Real-time Dashboard**
- Live status updates with WebSocket connections
- Interactive charts and statistics
- 24-hour uptime tracking
- Response time monitoring (ping)
- Visual heartbeat history

### 🔔 **Smart Notifications**
- **Multiple Channels**: Email (SMTP), Slack, Telegram, Webhooks
- **Intelligent Alerting**: Configurable retry logic before marking as down
- **Notification Control**: Set resend intervals to avoid spam
- **Important Events**: Only get notified when status actually changes

### 📄 **Status Pages**
- **Public Status Pages** - Share service status with your users

### 🛠 **Advanced Features**
- **Maintenance Windows** - Schedule maintenance to prevent false alerts
- **Proxy Support** - Route monitoring through HTTP proxies
- **Multi-user Authentication** - Secure login with 2FA support
- **Real-time Collaboration** - Multiple users can monitor simultaneously
- **Data Retention** - Automatic cleanup of old heartbeat data

### 🏗 **Technical Highlights**
- **Modern Stack**: Go backend, React frontend, MongoDB database
- **Cloud Native**: Docker support with docker-compose
- **API First**: Complete REST API with Swagger documentation
- **Real-time**: WebSocket connections for live updates
- **Scalable**: Architecture with dependency injection

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ and pnpm (for development)
- Go 1.24+ (for development)


## 🛠 Development Setup

### Full Stack Development
**Create environment file in root:**
```bash
cp .env.example .env
# Edit .env with your configuration
```

```bash
# Install all dependencies
pnpm install

# run turbo
turbo run dev docs:watch
```


## ⚙️ Configuration

### Environment Variables

```env
# Database
DB_USER=root
DB_PASSWORD=your-secure-password
DB_NAME=peekaping
DB_HOST=localhost
DB_PORT=6001
```

## 🔒 Security

### Authentication
- JWT-based authentication
- Optional 2FA with TOTP
- Secure session management

### Best Practices
- Use strong passwords and JWT secrets
- Enable 2FA for all users
- Regular security updates
- Secure your MongoDB instance
- Use HTTPS in production

## 🚀 Production Deployment

### Docker Production Setup

1. **Build production images:**
```bash
# Build server
cd apps/server
docker build -t peekaping/server:latest .

# Build web client
cd ../web
docker build -t peekaping/web:latest .
```

2. **Deploy with optimized compose:**
```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Reverse Proxy Setup

Example Nginx configuration included in `infra/nginx.conf`.

## 🤝 Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request


## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Uptime Kuma](https://github.com/louislam/uptime-kuma)
- Built with amazing open-source technologies
- Thanks to all contributors and users

## 📞 Support

- **Issues**: Report bugs and request features via GitHub Issues
---

**Made with ❤️ by the Peekaping team**
