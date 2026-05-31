# Overengineering as a Service (OaaS)

An overengineered web application that always ships as a single container.

## Roadmap

- [ ] Some sort of web server
- [ ] SPA UI
- [ ] More things

## Development

### Prerequisites

Development requires the following tools:

- Docker
- Go (plus `delve` for debugging)
- Git
- Node with `npm`
- Make

All of these prerequisites except for Docker are included if you use VS Code and
open the project with the included dev container.

### Commands

Common development commands:

| Command      | Purpose                                                                       |
|--------------|-------------------------------------------------------------------------------|
| `make build` | Package the application into a container image                                |
| `make dev`   | Start separate development servers for the UI and backend with live-reloading |

All available commands can be found in the `Makefile`.

### URLs

| Target                      | URL                   |
|-----------------------------|-----------------------|
| Application                 | http://localhost:8000 |
| Standalone UI (development) | http://localhost:5173 |

