package imgutil

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/match"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/partial"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"golang.org/x/sync/errgroup"

	cnbErrs "github.com/buildpacks/imgutil/errors"
)

type CNBIndex struct {
	// required
	v1.ImageIndex // The working Image Index

	// optional
	RegistrySetting
	IndexFormatOptions
	RepoName         string
	XdgPath          string
	KeyChain         authn.Keychain
	removedManifests []v1.Hash
	images           ImageHolder
}

type ImageHolder struct {
	mutex  sync.Mutex
	images map[v1.Hash]v1.Descriptor
}

func (i *ImageHolder) Set(key v1.Hash, value v1.Descriptor) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.images[key] = value
}

func (i *ImageHolder) Get(key v1.Hash) (desc v1.Descriptor, found bool) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for k, v := range i.images {
		if k == key {
			return v, true
		}
	}

	return v1.Descriptor{}, false
}

func (i *ImageHolder) Range(each func(key v1.Hash, value v1.Descriptor) error) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for k, v := range i.images {
		if err := each(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (i *ImageHolder) Size() (size int64) {
	i.Range(func(_ v1.Hash, _ v1.Descriptor) error {
		size++
		return nil
	})
	return size
}

func (h *CNBIndex) getHash(digest name.Digest) (hash v1.Hash, err error) {
	if hash, err = v1.NewHash(digest.Identifier()); err != nil {
		return hash, err
	}

	// if any image is removed with given hash return an error
	for _, h := range h.removedManifests {
		if h == hash {
			return hash, cnbErrs.NewDigestNotFoundError(h.String())
		}
	}

	return hash, nil
}

// OS returns `OS` of an existing Image.
func (h *CNBIndex) OS(digest name.Digest) (os string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return os, err
	}

	getOS := func(desc v1.Descriptor) (os string, err error) {
		if desc.Platform == nil {
			return os, cnbErrs.ErrPlatformUndefined
		}

		if desc.Platform.OS == "" {
			return os, cnbErrs.NewPlatformError(cnbErrs.OS, desc.MediaType, hash.String())
		}

		return desc.Platform.OS, nil
	}

	// return the OS of the added image(using ImageIndex#Add) if found
	if desc, ok := h.images.Get(hash); ok {
		return getOS(desc)
	}

	// check for the digest in the IndexManifest and return `OS` if found
	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return os, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getOS(desc)
		}
	}

	// when no image found with the given digest return an error
	return os, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetOS annotates existing Image by updating `OS` field in IndexManifest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetOS(digest name.Digest, os string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		desc.Platform.OS = os
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	desc.Platform.OS = os
	h.images.Set(hash, *desc)
	return err
}

// Architecture return the Architecture of an Image/Index based on given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) Architecture(digest name.Digest) (arch string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return arch, err
	}

	getArch := func(desc v1.Descriptor) (arch string, err error) {
		if desc.Platform == nil {
			return arch, cnbErrs.ErrPlatformUndefined
		}

		if desc.Platform.Architecture == "" {
			return arch, cnbErrs.NewPlatformError(cnbErrs.Arch, desc.MediaType, hash.String())
		}

		return desc.Platform.Architecture, nil
	}

	if desc, ok := h.images.Get(hash); ok {
		return getArch(desc)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return arch, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getArch(desc)
		}
	}

	return arch, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetArchitecture annotates the `Architecture` of an Image.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetArchitecture(digest name.Digest, arch string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		desc.Platform.Architecture = arch
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	desc.Platform.Architecture = arch
	h.images.Set(hash, *desc)
	return err
}

// Variant return the `Variant` of an Image.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) Variant(digest name.Digest) (osVariant string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return osVariant, err
	}

	getVariant := func(desc v1.Descriptor) (osVariant string, err error) {
		if desc.Platform == nil {
			return osVariant, cnbErrs.ErrPlatformUndefined
		}

		if desc.Platform.Variant == "" {
			return osVariant, cnbErrs.NewPlatformError(cnbErrs.Variant, desc.MediaType, hash.String())
		}

		return desc.Platform.Variant, nil
	}

	if desc, ok := h.images.Get(hash); ok {
		return getVariant(desc)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return osVariant, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getVariant(desc)
		}
	}

	return osVariant, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetVariant annotates the `Variant` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetVariant(digest name.Digest, variant string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		desc.Platform.Variant = variant
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	desc.Platform.Variant = variant
	h.images.Set(hash, *desc)
	return err
}

