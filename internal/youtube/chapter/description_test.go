// ./internal/youtube/description_chapter_test.go

package chapter

import (
	"testing"
	"reflect"
)

func TestExtractChapters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Chapters
	}{
		{
			name:  "empty description",
			input: "",
			expected: Chapters{},
		},
		{
			name: "valid chapters",
			input: `
0:00 Introduction
2:14 Setup
5:30 Conclusion
`,
			expected: Chapters{
				List: []Chapter{
					{"Introduction", "0:00", 0},
					{"Setup", "2:14", 134},
					{"Conclusion", "5:30", 330},
				},
				Valid: true,
			},
		},
		{
			name: "chapters with parentheses",
			input: `
Introduction (0:00)
Setup (2:14)
Conclusion (5:30)
`,
			expected: Chapters{
				List: []Chapter{
					{"Introduction", "0:00", 0},
					{"Setup", "2:14", 134},
					{"Conclusion", "5:30", 330},
				},
				Valid: true,
			},
		},
		{
			name: "description without chapters",
			input: `
This is a normal description.

Please subscribe.

https://example.com
`,
			expected: Chapters{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractChapters(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf(
					"expected %#v\ngot %#v",
					tt.expected,
					result,
				)
			}
		})
	}
}

func TestChapters_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		chapters Chapters
		valid    bool
	}{
		{
			name: "valid chapters",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
					{"Part 1", "2:00", 120},
					{"Part 2", "5:00", 300},
				},
			},
			valid: true,
		},
		{
			name: "less than three chapters",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
					{"Part 1", "2:00", 120},
				},
			},
			valid: false,
		},
		{
			name: "first chapter not zero",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "1:00", 60},
					{"Part 1", "2:00", 120},
					{"Part 2", "5:00", 300},
				},
			},
			valid: false,
		},
		{
			name: "timestamps not increasing",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
					{"Part 1", "5:00", 300},
					{"Part 2", "3:00", 180},
				},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chapters.IsValid()

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



func TestChapters_Text(t *testing.T) {
	tests := []struct {
		name     string
		chapters Chapters
		expected string
	}{
		{
			name:     "empty chapters",
			chapters: Chapters{},
			expected: "",
		},
		{
			name: "single chapter",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
				},
			},
			expected: "0:00 Intro",
		},
		{
			name: "multiple chapters",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
					{"Setup", "2:14", 134},
					{"Conclusion", "5:30", 330},
				},
			},
			expected: `0:00 Intro
2:14 Setup
5:30 Conclusion`,
		},
		{
			name: "chapter without title",
			chapters: Chapters{
				List: []Chapter{
					{"", "0:00", 0},
					{"Setup", "2:14", 134},
				},
			},
			expected: `0:00
2:14 Setup`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chapters.Text()

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


func TestChapters_Titles(t *testing.T) {
	tests := []struct {
		name     string
		chapters Chapters
		expected string
	}{
		{
			name:     "empty chapters",
			chapters: Chapters{},
			expected: "",
		},
		{
			name: "single chapter",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
				},
			},
			expected: "Intro",
		},
		{
			name: "multiple chapters",
			chapters: Chapters{
				List: []Chapter{
					{"Intro", "0:00", 0},
					{"Setup", "2:14", 134},
					{"Conclusion", "5:30", 330},
				},
			},
			expected: `Intro
Setup
Conclusion`,
		},
		{
			name: "chapter without title",
			chapters: Chapters{
				List: []Chapter{
					{"", "0:00", 0},
					{"Setup", "2:14", 134},
				},
			},
			expected: `
Setup`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chapters.Titles()

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