// Package standards provides functionality for loading, parsing, and validating standard files
// from the file system. It implements the StandardLoader interface and handles markdown
// frontmatter parsing, file validation, and security checks.
package standards

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

// frontmatterData represents the YAML frontmatter structure we expect
type frontmatterData struct {
	Description string `yaml:"description"`
}

const (
	// minimumFrontmatterLines is the minimum number of lines required for valid frontmatter
	minimumFrontmatterLines = 3
	// oneMB is the default maximum standard file size in bytes
	oneMB = 1024 * 1024
	// defaultMaxStandards is the default maximum number of standard files
	defaultMaxStandards = 100
)

// parseFrontmatter parses markdown content with optional YAML frontmatter.
// It extracts the description field from frontmatter and returns the description
// and content separately. If no frontmatter is present, description will be empty.
func parseFrontmatter(content string) (description string, parsedContent string, err error) {
	// Handle empty content
	if content == "" {
		return "", "", nil
	}

	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		// No frontmatter, return content as-is with empty description
		return "", content, nil
	}

	// Find the end of frontmatter
	lines := strings.Split(content, "\n")
	if len(lines) < minimumFrontmatterLines {
		// Not enough lines for proper frontmatter
		return "", content, nil
	}

	// Find the closing delimiter
	endIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		// No closing delimiter found, treat as no frontmatter
		return "", content, nil
	}

	// Extract frontmatter content
	frontmatterLines := lines[1:endIndex]
	frontmatterText := strings.Join(frontmatterLines, "\n")

	// Parse YAML frontmatter
	var fm frontmatterData
	err = yaml.Unmarshal([]byte(frontmatterText), &fm)
	if err != nil {
		return "", "", err
	}

	fm.Description = strings.TrimSpace(fm.Description)

	// Extract content after frontmatter
	var contentLines []string
	if endIndex+1 < len(lines) {
		contentLines = lines[endIndex+1:]
	}
	parsedContent = strings.TrimSpace(strings.Join(contentLines, "\n"))

	if fm.Description == "" {
		return "", "", errors.New("frontmatter 'description' cannot be empty")
	}

	if parsedContent == "" {
		return "", "", errors.New("standard content cannot be empty")
	}

	return fm.Description, parsedContent, nil
}
