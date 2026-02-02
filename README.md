# Go SaaS API

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/doc/devel/release)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

REST API yang dibangun dengan Go (Golang) menggunakan Gin framework. Project ini mengimplementasikan clean architecture dengan pemisahan concerns yang jelas untuk management User, Product, dan Place.

## 📋 Daftar Isi

- [Fitur](#fitur)
- [Teknologi](#teknologi)
- [Prasyarat](#prasyarat)
- [Instalasi](#instalasi)
- [Konfigurasi](#konfigurasi)
- [Cara Menjalankan](#cara-menjalankan)
- [Struktur Project](#struktur-project)
- [API Endpoints](#api-endpoints)
- [Database](#database)
- [Development](#development)

## ✨ Fitur

- **Authentication & Authorization** - JWT-based authentication dengan middleware
- **User Management** - CRUD operations untuk user dengan password encryption
- **Product Management** - Management produk lengkap dengan validasi
- **Place Management** - Management tempat/lokasi dengan schedule planning
- **Input Validation** - Validasi input menggunakan playground validator
- **Error Handling** - Consistent error responses
- **Database** - MySQL dengan prepared statements untuk security
- **Logging** - Comprehensive logging dengan Gin logger

## 🛠 Teknologi

| Teknologi | Versi | Deskripsi |
|-----------|-------|-----------|
| Go | 1.24.0 | Programming Language |
| Gin | 1.11.0 | Web Framework |
| MySQL | 5.7+ | Database |
| JWT | 5.3.1 | Authentication |
| SQLx | 1.4.0 | Database Abstraction |
| godotenv | 1.5.1 | Environment Variables |
| Validator | 10.27.0 | Input Validation |

## 📦 Prasyarat

- Go 1.24.0 atau lebih tinggi
- MySQL 5.7 atau lebih tinggi
- Make (opsional, untuk menggunakan Makefile)

## 🚀 Instalasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd golang-api
```

### 2. Install Dependencies

```bash
go mod download
go mod tidy
```

### 3. Setup Database

```bash
# Create database
mysql -u root -p < db/db.sql

# Atau via MySQL client
mysql -u root -p
mysql> CREATE DATABASE golang_api;
mysql> USE golang_api;
mysql> SOURCE db/db.sql;
```

### 4. Setup Environment Variables

Buat file `.env` di root project:

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_DSN=root:password@tcp(localhost:3306)/golang_api?parseTime=true

# JWT Configuration
JWT_SECRET=your-secret-key-change-this-in-production
```

**Penjelasan konfigurasi:**
- `PORT` - Port server (default: 8080)
- `DB_DSN` - Database connection string MySQL
- `JWT_SECRET` - Secret key untuk signing JWT token (gunakan string yang kuat di production)

## 🎯 Cara Menjalankan

### Development Mode

```bash
# Menggunakan make
make run

# Atau langsung dengan go
go run cmd/api/main.go
```

Server akan berjalan di `http://localhost:8080`

### Build Binary

```bash
# Menggunakan make
make build

# Atau langsung dengan go
go build -o bin/api cmd/api/main.go
```

Jalankan binary:

```bash
./bin/api
```

### Testing

```bash
make test
```

### Cleaning

```bash
make clean
```

### Help

```bash
make help
```

## 📁 Struktur Project

```
golang-api/
├── bin/                    # Build output
│   └── api                 # Binary executable
├── cmd/
│   └── api/
│       └── main.go         # Entry point aplikasi
├── db/
│   └── db.sql              # Database schema
├── internal/               # Private application code
│   ├── config/             # Configuration management
│   │   └── config.go       # Load env config
│   ├── database/           # Database connection
│   │   └── database.go     # MySQL connection setup
│   ├── middleware/         # Middleware
│   │   └── auth.go         # JWT authentication
│   ├── user/               # User module
│   │   ├── dto.go          # Data Transfer Objects
│   │   ├── handler.go      # HTTP handlers
│   │   ├── model.go        # Data models
│   │   ├── repository.go   # Database operations
│   │   ├── routes.go       # Route definitions
│   │   └── service.go      # Business logic
│   ├── product/            # Product module
│   │   ├── dto.go
│   │   ├── handler.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── routes.go
│   │   └── service.go
│   └── place/              # Place module
│       ├── dto.go
│       ├── handler.go
│       ├── model.go
│       ├── repository.go
│       ├── routes.go
│       └── service.go
├── pkg/                    # Reusable packages
│   └── response/           # Response utilities
│       └── response.go     # Standard response format
├── go.mod                  # Module definition
├── go.sum                  # Dependency checksums
├── Makefile                # Build commands
└── README.md               # Dokumentasi (file ini)
```

## 🔌 API Endpoints

### Health Check

```
GET /health
```

Response:
```json
{
  "status": "ok"
}
```

### User Module

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| POST | `/api/users/register` | Register user baru | ❌ |
| POST | `/api/users/login` | Login user | ❌ |
| GET | `/api/users/profile` | Get user profile | ✅ |
| PUT | `/api/users/:id` | Update user | ✅ |
| DELETE | `/api/users/:id` | Delete user | ✅ |

### Product Module

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| GET | `/api/products` | List semua products | ✅ |
| POST | `/api/products` | Create product baru | ✅ |
| GET | `/api/products/:id` | Get product detail | ✅ |
| PUT | `/api/products/:id` | Update product | ✅ |
| DELETE | `/api/products/:id` | Delete product | ✅ |

### Place Module

| Method | Endpoint | Deskripsi | Auth |
|--------|----------|-----------|------|
| GET | `/api/places` | List semua places | ✅ |
| POST | `/api/places` | Create place baru | ✅ |
| GET | `/api/places/:id` | Get place detail | ✅ |
| PUT | `/api/places/:id` | Update place | ✅ |
| DELETE | `/api/places/:id` | Delete place | ✅ |

## 🗄 Database

### Tables

#### users
- `id` - Primary key
- `email` - Unique email
- `password` - Hashed password
- `name` - User name
- `created_at` - Timestamp
- `updated_at` - Timestamp

#### products
- `id` - Primary key
- `user_id` - Foreign key to users
- `name` - Product name
- `description` - Product description
- `price` - Product price
- `created_at` - Timestamp
- `updated_at` - Timestamp

#### places
- `id` - Primary key
- `user_id` - Foreign key to users
- `name` - Place name
- `link` - Place URL/link
- `description` - Place description
- `go_at` - Planned visit date
- `go_at_time` - Planned visit time
- `status` - Place status
- `created_at` - Timestamp
- `updated_at` - Timestamp

### Relasi Database

- **User** → **Product**: 1 user memiliki banyak products
- **User** → **Place**: 1 user memiliki banyak places

## 👨‍💻 Development

### Project Architecture

Project ini mengikuti **Clean Architecture** pattern:

```
Presentation Layer (Handlers)
       ↓
Business Logic Layer (Services)
       ↓
Data Access Layer (Repositories)
       ↓
Database Layer
```

### Conventions

- **Naming**: CamelCase untuk functions/variables, snake_case untuk database columns
- **Error Handling**: Konsisten menggunakan error wrapping
- **Validation**: Gunakan struct tags untuk validasi
- **Authentication**: JWT token di header `Authorization: Bearer <token>`

### Adding New Module

1. Buat folder di `internal/<module-name>/`
2. Implementasi file standar:
   - `model.go` - Data structure
   - `dto.go` - Request/Response objects
   - `repository.go` - Database operations
   - `service.go` - Business logic
   - `handler.go` - HTTP handlers
   - `routes.go` - Route definition
3. Register di `cmd/api/main.go` dalam `setup<Module>Module()` function

### Common Commands

```bash
# Download dependencies
go mod download

# Tidy dependencies
make tidy

# Run tests
make test

# Clean build artifacts
make clean

# Format code
go fmt ./...

# Lint code
golangci-lint run ./...
```

## 🔐 Security Notes

- **JWT Secret**: Selalu gunakan secret yang kuat di production
- **Password**: Passwords di-hash menggunakan bcrypt
- **CORS**: Configure CORS sesuai kebutuhan production
- **Environment Variables**: Jangan commit `.env` file ke repository
- **Database Credentials**: Gunakan strong passwords dan limit user permissions

## 📝 License

MIT License - Lihat LICENSE file untuk detail

## 🤝 Contributing

1. Fork repository
2. Buat feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Open Pull Request

## 📧 Contact

Untuk pertanyaan atau support, silakan buka issue di repository.

---

**Last Updated**: February 3, 2026
