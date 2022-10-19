package storage

type SitemapRepository struct {
	store *Store
}

func NewSitemapRepository(store *Store) *SitemapRepository {
	return &SitemapRepository{
		store: store,
	}
}
