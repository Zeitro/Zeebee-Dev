package gormx

import "gorm.io/gorm"

type PageRequest struct {
	Page         int64    `form:"page,default=1" json:"page"`
	Limit        int64    `form:"limit,default=10" json:"limit"`
	CountColumn  string   `json:"count_column,omitempty"`
	SelectFields []string `json:"select_fields,omitempty"`
}

type PageResponse[T any] struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
	List      []T   `json:"list"`
}

func Page[T any](query *gorm.DB, pageReq PageRequest) (page PageResponse[T], err error) {
	if pageReq.Page <= 0 {
		pageReq.Page = 1
	}
	if pageReq.Limit <= 0 {
		pageReq.Limit = 20
	}
	var total int64

	countQuery := query.Session(&gorm.Session{})
	if len(pageReq.CountColumn) > 0 {
		countQuery = countQuery.Select(pageReq.CountColumn)
	}
	err = countQuery.Count(&total).Error
	if err != nil {
		return page, err
	}
	if total == 0 {
		return page, nil
	}
	if len(pageReq.SelectFields) > 0 {
		query = query.Select(pageReq.SelectFields)
	}
	err = query.Limit(int(pageReq.Limit)).
		Offset(int((pageReq.Page - 1) * pageReq.Limit)).
		Scan(&page.List).
		Error
	if err != nil {
		return page, err
	}
	page.Page = pageReq.Page
	page.Limit = pageReq.Limit
	page.Total = total
	page.TotalPage = total / pageReq.Limit
	if total%pageReq.Limit > 0 {
		page.TotalPage++
	}
	return page, nil
}
