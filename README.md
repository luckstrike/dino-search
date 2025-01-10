# Search Engine

A terminal-based search engine written in Go, with future web API support.

## Project Structure

\`\`\`
search-engine/
├── cmd/
│   ├── cli/                    # Command-line interface
│   └── api/                   # Future API server
├── internal/
│   ├── search/                # Core search functionality
│   ├── crawler/               # Web crawling functionality
│   ├── storage/               # Data persistence
│   └── common/                # Shared utilities
├── pkg/                       # Public packages
├── api/                       # API documentation
├── configs/                   # Configuration files
├── deployments/               # Deployment configs
├── docs/                      # Documentation
├── scripts/                   # Build scripts
└── test/                     # Integration tests
\`\`\`

## Getting Started

1. Build the project:
   \`\`\`bash
   make build
   \`\`\`

2. Run the search engine:
   \`\`\`bash
   make run
   \`\`\`

## Development

- \`make dev\`: Build and run the application
- \`make clean\`: Clean build artifacts

## License

TODO
