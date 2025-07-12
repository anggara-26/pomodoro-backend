# Pomodoro Backend API

A complete backend API for a Pomodoro timer application built with Go, Fiber, and MongoDB.

## Features

- **User Management**: Create, read, and update user profiles
- **Task Management**: Full CRUD operations for tasks with status tracking
- **Pomodoro Sessions**: Start and end Pomodoro sessions with automatic task progress tracking
- **Statistics**: Daily, weekly, and monthly productivity statistics
- **Firebase Authentication**: Integrated Firebase authentication middleware
- **API Documentation**: Swagger/OpenAPI documentation

## Tech Stack

- **Framework**: Fiber (Go web framework)
- **Database**: MongoDB
- **Authentication**: Firebase Auth
- **Documentation**: Swagger/OpenAPI
- **Validation**: go-playground/validator

## Getting Started

### Prerequisites

- Go 1.24.1 or later
- MongoDB instance (local or cloud)
- Firebase project with authentication enabled

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd pomodoro-backend
```

2. Install dependencies:

```bash
go mod tidy
```

3. Set up environment variables:
   Copy `.env.example` to `.env` and fill in your values:

```env
BINARY=pomodoro
MONGODB=pomodoro_db
MONGODB_USERNAME=your_username
MONGODB_PASSWORD=your_password
MONGODB_HOST=your_mongodb_host
PORT=8080
```

4. Add your Firebase service account key:
   Place your `serviceAccountKey.json` file in the root directory.

5. Build and run:

```bash
make build
make start
```

Or using Go directly:

```bash
go run ./cmd/api
```

## API Endpoints

### Health Check

- `GET /api/v1/healthcheck` - Check if the server is running

### Users

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user by ID

### Tasks

- `POST /api/v1/tasks` - Create a new task
- `GET /api/v1/tasks/{id}` - Get task by ID
- `PUT /api/v1/tasks/{id}` - Update task by ID
- `DELETE /api/v1/tasks/{id}` - Delete task by ID (soft delete)
- `GET /api/v1/tasks/user/{id}` - Get tasks by user ID with filters and pagination

### Pomodoro Sessions

- `POST /api/v1/sessions/start` - Start a new Pomodoro session
- `POST /api/v1/sessions/end/{id}` - End a Pomodoro session
  - Query parameter: `is_skip=true` to mark session as skipped

### Statistics

- `GET /api/v1/stats/daily` - Get daily statistics
  - Query parameters: `user_id` (required), `date` (optional, format: YYYY-MM-DD)
- `GET /api/v1/stats/weekly` - Get weekly statistics
  - Query parameters: `user_id` (required), `week_start` (optional, format: YYYY-MM-DD)
- `GET /api/v1/stats/monthly` - Get monthly statistics
  - Query parameters: `user_id` (required), `month` (optional, 1-12), `year` (optional)

## Data Models

### User

```json
{
  "id": "ObjectID",
  "firebase_uid": "string",
  "email": "string",
  "name": "string",
  "created_at": "timestamp"
}
```

### Task

```json
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "title": "string",
  "description": "string",
  "assigned_at": "timestamp",
  "status": "pending|in_progress|completed|deleted",
  "estimated_pomodoros": "number",
  "completed_pomodoros": "number",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "deleted_at": "timestamp"
}
```

### Session

```json
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "task_id": "ObjectID",
  "started_at": "timestamp",
  "ended_at": "timestamp",
  "duration": "number (minutes)",
  "type": "focus|short_break|long_break",
  "status": "active|break|skipped|completed"
}
```

## Task Status Flow

- **pending** → **in_progress** → **completed**
- Tasks can be soft-deleted (status: **deleted**)
- Completed pomodoros are automatically incremented when focus sessions are completed
- Tasks are automatically marked as completed when completed_pomodoros >= estimated_pomodoros

## Session Types and Durations

- **Focus**: 25 minutes (work session)
- **Short Break**: 5 minutes
- **Long Break**: 15-30 minutes

## Features

### Automatic Task Progress Tracking

When a focus session is completed (not skipped), the system automatically:

1. Increments the task's completed_pomodoros count
2. Marks the task as completed if completed_pomodoros >= estimated_pomodoros

### Statistics and Analytics

The API provides comprehensive statistics:

- Daily: Sessions completed, tasks finished, total focus time
- Weekly: Aggregated daily stats with daily breakdown
- Monthly: Aggregated weekly stats with weekly breakdown

### Filtering and Pagination

Tasks can be filtered by:

- Status
- Title (partial match)
- Date range (assigned_at)
- Pagination support (page, limit)

## API Documentation

When the server is running, visit:

- Swagger UI: `http://localhost:8080/swagger/`

## Development

### Project Structure

```
.
├── app/
│   ├── handler/          # HTTP handlers
│   └── model/            # Data models
├── cmd/
│   └── api/              # Application entry point
├── docs/                 # Generated Swagger docs
├── pkg/
│   ├── middleware/       # Custom middleware
│   └── router/          # Route definitions
├── platform/
│   ├── db/              # Database connection
│   └── firebase/        # Firebase configuration
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### Building

```bash
make build
```

### Running

```bash
make start
```

### Restart (build + start)

```bash
make restart
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
