# Helpdesk REST API

A production-ready REST API for IT Helpdesk management system built with Go, Echo framework, and PostgreSQL. The API follows clean architecture principles with repository, service, and handler layers.

**Project Structure:** [ERD Diagram](https://app.eraser.io/workspace/MCKUzCCls92JCU5rpuew?origin=share)

## Stack

- **Language:** Go 1.25+
- **Web Framework:** Echo v5
- **Database:** PostgreSQL
- **Query Builder:** sqlx
- **Database Migration:** Goose
- **Logging:** Structured logging (slog)
- **Containerization:** Docker

## Development (Hot Reload)

Use Air for hot reload:

```bash
air
```

This project intentionally builds to `tmp/main.exe` in [.air.toml](.air.toml).

- On Windows: CMD requires `.exe` to execute the binary.
- On macOS/Linux: `.exe` is only a file name suffix, and the binary still runs normally.
