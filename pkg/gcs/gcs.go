package gcs

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

// GCS represents a Google cloud storage
type GCS struct {
	client     *storage.Client
	bucketName string
}

// New builds new GCS
func New(bucketName string) (*GCS, error) {
	ctx := context.TODO()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GCS client: %v", err)
	}
	return &GCS{
		client:     client,
		bucketName: bucketName,
	}, nil
}

// Close closes the client connection
func (g *GCS) Close() error {
	return g.client.Close()
}

func (g *GCS) Put(ctx context.Context, data []byte, path string) error {
	w := g.client.Bucket(g.bucketName).Object(path).NewWriter(ctx)
	defer func() {
		if err := w.Close(); err != nil {
			fmt.Print(err)
		}
	}()

	if n, err := w.Write(data); err != nil {
		return fmt.Errorf("failed to write data: %v", err)
	} else if n != len(data) {
		return fmt.Errorf("failed to write data, got: %d, want: %d, err: %v", n, len(data), err)
	}

	return nil
}

func (g *GCS) Get(ctx context.Context, path string) ([]byte, error) {
	r, err := g.client.Bucket(g.bucketName).Object(path).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new object reader: %v", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			fmt.Print(err)
		}
	}()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %v", err)
	}
	return b, nil
}

func (g *GCS) HasObject(ctx context.Context, objectPath string) (bool, error) {
	obj := g.client.Bucket(g.bucketName).Object(objectPath)
	rc, err := obj.NewReader(ctx)
	if rc != nil {
		defer func() {
			if err := rc.Close(); err != nil {
				fmt.Print(err)
			}
		}()
	}
	if err == nil { // already exists
		return true, nil
	}

	// > ErrObjectNotExist will be returned if the object is not found.
	// https://godoc.org/cloud.google.com/go/storage#ObjectHandle.NewReader
	if err == storage.ErrObjectNotExist {
		return false, nil
	}

	return false, err
}
