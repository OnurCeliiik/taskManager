# Task Manager API

A production-ready REST API for managing tasks with secure user authentication, built with Go and Gin.

## Features

### User Management
- User registration and login with JWT authentication
- Secure password hashing with bcrypt
- Password reset flow with time-limited tokens (1 hour expiry)
- User profile management (get, update, delete)
- Role-based access control (user, admin)

### Task Management
- Full CRUD operations (Create, Read, Update, Delete)
- Task ownership enforcement - users can only manage their own tasks
- Task filtering by status (pending, in_progress, completed) and category
- Task pagination support
- Soft deletes with timestamps

### Security & Infrastructure
- JWT-based authentication (15-minute token expiry)
- Role-based authorization middleware
- Rate limiting (10 req/sec global, 5 req/sec for auth endpoints)
- Request logging with structured JSON output
- Password strength validation (min 6 chars, uppercase, lowercase, digit, special char)
- Email notifications (registration welcome, password reset) via Mailtrap/SMTP
- Docker containerization for easy deployment

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login and get JWT token |
| POST | `/api/v1/auth/forgot-password` | Request password reset |
| POST | `/api/v1/auth/reset-password` | Reset password with token |

### Users (requires authentication)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/:id` | Get user details |
| PUT | `/api/v1/users/:id` | Update user profile |
| DELETE | `/api/v1/users/:id` | Delete user account |

### Tasks (requires authentication)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tasks/create` | Create a new task |
| GET | `/api/v1/tasks` | List all tasks (with filters) |
| GET | `/api/v1/tasks/:id` | Get task details |
| PUT | `/api/v1/tasks/:id` | Update task |
| DELETE | `/api/v1/tasks/:id` | Delete task |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Health check |

## Authentication

All endpoints except `/healthz`, `/api/v1/auth/*` require a valid JWT token in the `Authorization` header:
