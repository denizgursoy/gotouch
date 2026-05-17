# Distributing Templates with Docker

The suggested way to use Gotouch across a team or organization is to build a **Docker image that bundles everything
together**: the Gotouch binary, your template projects, and your properties YAML. This gives you a single, versioned,
self-contained image that anyone can pull and run to create projects — no local setup required.

## Why Docker?

When you install Gotouch standalone, your properties YAML typically points to remote URLs for template projects.
This works, but it means:

- Template projects must be hosted somewhere accessible (GitHub, HTTP server)
- Network access is required every time a project is created
- Template versions are harder to pin — URLs can change

By building a Docker image with `localPath`, you solve all of these:

- **Self-contained** — templates, configuration, and the tool are all in one image
- **Versioned** — tag your image (`my-creator:1.0`, `my-creator:2.0`) and teams always get the right version
- **Offline** — no network needed at project creation time for `localPath` templates (templates using `url` still
  require network access)
- **Shareable** — push to any container registry (Docker Hub, GitHub Container Registry, private registry) and
  anyone with Docker can use it
- **Reproducible** — same image, same result, every time

## How It Works

1. Create your [template projects](./template-project) as local directories or point to remote git repositories
2. Write a [properties YAML](./properties-yaml) using `localPath` for bundled templates or `url` for remote ones
3. Build a Docker image based on `ghcr.io/denizgursoy/gotouch:latest` that copies in your templates and YAML
4. Distribute the image — your team runs one `docker run` command to create projects

## Full Example

This example includes three project templates that demonstrate every available feature in a properties YAML.
The first two use `localPath` to bundle templates inside the image, while the third uses `url` with
[git checkout](./project-structure) to clone a remote repository at project creation time — showing that you
can mix local and remote templates in the same properties YAML.

### Directory Structure

```
my-project-creator/
├── Dockerfile
├── properties.yaml
├── templates/
│   ├── rest-api/                  # Template as a directory
│   │   ├── main.go
│   │   ├── init.sh
│   │   ├── config/
│   │   │   └── config.go
│   │   └── handler/
│   │       └── handler.go
│   └── cli-app/                   # Template to be packaged as .tar.gz
│       ├── main.go
│       └── cmd/
│           └── root.go
│
│   # Web Application template is not stored locally —
│   # it is cloned from a git repository at runtime
│   # (see Project 3 in properties.yaml)
```

### Template Files

#### REST API Service

This template uses custom [delimiters](./project-structure) `<< >>` instead of the default `{{ }}`.
It includes an [init.sh](./init) script that runs after project creation.

::: code-group
```go [main.go]
package main

import (
	"fmt"
	"log"

	"<< .ModuleName >>/config"
	"<< .ModuleName >>/handler"
)

func main() {
	cfg := config.Load()
	router := handler.NewRouter(cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting << .ServiceName >> on %s", addr)
	log.Fatal(router.Start(addr))
}
```
```go [config/config.go]
package config

// Config holds the application configuration.
type Config struct {
	Port        string
	ServiceName string
	BaseURL     string
}

// Load returns the application configuration.
func Load() *Config {
	return &Config{
		Port:        "<< .Port >>",
		ServiceName: "<< .ServiceName >>",
		BaseURL:     "<< .BaseURL >>",
	}
}
```
```go [handler/handler.go]
package handler

import (
	"net/http"

	"<< .ModuleName >>/config"
)

// Router is the HTTP router.
type Router struct {
	cfg *config.Config
}

// NewRouter creates a new router.
func NewRouter(cfg *config.Config) *Router {
	return &Router{cfg: cfg}
}

// Start starts the HTTP server.
func (r *Router) Start(addr string) error {
	http.HandleFunc("<< .BaseURL >>/health", r.health)
	return http.ListenAndServe(addr, nil)
}

func (r *Router) health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
```
```sh [init.sh]
#!/bin/bash
echo "Setting up << .ServiceName >>..."

<< if .isDocker >>
echo "Building Docker image..."
docker build -t << .ServiceName >> .
<< end >>

<< if .isSwagger >>
echo "Generating Swagger docs..."
swag init
<< end >>

echo "Setup complete!"
```
:::

::: tip
The `init.sh` file is [templated](./init) with values before execution — you can use conditions and values inside it.
After execution, Gotouch deletes both `init.sh` and `init.bat` from the root folder.
:::

#### CLI Application

This template uses the default delimiters `{{ }}`.

::: code-group
```go [main.go]
package main

import (
	"fmt"
	"os"

	"{{.ModuleName}}/cmd"
)

var version = "{{.AppVersion}}"

func main() {
	if err := cmd.Execute(version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```
```go [cmd/root.go]
package cmd

import (
	"fmt"
)

// Execute runs the root command.
func Execute(version string) error {
	fmt.Println("{{.AppName}} version", version)
	fmt.Println("{{.Description}}")

	{{ if .configFormat }}
	fmt.Println("Config format:", "{{.configFormat}}")
	{{ end }}

	return nil
}
```
:::

### Properties YAML

This YAML demonstrates every available feature: [`localPath`](./project-structure),
[`initialModuleName`](./project-structure), [`delimiters`](./project-structure),
[`dependencies`](./dependencies), [`files`](./files) with inline content,
[`values` and `customValues`](./values), and all [question](./question) types
(yes/no, multiple choice, none of above, multiple select).

```yaml
# Project 1: REST API Service
# Uses localPath with a directory
- name: REST API Service
  reference: https://go.dev/
  initialModuleName: github.com/my-company/my-api
  localPath: /app/templates/rest-api
  language: go
  delimiters: "<< >>"
  dependencies:
    - go.uber.org/zap
  files:
    - content: |
        # << .ServiceName >>

        REST API service running on port << .Port >>.
      pathFromRoot: README.md
    - content: |
        bin/
        *.exe
        *.log
        .env
      pathFromRoot: .gitignore
  values:
    BaseURL: /api/v1
  customValues:
    Port: "8080"
    ServiceName: my-service
  questions:
    # Yes/No question — canSkip with one choice
    - direction: Do you want Dockerfile?
      canSkip: true
      choices:
        - choice: Yes
          values:
            isDocker: true
          files:
            - content: |
                FROM golang:1.25-alpine AS builder
                WORKDIR /app
                COPY . .
                RUN go build -o server .

                FROM alpine:latest
                WORKDIR /app
                COPY --from=builder /app/server .
                EXPOSE << .Port >>
                CMD ["./server"]
              pathFromRoot: Dockerfile

    # Multiple choice question — exactly one selection
    - direction: Which HTTP framework do you want to use?
      choices:
        - choice: Echo
          dependencies:
            - github.com/labstack/echo/v4
          values:
            framework: echo
        - choice: Gorilla Mux
          dependencies:
            - github.com/gorilla/mux
          values:
            framework: gorilla
        - choice: Gin
          dependencies:
            - github.com/gin-gonic/gin
          values:
            framework: gin

    # Multiple select question — zero or more selections
    - direction: Select additional features
      canSelectMultiple: true
      choices:
        - choice: Swagger Documentation
          dependencies:
            - github.com/swaggo/swag
            - github.com/swaggo/echo-swagger
          values:
            isSwagger: true
          files:
            - content: |
                // @title << .ServiceName >> API
                // @version 1.0
                // @description API documentation for << .ServiceName >>
                // @host localhost:<< .Port >>
                // @BasePath << .BaseURL >>
              pathFromRoot: docs/swagger.go
        - choice: Prometheus Metrics
          dependencies:
            - github.com/prometheus/client_golang
          values:
            isMetrics: true
          files:
            - content: |
                package metrics

                import "github.com/prometheus/client_golang/prometheus"

                var RequestCount = prometheus.NewCounterVec(
                    prometheus.CounterOpts{
                        Name: "<< .ServiceName >>_requests_total",
                        Help: "Total number of requests",
                    },
                    []string{"method", "path", "status"},
                )
              pathFromRoot: metrics/metrics.go
        - choice: JWT Authentication
          dependencies:
            - github.com/golang-jwt/jwt/v5
          values:
            isAuth: true
          files:
            - content: |
                package auth

                import "github.com/golang-jwt/jwt/v5"

                // Claims holds JWT claims.
                type Claims struct {
                    UserID string `json:"user_id"`
                    jwt.RegisteredClaims
                }
              pathFromRoot: auth/auth.go

# Project 2: CLI Application
# Uses localPath with a compressed file
- name: CLI Application
  reference: https://cobra.dev/
  initialModuleName: github.com/my-company/my-cli
  localPath: /app/templates/cli-app.tar.gz
  language: go
  dependencies:
    - github.com/spf13/cobra
  files:
    - content: |
        # {{.AppName}}

        {{.Description}}

        ## Installation

        ```bash
        go install {{.ModuleName}}@latest
        ```
      pathFromRoot: README.md
  values:
    AppVersion: 0.1.0
  customValues:
    AppName: mycli
    Description: A CLI application
  questions:
    # None of above question — canSkip with multiple choices
    - direction: Which configuration format do you want to use?
      canSkip: true
      choices:
        - choice: YAML
          dependencies:
            - gopkg.in/yaml.v3
          values:
            configFormat: yaml
          files:
            - content: |
                app:
                  name: "{{.AppName}}"
                  version: "{{.AppVersion}}"
              pathFromRoot: config.yaml
        - choice: TOML
          dependencies:
            - github.com/BurntSushi/toml
          values:
            configFormat: toml
          files:
            - content: |
                [app]
                name = "{{.AppName}}"
                version = "{{.AppVersion}}"
              pathFromRoot: config.toml
        - choice: JSON
          values:
            configFormat: json
          files:
            - content: |
                {
                  "app": {
                    "name": "{{.AppName}}",
                    "version": "{{.AppVersion}}"
                  }
                }
              pathFromRoot: config.json

    # Yes/No question
    - direction: Do you want shell completions?
      canSkip: true
      choices:
        - choice: Yes
          values:
            hasCompletions: true
          files:
            - content: |
                package cmd

                import (
                    "os"
                    "github.com/spf13/cobra"
                )

                func newCompletionCmd() *cobra.Command {
                    return &cobra.Command{
                        Use:   "completion [bash|zsh|fish]",
                        Short: "Generate shell completion scripts",
                        Args:  cobra.ExactArgs(1),
                        RunE: func(cmd *cobra.Command, args []string) error {
                            switch args[0] {
                            case "bash":
                                return cmd.Root().GenBashCompletion(os.Stdout)
                            case "zsh":
                                return cmd.Root().GenZshCompletion(os.Stdout)
                            case "fish":
                                return cmd.Root().GenFishCompletion(os.Stdout, true)
                            default:
                                return fmt.Errorf("unsupported shell: %s", args[0])
                            }
                        },
                    }
                }
              pathFromRoot: cmd/completion.go

# Project 3: Web Application
# Uses url with git checkout and branch
- name: Web Application
  reference: https://go.dev/
  initialModuleName: github.com/my-company/my-web
  url: https://github.com/denizgursoy/go-touch-projects.git
  branch: standard
  language: go
  dependencies:
    - github.com/labstack/echo/v4
  files:
    - content: |
        # {{.ProjectName}}

        Web application scaffolded with Gotouch.
      pathFromRoot: README.md
  values:
    BaseURL: /
  customValues:
    Port: "3000"
    AppTitle: My Web App
  questions:
    # Yes/No question
    - direction: Do you want a Docker Compose setup?
      canSkip: true
      choices:
        - choice: Yes
          values:
            isCompose: true
          files:
            - content: |
                version: "3.8"
                services:
                  app:
                    build: .
                    ports:
                      - "{{.Port}}:{{.Port}}"
                    environment:
                      - APP_TITLE={{.AppTitle}}
              pathFromRoot: docker-compose.yaml

    # Multiple choice question
    - direction: Which template engine do you want to use?
      choices:
        - choice: Go html/template
          values:
            templateEngine: html
        - choice: Templ
          dependencies:
            - github.com/a-h/templ
          values:
            templateEngine: templ
        - choice: None (JSON API only)
          values:
            templateEngine: none
```

### Dockerfile

The Dockerfile is simple — it extends the official Gotouch image and copies in your templates and configuration.
The Web Application project uses `url` with git checkout, so it doesn't need any files copied — Gotouch clones
the repository at runtime:

```dockerfile
FROM ghcr.io/denizgursoy/gotouch:latest

# Copy local templates
COPY templates/rest-api /app/templates/rest-api
COPY templates/cli-app.tar.gz /app/templates/cli-app.tar.gz

# Web Application template is cloned from git at runtime — nothing to copy

# Copy the properties YAML
COPY properties.yaml /app/properties.yaml

ENTRYPOINT ["gotouch", "-f", "/app/properties.yaml"]
```

## Build & Run

**1. Package the CLI Application template as `.tar.gz`:**

The REST API template is used as a directory directly, but the CLI Application template needs to be
compressed. Use the [package command](../commands#package) to create the archive:

```bash
gotouch package --source templates/cli-app --target templates/
```

This creates `templates/cli-app.tar.gz`.

**2. Build the Docker image:**

```bash
docker build -t my-project-creator .
```

Tag it with a version to keep things reproducible:

```bash
docker build -t my-project-creator:1.0 .
```

**3. Run the container:**

```bash
docker run -it -v $(pwd):/out --rm my-project-creator
```

Gotouch will prompt you to select a project, enter a module name, answer questions, and optionally
edit custom values. After all prompts complete, a progress view shows task execution logs until the
project is ready.

**4. Distribute:**

Push to a container registry so your team can use it:

```bash
docker tag my-project-creator:1.0 ghcr.io/my-org/project-creator:1.0
docker push ghcr.io/my-org/project-creator:1.0
```

Anyone on your team can now create projects with:

```bash
docker run -it -v $(pwd):/out --rm ghcr.io/my-org/project-creator:1.0
```
