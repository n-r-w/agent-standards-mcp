---
description: Golang coding guidelines and best practices.
---

# Golang Guidelines

🚨🚨 **CRITICAL INTERFACE USAGE RULES** 🚨🚨
- MUST define interfaces in consumer packages, not provider packages to enable reverse dependencies and avoid circular imports. E.g.:
BAD (defining interface in provider package):
    `internal/service/service.go` - implements `SomeInterface`
    `internal/service/interface.go` - defines `SomeInterface`
GOOD (defining interface in consumer package):
    `internal/service/service.go` - implements `SomeInterface`
    `internal/consumer/consumer.go` - defines `SomeInterface` and uses it

**General**
- Reuse buffers in pools with sufficient capacity; avoid returning same buffer twice
- Protect concurrent field access with atomics or mutexes covering shared state
- Avoid atomics if mutex already protects data
- Use sync.Pool pointers with slices to avoid value copying
- Limit sync.Pool to short-lived data; MUST NOT return pooled items to external callers
- Use range loops over integers: `for i := range n { ... }` instead of `for i := 0; i < n; i++ { ... }`
- Replace magic numbers with named constants: `const maxHotels = 1000; if c.MaxHotels > maxHotels {`
- Keep lines under 120 chars; break long function signatures with proper indentation
- Remove empty blocks; handle conditions properly or omit them
- Limit function returns to 3-4 values; use structs for complex returns
- Prefer typified variables over `any`
- Use `any` instead of `interface{}`
- SHOULD NOT Accept parameters suppression like `func(_ param)` if it's not really needed
- MUST NOT modify `.gitignore` without explicit approval

**Dependency Management**
For adding new dependencies:
1) Add import in code
2) Run `go mod tidy`
3) Verify `go.mod` and `go.sum` changes

🚨🚨 MUST NOT manually edit `go.mod` or `go.sum` 🚨🚨

**Linter**
- Use `//exhaustruct:enforce` to verify all struct fields are set
- External SDK types can use `//nolint:exhaustruct` (optional fields use zero values)
- Add reason to suppress linter warnings: `//nolint:gosec // brief reason`
- Suppress int64-to-int conversion warnings with `nolint:gosec` only when justified
- MUST NOT suppress linter warnings without explicit approval
- `G304: Potential file inclusion via variable (gosec)` -> use `filepath.Clean`

**Naming**
- Use uppercase 'ID' for identifiers: HotelIDs []int64
- Avoid similar names for public/private functions; use distinct terms like NewMetrics()/createMetrics()
- Omit package name repetition: struct `Factory` in `cache` package, not `CacheFactory`
- Use '_' or omit name for unused receivers (e.g. `func (_ *GzipCompressor) Method() {}`)
- Use lowercase folders/package names without underscores
- Keep filenames short; use underscores for very long names: `some.go` or `some_very_long_file_name.go`

**Errors**
- Check all errors, including Close(): `if err := writer.Close(); err != nil { handle }`
- Check bounds before casting to prevent overflow: `maxTries := uint(1); if config.MaxRetries > 0 { maxTries = uint(config.MaxRetries) }`

**Interfaces**
- Verify interface implementation: `var _ somepkg.SomeInterface = (*SomeService)(nil)` with comment `// SomeService implements somepkg.SomeInterface`

**Context**
- Use `context.Background()` ONLY in main(), init(), and tests
- MUST NOT store context in a struct, ALWAYS pass context as a parameter
- If you need a non-cancellable context, use context.WithoutCancel() at the call site

**Testing**
- Avoid verifying APIs by executing services; use integration tests where appropriate
- Use specific testify assertions: `assert.Positive(t, value)` instead of `assert.Greater(t, value, 0)`
- Use epsilon tolerance for float comparisons: `assert.InEpsilon(t, 1.0, ratio, 0.01)`
- Use require only in main test goroutine; use assert in concurrent goroutines
- Use require for errors that should fail test: `require.Error(t, err)`
- Use assert for value comparisons: `assert.Equal(t, expected, actual)`
- Put unit tests in the same package as the code
- MUST NOT use `ctrl.Finish` in tests, it executed automatically by gomock

**Performance**
- Pre-allocate slices with known size: `keys := make([]string, 0, expectedSize)`
- Prefer `fmt.Fprintf(builder, "%.2f", value)` over `WriteString` with `fmt.Sprintf` for string building
- Prefer values over pointers for small structs to reduce allocations

**Security**
- Validate all user input; sanitize before use
- Check vulnerabilities (govulncheck, etc) before task completion
- Avoid `ioutil.ReadAll` for large requests; use streaming with `io.Copy`

**Docs**
- Comment exported constants with purpose: `// VatRateUnspecified is unspecified VAT rate. const VatRateUnspecified VatRate = 0`
- Avoid variable shadowing; use unique names in inner scopes: `if bytes, ok := value.([]byte); ok { ... }`

**Brace fixing**
- Read full block (function/method/scope) to understand structure when fixing brace errors
- Creating temporary examples: always run single file, absolute path (e.g. `go run /tmp/test.go`)
