package anime

type AnimeService struct{}

func NewAnimeService() *AnimeService {
	return &AnimeService{}
}

func (s *AnimeService) GetTest() string {
	return "Anime test successful!"
}
