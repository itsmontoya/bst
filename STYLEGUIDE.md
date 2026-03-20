# Go Style Guide

This document defines the coding standards for this repository.

The goal is clarity, consistency, and long-term maintainability.

We optimize for:

* Explicitness over cleverness
* Small, focused functions
* Predictable structure
* Easy code review

# Core Principles

* Write code that is obvious to the next engineer.
* Keep functions small and single-purpose.
* Prefer explicit declarations over inference.
* Avoid magic and hidden control flow.
* Structure the repository so it is easy to navigate.

# Project Structure

## File Per Type

Each primary type gets its own file.

All methods belonging to a type must live in the same file as the type definition.

## Constructor Placement

Place constructor/initializer functions directly above the type they construct.

### Why

* Keeps navigation predictable.
* Makes type entry points obvious during review.

### Preferred

```go
func NewService(store Store) (s *Service) { ... }

type Service struct { ... }
```

### Avoid

```go
type Service struct { ... }

func NewService(store Store) (s *Service) { ... }
```

### Why

* Makes navigation trivial.
* Keeps related behavior grouped together.
* Prevents method scattering across the codebase.
* Makes refactoring safer.

### Example

service.go
client.go
parser.go
store.go

### Preferred

```go
// service.go

// NewService constructs a Service with the provided Store.
func NewService(store Store) (s *Service) {
    s = &Service{
        store: store,
    }
    return s
}

// Service coordinates request handling and persistence.
type Service struct {
    store Store
}

// Handle processes the incoming request.
func (s *Service) Handle(...) (...) {
    ...
}

func (s *Service) validate(...) error {
    ...
}
```

### Avoid

* Defining a type in one file and spreading its methods across multiple files.
* Creating service_helpers.go, service_utils.go, etc. for methods of the same type.

If behavior grows too large, extract a new type instead of splitting methods across files.

# Formatting

* Always run gofmt.
* Follow standard Go conventions.
* Keep vertical noise low.
* Avoid unnecessary blank lines.

# Naming

* Use descriptive, intention-revealing names.
* Avoid unnecessary abbreviations.
* Use verbs for functions (Build, Parse, Fetch).
* Use nouns for types (Parser, Client, Store).

Avoid stutter:

* Preferred: type Client struct{}
* Avoid: type MyProjectClient struct{}

# Variable Declarations

## Prefer var over := (with narrow exceptions)

We prefer explicit variable declarations.

### Why

* Clear scope and type.
* Reduces accidental shadowing.
* Encourages intentional initialization.

### Preferred

```go
var (
    count int
    user  User
    ids   []string
    err   error
)

timeout := 5 * time.Second
```

Use `:=` only when it meaningfully improves clarity in tight scopes and the declaration is simple.

### Use := for simple local declarations

`:=` is acceptable when all of the following are true:

* The declaration is local and close to first use.
* The right-hand side is short and obvious.
* The variable has no ambiguity in meaning or type.
* The statement does not risk shadowing an existing variable.

Preferred examples:

```go
b := newMemory()
done := make(chan struct{}, 1)
if v, ok := m[key]; ok {
    return v, nil
}
```

### Use var for shared, reused, or stateful values

Use `var` when the variable is reused across branches, needs a zero-value declaration, or benefits from explicit grouping.

Preferred examples:

```go
var (
    b   backend
    err error
)

if b, err = newFile(filepath); err != nil {
    return nil, err
}
```

### Acceptable

```go
if v, ok := m[key]; ok {
    return v, nil
}
```

## Avoid Shadowing

Do not redeclare variables in nested scopes.

Avoid:

```go
var err error
err = doThing()

if !cond {
    return nil
}

if err = doOtherThing(); err != nil {
    return err
}
```

# Named Returns

## Named Returns Are Encouraged

Named return values improve readability and self-documentation.

Use them especially for:

* Multiple return values
* Public APIs
* Non-obvious outputs

Preferred:

```go
func Parse(input string) (result Result, err error) {
    ...
    return result, err
}
```

## No Naked Returns

Even with named returns, never use naked return.

Scope note:

* This rule applies to functions with named return values.
* Early bare `return` is acceptable in `func(...){}` blocks that return no values (common in tests).

### Why

* Control flow remains explicit.
* Safer during refactors.
* Easier to review.

Avoid:

```go
func Parse(input string) (result Result, err error) {
    if input == "" {
        err = errors.New("empty")
        return
    }
    return
}
```

Prefer:

```go
func Parse(input string) (result Result, err error) {
    if input == "" {
        err = errors.New("empty")
        return result, err
    }

    return result, err
}
```

# Function Design

## Keep Functions Small

Guidelines:

* Should fit comfortably on one screen.
* Should do one thing.
* If it has multiple responsibilities, extract helpers.

## Use Early Returns

Reduce nesting and keep control flow flat.

```go
if err != nil {
    return resp, err
}
```

## Prefer Composition Over Large Methods

If a method grows too large:

* Extract smaller private methods.
* Or extract a new type.

Do not split methods across files. Extract types instead.

# Error Handling

* Handle errors immediately.
* For external/operation errors (IO, filesystem, network, system calls), context is required.
* Errors from calls like `os.Open`, `os.OpenFile`, `Read`, `Write`, and `Close` must be wrapped with context via `%w` when returned.
* Return sentinel contract/state errors directly when they are already descriptive.

```go
if err != nil {
    return out, fmt.Errorf("load user %q: %w", id, err)
}
```

Avoid swallowing errors.

Sentinel example:

```go
if closed {
    return ErrIsClosed
}
```

# Interfaces

* Accept interfaces.
* Return concrete types.
* Keep interfaces small and behavior-focused.
* Define interfaces near where they are used.

# Receivers

* Use pointer receivers for large structs or when mutating state.
* Keep receiver names short (s, c, p).
* Never use this.

# Comments

* All exported types, functions, methods, and interfaces must have GoDoc comments.
* Comments must begin with the name of the exported symbol.
* Explain why, not what.
* Keep comments concise. Err on the side of brevity.
* Avoid restating obvious code.
* Private functions may be commented consistently for navigation and readability.
* Private-function comments must stay concise and non-redundant.
* Comments should rarely exceed a few short lines unless documenting complex behavior.

### GoDoc Example

```go
// Service coordinates request handling and persistence.
type Service struct {
    store Store
}

// NewService constructs a Service with the provided Store.
func NewService(store Store) (s *Service) {
    s = &Service{
        store: store,
    }
    return s
}
```

### Avoid

```go
// Handle handles the request.
func (s *Service) Handle(...) { ... }

// i increments by 1.
i++
```

# Tests

* Prefer table-driven tests.
* Keep setup explicit.
* Avoid clever test abstractions.

```go
func Test_Buffer_Write(t *testing.T) {
    type testcase struct {
        name string
        // Additional needed fields here.
    }

    testInput := []byte("hello world")
    tests := []testcase{
        // Test cases go here.
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic goes here.
        })
    }
}
```

# PR Checklist

Before opening a PR:

* [ ] gofmt applied
* [ ] One file per primary type
* [ ] All methods for a type live in the same file
* [ ] Functions are small and focused
* [ ] Named returns used appropriately
* [ ] No naked returns
* [ ] Prefer var over :=
* [ ] No shadowing
* [ ] External/operation errors include context
* [ ] Sentinel/state errors remain direct
* [ ] Exported symbols have GoDoc comments
* [ ] Tests cover behavior
