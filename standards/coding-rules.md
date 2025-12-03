---
description: Coding standards for consistent, high-quality software development.
---

# Code Quality Rules

**Software Architecture Principles**
- Separate concerns into clear layers (presentation, application, domain, infrastructure)
- Apply SOLID principles and the 12‑factor app guidelines
- Keep it simple; choose straightforward, maintainable solutions (YAGNI)
- Prefer small, cohesive modules with explicit interfaces and boundaries
- Prioritize clarity and correctness over micro‑optimizations; optimize with evidence
- Favor functional core, imperative shell; keep side effects at the edges
- Design for idempotency and safe retries for external interactions
- MUST NOT duplicate business logic; extract shared behavior into utilities with clear ownership
- MUST NOT use json/bson/yaml/etc tags for objects in domain or application layer. Use DTO for serialization

**Coding Standards**
- Use dependency injection and composition over inheritance
- Define interfaces where consumed; depend on abstractions, not concrete implementations
- Prefer typed APIs/annotations and validate inputs at boundaries
- Use proven libraries for logging, configuration, and runtime essentials
- Write self‑documenting code; choose clear names; add comments for intent/“why”, not “what”
- Use Makefile/Taskfile for build, test, and common workflows
- Use named constants; MUST NOT use magic numbers or strings
- Enforce linters/formatters; keep consistent project style
- Handle errors explicitly; return actionable messages; avoid silent failures
- Secure coding: MUST NOT hardcode secrets; load via env/secret manager; least privilege for creds
- Add `TODO` comments for technical debt with context and plan to address
- MUST NOT add comments like `In a real implementation...`. Your implementation ALWAYS REAL!
- MUST NOT use identifiers as metric labels; use constant labels
- MUST NOT delete or change business logic without user confirmation or a plan item
- MUST NOT guess requirements; verify against codebase, docs, or approved sources
- MUST NOT implement backward compatibility, auth, or authorization unless requested
- MUST NOT return interfaces from constructors, use concrete types instead

**Testing Standards**
- Analyze existing tests and coverage before adding new ones
- Use fast, deterministic unit tests; isolate with mocks/fakes where appropriate
- Extract common setup/teardown; use factories/builders for complex data
- Not mark testing tasks complete until all tests pass
- Add property‑based and edge‑case tests when value is high
- Not add integration or API‑layer tests unless requested
- Test MUST create temporary files ONLY in temporary folders (e.g., /tmp, /var/tmp, etc.)
- MUST use web or tools like `perplexity` in case of uncertainty about testing best practices or stuck on a problem
- MUST NOT test trivial constructors/getters/setters
- MUST NOT modify tests to make implementation pass; tests define correct behavior
- MUST NOT delete or skip tests to ship; fix code or adjust plan
- MUST NOT write tests for metrics/logging unless explicitly requested
- MUST NOT write tests for interfaces itself; test implementations instead

**Debugging Standards**
- Reproduce issues and, when feasible, add a failing test before fixing
- Follow TDD mindset: code satisfies tests; do not bend tests to code
- MUST NOT assume code is correct because tests fail; investigate root cause
- MUST NOT modify tests to match broken behavior without prior analysis and sign‑off

**Mocking Guidelines**
- Confirm with user before introducing new mocks that change contract surface
- Mock external systems (DBs, APIs, FS), slow operations, and nondeterminism (time, random)
- Mock complex dependencies not under test; prefer thin seams
- ALWAYS use special mocking libraries (e.g., mockery, gomock) to generate mocks
- MUST NOT create mocks manually
- ALWAYS create mocks in place of use, not in a separate package
- ALWAYS use `//go:generate ...` for generating mocks
- MUST NOT mock logging; use a real logger

**Comment Standards**
- Keep comments short, precise, and purposeful
- Focus on intent and constraints
  * BAD: "New functionality"
  * GOOD: "Customer filtering"
- [MUST NOT] Document removed/non‑existent logic
  * BAD: "Order processing removed"
  * GOOD: Omit the comment

# Business Logic Disabling Protocol
Use when temporarily disabling functionality during current stage

Workflow:
1) Disable the logic behind a clear feature flag or guard
2) Add a TODO with re‑enable context (e.g., `// TODO: Re‑enable after step 5 in plan.md`)
3) Add a matching plan/todo item to track re‑enable (e.g., `Re‑enable customer filtering after step 5 of plan.md`)
