package auth

import "testing"

func TestHashPassword(t *testing.T) {
	// Arrange
	password := "securepassword"

	// Act & Assert
	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == password {
		t.Fatalf("expected hashed password to differ from original password")
	}
	hash2, err := hashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hash == hash2 {
		t.Fatalf("expected different hashes for the same password due to salting")
	}
}
