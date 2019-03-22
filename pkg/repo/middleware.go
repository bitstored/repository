package repo

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

// NamespaceRepo returnes a repo which namespaces on top of the the given
// repository by prefixing all keys.
func NamespaceRepo(prefix string) Middleware {
	return func(r Repository) Repository {
		return &repo{prefix: prefix, next: r}
	}
}

type repo struct {
	prefix string
	next   Repository
}

func (r *repo) Create(ctx context.Context, key string, doc interface{}) (uint, error) {
	return r.next.Create(ctx, r.prefix+key, doc)
}

func (r *repo) Read(ctx context.Context, key string, doc interface{}) (uint, error) {
	return r.next.Read(ctx, r.prefix+key, doc)
}

func (r *repo) Update(ctx context.Context, key string, cas uint, doc interface{}) (uint, error) {
	return r.next.Update(ctx, r.prefix+key, cas, doc)
}

func (r *repo) Delete(ctx context.Context, key string, cas uint) (uint, error) {
	return r.next.Delete(ctx, r.prefix+key, cas)
}

// EncrypterDecrypter is a interface that defines the basic methods to
// encrypt and decrypt documents.
type EncrypterDecrypter interface {
	Encrypt(io.Writer, io.Reader) error
	Decrypt(io.Writer, io.Reader) error
}

// EncryptionRepo returnes a Repository which uses the given SecurityDevice to
// encrypt stored and decrypt retrieved documents.
func EncryptionRepo(device EncrypterDecrypter) Middleware {
	return func(r Repository) Repository {
		return &encRepo{
			device: device,
			next:   r,
		}
	}
}

type encRepo struct {
	device EncrypterDecrypter
	next   Repository
}

func (r *encRepo) Create(ctx context.Context, key string, doc interface{}) (uint, error) {
	var plaintext, ciphertext bytes.Buffer
	if err := json.NewEncoder(&plaintext).Encode(doc); err != nil {
		return 0, err
	}
	if err := r.device.Encrypt(&ciphertext, &plaintext); err != nil {
		return 0, err
	}
	return r.next.Create(ctx, key, base64.StdEncoding.EncodeToString(ciphertext.Bytes()))
}

func (r *encRepo) Read(ctx context.Context, key string, doc interface{}) (uint, error) {
	var cipherBase64 string
	cas, err := r.next.Read(ctx, key, &cipherBase64)
	if err != nil {
		return 0, err
	}
	ciphertext := base64.NewDecoder(base64.StdEncoding, strings.NewReader(cipherBase64))

	var plaintext bytes.Buffer
	if err := r.device.Decrypt(&plaintext, ciphertext); err != nil {
		return 0, err
	}
	if err = json.NewDecoder(&plaintext).Decode(doc); err != nil {
		return 0, err
	}
	return cas, nil
}

// Update replaces a existing document.
func (r *encRepo) Update(ctx context.Context, key string, cas uint, doc interface{}) (uint, error) {
	var plaintext, ciphertext bytes.Buffer
	if err := json.NewEncoder(&plaintext).Encode(doc); err != nil {
		return 0, err
	}
	if err := r.device.Encrypt(&ciphertext, &plaintext); err != nil {
		return 0, err
	}
	return r.next.Update(ctx, key, cas, base64.StdEncoding.EncodeToString(ciphertext.Bytes()))
}

// Delete deletes the document with the given key.
func (r *encRepo) Delete(ctx context.Context, key string, cas uint) (uint, error) {
	return r.next.Delete(ctx, key, cas)
}

// NewAESPassphrase returnes a AESPassphares.
func NewAESPassphrase(passphrase string) *AESPassphrase {
	hasher := sha512.New()
	hasher.Write([]byte(passphrase))
	hasher.Write([]byte(strconv.Itoa(len(passphrase))))
	return &AESPassphrase{key: hasher.Sum(nil)[:32]}
}

// AESPassphrase implements the EncrypterDecrypter interface.
type AESPassphrase struct {
	key []byte
}

// Encrypt implements the encryption for the EncrypterDecrypter interface.
func (p *AESPassphrase) Encrypt(w io.Writer, r io.Reader) error {
	return nil
}

// Decrypt implements the decryption for the EncrypterDecrypter interface.
func (p *AESPassphrase) Decrypt(w io.Writer, r io.Reader) error {
	return nil
}
