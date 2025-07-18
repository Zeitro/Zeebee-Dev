package gormx

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BaseRepo[T schema.Tabler] struct {
	db *gorm.DB
	// model T
}

func NewBaseRepo[T schema.Tabler](db *gorm.DB) *BaseRepo[T] {
	// m := new(T)
	return &BaseRepo[T]{
		db: db,
		// model: *m,
	}
}

// Client return gorm db
func (repo *BaseRepo[T]) Client(ctx context.Context) *gorm.DB {
	return repo.db.Model(new(T)).WithContext(ctx)
}

// BuildQuery build query
func (repo *BaseRepo[T]) BuildQuery(query *gorm.DB, opt QueryOption) *gorm.DB {
	for _, whereRaw := range opt.Filter {
		query = query.Where(whereRaw.Raw, whereRaw.Args...)
	}
	for _, order := range opt.Order {
		query = query.Order(order)
	}

	for _, preload := range opt.Preloads {
		query = query.Preload(preload)
	}

	if len(opt.Fields) > 0 {
		query = query.Select(opt.Fields)
	}

	return query
}

// Transaction transaction
func (repo *BaseRepo[T]) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return repo.db.WithContext(ctx).Transaction(fn)
}

// IsExisted check is exist
func (repo *BaseRepo[T]) IsExisted(ctx context.Context, opt QueryOption) (bool, error) {
	query := repo.Client(ctx)
	query = repo.BuildQuery(query, opt)
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *BaseRepo[T]) IsExistedByColumns(ctx context.Context, columns map[string]interface{}) (bool, error) {
	query := repo.Client(ctx)
	for k, v := range columns {
		query = query.Where(k+" = ?", v)
	}
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create create
func (repo *BaseRepo[T]) Create(ctx context.Context, data T) (T, error) {
	err := repo.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// CreateWithTx create with tx
func (repo *BaseRepo[T]) CreateWithTx(ctx context.Context, tx *gorm.DB, data T) (T, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// FindOneById find one by id
func (repo *BaseRepo[T]) FindOneById(ctx context.Context, id int64) (T, error) {
	var data T
	var err error
	filters := []WhereRaw{
		{
			Raw:  "id = ?",
			Args: []interface{}{id},
		},
	}
	data, err = repo.FindOne(ctx, NewQueryOption(WithFilterOpt(filters)))
	return data, err
}

// FindOne find one
func (repo *BaseRepo[T]) FindOne(ctx context.Context, opt QueryOption) (T, error) {
	query := repo.Client(ctx)

	query = repo.BuildQuery(query, opt)
	var data T
	err := query.First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// FindOneByColumns find one by columns
func (repo *BaseRepo[T]) FindOneByColumns(ctx context.Context, columns map[string]interface{}) (T, error) {
	query := repo.Client(ctx)
	for k, v := range columns {
		query = query.Where(k+" = ?", v)
	}
	var data T
	err := query.First(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// DeleteById delete one
func (repo *BaseRepo[T]) DeleteByID(ctx context.Context, id int64) (effectRaws int64, err error) {
	filters := []WhereRaw{
		{
			Raw:  "id = ?",
			Args: []interface{}{id},
		},
	}
	return repo.Delete(ctx, NewQueryOption(WithFilterOpt(filters)))
}

// Delete delete
func (repo *BaseRepo[T]) Delete(ctx context.Context, opt QueryOption) (effectRaws int64, err error) {
	if len(opt.Filter) == 0 {
		return 0, errors.New("filters is empty")
	}
	query := repo.Client(ctx)
	query = repo.BuildQuery(query, opt)
	err = query.Delete(new(T)).Error
	if err != nil {
		return 0, err
	}
	return query.RowsAffected, nil
}

// Page page
func (repo *BaseRepo[T]) Page(query *gorm.DB, pageReq PageRequest) (page PageResponse[T], err error) {
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
	err = query.Limit(pageReq.Limit).Offset((pageReq.Page - 1) * pageReq.Limit).Find(&page.List).Error
	if err != nil {
		return page, err
	}
	page.Page = pageReq.Page
	page.Limit = pageReq.Limit
	page.Total = int(total)
	page.TotalPage = int(total) / pageReq.Limit
	if int(total)%pageReq.Limit != 0 {
		page.TotalPage++
	}
	return page, nil
}

// FindAll find all
func (repo *BaseRepo[T]) FindAll(ctx context.Context, opt QueryOption) ([]T, error) {
	query := repo.Client(ctx)
	query = repo.BuildQuery(query, opt)
	var data []T
	err := query.Find(&data).Error
	return data, err
}
