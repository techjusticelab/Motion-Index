# Repository Guidelines

## Project Structure & Module Organization
Motion-Index spans three active workspaces. `motion-index-fiber/` hosts the Go Fiber API, with `cmd/` entrypoints, `internal/` service logic, `pkg/` reusable clients, and `test/` integration fixtures. `Web/` contains the SvelteKit frontend; UI code lives in `src/lib` and `src/routes`, static assets in `static/`, and shared config in `svelte.config.js`. `Fiber/` provides architecture briefs and should be updated when major flows change. Keep generated binaries in `motion-index-fiber/bin/` and never commit local `Web/node_modules/`.

## Build, Test, and Development Commands
Backend:
- `cd motion-index-fiber && go run cmd/server/main.go` starts the API with the local `.env`.
- `go test ./...` runs unit and integration tests; append `-run` to target suites.
- `go build -o bin/server cmd/server/main.go` produces the deployable binary.
Frontend:
- `cd Web && npm install` ensures dependencies align with `package.json`.
- `npm run dev` launches Vite at http://localhost:5173.
- `npm run check` runs Svelte type analysis; `npm run lint` enforces formatting; `npm run build` emits the production bundle.

## Coding Style & Naming Conventions
Run `gofmt` (or `goimports`) before committing Go code; stick to idiomatic camelCase identifiers and package-scoped constants in ALL_CAPS only when immutable. HTTP handlers live under `internal/http` and should use descriptive Fiber route names. Svelte components are PascalCase files under `src/lib`, route modules stay lowercase with `+page.svelte` naming. Use 2-space indentation, Prettier formatting (`npm run format`), and keep Tailwind classes sorted by Prettier’s plugin.

## Testing Guidelines
Place Go tests alongside the code with `_test.go` suffixes; broader scenarios belong in `motion-index-fiber/test/` and should mock external services before hitting real DO resources. Maintain assertions for search relevance and storage fallbacks. For Svelte, validate new UI logic with `npm run check` and extend Playwright or Vitest suites if introduced; document any manual QA steps in PRs.

## Commit & Pull Request Guidelines
Write concise, present-tense commit titles (e.g., “tune search faceting”); group related changes rather than omnibus commits. PRs must include context, screenshots for UI tweaks, test evidence (`go test` / `npm run check` logs), and call out migrations or environment variables. Link GitHub issues when available and update relevant docs (`Fiber/*.md`, `README.md`) alongside functional changes.

## Environment & Security Notes
Copy provided templates (`motion-index-fiber/.env.example`, `Web/template.env`) and never commit filled secrets. When adjusting DigitalOcean or Supabase keys, rotate credentials in cloud consoles and record the change in deployment docs. Verify Elasticsearch or Spaces configuration locally with the scripts under `motion-index-fiber/deployments/` before merging.
