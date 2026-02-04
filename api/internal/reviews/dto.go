package reviews

type ReviewListingUriRequest struct {
	BookID string `uri:"id" binding:"required"`
}
