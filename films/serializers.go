package films

import "github.com/gin-gonic/gin"

type FilmSerializer struct {
	C *gin.Context
	FilmModel
}

type FilmResponse struct {
	ID             uint                  `json:"-"`
	Title          string                `json:"title"`
}

func (s *FilmSerializer) Response() FilmResponse {
	response := FilmResponse{
		ID:          s.ID,
		Title:       s.Title,
	}
	return response
}