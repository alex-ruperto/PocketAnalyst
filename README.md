# GoSimpleREST

GoSimpleREST is exactly as the title says; a simple REST API written in Go.

## Architecture

GoSimpleREST uses MVC (Model-View-Controller) architecture.

### Project Structure

- README.md ------------- Project Documentation
- api/
  - app.go -------------- Application Setup
  - handlers.go --------- Request Handlers (CONTROLLER)
  - responses.go -------- Response Utilities (VIEW)
- go.mod ---------------- Go Module file
- main.go --------------- Entry Point
- models/
  - item.go ------------- Data Model (MODEL)
- store/
  - memstore.go --------- In-memory data store
