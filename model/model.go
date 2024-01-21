package model

type Artist struct {
	ID    string   `json:"id" validate:"required"`                   // id коллектива
	Name  string   `json:"name" validate:"required"`                 // название группы
	Born  string   `json:"born" validate:"required,gt=1900,lt=2024"` // год основания группы
	Genre string   `json:"genre" validate:"required"`                // жанр
	Songs []string `json:"songs" validate:"required"`                // популярные песни, это слайс строк, так как песен может быть несколько
}

type Filter struct {
	Genre  string `json:"genre"`
	Born   string `json:"born"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
