package couchbaserepo

import (
	"context"

	"github.com/couchbase/gocb"
)

// BucketFunc returns a gocb.Bucket or a errors.
type BucketFunc func() (*gocb.Bucket, error)

// Repository is a couchbase implementation to of the Repository pattern.
type Repository struct {
	bucket BucketFunc
	TTL    uint32
}

// NewRepository returnes a new Repository.
func NewRepository(f BucketFunc) *Repository {
	return &Repository{
		bucket: f,
	}
}

// Create creates a new document. It returns a DocumentExistsError if the key already exists.
func (r *Repository) Create(ctx context.Context, key string, doc interface{}) (uint, error) {
	b, err := r.bucket()
	if err != nil {
		return 0, err
	}
	gocbCas, err := b.Insert(key, doc, r.TTL)
	return uint(gocbCas), err
}

// Update replaces a existing document.
func (r *Repository) Update(ctx context.Context, key string, cas uint, doc interface{}) (uint, error) {
	b, err := r.bucket()
	if err != nil {
		return 0, err
	}
	gocbCas, err := b.Replace(key, doc, gocb.Cas(cas), r.TTL)
	return uint(gocbCas), err
}

// Read reads the document from the key.
func (r *Repository) Read(ctx context.Context, key string, doc interface{}) (uint, error) {
	b, err := r.bucket()
	if err != nil {
		return 0, err
	}
	gocbCas, err := b.Get(key, &doc)
	return uint(gocbCas), err
}

// Delete deletes the document with the given key.
func (r *Repository) Delete(ctx context.Context, key string, cas uint) (uint, error) {
	b, err := r.bucket()
	if err != nil {
		return 0, err
	}
	cas1, err := b.Remove(key, gocb.Cas(cas))
	return uint(cas1), err
}
