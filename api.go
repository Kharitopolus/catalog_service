package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Storage interface {
	CreateProduct(*Product) error
	UpdateProduct(*Product) error
	DeleteProduct(uuid.UUID) error
	GetProducts() ([]*Product, error)
	GetProductByID(uuid.UUID) (*Product, error)
	GetCategories() ([]*Category, error)
	CreateCategory(*Category) error
}

type APIServer struct {
	listerAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listerAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /products/", Make(s.handleGetProducts))
	mux.HandleFunc("GET /products/{id}", Make(s.handleGetProductsByID))
	mux.HandleFunc("POST /products/", Make(s.handleCreateProducts))
	mux.HandleFunc("PUT /products/{id}", Make(s.handleUpdateProducts))
	mux.HandleFunc(
		"DELETE /products/{id}",
		Make(s.handleDeleteProducts),
	)

	mux.HandleFunc("GET /categories/", Make(s.handleGetCategories))
	mux.HandleFunc("POST /categories/", Make(s.handleCreateCategories))

	http.ListenAndServe(s.listerAddr, mux)
}

func (s *APIServer) handleGetProducts(
	w http.ResponseWriter,
	r *http.Request,
) error {
	products, err := s.store.GetProducts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, products)
}

func (s *APIServer) handleGetProductsByID(
	w http.ResponseWriter,
	r *http.Request,
) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	product, err := s.store.GetProductByID(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleCreateProducts(
	w http.ResponseWriter,
	r *http.Request,
) error {
	req := &CreateProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return InvalidJSON()
	}
	defer r.Body.Close()
	if errors := req.validate(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}
	product := NewProduct(req.Name, req.Description, req.Price, req.CategoryID)
	if err := s.store.CreateProduct(product); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, product)
}

func (s *APIServer) handleUpdateProducts(
	w http.ResponseWriter,
	r *http.Request,
) error {
	return nil
}

func (s *APIServer) handleDeleteProducts(
	w http.ResponseWriter,
	r *http.Request,
) error {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		return err
	}
	if err := s.store.DeleteProduct(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]uuid.UUID{"deleted": id})
}

func (s *APIServer) handleGetCategories(
	w http.ResponseWriter,
	r *http.Request,
) error {
	categories, err := s.store.GetCategories()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, categories)
}

func (s *APIServer) handleCreateCategories(
	w http.ResponseWriter,
	r *http.Request,
) error {
	req := &CreateCategoryRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return InvalidJSON()
	}
	defer r.Body.Close()
	if errors := req.validate(); len(errors) > 0 {
		return InvalidRequestData(errors)
	}
	c := NewCategory(req.Name)
	if err := s.store.CreateCategory(c); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, c)
}
