# Otomo - Architectural Principles

This document outlines the core architectural principles of the Otomo project. All developers MUST adhere to these principles to maintain code quality, robustness, and consistency.

## 1. Clean Architecture

The project follows Clean Architecture principles with clear separation of concerns across multiple layers. The dependency direction must always point inwards (from Infrastructure/Presentation through Application to Domain).

```
┌─────────────────────────────────────────────────────────────┐
│                     Infrastructure Layer                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │  AWS Bedrock    │  │   Slack API     │  │   Storage    │ │
│  │   (AI Brain)    │  │                 │  │              │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                        Presentation Layer                    │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   HTTP Router   │  │   CLI Commands  │  │   Lambda     │ │
│  │   (Slack API)   │  │                 │  │   Handler    │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                       Application Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │    Use Cases    │  │    Services     │  │   Event      │ │
│  │                 │  │                 │  │   Handlers   │ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
                                │
┌─────────────────────────────────────────────────────────────┐
│                        Domain Layer                          │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │    Entities     │  │  Value Objects  │  │  Repositories│ │
│  │                 │  │                 │  │  (Interfaces)│ │
│  └─────────────────┘  └─────────────────┘  └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Core Components & Layer Responsibilities

#### Domain Layer
The core business logic and models. This layer is completely isolated from external frameworks, databases, or APIs.
- **Entities**: Objects with distinct identities (e.g., `Otomo`, `Thread`).
- **Value Objects**: Immutable data containers with no identity (e.g., `Message`, `Prompt`, `Reply`).
- **Repositories (Interfaces)**: Interfaces defining data access contracts, implemented in the Infrastructure layer.

- **Packages**:
  - `core`: System-wide core concepts (Prompt, Event, Message).
  - `reasoning`: Reasoning and AI model context (Brain, Context, Answer).
  - `chat`: Chat conversation context (Otomo, Thread, Reply).

#### Application Layer
Coordinates the flow of data and defines use cases. It depends only on the Domain Layer.
- **Use Cases**: Implements the business orchestrations (e.g., `ReplyToUser`, `ClassifySlackEventAndPublish`).
- **Services (Interfaces)**: Interfaces for system-level actions (e.g., `Messenger`).

#### Presentation Layer
Handles input delivery from external sources and formats output.
- **UI / HTTP**: Handles HTTP endpoints (chi router), Slack events, and CLI commands.

#### Infrastructure Layer
Implements interfaces defined in Domain and Application layers for external tools and services.
- **External Integrations**: AWS Bedrock (Claude), Slack API, event publishing.

---

## 2. Always Valid Domain Model

We enforce the **Always Valid Domain Model** pattern as our primary rule for domain layer design. A domain object (Entity or Value Object) must never exist in an invalid state. 

### Core Guidelines

1. **Strict Encapsulation (Private Fields)**
   All struct fields in the domain layer (`internal/domain/...`) MUST be private (unexported, starting with a lowercase letter). This prevents code outside the package from instantiating objects with zero values or invalid states.

2. **Error-Returning Factory Functions**
   Instantiating domain objects is only allowed through factory functions (constructors) prefixed with `New` (e.g., `NewMessage`, `NewThread`). These functions MUST validate their inputs and return an `error` if any business rules are violated.
   - Signature: `func NewXxx(...) (*Xxx, error)` or `func NewXxx(...) (Xxx, error)`

3. **Handcrafted Inline Validation**
   Validation logic must be written explicitly using simple Go code (such as `if` statements). Do not use external validation libraries (e.g., struct tags with `go-playground/validator`) in the domain models. The validation rules are part of the business logic and should be readable as pure Go code.

4. **Immutable Value Objects**
   Value Objects should generally be immutable. Do not provide setter methods. If a value needs to be updated, provide a method that returns a new instance of the Value Object, performing validation again.

5. **Mutator Validation**
   If an Entity has mutator methods that change its state (e.g., `AddMessage` on `Thread`), those methods must also enforce that the transition leads to a valid state.

6. **Avoid Primitive Obsession (Domain-Specific Types)**
   Do not use Go primitive types (like `string`, `int`, `bool`) directly for domain values or function signatures. Every business value MUST have its own defined type to make the domain model self-documenting and prevent bugs caused by accidentally swapping parameters of the same underlying type.

   ```go
   // Avoid:
   func SendMessage(channelID string, userID string)

   // Prefer:
   type ChannelID string
   type UserID string
   func SendMessage(channelID ChannelID, userID UserID)
   ```

### Code Examples

#### Before (Invalid state allowed)
```go
package core

// MESSAGE CAN BE INITIALIZED WITH INVALID VALUES (e.g., Role: "", Body: "")
type Message struct {
	Role MessageRole
	User string
	Body string
}
```

#### After (Always Valid)
```go
package core

import "fmt"

type Message struct {
	role MessageRole
	user string
	body string
}

func NewMessage(role MessageRole, user string, body string) (*Message, error) {
	if role != RoleSystem && role != RoleUser && role != RoleAssistant {
		return nil, fmt.Errorf("invalid message role: %s", role)
	}
	if body == "" {
		return nil, fmt.Errorf("message body cannot be empty")
	}
	return &Message{
		role: role,
		user: user,
		body: body,
	}, nil
}

func (m *Message) Role() MessageRole { return m.role }
func (m *Message) User() string      { return m.user }
func (m *Message) Body() string      { return m.body }
```

### Testing Guidelines

When writing tests, always use the factory functions. If you need helper fixtures for tests, create dedicated test helpers that return valid models, rather than exposing fields or bypassing constructors.
