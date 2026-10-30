# CapIt

**Your content. The perfect caption.**

CapIt is an AI-powered social media caption generator designed to help content creators, social media managers, digital marketers, photographers, and communication teams turn rough content descriptions into publish-ready captions.

The project uses a lightweight frontend built with **HTML, CSS, and vanilla JavaScript**, backed by **Go**, `net/http`, and **Chi** for HTTP routing.

## Goals

- Generate captions for different social media platforms.
- Customize captions by content type, tone, language, and length.
- Edit, copy, and regenerate captions.
- Keep the frontend simple and responsive.
- Keep AI provider credentials and integrations on the backend.
- Build incrementally without unnecessary frameworks or dependencies.

## Technology Stack

| Area | Technology |
|---|---|
| Frontend | HTML5, CSS3, vanilla JavaScript |
| Backend | Go |
| HTTP routing | Chi |
| HTTP client/server communication | Fetch API and JSON |
| Initial caption history | Browser `localStorage` |
| Future database | PostgreSQL, when persistence is needed |
| AI integration | Provider adapter behind a Go interface |

Templ, HTMX, and Tailwind CSS are optional and are not required for the initial version.

## Project Structure

```text
capit/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
│
├── internal/
│   ├── config/
│   │   └── config.go               # Environment and application configuration
│   │
│   ├── domain/
│   │   └── caption.go              # Caption request/response data types
│   │
│   ├── handler/
│   │   ├── caption.go              # Caption API handlers
│   │   └── health.go               # Health-check endpoint
│   │
│   ├── service/
│   │   └── caption.go              # Caption-generation business logic
│   │
│   ├── ai/
│   │   ├── provider.go             # AI provider interface
│   │   └── mock.go                 # Mock provider for development
│   │
│   ├── middleware/
│   │   └── logging.go              # HTTP request logging
│   │
│   └── repository/
│       └── caption.go              # Caption persistence (add when needed)
│
├── web/
│   ├── index.html                  # Main frontend page
│   └── assets/
│       ├── css/
│       │   └── style.css           # Application styles
│       ├── js/
│       │   └── app.js              # Frontend interactions and API calls
│       └── images/
│           └── logo.svg            # CapIt logo
│
├── migrations/                     # Database migrations (add when needed)
├── .env.example                    # Example environment variables
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

> Start with the files needed for the first working feature. Directories such as `repository/` and `migrations/` can be added when persistent storage is introduced.

## Architecture

```text
                         USER
                          |
                          v
                 +-------------------+
                 |     CapIt UI      |
                 | index.html        |
                 | style.css         |
                 | app.js            |
                 +---------+---------+
                           |
                           | HTTP + JSON
                           v
                 +-------------------+
                 | Go HTTP Server    |
                 | Chi Router        |
                 +---------+---------+
                           |
                           v
                 +-------------------+
                 | HTTP Handlers     |
                 +---------+---------+
                           |
                           v
                 +-------------------+
                 | Caption Service   |
                 | Validation        |
                 | Prompt preparation|
                 | Generation logic  |
                 +---------+---------+
                           |
                           v
                 +-------------------+
                 | AI Provider       |
                 | Adapter/interface |
                 +---------+---------+
                           |
                           v
                 +-------------------+
                 | AI Model or API   |
                 +-------------------+
