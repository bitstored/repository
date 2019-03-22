package repo

import (
	"context"
)

type Repository interface {
	// Create method writes the obj to database with the prefix equal with key
	Create(ctx context.Context, key string, obj interface{}) (uint, error)
	// Read method reads the object from the database, where prefix is equal to key and the selction criteria is cas
	Read(ctx context.Context, key string, obj interface{}) (uint, error)
	// Update method updates the object at the given prefix with the given criteria to the new value.
	Update(ctx context.Context, key string, cas uint, obj interface{}) (uint, error)
	// Delete method sets the deleted field to true and the entity will not be available anymore.
	Delete(ctx context.Context, key string, cas uint) (uint, error)
}

type Middleware func(Repository) Repository

// UseMiddleware returnes the given repo wrapped into the given middlewares
func UseMiddleware(repo Repository, middlewares ...Middleware) Repository {
	r := repo
	for _, mw := range middlewares {
		r = mw(r)
	}
	return r
}
