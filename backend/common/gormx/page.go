package gormx

import "gorm.io/gorm"

type PageRequest struct {
	Page  int `form:"page,default=1" json:"page"`
	Limit int `form:"limit,default=10" json:"limit"`
}

type PageResponse[T any] struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
	List      []T `json:"list"`
}

func Page[T any](query *gorm.DB, pageReq PageRequest) (page PageResponse[T], err error) {
	if pageReq.Page <= 0 {
		pageReq.Page = 1
	}
	if pageReq.Limit <= 0 {
		pageReq.Limit = 20
	}
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		return page, err
	}
	if total == 0 {
		return page, nil
	}
	err = query.Limit(pageReq.Limit).
		Offset((pageReq.Page - 1) * pageReq.Limit).
		Scan(&page.List).
		Error
	if err != nil {
		return page, err
	}
	page.Page = pageReq.Page
	page.Limit = pageReq.Limit
	page.Total = int(total)
	page.TotalPage = int(total) / pageReq.Limit
	if int(total)%pageReq.Limit > 0 {
		page.TotalPage++
	}
	return page, nil
}
