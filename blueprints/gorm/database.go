package gorm

import (
	"context"
	"db_blueprints/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DatabaseTimeout = time.Second * 5
)

type IDatabase interface {
	GetDB() *gorm.DB
	AutoMigrate(models ...any) error
	WithTransaction(function func() error) error
	Create(ctx context.Context, doc any) error
	CreateInBatches(ctx context.Context, docs any, batchSize int) error
	Update(ctx context.Context, doc any) error
	Delete(ctx context.Context, value any, opts ...FindOption) error
	FindById(ctx context.Context, id int64, result any) error
	FindOne(ctx context.Context, result any, opts ...FindOption) error
	Find(ctx context.Context, result any, opts ...FindOption) error
	Count(ctx context.Context, model any, total *int64, opts ...FindOption) error
}

type Database struct {
	db *gorm.DB
}

func NewDatabase(config *config.Config) (*Database, error) {
	// 1. Construct the Data Source Name (DSN) string from the config
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)

	// 2. Open the database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 3. If connection fails, return the error
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 4. On success, wrap the connection in your struct and return it
	gormDB := &Database{
		db: db,
	}

	fmt.Println("Successfully connected to the database!")
	return gormDB, nil
}

func (d *Database) AutoMigrate(models ...any) error {
	return d.db.AutoMigrate(models...)
}

func (d *Database) WithTransaction(function func() error) error {
	callback := func(db *gorm.DB) error {
		return function()
	}

	tx := d.db.Begin()
	if err := callback(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Preload(query string, args ...interface{}) IDatabase {
	d.db.Preload(query, args...)
	return d
}

func (d *Database) Create(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.Create(doc).Error
}

func (d *Database) CreateInBatches(ctx context.Context, docs any, batchSize int) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.CreateInBatches(docs, batchSize).Error
}

func (d *Database) Update(ctx context.Context, doc any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	return d.db.Save(doc).Error
}

func (d *Database) Delete(ctx context.Context, value any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	return query.Delete(value).Error
}

func (d *Database) FindById(ctx context.Context, id int64, result any) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	if err := d.db.Where("id = ? ", id).First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) FindOne(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.First(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Find(ctx context.Context, result any, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.Find(result).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) Count(ctx context.Context, model any, total *int64, opts ...FindOption) error {
	ctx, cancel := context.WithTimeout(ctx, DatabaseTimeout)
	defer cancel()

	query := d.applyOptions(opts...)
	if err := query.Model(model).Count(total).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func (d *Database) applyOptions(opts ...FindOption) *gorm.DB {
	query := d.db

	opt := getOption(opts...)

	if len(opt.preloads) != 0 {
		for _, preload := range opt.preloads {
			query = query.Preload(preload)
		}
	}

	if opt.query != nil {
		for _, q := range opt.query {
			query = query.Where(q.Query, q.Args)
		}
	}

	if opt.order != "" {
		query = query.Order(opt.order)
	}

	if opt.offset != 0 {
		query = query.Offset(opt.offset)
	}

	if opt.limit != 0 {
		query = query.Limit(opt.limit)
	}

	return query
}
