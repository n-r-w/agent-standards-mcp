// Package domain contains core business entities without any external dependencies.
package domain

// StandardInfo represents basic information about a standard.
// This is a pure domain entity without any serialization tags.
type StandardInfo struct {
	Name        string
	Description string
}

// Standard represents the full content of a standard.
// This is a pure domain entity without any serialization tags.
type Standard struct {
	Name        string
	Description string
	Content     string
}
