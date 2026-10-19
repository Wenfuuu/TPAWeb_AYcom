# AYcom

AYcom is a Twitter-style social media web app built as a set of Go microservices behind a single API gateway, with a Svelte frontend.
Each service follows a layered architecture (handler, server, repository, model) and talks to the others over gRPC.

## Features

- **Threads** with posting, replies, reposts, likes, bookmarks, media attachments, hashtags and categories.
- **Feeds** for all threads, followed users, communities and top hashtags, plus fuzzy thread search using Damerau-Levenshtein distance.
- **Communities** with create and join requests, member lists, community media and categories managed by admins.
- **Profiles and follows** with followers, following, who-to-follow suggestions and profile image and video uploads.
- **Realtime direct and group messaging** over WebSocket, with mention and activity notifications.
- **Auth** with JWT access and refresh tokens, email verification codes, security-question password reset and Google reCAPTCHA.
- **Premium and moderation** flows where users request premium or report others, and admins accept, reject or ban from the admin page.
- **Voice-to-text input** through a Whisper model fine-tuned on the PolyAI MINDS-14 dataset.

## Tech Stack

| Layer | Tools |
| --- | --- |
| Frontend | Svelte 5, TypeScript, Vite, svelte-routing, Axios |
| API gateway | Go, Gin, Swagger (swaggo) |
| Services | Go, gRPC, GORM, Zap logging, gomail |
| AI service | Python, Flask, Hugging Face Transformers, PyTorch, librosa |
| Data and messaging | PostgreSQL, Redis, RabbitMQ, Supabase Storage |
| Tooling | Docker Compose, Air (hot reload) |

## Architecture

The frontend calls the API gateway over REST, and the gateway forwards each request to the right service over gRPC.
Services that need to send email publish to RabbitMQ, and the email service consumes those messages and sends them over SMTP.
Uploaded images and videos go to Supabase Storage.

| Service | Port | Responsibility |
| --- | --- | --- |
| `frontend` | 5173 | Svelte web client |
| `api-gateway` | 8080 | REST entry point, auth middleware, Swagger docs |
| `service-user` | 50051 | Users, follows, uploads, premium, reports, bans |
| `service-auth` | 50052 | Login, register, JWT tokens, verification codes (Redis) |
| `service-email` | 50053 | Email delivery from RabbitMQ messages |
| `service-thread` | 50054 | Threads, communities, categories, search |
| `service-chat` | 50055, 50056 | Conversations and messages over gRPC and WebSocket |
| `service-ai` | 5000 | Speech-to-text transcription |
| `db-user`, `db-chat` | 5432, 5433 | PostgreSQL databases |
| `rabbitmq`, `redis` | 5672, 6379 | Message broker and cache |

```
.
├── tpa-frontend/          Svelte client (pages, components, stores)
├── tpa-backend/
│   ├── api-gateway/       Gin REST gateway and Swagger docs
│   ├── pkg/proto/         Shared generated gRPC code
│   └── service-*/         One folder per microservice
├── database/              Postgres init scripts
├── ai/                    Whisper fine-tuning notebook and model link
└── docker-compose.yml
```

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) with Docker Compose
- A Supabase project with a storage bucket
- An SMTP account for sending verification emails
- A Google reCAPTCHA site key

### Setup

1. Clone the repo.

   ```bash
   git clone https://github.com/Wenfuuu/TPAWeb_AYcom.git
   cd TPAWeb_AYcom
   ```

2. Fill in the `.env` file in each `tpa-backend/service-*` and `tpa-backend/api-gateway` folder with your own database, Redis, RabbitMQ, Supabase, SMTP and JWT values.

3. Download the fine-tuned Whisper model from the link in `ai/model_link.txt` into `tpa-backend/service-ai/model`.
   You can also train your own with `ai/tpa-web.ipynb`.

4. Start everything.

   ```bash
   docker compose up --build
   ```

5. Open the app at `http://localhost:5173` and the API docs at `http://localhost:8080/swagger/index.html`.

## Course Context

This app was built for the Web Technical Path Assessment (TPA) at Software Laboratory Center, BINUS University.