// OSVersion returns the `OSVersion` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) OSVersion(digest name.Digest) (osVersion string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return osVersion, err
	}

	getOSVersion := func(desc v1.Descriptor) (osVersion string, err error) {
		if desc.Platform == nil {
			return osVersion, cnbErrs.ErrPlatformUndefined
		}

		if desc.Platform.OSVersion == "" {
			return osVersion, cnbErrs.NewPlatformError(cnbErrs.OSVersion, desc.MediaType, hash.String())
		}

		return desc.Platform.OSVersion, nil
	}

	if desc, ok := h.images.Get(hash); ok {
		return getOSVersion(desc)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return osVersion, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getOSVersion(desc)
		}
	}

	return osVersion, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetOSVersion annotates the `OSVersion` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetOSVersion(digest name.Digest, osVersion string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		desc.Platform.OSVersion = osVersion
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	desc.Platform.OSVersion = osVersion
	h.images.Set(hash, *desc)
	return err
}

// Features returns the `Features` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) Features(digest name.Digest) (features []string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return features, err
	}

	getFeatures := func(desc v1.Descriptor) (features []string, err error) {
		if desc.Platform == nil {
			return features, cnbErrs.ErrPlatformUndefined
		}

		if len(desc.Platform.Features) == 0 {
			return features, cnbErrs.NewPlatformError(cnbErrs.Features, desc.MediaType, hash.String())
		}

		return desc.Platform.Features, nil
	}

	if desc, ok := h.images.Get(hash); ok {
		return getFeatures(desc)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return features, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getFeatures(desc)
		}
	}

	return features, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetFeatures annotates the `Features` of an Image with given Digest by appending to existsing Features if any.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetFeatures(digest name.Digest, features []string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		strSet := NewStringSet()
		for _, feat := range append(desc.Platform.Features, features...) {
			strSet.Add(feat)
		}
		desc.Platform.Features = strSet.StringSlice()
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	strSet := NewStringSet()
	for _, feat := range append(desc.Platform.Features, features...) {
		strSet.Add(feat)
	}
	desc.Platform.Features = strSet.StringSlice()
	h.images.Set(hash, *desc)
	return err
}

// OSFeatures returns the `OSFeatures` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) OSFeatures(digest name.Digest) (osFeatures []string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return osFeatures, err
	}

	getOSFeatures := func(desc v1.Descriptor) (osFeatures []string, err error) {
		if desc.Platform == nil {
			return osFeatures, cnbErrs.ErrPlatformUndefined
		}

		if len(desc.Platform.OSFeatures) == 0 {
			return osFeatures, cnbErrs.NewPlatformError(cnbErrs.OSFeatures, desc.MediaType, digest.Identifier())
		}

		return desc.Platform.OSFeatures, nil
	}

	if desc, ok := h.images.Get(hash); ok {
		return getOSFeatures(desc)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return osFeatures, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getOSFeatures(desc)
		}
	}

	return osFeatures, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetOSFeatures annotates the `OSFeatures` of an Image with given Digest by appending to existsing OSFeatures if any.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetOSFeatures(digest name.Digest, osFeatures []string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		strSet := NewStringSet()
		for _, feat := range append(desc.Platform.OSFeatures, osFeatures...) {
			strSet.Add(feat)
		}
		desc.Platform.OSFeatures = strSet.StringSlice()
		h.images.Set(hash, desc)
		return nil
	}

	img, err := h.ImageIndex.Image(hash)
	if err != nil {
		return cnbErrs.NewDigestNotFoundError(digest.Identifier())
	}

	desc, err := partial.Descriptor(img)
	if desc.Platform == nil {
		desc.Platform = &v1.Platform{}
	}

	strSet := NewStringSet()
	for _, feat := range append(desc.Platform.OSFeatures, osFeatures...) {
		strSet.Add(feat)
	}
	desc.Platform.OSFeatures = strSet.StringSlice()
	h.images.Set(hash, *desc)
	return err
}

