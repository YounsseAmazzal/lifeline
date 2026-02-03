package dto

type PaginationHeader struct {
	PageNumber int   `json:"pageNumber"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	TotalCount int64 `json:"totalCount"`
}

type UserParams struct {
	PageNumber int    `query:"pageNumber"`
	PageSize   int    `query:"pageSize"`
	
	CurrentUserName string 
	Gender          string `query:"gender"`
	BloodGroup      string `query:"bloodGroup"`
	MinAge int `query:"minAge"`
	MaxAge int `query:"maxAge"`
	
	OrderBy string `query:"orderBy"` 
}

func (p *UserParams) SetDefaults() {
	if p.PageNumber < 1 {
		p.PageNumber = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 50 {
		p.PageSize = 50 
	}
	if p.MinAge == 0 {
		p.MinAge = 18 
	}
	if p.MaxAge == 0 {
		p.MaxAge = 65 
		}
}