```

### Request flow

1. The user enters a content description and selects caption options.
2. `app.js` sends a JSON request to the Go API using the Fetch API.
3. Chi routes the request to the relevant handler.
4. The handler validates the request and calls the caption service.
5. The service prepares the generation request and calls the configured AI provider.
6. The handler returns the result as JSON.
7. The frontend displays the caption so the user can edit or copy it.

The frontend should not depend on a particular AI provider. The Go backend owns provider integration so the provider can be changed without redesigning the UI.

## Initial API

| Method | Endpoint | Purpose |
|---|---|---|
| `GET` | `/` | Serves the CapIt frontend |
| `GET` | `/assets/*` | Serves CSS, JavaScript, and images |
| `GET` | `/api/v1/health` | Reports whether the backend is running |
| `POST` | `/api/v1/captions/generate` | Generates a caption |

### Example request

`POST /api/v1/captions/generate`

```json
{
  "content": "A community meeting discussing access to clean water.",
  "platform": "Facebook",
  "contentType": "Community Engagement",
  "tone": "Professional",
  "language": "English",
  "length": "Short",
  "instructions": "Keep it concise and engaging."
}
```

### Example response

```json
{
  "caption": "Building stronger communities starts with listening. Today, residents came together to discuss access to clean water and practical solutions for a healthier future."
}
```

The API contract can be extended as the application grows, for example to return alternative captions or metadata.

## Backend Responsibilities

### `cmd/server/`

Starts the application, loads configuration, sets up routes, and starts the HTTP server.

### `internal/handler/`

Handles HTTP concerns:
- Decode incoming JSON.
- Validate request structure.
- Call the appropriate service.
- Return consistent JSON responses and HTTP status codes.

### `internal/service/`

Contains caption-generation business logic:
- Validate supported options.
- Prepare the generation request.
- Call the AI provider.
- Handle service-level errors.

### `internal/ai/`

Contains the AI provider interface and implementations. A mock provider can be used to test the request-response flow before connecting a real model.

Example interface:

```go
package ai

import "context"

type Provider interface {
    GenerateCaption(ctx context.Context, prompt string) (string, error)
}
```

### `internal/domain/`

Defines core data types such as caption-generation requests and responses.

### `internal/config/`

Loads configuration from environment variables and applies safe defaults where appropriate.

### `internal/middleware/`

Holds shared HTTP middleware, such as request logging. Add middleware only when it serves a clear purpose.

### `internal/repository/`

Provides persistence when server-side caption history or user data is needed. It is not required for the initial prototype.

## Frontend Responsibilities

### `web/index.html`

Provides the semantic HTML structure and form controls.

### `web/assets/css/style.css`

Contains responsive styling, layout, colors, spacing, and focus states.

### `web/assets/js/app.js`

Handles:
- Form submission and validation.
- Fetch requests to the Go API.
- Loading, success, and error states.
- Caption editing and copying.
- Character counts.
- Caption history in `localStorage`.
- Regeneration and alternative captions when supported.

Do not put AI provider API keys in frontend JavaScript or browser storage.

## Development Roadmap

### Stage 1: Working prototype

- [ ] Create the HTML, CSS, and JavaScript frontend.
- [ ] Set up the Go server and Chi router.
- [ ] Serve the frontend from Go.
- [ ] Implement `GET /api/v1/health`.
- [ ] Implement `POST /api/v1/captions/generate`.
- [ ] Add a mock AI provider.
- [ ] Display, edit, and copy generated captions.
- [ ] Store recent captions in browser `localStorage`.

**Goal:** Prove that the frontend-to-backend flow works before spending money on an AI API.

### Stage 2: Real AI integration

- [ ] Implement a real AI provider behind the provider interface.
- [ ] Load provider credentials from environment variables.
- [ ] Add request timeouts and robust error handling.
- [ ] Add suitable logging and basic rate limiting.
- [ ] Test malformed requests, provider failures, and network errors.

### Stage 3: Persistence and accounts

When the product needs these features:

- [ ] Add PostgreSQL.
- [ ] Introduce database migrations.
- [ ] Save caption history on the server.
- [ ] Add authentication and authorization.
- [ ] Add user-specific settings and usage limits.
- [ ] Consider subscriptions and payments if the business model requires them.

Avoid introducing a database, authentication, or payment features before the core caption-generation flow is reliable.

## Configuration and Secrets

Use `.env.example` to document the environment variables required for local development. Keep real secrets in a local environment file or a secrets manager, and ensure secret files are ignored by Git.

For example, a future AI integration might need variables such as:

```text
PORT=3000
AI_PROVIDER=
AI_API_KEY=
```

These are illustrative configuration names. The application must explicitly load and validate whichever variables it actually uses. Never commit real API keys or expose them to the browser.

## Running the Project

The exact commands depend on the server setup implemented in `cmd/server/main.go`.

A typical Go development workflow is:

1. Install a supported Go version.
2. Initialize the module if it has not already been created.
3. Add Chi as a dependency if using Chi for routing.
4. Implement the server and route setup.
5. Start the server with `go run ./cmd/server`.
6. Open the local address configured by the server, for example `http://localhost:3000`.

Example setup commands for a new repository:

```bash
go mod init example.com/capit
go get github.com/go-chi/chi/v5
go run ./cmd/server
```

Replace `example.com/capit` with the module path you choose. If `go.mod` already exists, do not initialize it again.

## Design Principles

- **Keep it simple:** use vanilla HTML, CSS, and JavaScript.
- **Separate responsibilities:** HTTP handling, business logic, and AI integration belong in different layers.
- **Avoid premature complexity:** add packages and directories as real needs emerge.
- **Keep provider credentials private:** all AI calls should go through the backend.
- **Make failures clear:** show useful messages for validation, network, and provider errors.
- **Build incrementally:** get one end-to-end feature working before adding more.
- **Stay maintainable:** prefer readable Go code and small, focused functions.

## Future Ideas

Possible features after the initial release:

- Multiple caption variations.
- Platform-specific formatting.
- English, Kiswahili, and Sheng support.
- Tone and length presets.
- Saved caption collections.
- User accounts and cross-device history.
- Usage limits and analytics.

## License

Choose and add a license before distributing or accepting external contributions.