var (
	DockerMediaTypes = map[types.MediaType]bool{
		types.DockerManifestList:          true,
		types.DockerManifestSchema2:       true,
		types.DockerManifestSchema1:       true,
		types.DockerManifestSchema1Signed: true,
	}
	OCIMediaTypes = map[types.MediaType]bool{
		types.OCIImageIndex:      true,
		types.OCIManifestSchema1: true,
	}
)

// Annotations return the `Annotations` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
// For Docker images and Indexes it returns an error.
func (h *CNBIndex) Annotations(digest name.Digest) (annotations map[string]string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return annotations, err
	}

	getAnnotations := func(annos map[string]string, format types.MediaType) (map[string]string, error) {
		var (
			_, dockerMediaType = DockerMediaTypes[format]
			_, ociMediaType    = OCIMediaTypes[format]
		)
		switch {
		case dockerMediaType:
			// Docker Manifest doesn't support annotations
			return nil, cnbErrs.NewPlatformError(cnbErrs.Annotations, format, digest.Identifier())
		case ociMediaType:
			if len(annos) == 0 {
				return nil, cnbErrs.NewPlatformError(cnbErrs.Annotations, format, digest.Identifier())
			}

			return annos, nil
		default:
			return annos, cnbErrs.NewUnknownMediaTypeError(format)
		}
	}

	if desc, ok := h.images.Get(hash); ok {
		return getAnnotations(desc.Annotations, desc.MediaType)
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return annotations, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			return getAnnotations(desc.Annotations, desc.MediaType)
		}
	}

	return annotations, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetAnnotations annotates the `Annotations` of an Image with given Digest by appending to existing Annotations if any.
//
// Returns an error if no Image/Index found with given Digest.
//
// For Docker images and Indexes it ignores updating Annotations.
func (h *CNBIndex) SetAnnotations(digest name.Digest, annotations map[string]string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		for key, value := range annotations {
			desc.Annotations[key] = value
		}
		h.images.Set(hash, desc)
		return nil
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest != hash {
			continue
		}

		if len(desc.Annotations) == 0 {
			desc.Annotations = make(map[string]string)
		}

		for key, value := range desc.Annotations {
			desc.Annotations[key] = value
		}
		h.images.Set(hash, desc)
		return nil
	}
	return cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// URLs returns the `URLs` of an Image with given Digest.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) URLs(digest name.Digest) (urls []string, err error) {
	hash, err := h.getHash(digest)
	if err != nil {
		return urls, err
	}

	if urls, err = h.getIndexURLs(hash); err == nil || cnbErrs.IsPlatformError(err, hash.String()) {
		return urls, nil
	}

	// OCI ImageIndex can have Docker Images and vice versa
	// Check If it is Platform error, if true return error.
	if urls, _, err = h.getImageURLs(hash); err == nil || cnbErrs.IsPlatformError(err, hash.String()) {
		return urls, err
	}

	return urls, cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// SetURLs annotates the `URLs` of an Image with given Digest by appending to existsing URLs if any.
// Returns an error if no Image/Index found with given Digest.
func (h *CNBIndex) SetURLs(digest name.Digest, urls []string) error {
	hash, err := h.getHash(digest)
	if err != nil {
		return err
	}

	// set the `OS` of an Image added to ImageIndex if found
	if desc, ok := h.images.Get(hash); ok {
		if desc.Platform == nil {
			desc.Platform = &v1.Platform{}
		}

		strSet := NewStringSet()
		for _, feat := range append(desc.URLs, urls...) {
			strSet.Add(feat)
		}
		desc.URLs = strSet.StringSlice()
		h.images.Set(hash, desc)
		return nil
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest != hash {
			continue
		}
		strSet := NewStringSet()
		for _, feat := range append(desc.URLs, urls...) {
			strSet.Add(feat)
		}
		desc.URLs = strSet.StringSlice()
		h.images.Set(hash, desc)
		return nil
	}
	return cnbErrs.NewDigestNotFoundError(digest.Identifier())
}

