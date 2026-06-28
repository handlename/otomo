## Architecture
- Adopt a Layered Architecture divided into 4 layers: Domain layer/Application layer/Presentation layer/Infrastructure layer.
- Enforce the **Always Valid Domain Model** principle for all domain layer designs (strict encapsulation, validation in constructors).
- Refer to [ARCHITECTURE.md](./ARCHITECTURE.md) for core architectural principles and patterns.
- Refer to [DESIGN.md](./DESIGN.md) for specific configuration.

## Coding
- Follow general Go best practices
- Define appropriate types as needed and avoid using primitive types directly as much as possible
- Prioritize clarity of processing over technical approaches

## Comments
- Write comments in English
- Write comments only when structures or operations are complex
- Write git commit messages, GitHub issue descriptions, and Pull Request summaries in English

## Testing
- Use the command `go test -v {package path} {test function name}` when running tests
- Write test cases in table-based test format. However, consider splitting test functions when conditional branching is needed for each test case
- Write only necessary and sufficient test cases
- Actively use `github.com/stretchr/testify/require` and `github.com/stretchr/testify/assert` for value comparisons in test code

## Glossary
- Refer to [GLOSSARY.md](GLOSSARY.md) for definitions of key terms, domain concepts, and technical terminology used in this project

## Documentation
- Update or modify related documentation (e.g., [ARCHITECTURE.md](./ARCHITECTURE.md), [DESIGN.md](./DESIGN.md), [GLOSSARY.md](./GLOSSARY.md)) whenever features in `otomo` are added or changed.

## Tracing
- When adding or modifying operations in the Presentation, Application, or Infrastructure layers that involve network calls, latency-sensitive actions, or significant use case steps, instrument them using OpenTelemetry spans.
- Use `otel.Tracer("otomo").Start(ctx, "Span Name")` to create and context-propagate a span, and defer its closure via `defer span.End()`.
- Use the central helper `trace.RecordError(span, err)` from `github.com/handlename/otomo/internal/infra/trace` on error return paths to record errors. Do not use named return variables solely for deferred span logging.
- **Strictly prohibit** any OpenTelemetry imports or tracing dependencies inside the Domain layer (`internal/domain/...`). The Domain layer must remain pure.

