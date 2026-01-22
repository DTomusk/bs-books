package authors

import "testing"

func TestNormaliseAuthorName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"J.K. Rowling", "j k rowling"},
		{"george r.r. martin", "george r r martin"},
		{"  Agatha Christie  ", "agatha christie"},
		{"Émile Zola", "emile zola"},
		{"Gabriel García Márquez", "gabriel garcia marquez"},
		{"Mary-Sue O'Connor", "mary sue o connor"},
		{"Jean-Luc Picard", "jean luc picard"},
		{" Hồ Chí Minh ", "ho chi minh"},
		{"Renée Zellweger", "renee zellweger"},
		{"Björk Guðmundsdóttir", "bjork guðmundsdottir"},
		{" ROSALÍA ", "rosalia"},
		{"黒澤 明", "黒澤 明"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normaliseAuthorName(tt.name); got != tt.expected {
				t.Errorf("normaliseAuthorName() = %v, want %v", got, tt.expected)
			}
		})
	}
}
