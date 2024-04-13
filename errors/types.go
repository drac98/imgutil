package errors

import "github.com/google/go-containerregistry/pkg/v1/types"

type Platform string

const (
	OS          Platform = "os"
	Arch        Platform = "architecture"
	Variant     Platform = "variant"
	OSVersion   Platform = "os-version"
	Features    Platform = "features"
	OSFeatures  Platform = "os-features"
	URLs        Platform = "urls"
	Annotations Platform = "annotations"
)

type platformUndefined struct {
	format   types.MediaType
	digest   string
	platform Platform
}

type digestNotFound struct {
	digest string
}

type unknownMediaType struct {
	format types.MediaType
}
