package repository

const (
	componentRepo = "Repository"
)

type IdGenerator interface {
	NewID() string
}

func New(idgen IdGenerator) *Repository {
	return &Repository{
		idgen: idgen,
	}
}

type Repository struct {
	idgen IdGenerator
}
