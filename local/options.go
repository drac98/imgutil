package local

import (
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"

	"github.com/buildpacks/imgutil"
)

// FIXME: the following functions are defined in this package for backwards compatibility,
// and should eventually be deprecated.

func FromBaseImage(name string) func(*imgutil.ImageOptions) {
	return imgutil.FromBaseImage(name)
}

func WithConfig(c *v1.Config) func(*imgutil.ImageOptions) {
	return imgutil.WithConfig(c)
}

func WithCreatedAt(t time.Time) func(*imgutil.ImageOptions) {
	return imgutil.WithCreatedAt(t)
}

func WithDefaultPlatform(p imgutil.Platform) func(*imgutil.ImageOptions) {
	return imgutil.WithDefaultPlatform(p)
}

func WithHistory() func(*imgutil.ImageOptions) {
	return imgutil.WithHistory()
}

func WithMediaTypes(m imgutil.MediaTypes) func(*imgutil.ImageOptions) {
	return imgutil.WithMediaTypes(m)
}

func WithPreviousImage(name string) func(*imgutil.ImageOptions) {
	return imgutil.WithPreviousImage(name)
}

type Option func(options *imgutil.IndexOptions) error
type PushOption func(*imgutil.IndexPushOptions) error
type AddOption func(*imgutil.IndexAddOptions) error

// WithKeychain fetches Index from registry with keychain
func WithKeychain(keychain authn.Keychain) Option {
	return imgutil.WithKeychain(keychain)
}

// WithXDGRuntimePath Saves the Index to the '`xdgPath`/manifests'
func WithXDGRuntimePath(xdgPath string) Option {
	return imgutil.WithXDGRuntimePath(xdgPath)
}

// PullInsecure If true, pulls images from insecure registry
func PullInsecure() Option {
	return imgutil.PullInsecure()
}

// Push index to Insecure Registry
func WithInsecure(insecure bool) PushOption {
	return imgutil.WithInsecure(insecure)
}

// Others

// Add all images within the index
func WithAll(all bool) AddOption {
	return imgutil.WithAll(all)
}

// Add a single image from index with given OS
func WithOS(os string) AddOption {
	return imgutil.WithOS(os)
}

// Add a Local image to Index
func WithLocalImage(image imgutil.EditableImage) AddOption {
	return imgutil.WithLocalImage(image)
}

// Add a single image from index with given Architecture
func WithArchitecture(arch string) AddOption {
	return imgutil.WithArchitecture(arch)
}

// Add a single image from index with given Variant
func WithVariant(variant string) AddOption {
	return imgutil.WithVariant(variant)
}

// Add a single image from index with given OSVersion
func WithOSVersion(osVersion string) AddOption {
	return imgutil.WithOSVersion(osVersion)
}

// Add a single image from index with given Features
func WithFeatures(features []string) AddOption {
	return imgutil.WithFeatures(features)
}

// Add a single image from index with given OSFeatures
func WithOSFeatures(osFeatures []string) AddOption {
	return imgutil.WithOSFeatures(osFeatures)
}

// Add a single image from index with given Annotations
func UsingAnnotations(annotations map[string]string) AddOption {
	return imgutil.WithAnnotations(annotations)
}

// If true, Deletes index from local filesystem after pushing to registry
func WithPurge(purge bool) PushOption {
	return imgutil.WithPurge(purge)
}

// Push the Index with given format
func WithTags(tags ...string) PushOption {
	return imgutil.WithTags(tags...)
}

// Push Index with format
func EnsureFormat(format types.MediaType) PushOption {
	return imgutil.UsingFormat(format)
}
