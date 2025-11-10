package test

// DefaultStandardFiles returns the default set of test standard files
func DefaultStandardFiles() map[string]string {
	return map[string]string{
		"standard1.md": `---
description: "A test standard for basic functionality"
---
This is the content of standard1.
It contains multiple lines of text.
And demonstrates basic standard loading.
`,
		"standard2.md": `---
description: "Another test standard with different content"
---
Standard 2 content here.
This standard has different content to test variety.
`,
		"standard3.md": `---
description: "A third standard for testing"
---
Content for standard 3.
Used to test multiple standards scenario.
`,
		"no-description.md": `This standard has no frontmatter description.
It should still be loadable but with empty description.
`,
		"complex-standard.md": `---
description: "A more complex standard with advanced features"
other_field: "this field should be ignored"
author: "test author"
---
# Complex Standard

This standard demonstrates:
- Advanced formatting
- Multiple sections
- Various content types

## Usage

Use this standard for testing complex scenarios.
`,
	}
}

// EmptyStandardFiles returns an empty set of standard files
func EmptyStandardFiles() map[string]string {
	return map[string]string{}
}

// CustomStandardFiles returns a custom set of standard files for specific tests
func CustomStandardFiles() map[string]string {
	return map[string]string{
		"custom1.md": `---
description: "Custom standard 1"
---
Custom content 1
`,
		"custom2.md": `---
description: "Custom standard 2"
---
Custom content 2
`,
	}
}

// NonExistentStandardFile returns a standard file for testing non-existent standards
func NonExistentStandardFile() map[string]string {
	return map[string]string{
		"real-standard.md": `---
description: "A real standard"
---
This standard actually exists.
`,
	}
}
