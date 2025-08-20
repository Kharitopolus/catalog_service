package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateProductRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  uuid.UUID `json:"categoryID"`
}

func (r CreateProductRequest) validate() map[string]string {
	errors := map[string]string{}
	if len(r.Name) > 256 || len(r.Name) < 1 {
		errors["name"] = "The length of the name should be between 1 and 256 characters long"
	}
	if len(r.Description) > 512 {
		errors["description"] = "The length of the description should shorter than 512 characters long"
	}
	if r.Price <= 0 {
		errors["price"] = "Price should be greater than 0"
	}
	return errors
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (r CreateCategoryRequest) validate() map[string]string {
	errors := map[string]string{}
	if len(r.Name) > 256 || len(r.Name) < 1 {
		errors["name"] = "The length of the name should be between 1 and 256 characters long"
	}
	return errors
}

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	CategoryID  uuid.UUID `json:"categoryID"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewProduct(
	name, description string,
	price int,
	categoryID uuid.UUID,
) *Product {
	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		CreatedAt:   time.Now().UTC(),
	}
}

type Category struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewCategory(name string) *Category {
	return &Category{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}
}
