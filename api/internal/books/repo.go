package books

type booksRepo struct{}

func NewBooksRepo() *booksRepo {
	return &booksRepo{}
}
