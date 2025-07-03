#  lite-calendar

A lightweight calendar backend built in Go — fast, modular, containerized. Designed for clarity, testability, and extensibility (Kafka, workers, CI/CD).
Clean layered architecture: `handler → service → repository`
---

##  Tech Stack

- **Language:** Go 1.24+
- **Database:** PostgreSQL
- **Router:** chi
- **DB Layer:** sqlx
- **Migrations:** Goose
- **Logging:** slog
- **Containerization:** Docker + Docker Compose

---

## 🚀 Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/YOUR_USERNAME/lite-calendar-v2.git
   cd lite-calendar-v2
