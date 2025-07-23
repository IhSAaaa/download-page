# URL Shortener Service

A high-performance URL shortening service built with Go (Golang) and Gin framework. Features include URL shortening, QR code generation, click analytics, and a modern web interface.

## 🚀 Features

- **URL Shortening**: Create short URLs with custom codes
- **QR Code Generation**: Automatic QR code generation for short URLs
- **Click Analytics**: Track clicks, countries, browsers, and devices
- **Expiration Dates**: Set expiration dates for URLs
- **Modern UI**: Responsive web interface with Vue.js
- **RESTful API**: Complete API for integration
- **High Performance**: Built with Go for optimal performance
- **Database Support**: PostgreSQL with proper indexing

## 🛠 Tech Stack

### Backend
- **Language**: Go (Golang) 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **QR Codes**: go-qrcode
- **Containerization**: Docker

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **State Management**: React Query
- **Routing**: React Router DOM
- **Icons**: Lucide React
- **Notifications**: React Hot Toast
- **Charts**: Recharts

## 📋 Prerequisites

### Backend
- Go 1.21 or higher
- PostgreSQL 12 or higher

### Frontend
- Node.js 18 or higher
- npm or yarn

### General
- Git
- Docker (optional)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd url-shortener
```

### 2. Backend Setup

```bash
# Set up environment
cp env.example .env
# Edit .env with your database credentials

# Install Go dependencies
go mod download

# Set up PostgreSQL database
CREATE DATABASE url_shortener;

# Run the backend
go run main.go
```

The backend will be available at `http://localhost:8080`

### 3. Frontend Setup

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

The frontend will be available at `http://localhost:3000`

### 4. Using Docker (Recommended)

```bash
# Make scripts executable (first time only)
chmod +x scripts/*.sh

# Start all services with Docker
./scripts/start.sh

# Or use Docker Compose directly
docker-compose up --build -d
```

### 5. Development Mode

```bash
# Initialize dependencies using Docker (no Go/Node.js required on host)
./scripts/init-all-docker.sh

# Start backend and frontend with Docker (no Go/Node.js required on host)
./scripts/dev-docker.sh

# Or start backend and frontend separately (requires Go/Node.js on host)
./scripts/dev.sh

# Or start individual services
./scripts/dev.sh backend
./scripts/dev.sh frontend
```

## 📖 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### Create Short URL
```http
POST /shorten
Content-Type: application/json

{
  "original_url": "https://example.com/very-long-url",
  "custom_code": "my-link",
  "title": "My Awesome Link",
  "description": "Description of the link",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

#### Get All URLs
```http
GET /urls?page=1&limit=10
```

#### Get URL by ID
```http
GET /urls/{id}
```

#### Delete URL
```http
DELETE /urls/{id}
```

#### Get URL Analytics
```http
GET /analytics/{id}
```

#### Get All Analytics
```http
GET /analytics
```

#### Redirect to Original URL
```http
GET /{shortCode}
```

## 🏗 Project Structure

```
url-shortener/
├── main.go                 # Backend entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksum
├── .env                    # Environment variables
├── env.example            # Example environment file
├── README.md              # Project documentation
├── Dockerfile             # Backend Docker configuration
├── docker-compose.yml     # Multi-container setup
├── scripts/               # Automation scripts
│   ├── setup.sh           # 🚀 Complete setup (init + build + start)
│   ├── start.sh           # Start all services (production)
│   ├── stop.sh            # Stop all services
│   ├── dev.sh             # Development mode (requires Go/Node.js)
│   ├── dev-docker.sh      # Development mode (Docker only)
│   ├── init-simple.sh     # Initialize all dependencies (Docker only)
│   ├── init-all-docker.sh # Initialize dependencies (Docker only, alternative)
│   ├── init-go-docker.sh  # Initialize Go dependencies (Docker only)
│   └── init-frontend-docker.sh # Initialize frontend dependencies (Docker only)
├── models/                # Data models
│   └── url.go
├── handlers/              # HTTP handlers
│   └── url_handlers.go
├── database/              # Database configuration
│   └── database.go
├── templates/             # Legacy HTML template
│   └── index.html
└── frontend/              # React frontend
    ├── package.json       # Frontend dependencies
    ├── vite.config.ts     # Vite configuration
    ├── tailwind.config.js # Tailwind CSS config
    ├── tsconfig.json      # TypeScript config
    ├── Dockerfile         # Frontend Docker configuration
    ├── nginx.conf         # Nginx configuration
    ├── index.html         # Frontend entry point
    ├── README.md          # Frontend documentation
    └── src/               # Source code
        ├── components/    # React components
        ├── pages/         # Page components
        ├── services/      # API services
        ├── types/         # TypeScript types
        ├── App.tsx        # Main app component
        ├── main.tsx       # Entry point
        └── index.css      # Global styles
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (debug/release) | `debug` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `url_shortener` |
| `DB_SSLMODE` | Database SSL mode | `disable` |

## 🐳 Docker Support

### Docker-Only Approach (No Go/Node.js Required)

If you don't have Go or Node.js installed on your host machine, you can use the Docker-only approach:

```bash
# 🚀 ONE COMMAND SETUP: Initialize, build, and start everything
./scripts/setup.sh

# Or step by step:
# Initialize all dependencies using Docker
./scripts/init-simple.sh

# Start all services
./scripts/start.sh

# Development mode with Docker
./scripts/dev-docker.sh
```

**Note:** `setup.sh` is the ultimate one-command solution that does everything. `init-simple.sh` is recommended for just initialization as it uses Docker volumes to avoid permission issues.

### Quick Start with Scripts

```bash
# Start all services
./scripts/start.sh

# Stop all services
./scripts/stop.sh

# View logs
./scripts/start.sh logs

# Check status
./scripts/start.sh status
```

### Manual Docker Commands

```bash
# Build and start all services
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild specific service
docker-compose up --build -d backend
```

### Service URLs (Docker)

- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Database**: localhost:5432

## 📊 Features in Detail

### URL Shortening
- Generate random 6-character codes
- Support for custom codes (3-50 characters)
- URL validation and sanitization
- Duplicate custom code prevention

### Analytics
- Click tracking with IP addresses
- User agent parsing
- Geographic location (future enhancement)
- Device and browser detection
- Click timeline (last 30 days)
- Unique vs total clicks

### QR Codes
- Automatic QR code generation
- Base64 encoded PNG format
- Medium error correction level
- 256x256 pixel size

### Web Interface
- Modern, responsive design
- Real-time form validation
- Copy-to-clipboard functionality
- URL management dashboard
- Analytics visualization

## 🔒 Security Features

- Input validation and sanitization
- SQL injection prevention
- CORS configuration
- Rate limiting (future enhancement)
- HTTPS support (production)

## 🚀 Performance Optimizations

- Database connection pooling
- Proper indexing on frequently queried columns
- Asynchronous click tracking
- Efficient QR code generation
- Optimized database queries

## 🧪 Testing

Run tests:
```bash
go test ./...
```

## 📈 Monitoring

### Health Check
```http
GET /health
```

### Metrics (Future Enhancement)
- Prometheus metrics
- Grafana dashboards
- Application performance monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue on GitHub
- Check the documentation
- Review the API examples

## 🔮 Future Enhancements

- [ ] User authentication and authorization
- [ ] Rate limiting and API quotas
- [ ] Geographic location tracking
- [ ] Advanced analytics and reporting
- [ ] Bulk URL operations
- [ ] API rate limiting
- [ ] Webhook notifications
- [ ] Mobile application
- [ ] Social media integration
- [ ] Custom domains support

---

**Built with ❤️ using Go and Gin** 