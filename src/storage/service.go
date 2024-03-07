package storage

import (
	"context"
	kithttp "github.com/go-kit/kit/transport/http"
	"mime/multipart"
)

// CreationOptions is the options for the faceswap service.
type CreationOptions struct {
	httpClientOpts []kithttp.ClientOption
	storageType    string
	endpoint       string
	region         string
	secretId       string
	secretKey      string
}

// CreationOption is a creation option for the faceswap service.
type CreationOption func(*CreationOptions)

// WithEndpoint returns a CreationOption that sets the base url.
func WithEndpoint(s3url string) CreationOption {
	return func(co *CreationOptions) {
		co.endpoint = s3url
	}
}

// WithClientHttpOptions returns a CreationOption that sets the http client options.
func WithClientHttpOptions(opts ...kithttp.ClientOption) CreationOption {
	return func(co *CreationOptions) {
		var options []kithttp.ClientOption
		for _, opt := range opts {
			options = append(options, opt)
		}
		co.httpClientOpts = options
	}
}

// WithRegion returns a CreationOption that sets the region.
func WithRegion(region string) CreationOption {
	return func(co *CreationOptions) {
		co.region = region
	}
}

// WithSecretId returns a CreationOption that sets the secret id.
func WithSecretId(secretId string) CreationOption {
	return func(co *CreationOptions) {
		co.secretId = secretId
	}
}

// WithSecretKey returns a CreationOption that sets the secret key.
func WithSecretKey(secretKey string) CreationOption {
	return func(co *CreationOptions) {
		co.secretKey = secretKey
	}
}

type Service interface {
	SaveFile(ctx context.Context, file multipart.File)
}
