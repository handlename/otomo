## Architecture
- Adopt a Layered Architecture divided into 4 layers: Domain layer/Application layer/Presentation layer/Infrastructure layer
- The Domain layer must not depend on any other layer
- The Application layer depends only on the Domain layer
- The Presentation layer depends on the Domain layer and Application layer
- Other layers must not depend on the Infrastructure layer
- Define interfaces when necessary to invert dependencies
- Refer to DESIGN.md for specific configuration

## Coding
- Follow general Go best practices
- Define appropriate types as needed and avoid using primitive types directly as much as possible
- Prioritize clarity of processing over technical approaches

## Comments
- Write comments in English
- Write comments only when structures or operations are complex

## Testing
- Use the command `go test -v {package path} {test function name}` when running tests
- Write test cases in table-based test format. However, consider splitting test functions when conditional branching is needed for each test case
- Write only necessary and sufficient test cases
- Actively use `github.com/stretchr/testify/require` and `github.com/stretchr/testify/assert` for value comparisons in test code

## Glossary
- Refer to [GLOSSARY.md](GLOSSARY.md) for definitions of key terms, domain concepts, and technical terminology used in this project