// Add the ImageIndex from the registry with the given Reference.
//
// If referencing an ImageIndex, will add Platform Specific Image from the Index.
// Use IndexAddOptions to alter behaviour for ImageIndex Reference.
func (h *CNBIndex) Add(name string, ops ...func(*IndexAddOptions) error) error {
	var addOps = &IndexAddOptions{}
	for _, op := range ops {
		if err := op(addOps); err != nil {
			return err
		}
	}

	layoutPath := filepath.Join(h.XdgPath, MakeFileSafeName(h.RepoName))
	path, pathErr := layout.FromPath(layoutPath)
	if addOps.Local {
		if pathErr != nil {
			return cnbErrs.ErrIndexNeedToBeSaved
		}
		img := addOps.Image
		desc, err := partial.Descriptor(img.UnderlyingImage())
		if err != nil {
			return err
		}

		return path.AppendDescriptor(*desc)
	}

	ref, auth, err := referenceForRepoName(h.KeyChain, name, h.Insecure)
	if err != nil {
		return err
	}

	// Fetch Descriptor of the given reference.
	//
	// This call is returns a v1.Descriptor with `Size`, `MediaType`, `Digest` fields only!!
	// This is a lightweight call used for checking MediaType of given Reference
	desc, err := remote.Head(
		ref,
		remote.WithAuth(auth),
	)
	if err != nil {
		return err
	}

	if desc == nil {
		return cnbErrs.ErrManifestUndefined
	}

	switch {
	case desc.MediaType.IsImage():
		// Get the Full Image from remote if the given Reference refers an Image
		img, err := remote.Image(
			ref,
			remote.WithAuth(auth),
		)
		if err != nil {
			return err
		}

		desc, err := partial.Descriptor(img)
		if err != nil {
			return err
		}
		if desc == nil {
			return cnbErrs.ErrManifestUndefined
		}

		if len(desc.Annotations) == 0 {
			desc.Annotations = make(map[string]string)
		}

		for k, v := range addOps.Annotations {
			desc.Annotations[k] = v
		}

		if pathErr != nil {
			if path, err = layout.Write(layoutPath, h.ImageIndex); err != nil {
				return err
			}
		}

		// Append Image to V1.ImageIndex with the Annotations if any
		return path.AppendDescriptor(*desc)
	case desc.MediaType.IsIndex():
		switch {
		case addOps.All:
			idx, err := remote.Index(
				ref,
				remote.WithAuthFromKeychain(h.KeyChain),
				remote.WithTransport(GetTransport(h.Insecure)),
			)
			if err != nil {
				return err
			}
			// Add all the images from Nested ImageIndexes
			if err = h.addAllImages(idx, addOps.Annotations); err != nil {
				return err
			}

			if pathErr != nil {
				// if the ImageIndex is not saved till now for some reason Save the ImageIndex locally to append images
				if err = h.Save(); err != nil {
					return err
				}
			}

			return err
		case !addOps.Platform.Satisfies(v1.Platform{}):

			// Add an Image from the ImageIndex with the given Platform
			return h.addPlatformSpecificImages(ref, addOps.Platform, addOps.Annotations)
		default:
			platform := v1.Platform{
				OS:           runtime.GOOS,
				Architecture: runtime.GOARCH,
			}

			// Add the Image from the ImageIndex with current Device's Platform
			return h.addPlatformSpecificImages(ref, platform, addOps.Annotations)
		}
	default:
		// return an error if the Reference is neither an Image not an Index
		return cnbErrs.NewUnknownMediaTypeError(desc.MediaType)
	}
}

func (h *CNBIndex) addAllImages(idx v1.ImageIndex, annotations map[string]string) error {
	mfest, err := getIndexManifest(idx)
	if err != nil {
		return err
	}

	var errs, _ = errgroup.WithContext(context.Background())
	for _, desc := range mfest.Manifests {
		desc := desc
		errs.Go(func() error {
			return h.addIndexAddendum(annotations, desc, idx)
		})
	}

	return errs.Wait()
}

func (h *CNBIndex) addIndexAddendum(annotations map[string]string, desc v1.Descriptor, idx v1.ImageIndex) error {
	switch {
	case desc.MediaType.IsIndex():
		ii, err := idx.ImageIndex(desc.Digest)
		if err != nil {
			return err
		}

		return h.addAllImages(ii, annotations)
	case desc.MediaType.IsImage():
		img, err := idx.Image(desc.Digest)
		if err != nil {
			return err
		}

		desc, err := partial.Descriptor(img)
		if err != nil {
			return err
		}

		if len(desc.Annotations) == 0 {
			desc.Annotations = make(map[string]string)
		}

		h.images.Set(desc.Digest, *desc)
		return nil
	default:
		return cnbErrs.NewUnknownMediaTypeError(desc.MediaType)
	}
}

