package errors

import "github.com/google/go-containerregistry/pkg/v1/types"

func NewPlatformError(platform Platform, format types.MediaType, digest string) error {
	return &platformUndefined{
		format:   format,
		digest:   digest,
		platform: platform,
	}
}

func NewDigestNotFoundError(digest string) error {
	return &digestNotFound{
		digest: digest,
	}
}

func NewUnknownMediaTypeError(format types.MediaType) error {
	return &unknownMediaType{
		format: format,
	}
}
