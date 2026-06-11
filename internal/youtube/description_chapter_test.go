// ./internal/youtube/description_chapter_test.go

package youtube

import "testing"

func TestExtractChapterText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty description",
			input:    "",
			expected: "",
		},
		{
			name: "valid chapters",
			input: `
0:00 Introduction
2:14 Setup
5:30 Conclusion
`,
			expected: `0:00 Introduction
2:14 Setup
5:30 Conclusion`,
		},
		{
			name: "chapters with parentheses",
			input: `
Introduction (0:00)
Setup (2:14)
Conclusion (5:30)
`,
			expected: `0:00 Introduction
2:14 Setup
5:30 Conclusion`,
		},
		{
			name: "less than three chapters",
			input: `
0:00 Intro
1:00 Setup
`,
			expected: "",
		},
		{
			name: "timestamps not increasing",
			input: `
0:00 Intro
5:00 Advanced
3:00 Basics
`,
			expected: "",
		},
		{
			name: "first chapter not zero",
			input: `
1:00 Intro
2:00 Setup
3:00 End
`,
			expected: "",
		},
		{
			name: "description without chapters",
			input: `
This is a normal description.

Please subscribe.

https://example.com
`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractChapterText(tt.input)

			if result != tt.expected {
				t.Fatalf(
					"expected:\n%s\n\ngot:\n%s",
					tt.expected,
					result,
				)
			}
		})
	}
}

func TestValidateChapters(t *testing.T) {
	tests := []struct {
		name     string
		chapters []Chapter
		valid    bool
	}{
		{
			name: "valid chapters",
			chapters: []Chapter{
				{"Intro", "0:00", 0},
				{"Part 1", "2:00", 120},
				{"Part 2", "5:00", 300},
			},
			valid: true,
		},
		{
			name: "less than three chapters",
			chapters: []Chapter{
				{"Intro", "0:00", 0},
				{"Part 1", "2:00", 120},
			},
			valid: false,
		},
		{
			name: "first chapter not zero",
			chapters: []Chapter{
				{"Intro", "1:00", 60},
				{"Part 1", "2:00", 120},
				{"Part 2", "5:00", 300},
			},
			valid: false,
		},
		{
			name: "timestamps not increasing",
			chapters: []Chapter{
				{"Intro", "0:00", 0},
				{"Part 1", "5:00", 300},
				{"Part 2", "3:00", 180},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateChapters(tt.chapters)

			if result != tt.valid {
				t.Fatalf(
					"expected %v, got %v",
					tt.valid,
					result,
				)
			}
		})
	}
}