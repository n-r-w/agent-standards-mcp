package standards

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantDesc    string
		wantContent string
		wantErr     bool
	}{
		{
			name: "valid frontmatter with description",
			content: `---
description: "A test standard for doing X"
---
This is the standard content.
It can have multiple lines.`,
			wantDesc:    "A test standard for doing X",
			wantContent: "This is the standard content.\nIt can have multiple lines.",
			wantErr:     false,
		},
		{
			name: "no frontmatter",
			content: `This is just content without frontmatter.
No YAML header here.`,
			wantDesc:    "",
			wantContent: "This is just content without frontmatter.\nNo YAML header here.",
			wantErr:     false, // No frontmatter means no validation
		},
		{
			name: "empty frontmatter",
			content: `---
---
Just content after empty frontmatter.`,
			wantDesc:    "",
			wantContent: "",   // Will be empty due to validation error
			wantErr:     true, // Empty description causes validation error
		},
		{
			name: "frontmatter without description",
			content: `---
other: "value"
---
Content here.`,
			wantDesc:    "",
			wantContent: "",   // Will be empty due to validation error
			wantErr:     true, // Empty description causes validation error
		},
		{
			name: "malformed frontmatter",
			content: `---
description: "unclosed quote
---
Content`,
			wantDesc:    "",
			wantContent: "",
			wantErr:     true,
		},
		{
			name:        "empty file",
			content:     "",
			wantDesc:    "",
			wantContent: "",
			wantErr:     false, // Empty file with no frontmatter is allowed
		},
		{
			name: "only frontmatter",
			content: `---
description: "Only description"
---`,
			wantDesc:    "",
			wantContent: "",   // Will be empty due to validation error
			wantErr:     true, // Empty content causes validation error
		},
		{
			name: "description with special characters",
			content: `---
description: "Standard with apostrophes 'and' \"quotes\" & symbols"
---
Content with symbols`,
			wantDesc:    "Standard with apostrophes 'and' \"quotes\" & symbols",
			wantContent: "Content with symbols",
			wantErr:     false,
		},
		{
			name: "multiline description",
			content: `---
description: |
  This is a multiline description
  with multiple lines
  and indentation
---
Content here`,
			wantDesc:    "This is a multiline description\nwith multiple lines\nand indentation",
			wantContent: "Content here",
			wantErr:     false,
		},
		{
			name: "whitespace only description",
			content: `---
description: "   \n\t   "
---
Valid content`,
			wantDesc:    "",
			wantContent: "",   // Will be empty due to validation error
			wantErr:     true, // Should fail after trimming description
		},
		{
			name: "whitespace only content",
			content: `---
description: "Valid description"
---
		 	 
		 	 
		 `,
			wantDesc:    "",
			wantContent: "",   // Will be empty due to validation error
			wantErr:     true, // Should fail after trimming content
		},
		{
			name: "description and content with surrounding whitespace",
			content: `---
description: "   Valid description with spaces   "
---
		 	 
		  Valid content with spaces
		 	 
`,
			wantDesc:    "Valid description with spaces",
			wantContent: "Valid content with spaces",
			wantErr:     false, // Should succeed after trimming
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDesc, gotContent, err := parseFrontmatter(tt.content)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotDesc != tt.wantDesc {
				t.Errorf("ParseFrontmatter() gotDesc = %q, wantDesc %q", gotDesc, tt.wantDesc)
			}

			if gotContent != tt.wantContent {
				t.Errorf("ParseFrontmatter() gotContent = %q, wantContent %q", gotContent, tt.wantContent)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Set up environment variables for testing
	originalMaxStandards, hasMaxStandards := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARDS")
	originalMaxStandardSize, hasMaxStandardSize := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE")
	defer func() {
		if hasMaxStandards {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", originalMaxStandards); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARDS"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		}
		if hasMaxStandardSize {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", originalMaxStandardSize); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		}
	}()

	// Set test values
	if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", "10"); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
	}
	if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", "1024"); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
	}

	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid file within size limit",
			setup: func() string {
				path := filepath.Join(tempDir, "valid.md")
				content := "This is a valid standard file with acceptable content."
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "file too large",
			setup: func() string {
				path := filepath.Join(tempDir, "large.md")
				// Create content larger than 1024 bytes
				content := string(make([]byte, 2000))
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			wantErr: true,
			errMsg:  "file size exceeds maximum limit",
		},
		{
			name: "path traversal attack - relative path",
			setup: func() string {
				return "../../../etc/passwd"
			},
			wantErr: true,
			errMsg:  "path traversal detected",
		},
		{
			name: "path traversal attack - absolute path outside allowed",
			setup: func() string {
				return "/etc/passwd"
			},
			wantErr: true,
			errMsg:  "path traversal detected",
		},
		{
			name: "file does not exist",
			setup: func() string {
				return filepath.Join(tempDir, "nonexistent.md")
			},
			wantErr: true,
			errMsg:  "file does not exist",
		},
		{
			name: "directory instead of file",
			setup: func() string {
				path := filepath.Join(tempDir, "not_a_file")
				if err := os.Mkdir(path, 0755); err != nil {
					t.Fatalf("Failed to create test directory: %v", err)
				}
				return path
			},
			wantErr: true,
			errMsg:  "path is not a file",
		},
		{
			name: "valid file with markdown extension",
			setup: func() string {
				path := filepath.Join(tempDir, "standard.md")
				content := "---\ndescription: test\n---\ncontent"
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			wantErr: false,
			errMsg:  "",
		},
		{
			name: "valid file with non-markdown extension",
			setup: func() string {
				path := filepath.Join(tempDir, "standard.txt")
				content := "Just a text file"
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return path
			},
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPath := tt.setup()

			err := validateFile(testPath, tempDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateFile() error = %v, expected to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateStandardFiles(t *testing.T) {
	tempDir := t.TempDir()

	originalMaxStandards, hasMaxStandards := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARDS")
	defer func() {
		if hasMaxStandards {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", originalMaxStandards); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARDS"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		}
	}()

	tests := []struct {
		name         string
		setup        func() []string
		maxStandards string
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid number of files",
			setup: func() []string {
				var paths []string
				for i := 0; i < 3; i++ {
					path := filepath.Join(tempDir, fmt.Sprintf("standard%d.md", i))
					content := "---\ndescription: test\n---\ncontent"
					if err := os.WriteFile(path, []byte(content), 0644); err != nil {
						t.Fatalf("Failed to write test file: %v", err)
					}
					paths = append(paths, path)
				}
				return paths
			},
			maxStandards: "5",
			wantErr:      false,
			errMsg:       "",
		},
		{
			name: "too many files",
			setup: func() []string {
				var paths []string
				for i := 0; i < 8; i++ {
					path := filepath.Join(tempDir, fmt.Sprintf("standard%d.md", i))
					content := "---\ndescription: test\n---\ncontent"
					if err := os.WriteFile(path, []byte(content), 0644); err != nil {
						t.Fatalf("Failed to write test file: %v", err)
					}
					paths = append(paths, path)
				}
				return paths
			},
			maxStandards: "5",
			wantErr:      true,
			errMsg:       "number of files exceeds maximum limit",
		},
		{
			name: "empty file list",
			setup: func() []string {
				return []string{}
			},
			maxStandards: "5",
			wantErr:      false,
			errMsg:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", tt.maxStandards); err != nil {
				t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
			paths := tt.setup()

			err := validateStandardFiles(paths, tempDir)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStandardFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateStandardFiles() error = %v, expected to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 1; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

func TestFileStandardLoader_ListStandards(t *testing.T) {
	tempDir := t.TempDir()

	// Set up environment variables
	originalFolder, hasFolder := os.LookupEnv("AGENT_STANDARDS_MCP_FOLDER")
	originalMaxStandards, hasMaxStandards := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARDS")
	originalMaxStandardSize, hasMaxStandardSize := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE")
	defer func() {
		if hasFolder {
			if err := os.Setenv("AGENT_STANDARDS_MCP_FOLDER", originalFolder); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_FOLDER: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_FOLDER"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_FOLDER: %v", err)
			}
		}
		if hasMaxStandards {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", originalMaxStandards); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARDS"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
			}
		}
		if hasMaxStandardSize {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", originalMaxStandardSize); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		}
	}()

	if err := os.Setenv("AGENT_STANDARDS_MCP_FOLDER", tempDir); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_FOLDER: %v", err)
	}
	if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARDS", "10"); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARDS: %v", err)
	}
	if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", "1024"); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
	}

	tests := []struct {
		name        string
		setup       func()
		wantErr     bool
		expectedLen int
	}{
		{
			name: "empty directory",
			setup: func() {
				// No files created
			},
			wantErr:     false,
			expectedLen: 0,
		},
		{
			name: "directory with valid standard files",
			setup: func() {
				// Create standard1.md
				standard1Path := filepath.Join(tempDir, "standard1.md")
				standard1Content := `---
description: "First standard for testing"
---
This is the content of standard 1.`
				if err := os.WriteFile(standard1Path, []byte(standard1Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}

				// Create standard2.md
				standard2Path := filepath.Join(tempDir, "standard2.md")
				standard2Content := `---
description: "Second standard for testing"
---
This is the content of standard 2.`
				if err := os.WriteFile(standard2Path, []byte(standard2Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}

				// Create no_frontmatter.md
				standard3Path := filepath.Join(tempDir, "no_frontmatter.md")
				standard3Content := `This standard has no frontmatter.
Just plain content.`
				if err := os.WriteFile(standard3Path, []byte(standard3Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			},
			wantErr:     false,
			expectedLen: 3,
		},
		{
			name: "directory with mixed file types",
			setup: func() {
				// Create valid standard file
				standardPath := filepath.Join(tempDir, "standard.md")
				standardContent := `---
description: "Valid standard"
---
Content here.`
				if err := os.WriteFile(standardPath, []byte(standardContent), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}

				// Create non-markdown file (should be ignored)
				txtPath := filepath.Join(tempDir, "readme.txt")
				txtContent := "This is not a standard file"
				if err := os.WriteFile(txtPath, []byte(txtContent), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}

				// Create hidden file (should be ignored)
				hiddenPath := filepath.Join(tempDir, ".hidden.md")
				if err := os.WriteFile(hiddenPath, []byte("hidden"), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			},
			wantErr:     false,
			expectedLen: 1,
		},
		{
			name: "directory with malformed frontmatter",
			setup: func() {
				// Create file with bad frontmatter
				badPath := filepath.Join(tempDir, "bad.md")
				badContent := `---
description: "unclosed quote
---
Some content`
				if err := os.WriteFile(badPath, []byte(badContent), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}

				// Create valid file
				goodPath := filepath.Join(tempDir, "good.md")
				goodContent := `---
description: "Good standard"
---
Good content`
				if err := os.WriteFile(goodPath, []byte(goodContent), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			},
			wantErr:     true, // Should fail due to malformed frontmatter
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean temp dir
			for _, f := range []string{"standard1.md", "standard2.md", "no_frontmatter.md", "standard.md", "readme.txt", ".hidden.md", "bad.md", "good.md"} {
				_ = os.Remove(filepath.Join(tempDir, f)) // Ignore error - cleanup may fail if file doesn't exist
			}

			tt.setup()

			loader := NewFileStandardLoader()
			got, err := loader.ListStandards(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("FileStandardLoader.ListStandards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != tt.expectedLen {
					t.Errorf("FileStandardLoader.ListStandards() returned %d standards, expected %d", len(got), tt.expectedLen)
				}

				// Verify that all returned standards have names
				for _, standard := range got {
					if standard.Name == "" {
						t.Errorf("FileStandardLoader.ListStandards() returned standard with empty name")
					}
				}
			}
		})
	}
}

func TestFileStandardLoader_GetStandards(t *testing.T) {
	tempDir := t.TempDir()

	// Set up environment variables
	originalFolder, hasFolder := os.LookupEnv("AGENT_STANDARDS_MCP_FOLDER")
	originalMaxStandardSize, hasMaxStandardSize := os.LookupEnv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE")
	defer func() {
		if hasFolder {
			if err := os.Setenv("AGENT_STANDARDS_MCP_FOLDER", originalFolder); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_FOLDER: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_FOLDER"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_FOLDER: %v", err)
			}
		}
		if hasMaxStandardSize {
			if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", originalMaxStandardSize); err != nil {
				t.Logf("Warning: failed to restore AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		} else {
			if err := os.Unsetenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE"); err != nil {
				t.Logf("Warning: failed to unset AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
			}
		}
	}()

	if err := os.Setenv("AGENT_STANDARDS_MCP_FOLDER", tempDir); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_FOLDER: %v", err)
	}
	if err := os.Setenv("AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE", "1024"); err != nil {
		t.Fatalf("Failed to set AGENT_STANDARDS_MCP_MAX_STANDARD_SIZE: %v", err)
	}

	tests := []struct {
		name          string
		setup         func() map[string]string // standardName -> filePath
		standardNames []string
		wantErr       bool
		expected      int // number of standards expected to be returned
	}{
		{
			name: "get existing standards",
			setup: func() map[string]string {
				files := make(map[string]string)

				// Create standard1.md
				standard1Path := filepath.Join(tempDir, "standard1.md")
				standard1Content := `---
description: "First standard"
---
Content of standard 1`
				if err := os.WriteFile(standard1Path, []byte(standard1Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				files["standard1"] = standard1Path

				// Create standard2.md
				standard2Path := filepath.Join(tempDir, "standard2.md")
				standard2Content := `---
description: "Second standard"
---
Content of standard 2`
				if err := os.WriteFile(standard2Path, []byte(standard2Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				files["standard2"] = standard2Path

				return files
			},
			standardNames: []string{"standard1", "standard2"},
			wantErr:       false,
			expected:      2,
		},
		{
			name: "get non-existent standard",
			setup: func() map[string]string {
				// Create only standard1.md
				standard1Path := filepath.Join(tempDir, "standard1.md")
				standard1Content := `---
description: "First standard"
---
Content`
				if err := os.WriteFile(standard1Path, []byte(standard1Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				return map[string]string{"standard1": standard1Path}
			},
			standardNames: []string{"standard1", "nonexistent"},
			wantErr:       false, // Should not fail - missing standards are just skipped
			expected:      1,     // Should return only the existing standard
		},
		{
			name: "get standards with no frontmatter",
			setup: func() map[string]string {
				files := make(map[string]string)

				// Create standard with frontmatter
				standard1Path := filepath.Join(tempDir, "standard1.md")
				standard1Content := `---
description: "With frontmatter"
---
Content 1`
				if err := os.WriteFile(standard1Path, []byte(standard1Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				files["standard1"] = standard1Path

				// Create standard without frontmatter
				standard2Path := filepath.Join(tempDir, "standard2.md")
				standard2Content := `Just content without frontmatter`
				if err := os.WriteFile(standard2Path, []byte(standard2Content), 0644); err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
				files["standard2"] = standard2Path

				return files
			},
			standardNames: []string{"standard1", "standard2"},
			wantErr:       false,
			expected:      2,
		},
		{
			name: "empty standard names list",
			setup: func() map[string]string {
				return make(map[string]string)
			},
			standardNames: []string{},
			wantErr:       false,
			expected:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean temp dir
			for _, f := range []string{"standard1.md", "standard2.md"} {
				_ = os.Remove(filepath.Join(tempDir, f)) // Ignore error - cleanup may fail if file doesn't exist
			}

			tt.setup()

			loader := NewFileStandardLoader()
			got, err := loader.GetStandards(context.Background(), tt.standardNames)

			if (err != nil) != tt.wantErr {
				t.Errorf("FileStandardLoader.GetStandards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != tt.expected {
					t.Errorf("FileStandardLoader.GetStandards() returned %d standards, expected %d", len(got), tt.expected)
				}

				// Verify returned standards
				for i, standard := range got {
					if standard.Name == "" {
						t.Errorf("FileStandardLoader.GetStandards() returned standard with empty name at index %d", i)
					}
					if standard.Content == "" {
						t.Errorf("FileStandardLoader.GetStandards() returned standard with empty content at index %d", i)
					}
				}
			}
		})
	}
}

func TestExtractStandardName(t *testing.T) {
	tests := []struct {
		filePath string
		expected string
	}{
		{
			filePath: "/path/to/standard.md",
			expected: "standard",
		},
		{
			filePath: "/path/to/complex-standard-name.md",
			expected: "complex-standard-name",
		},
		{
			filePath: "simple.txt",
			expected: "simple",
		},
		{
			filePath: "/path/to/.hidden.md",
			expected: ".hidden",
		},
		{
			filePath: "/path/to/no-extension",
			expected: "no-extension",
		},
		{
			filePath: "/path/to/multiple.dots.in.name.md",
			expected: "multiple.dots.in.name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			got := extractStandardName(tt.filePath)
			if got != tt.expected {
				t.Errorf("extractStandardName(%s) = %s, expected %s", tt.filePath, got, tt.expected)
			}
		})
	}
}
