package imgutil

import (
	"encoding/json"
	"slices"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/partial"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/pkg/errors"
)

func GetConfigFile(image v1.Image) (*v1.ConfigFile, error) {
	configFile, err := image.ConfigFile()
	if err != nil {
		return nil, err
	}
	if configFile == nil {
		return nil, errors.New("missing config file")
	}
	return configFile, nil
}

func GetManifest(image v1.Image) (*v1.Manifest, error) {
	manifest, err := image.Manifest()
	if err != nil {
		return nil, err
	}
	if manifest == nil {
		return nil, errors.New("missing manifest")
	}
	return manifest, nil
}

// TaggableIndex any ImageIndex with RawManifest method.
type TaggableIndex struct {
	*v1.IndexManifest
}

// RawManifest returns the bytes of IndexManifest.
func (t *TaggableIndex) RawManifest() ([]byte, error) {
	return json.Marshal(t.IndexManifest)
}

// Digest returns the Digest of the IndexManifest if present.
// Else generate a new Digest.
func (t *TaggableIndex) Digest() (v1.Hash, error) {
	if t.IndexManifest.Subject != nil && t.IndexManifest.Subject.Digest != (v1.Hash{}) {
		return t.IndexManifest.Subject.Digest, nil
	}

	return partial.Digest(t)
}

// MediaType returns the MediaType of the IndexManifest.
func (t *TaggableIndex) MediaType() (types.MediaType, error) {
	return t.IndexManifest.MediaType, nil
}

// Size returns the Size of IndexManifest if present.
// Calculate the Size of empty.
func (t *TaggableIndex) Size() (int64, error) {
	if t.IndexManifest.Subject != nil && t.IndexManifest.Subject.Size != 0 {
		return t.IndexManifest.Subject.Size, nil
	}

	return partial.Size(t)
}

type StringSet struct {
	items map[string]bool
}

func NewStringSet() *StringSet {
	return &StringSet{items: make(map[string]bool)}
}

func (s *StringSet) Add(str string) {
	if s == nil {
		s = &StringSet{items: make(map[string]bool)}
	}

	s.items[str] = true
}

func (s *StringSet) Remove(str string) {
	if s == nil {
		s = &StringSet{items: make(map[string]bool)}
	}

	s.items[str] = false
}

func (s *StringSet) StringSlice() (slice []string) {
	if s == nil {
		s = &StringSet{items: make(map[string]bool)}
	}

	for i, ok := range s.items {
		if ok {
			slice = append(slice, i)
		}
	}

	return slice
}

func MapContains(src, target map[string]string) bool {
	for targetKey, targetValue := range target {
		if value := src[targetKey]; targetValue == value {
			continue
		}
		return false
	}
	return true
}

func SliceContains(src, target []string) bool {
	for _, value := range target {
		if ok := slices.Contains(src, value); !ok {
			return false
		}
	}
	return true
}

// MakeFileSafeName Change a reference name string into a valid file name
// Ex: cnbs/sample-package:hello-multiarch-universe
// to cnbs_sample-package-hello-multiarch-universe
func MakeFileSafeName(ref string) string {
	fileName := strings.ReplaceAll(ref, ":", "-")
	return strings.ReplaceAll(fileName, "/", "_")
}

func NewEmptyDockerIndex() v1.ImageIndex {
	idx := empty.Index
	return mutate.IndexMediaType(idx, types.DockerManifestList)
}

func ValidateRepoName(repoName string, o *IndexOptions) (err error) {
	if o.Insecure {
		_, err = name.ParseReference(repoName, name.Insecure, name.WeakValidation)
	} else {
		_, err = name.ParseReference(repoName, name.WeakValidation)
	}
	o.BaseImageIndexRepoName = repoName
	return err
}
