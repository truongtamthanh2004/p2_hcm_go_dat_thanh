# 🏢 Co-working Space Booking API

This is a backend system built with Go (Gin framework) following a **microservices architecture**, designed for booking and managing co-working spaces. It supports features like authentication, real-time communication, payments, email/SMS notifications, statistics, and more.

## 🎯 Features

- ✅ User registration, login, and email verification (normal & OAuth2)
- 🧑‍💼 Role-based access (Customer, Moderator, Admin, System)
- 📍 Location & room management (availability, filters, images, amenities)
- 📅 Booking system with check-in/check-out & status tracking
- 💳 Payment integration (e.g., VNPAY, Internet Banking)
- 📈 Admin statistics (revenue, users, bookings...)
- 💬 Real-time chat using WebSocket
- 🌐 Multi-language support (i18n)
- 📧 Email & SMS notifications
- 🗺️ Map view for location search
- ⏳ Job scheduling (e.g., daily reminders, cron tasks)
- 🔁 Message queue & rate limiting support

## 🗂️ Project Structure

```
/services
  /user-service           # Handles user registration, authentication, profile, roles (Customer, Moderator, Admin)
  /location-service       # Manages coworking spaces, rooms, availability, location search
  /booking-service        # Manages bookings, calendar, check-in/check-out, payment integration
  /payment-service        # Handles online payments, transactions, VNPAY or other gateways
  /auth-service           # JWT token management, email verification, password reset
  /chat-service           # WebSocket-based chat system for real-time communication
  /admin-service          # Admin-specific dashboards, statistics, user management
  /notification-service   # Email & SMS notifications, async job processing
/pkg
  /utils                  # Common utilities used across services (e.g., password hash, response format)
  /middlewares            # Auth middlewares, rate limiting, logging
  /config                 # Centralized configuration loading (YAML, ENV)
  /docs/openapi           # OpenAPI (Swagger) specifications for services
```

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/your-username/co-working-booking-api.git
cd co-working-booking-api
```

### 2. Spin up services using Docker Compose

```bash
docker-compose up --build
```

> Note: Ensure Docker and Docker Compose are installed.

### 3. Run unit tests

```bash
go test ./...
```

## 📚 Documentation

- Each service has its own Swagger/OpenAPI docs in `/docs/openapi`
- Run Swagger UI with Docker or use [Swagger Editor](https://editor.swagger.io/) to visualize the specs

## 👥 Contributing

1. Fork the repo
2. Create a new branch (`git checkout -b feature/awesome-feature`)
3. Commit your changes (`git commit -am 'Add awesome feature'`)
4. Push to the branch (`git push origin feature/awesome-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.
