
# CapIt

**Your content. The perfect caption.**

CapIt is an AI-powered social media caption generator designed to help content creators, businesses, photographers, and digital marketers create engaging captions for their social media posts.

By providing a topic, choosing a social media platform, and selecting a preferred tone, users can generate captions using Google's Gemini AI.

## Features

- **AI-powered caption generation:** Generate captions using Google's Gemini API.
- **Multiple social media platforms:** Supports Instagram, Facebook, TikTok, and LinkedIn.
- **Customizable tone:** Specify the tone that best suits your content.
- **REST API:** Generate captions through a Go-powered HTTP API.
- **Input validation:** Validate incoming requests and reject unsupported platforms.
- **Local caption storage:** Save generated captions in the browser using localStorage.
- **Mock generator:** Retain a mock implementation for development and testing.

## Technology Stack

| Component | Technology |
|---|---|
| Backend | Go |
| HTTP Router | Chi |
| AI Provider | Google Gemini |
| AI SDK | Google GenAI Go SDK |
| Frontend | HTML, CSS, JavaScript |
| Configuration | Environment variables |
| Browser Storage | localStorage |
| Testing | Go testing package |

## Project Structure

```text
CapIt/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── captions/
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   └── mock.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── providers/
│       └── gemini.go
│
├── web/
│   ├── assets/
│   │   ├── css/
│   │   └── js/
│   │
│   └── index.html
│
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

## Architecture Overview

CapIt follows a modular backend architecture.

- **cmd/server:** Application entry point, environment loading, dependency initialization, routing, and HTTP server startup.
- **internal/captions:** Caption request handling, input validation, response formatting, and the generator interface.
- **internal/config:** Application configuration and environment variable management.
- **internal/providers:** AI provider implementations, including the Gemini integration.
- **web:** Frontend assets and user interface.

The application uses dependency injection to separate caption generation from HTTP request handling.

The `Generator` interface allows different implementations, such as the mock generator and Gemini generator, to be used interchangeably.

## Getting Started

### Prerequisites

Make sure you have installed:

- Go
- Git
- A Google Gemini API key
- A modern web browser

### 1. Clone the repository

```bash
git clone https://github.com/nyae44/CapIt.git
cd CapIt
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Configure environment variables

Create a `.env` file in the project root.

```env
PORT=3000
APP_ENV=development

GEMINI_API_KEY=your_gemini_api_key
GEMINI_MODEL=your_supported_gemini_model
```

Replace the placeholders with your actual Gemini API credentials and a model available to your account.

**Security:** Never commit your `.env` file or expose your Gemini API key in frontend JavaScript.

### 4. Run the application

```bash
go run ./cmd/server
```

The application will be available at:

```text
http://localhost:3000
```

## API Documentation

### Health Check

Check whether the application is running.

**Endpoint**

```http
GET /api/v1/health
```

**Response**

```json
{
  "status": "ok"
}
```

### Generate Caption

Generate a caption using the configured AI provider.

**Endpoint**

```http
POST /api/v1/captions/generate
```

**Content-Type**

```text
application/json
```

**Request Body**

```json
{
  "platform": "instagram",
  "topic": "A beautiful coastal wedding in Kilifi",
  "tone": "professional"
}
```

**Supported Platforms**

- Instagram
- Facebook
- TikTok
- LinkedIn

**Response**

```json
{
  "caption": "Celebrate love, laughter, and unforgettable moments...",
  "platform": "instagram"
}
```

The generated caption will vary depending on the topic, platform, and selected tone.

### Error Handling

The API validates incoming requests and returns appropriate HTTP error responses.

Examples include:

- `400 Bad Request`: Invalid JSON, missing fields, or unsupported platform.
- `500 Internal Server Error`: Caption generation failure.

Detailed generation errors are logged on the backend rather than exposed to API clients.

## Testing

Run all project tests with:

```bash
go test ./...
```

The mock generator allows the caption handler to be tested without making actual requests to the Gemini API.

## AI Integration

CapIt uses Google's Gemini API through the official Go SDK.

The backend communicates with Gemini to generate captions based on user-provided content.

The AI provider is separated from the HTTP handler through the `Generator` interface, making it easier to introduce additional providers in the future.

AI services may occasionally experience temporary availability issues or rate limits.

## Development Roadmap

- [ ] Complete frontend integration with the Go caption API.
- [ ] Multiple caption variations per request.
- [ ] Improve error handling and retry mechanisms.
- [ ] Additional AI providers.
- [ ] User authentication and accounts.
- [ ] Cloud-based caption history.
- [ ] Caption editing and regeneration.
- [ ] Docker support.
- [ ] Deployment to production infrastructure.

## Contributing

Contributions, ideas, bug reports, and suggestions are welcome.

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Run the tests.
5. Submit a pull request.

## License

A license has not yet been specified.

---

**CapIt**

*Your content. The perfect caption.*
