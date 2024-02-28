package bucket

import (
	"fmt"
	"io"
	"reflect"
)

const (
	AwsProvider BucketType = iota
	MockProvider
)

type BucketType int

func New(bt BucketType, cfg any) (b *Bucket, err error) {
	b = new(Bucket)

	rt := reflect.TypeOf(cfg)

	switch bt {
	case AwsProvider:
		if rt.Name() != "AwsConfig" {
			return nil, fmt.Errorf("Config need's to be of type AwsConfig")
		}

		b.p = newAwsSession(cfg.(AwsConfig))
	case MockProvider:
		b.p = &MockBucket{
			content: make(map[string][]byte),
		}
	default:
		return nil, fmt.Errorf("type not implemented")
	}

	return
}

type BucketInterface interface {
	Upload(io.Reader, string) error
	Download(string, string) error
	Delete(string) error
}

type Bucket struct {
	p BucketInterface
}

func (b *Bucket) Upload(file io.Reader, key string) error {
	return b.p.Upload(file, key)
}

func (b *Bucket) Download(src, dst string) error {
	return b.p.Download(src, dst)
}

func (b *Bucket) Delete(key string) error {
	return b.p.Delete(key)
}
