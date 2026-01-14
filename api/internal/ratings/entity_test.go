package ratings

import "testing"

func TestNewRating(t *testing.T) {
	_, err := newRating("book-id", 4.0, 3.0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNewRating_NegativeScore(t *testing.T) {
	_, err := newRating("book-id", -1.0, 3.0)
	if err != ErrNegativeScore {
		t.Fatalf("expected ErrNegativeScore, got %v", err)
	}
	_, err = newRating("book-id", 4.0, -2.0)
	if err != ErrNegativeScore {
		t.Fatalf("expected ErrNegativeScore, got %v", err)
	}
}

func TestNewRating_LargeScore(t *testing.T) {
	_, err := newRating("book-id", 6.0, 3.0)
	if err != ErrLargeScore {
		t.Fatalf("expected ErrLargeScore, got %v", err)
	}
	_, err = newRating("book-id", 4.0, 7.0)
	if err != ErrLargeScore {
		t.Fatalf("expected ErrLargeScore, got %v", err)
	}
}
