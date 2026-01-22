package authors

type AuthorsService struct{}

func NewAuthorsService() *AuthorsService {
	return &AuthorsService{}
}

func (s *AuthorsService) ProcessAuthors(authorNames []string) map[string]string {
	// For each author:
	// Check for exact match in db
	// If exists, return that ID
	// If not, check for exact match in alias table
	// If exists, return that ID
	// If not, normalise name and check against normalised name in author table
	// If exists, add name to alias table and return that author's ID
	// If it is normalised versions are highly similar, create new author entry flagged as possible duplicate
	// Keep track of all new authors to batch insert at the end
	// Return map of author names to IDs
	return nil
}
