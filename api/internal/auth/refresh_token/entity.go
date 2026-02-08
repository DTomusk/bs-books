package refresh_token

type RefreshToken struct {
	ID string
	// Raw token used in cookie
	Token string
	// Persisted hash
	TokenHash string
	IsRevoked bool
	// ExpiresAt is a unix timestamp in seconds
	ExpiresAt    int64
	UserID       string
	IPAddress    string
	FamilyID     string
	ReplacedByID *string
}