func (h *CNBIndex) addPlatformSpecificImages(ref name.Reference, platform v1.Platform, annotations map[string]string) error {
	if platform.OS == "" || platform.Architecture == "" {
		return cnbErrs.ErrInvalidPlatform
	}

	img, err := remote.Image(
		ref,
		remote.WithAuthFromKeychain(h.KeyChain),
		remote.WithTransport(GetTransport(true)),
		remote.WithPlatform(platform),
	)
	if err != nil {
		return err
	}

	desc, err := partial.Descriptor(img)
	if desc == nil {
		return cnbErrs.ErrManifestUndefined
	}

	if len(desc.Annotations) == 0 {
		desc.Annotations = make(map[string]string)
	}

	for k, v := range annotations {
		desc.Annotations[k] = v
	}

	h.images.Set(desc.Digest, *desc)
	return err
}

// Save IndexManifest locally.
// Use it save manifest locally iff the manifest doesn't exist locally before
func (h *CNBIndex) save(layoutPath string) (path layout.Path, err error) {
	// If the ImageIndex is not saved before Save the ImageIndex
	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return path, err
	}

	// Initially write an empty IndexManifest with expected MediaType
	if mfest.MediaType == types.OCIImageIndex {
		if path, err = layout.Write(layoutPath, empty.Index); err != nil {
			return path, err
		}
	} else {
		if path, err = layout.Write(layoutPath, NewEmptyDockerIndex()); err != nil {
			return path, err
		}
	}

	// loop over each digest and append Image/ImageIndex
	for _, d := range mfest.Manifests {
		switch {
		case d.MediaType.IsIndex(), d.MediaType.IsImage():
			if err = path.AppendDescriptor(d); err != nil {
				return path, err
			}
		default:
			return path, cnbErrs.NewUnknownMediaTypeError(d.MediaType)
		}
	}

	return path, nil
}

// Save will locally save the given ImageIndex.
func (h *CNBIndex) Save() error {
	layoutPath := filepath.Join(h.XdgPath, MakeFileSafeName(h.RepoName))
	path, err := layout.FromPath(layoutPath)
	if err != nil {
		// Initially write index to disk then process the changes made to the current index.
		if path, err = h.save(layoutPath); err != nil {
			return err
		}
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return err
	}
	return h.images.Range(func(key v1.Hash, value v1.Descriptor) error {
		for _, desc := range mfest.Manifests {
			if desc.Digest == key {
				if err := path.RemoveDescriptors(match.Digests(key)); err != nil {
					return err
				}
			}
		}
		return path.AppendDescriptor(value)
	})
}

// Push Publishes ImageIndex to the registry assuming every image it referes exists in registry.
//
// It will only push the IndexManifest to registry.
func (h *CNBIndex) Push(ops ...func(*IndexPushOptions) error) error {
	if h.images.Size() != 0 {
		return cnbErrs.ErrIndexNeedToBeSaved
	}

	var pushOps = &IndexPushOptions{}
	for _, op := range ops {
		if err := op(pushOps); err != nil {
			return err
		}
	}

	if pushOps.Format != types.MediaType("") {
		h.ImageIndex = mutate.IndexMediaType(h.ImageIndex, pushOps.Format)
		if err := h.Save(); err != nil {
			return err
		}
	}

	refOps := []name.Option{name.WeakValidation}
	if pushOps.Insecure {
		refOps = append(refOps, name.Insecure)
	}

	ref, err := name.ParseReference(
		h.RepoName,
		refOps...,
	)
	if err != nil {
		return err
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return err
	}

	var taggableIndex = NewTaggableIndex(mfest)
	multiWriteTagables := map[name.Reference]remote.Taggable{
		ref: taggableIndex,
	}
	for _, tag := range pushOps.Tags {
		multiWriteTagables[ref.Context().Tag(tag)] = taggableIndex
	}

	// Note: It will only push IndexManifest, assuming all the images it refers exists in the registry
	err = remote.MultiWrite(
		multiWriteTagables,
		remote.WithAuthFromKeychain(h.KeyChain),
		remote.WithTransport(GetTransport(pushOps.Insecure)),
	)

	if pushOps.Purge {
		return h.Delete()
	}

	return err
}

