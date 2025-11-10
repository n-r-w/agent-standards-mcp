// Package prompt provides system and tool prompts for the agent-standards-mcp server.
package prompt

import _ "embed"

//go:embed get-standards-prompt.txt
var getStandardsPrompt []byte

//go:embed list-standards-prompt.txt
var listStandardsPrompt []byte

//go:embed system-prompt.txt
var systemPrompt []byte

//go:embed load-relevant-standards-prompt.txt
var loadRelevantStandardsPrompt []byte

//go:embed follow-standards-prompt.txt
var followStandardsPrompt []byte

// SystemPrompt returns the system prompt as a string.
func SystemPrompt() string {
	return string(systemPrompt)
}

// GetStandardsPrompt returns the get standards prompt as a string.
func GetStandardsPrompt() string {
	return string(getStandardsPrompt)
}

// ListStandardsPrompt returns the list standards prompt as a string.
func ListStandardsPrompt() string {
	return string(listStandardsPrompt)
}

// LoadRelevantStandardsPrompt returns the load relevant standards prompt as a string.
func LoadRelevantStandardsPrompt() string {
	return string(loadRelevantStandardsPrompt)
}

// FollowStandardsPrompt returns the follow standards prompt as a string.
func FollowStandardsPrompt() string {
	return string(followStandardsPrompt)
}
