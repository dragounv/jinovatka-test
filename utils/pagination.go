package utils

const DefaultLinesPerPage = 20

func NewPagination(page, noPages, linesPerPage int) Pagination {
	return Pagination{
		Page:         page,
		NoPages:      noPages,
		LinesPerPage: linesPerPage,
	}
}

type Pagination struct {
	Page         int
	NoPages      int
	LinesPerPage int
}
