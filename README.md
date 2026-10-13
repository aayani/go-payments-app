## Payments App

This Go-based payments MVP demonstrates the application lifecycle from development through deployment, with observability and monitoring.

```text
Caveat: This is a local demo, not a production payments system. The API has no authentication and represents money with floating-point values; do not expose it publicly or use it to process real funds. Transaction confirmation IDs are single-use and expire after five minutes. A confirmation ID is consumed before the database write, so a failed write requires initializing a new transaction.
```

### Prerequisites

- Go 1.27 or later
- Terraform 1.16 or later
- Docker with Compose
- `make`

Go and Terraform run inside Docker during setup.

### Local Setup

Start the local stack from the repository root:

```sh
make up
```

This starts MiniStack and the supporting UIs, builds the payments image, and applies Terraform to run one ECS Fargate task behind an ALB. On first run, Docker downloads the required images and Go dependencies. If `.env` is missing, the setup copies `.env.example` to create it.

### Local URLs

| Service                        | URL                                        | Login                         |
| ------------------------------ | ------------------------------------------ | ----------------------------- |
| Payments API documentation     | http://localhost:4566/\_alb/payments/docs/ | None                          |
| StackPort AWS resource browser | http://localhost:8081                      | None                          |
| Grafana                        | http://localhost:3000                      | `admin` / `admin`             |
| pgAdmin                        | http://localhost:8888                      | `admin@example.com` / `admin` |

Check that the API is responding through the ALB:

```sh
curl http://localhost:4566/_alb/payments/ping
```

Stop the local containers with:

```sh
make down
```

Named Docker volumes are retained, so MiniStack state and UI data survive `make down`.
