package books

import "testing"

func TestNormaliseBookTitle(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{" The Great Gatsby ", "great gatsby"},
		{"Moby-Dick", "moby dick"},
		{"To Kill a Mockingbird!", "to kill a mockingbird"},
		{"1984", "1984"},
		{"Pride & Prejudice", "pride & prejudice"},
		{" War and Peace ", "war and peace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormaliseBookTitle(tt.name); got != tt.expected {
				t.Errorf("NormaliseBookTitle() = %v, want %v", got, tt.expected)
			}
		})
	}
}
