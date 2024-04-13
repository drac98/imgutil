package errors

import (
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/v1/types"
)

var (
	ErrConfigFilePlatformUndefined = errors.New("unable to determine image platform: ConfigFile's platform is nil")
	ErrManifestUndefined           = errors.New("encountered unexpected error while parsing image: manifest or index manifest is nil")
	ErrPlatformUndefined           = errors.New("unable to determine image platform: platform is nil")
	ErrInvalidPlatform             = errors.New("unable to determine image platform: platform's 'OS' or 'Architecture' field is nil")
	ErrConfigFileUndefined         = errors.New("unable to access image configuration: ConfigFile is nil")
	ErrIndexNeedToBeSaved          = errors.New(`unable to perform action: ImageIndex requires local storage before proceeding.
	Please use '#Save()' to save the image index locally before attempting this operation`)
	ErrNoImageFoundWithGivenPlatform = errors.New("no image found for specified platform")
)

func (p *platformUndefined) Error() string {
	return fmt.Sprintf("image %s is undefined for %s ImageIndex (digest: %s)", p.platform, indexMediaType(p.format), p.digest)
}

func (d *digestNotFound) Error() string {
	return fmt.Sprintf(`no image or image index found for digest "%s"`, d.digest)
}

func (f *unknownMediaType) Error() string {
	return fmt.Sprintf("unsupported media type encountered in image: '%s'", f.format)
}

func indexMediaType(format types.MediaType) string {
	switch format {
	case types.DockerManifestList, types.DockerManifestSchema2:
		return "Docker"
	case types.OCIImageIndex, types.OCIManifestSchema1:
		return "OCI"
	default:
		return "UNKNOWN"
	}
}
