package books

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
}

func (h *BookHandler) GetBooks() []string {
	return []string{"Book 1", "Book 2", "Book 3"}
}
