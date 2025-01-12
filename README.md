# Domain Management System

A modern web application for managing domains and email accounts with real-time notifications.

## Features

- Domain Management (CRUD operations)
- Email Account Management
- Real-time WebSocket Notifications
- User Authentication & Authorization
- Responsive UI with Dynamic Updates
- Bulk Email Account Creation
- Domain Status Monitoring

## Tech Stack

### Backend

- **Go 1.21**
- **Chi Router** - Lightweight HTTP routing
- **SQLC** - Type-safe SQL in Go
- **PostgreSQL** - Primary database
- **WebSocket** - Real-time notifications
- **JWT** - Authentication
- **Templ** - Type-safe templates

### Frontend

- **HTMX** - Dynamic UI updates
- **TailwindCSS** - Styling
- **Font Awesome** - Icons

### Development Tools

- **Docker** - Containerization
- **Make** - Build automation
- **Air** - Live reload
- **Testify** - Testing framework

### Project Structure

```
├── assets/ # Static files (JS, CSS)
├── cmd/ # Application entrypoints
├── internal/
│ ├── api/ # HTTP handlers
│ ├── auth/ # Authentication
│ ├── db/ # Database layer
│ │ ├── queries/ # SQL queries
│ │ └── sqlc/ # Generated code
│ ├── models/ # Data models
│ ├── services/ # Business logic
│ └── views/ # UI templates
├── migrations/ # SQL migrations
└── scripts/ # Build scripts
```

## Testing

### Key Features Implementation

#### WebSocket Notifications

- Real-time notifications for domain status changes
- Automatic reconnection on connection loss
- Notification badge updates
- Toast notifications for immediate feedback

#### Domain Management

- Create and manage domains
- Update domain status
- Bulk email account creation
- Status monitoring and notifications

#### Authentication

- JWT-based authentication
- Secure session management
- Protected API endpoints