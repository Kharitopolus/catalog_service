package main

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) CreateProduct(p *Product) error {
	_, err := s.db.Exec(
		`INSERT INTO products (id, name, description, price, category_id,created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		p.ID,
		p.Name,
		p.Description,
		p.Price,
		p.CategoryID,
		p.CreatedAt,
	)
	return err
}

func (s *PostgresStore) UpdateProduct(prod *Product) error {
	return nil
}

func (s *PostgresStore) DeleteProduct(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

func (s *PostgresStore) GetProducts() ([]*Product, error) {
	rows, err := s.db.Query(
		`SELECT id, name, description, price, category_id, created_at
		FROM products`,
	)
	if err != nil {
		return nil, err
	}
	products := []*Product{}
	for rows.Next() {
		p := &Product{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (s *PostgresStore) GetProductByID(id uuid.UUID) (*Product, error) {
	p := &Product{}
	if err := s.db.QueryRow(
		`SELECT id, name, description, price, category_id, created_at
		FROM products WHERE id = $1`, id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PostgresStore) GetCategories() ([]*Category, error) {
	rows, err := s.db.Query("SELECT id, name, created_at FROM categories")
	if err != nil {
		return nil, err
	}
	categories := []*Category{}
	for rows.Next() {
		c := &Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (s *PostgresStore) CreateCategory(c *Category) error {
	_, err := s.db.Exec(
		`INSERT INTO categories (id, name, created_at)
		VALUES ($1, $2, $3)`,
		c.ID,
		c.Name,
		c.CreatedAt,
	)
	return err
}

func NewPostgresStore(url string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	_, err := s.db.Exec(
		`CREATE TABLE IF NOT EXISTS categories (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS products (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			price INT NOT NULL,
			category_id UUID REFERENCES categories(id),
			created_at TIMESTAMP
		);`,
	)
	return err
}
