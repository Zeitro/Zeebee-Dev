package gormx

import (
	"context"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BaseRepo[T schema.Tabler] struct {
	db *gorm.DB
}

func NewBaseRepo[T schema.Tabler](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{
		db: db,
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

// IsExistedByColumns check is exist by columns
func (repo *BaseRepo[T]) IsExistedByColumns(ctx context.Context, columns map[string]interface{}) (bool, error) {
	query := repo.Client(ctx)
	for k, v := range columns {
		query = query.Where(k+" = ?", v)
	}
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false, errors.WithStack(err)
	}
	return count > 0, nil
}

// Create create
func (repo *BaseRepo[T]) Create(ctx context.Context, data T) (T, error) {
	err := repo.db.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, errors.WithStack(err)
	}
	return data, nil
}

// CreateWithTx create with tx
func (repo *BaseRepo[T]) CreateWithTx(ctx context.Context, tx *gorm.DB, data T) (T, error) {
	err := tx.WithContext(ctx).Create(&data).Error
	if err != nil {
		return data, errors.WithStack(err)
	}
	return data, nil
}

// FindOneById find one by id
func (repo *BaseRepo[T]) FindOneById(ctx context.Context, id int64) (T, error) {
	var data T
	var err error
	data, err = repo.FindOneByColumns(ctx, map[string]interface{}{
		"id": id,
	})
	return data, errors.WithStack(err)
}

func (repo *BaseRepo[T]) FindOneByIdWithTx(ctx context.Context, tx *gorm.DB, id int64) (T, error) {
	var data T
	var err error
	data, err = repo.FindOneByColumnsWithTx(ctx, tx, map[string]interface{}{
		"id": id,
	})
	return data, errors.WithStack(err)
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
		return data, errors.WithStack(err)
	}
	return data, nil
}

func (repo *BaseRepo[T]) FindOneByColumnsWithTx(ctx context.Context, tx *gorm.DB, columns map[string]interface{}) (T, error) {
	query := tx.WithContext(ctx)
	for k, v := range columns {
		query = query.Where(k+" = ?", v)
	}
	var data T
	err := query.First(&data).Error
	if err != nil {
		return data, errors.WithStack(err)
	}
	return data, nil
}

func (repo *BaseRepo[T]) UpdateOne(ctx context.Context, data T) (T, error) {
	err := repo.Client(ctx).WithContext(ctx).Save(data).Error
	if err != nil {
		return data, errors.WithStack(err)
	}
	return data, nil
}

func (repo *BaseRepo[T]) UpdateOneWithTx(ctx context.Context, tx *gorm.DB, data T) (T, error) {
	err := tx.WithContext(ctx).Save(data).Error
	if err != nil {
		return data, errors.WithStack(err)
	}
	return data, nil
}

// DeleteById delete one
func (repo *BaseRepo[T]) DeleteByID(ctx context.Context, id int64) (effectRaws int64, err error) {

	return repo.Delete(ctx, map[string]interface{}{
		"id": id,
	})
}

// Delete delete
func (repo *BaseRepo[T]) Delete(ctx context.Context, filters map[string]interface{}) (effectRaws int64, err error) {
	if len(filters) == 0 {
		return 0, errors.New("filters is empty")
	}
	query := repo.Client(ctx)
	for k, v := range filters {
		query = query.Where(k+" = ?", v)
	}
	err = query.Delete(new(T)).Error
	if err != nil {
		return 0, errors.WithStack(err)
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
		return page, errors.WithStack(err)
	}
	if total == 0 {
		return page, nil
	}
	err = query.Limit(pageReq.Limit).Offset((pageReq.Page - 1) * pageReq.Limit).Find(&page.List).Error
	if err != nil {
		return page, errors.WithStack(err)
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
func (repo *BaseRepo[T]) FindAll(ctx context.Context, filter map[string]interface{}) ([]T, error) {
	query := repo.Client(ctx)
	for k, v := range filter {
		query = query.Where(k+" = ?", v)
	}

	var data []T
	err := query.Find(&data).Error
	return data, errors.WithStack(err)
}