// Inspect Displays IndexManifest.
func (h *CNBIndex) Inspect() (string, error) {
	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		return "", err
	}

	if h.images.Size() != 0 {
		return "", cnbErrs.ErrIndexNeedToBeSaved
	}

	mfestBytes, err := json.MarshalIndent(mfest, "", "	")
	if err != nil {
		return "", err
	}

	return string(mfestBytes), nil
}

// Removes Image/Index from ImageIndex.
func (h *CNBIndex) Remove(repoName string) (err error) {
	ref, auth, err := referenceForRepoName(h.KeyChain, repoName, h.Insecure)
	if err != nil {
		return err
	}

	hash, err := parseReferenceToHash(ref, auth)
	h.ImageIndex = mutate.RemoveManifests(h.ImageIndex, match.Digests(hash))
	return err
}

// Delete removes ImageIndex from local filesystem if exists.
func (h *CNBIndex) Delete() error {
	layoutPath := filepath.Join(h.XdgPath, MakeFileSafeName(h.RepoName))
	if _, err := os.Stat(layoutPath); err != nil {
		return err
	}

	return os.RemoveAll(layoutPath)
}

func (h *CNBIndex) getIndexURLs(hash v1.Hash) (urls []string, err error) {
	idx, err := h.ImageIndex.ImageIndex(hash)
	if err != nil {
		return urls, err
	}

	desc, err := partial.Descriptor(idx)
	if desc == nil {
		return urls, cnbErrs.ErrManifestUndefined
	}

	return desc.URLs, err
}

func (h *CNBIndex) getImageURLs(hash v1.Hash) (urls []string, format types.MediaType, err error) {
	if desc, ok := h.images.Get(hash); ok {
		if len(desc.URLs) == 0 {
			return urls, desc.MediaType, cnbErrs.NewPlatformError(cnbErrs.URLs, desc.MediaType, hash.String())
		}

		return desc.URLs, desc.MediaType, nil
	}

	mfest, err := getIndexManifest(h.ImageIndex)
	if err != nil {
		// Return Non-Image and Non-Index mediaType
		return urls, types.DockerConfigJSON, err
	}

	for _, desc := range mfest.Manifests {
		if desc.Digest == hash {
			if len(desc.URLs) == 0 {
				return urls, desc.MediaType, cnbErrs.NewPlatformError(cnbErrs.URLs, desc.MediaType, hash.String())
			}

			return desc.URLs, desc.MediaType, nil
		}
	}

	return urls, mfest.MediaType, cnbErrs.NewDigestNotFoundError(hash.String())
}

// Returns v1.Hash from the name.Digest.
// If name.Reference is name.Tag, returns the v1.Hash of the Tag.
func parseReferenceToHash(ref name.Reference, auth authn.Authenticator) (hash v1.Hash, err error) {
	switch v := ref.(type) {
	case name.Tag:
		desc, err := remote.Head(
			v,
			remote.WithAuth(auth),
		)
		if err != nil {
			return hash, err
		}

		if desc == nil {
			return hash, cnbErrs.ErrManifestUndefined
		}

		hash = desc.Digest
	default:
		hash, err = v1.NewHash(v.Identifier())
		if err != nil {
			return hash, err
		}
	}

	return hash, nil
}

func getIndexManifest(ii v1.ImageIndex) (mfest *v1.IndexManifest, err error) {
	mfest, err = ii.IndexManifest()
	if mfest == nil {
		return mfest, cnbErrs.ErrManifestUndefined
	}

	return mfest, err
}

// TODO this method is duplicated from remote.new file
// referenceForRepoName
func referenceForRepoName(keychain authn.Keychain, ref string, insecure bool) (name.Reference, authn.Authenticator, error) {
	var auth authn.Authenticator
	opts := []name.Option{name.WeakValidation}
	if insecure {
		opts = append(opts, name.Insecure)
	}
	r, err := name.ParseReference(ref, opts...)
	if err != nil {
		return nil, nil, err
	}

	auth, err = keychain.Resolve(r.Context().Registry)
	if err != nil {
		return nil, nil, err
	}
	return r, auth, nil
}
