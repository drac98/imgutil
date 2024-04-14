package index_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil/index"
	h "github.com/buildpacks/imgutil/testhelpers"
)

var (
	repoName = "some/index"
)

func TestIndexNew(t *testing.T) {
	spec.Run(t, "IndexNew", testIndexNew, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testIndexNew(t *testing.T, when spec.G, it spec.S) {
	when("#NewIndex", func() {
		var (
			xdgPath string
			err     error
		)
		it.Before(func() {
			// creates the directory to save all the OCI images on disk
			xdgPath, err = os.MkdirTemp("", "image-indexes")
			h.AssertNil(t, err)
		})

		it.After(func() {
			path, err := filepath.Abs(xdgPath)
			h.AssertNil(t, err)

			err = os.RemoveAll(path)
			h.AssertNil(t, err)
		})
		it("should return an error when invalid format passed", func() {
			_, err = index.NewIndex(
				repoName,
				index.WithFormat(types.OCIConfigJSON),
				index.WithXDGRuntimePath(xdgPath),
			)
			h.AssertNotNil(t, err)
		})
		it("Docker ImageIndex", func() {
			_, err = index.NewIndex(
				repoName,
				index.WithFormat(types.DockerManifestList),
				index.WithXDGRuntimePath(xdgPath),
			)
			h.AssertNil(t, err)
		})
		it("OCI ImageIndex", func() {
			_, err = index.NewIndex(
				repoName,
				index.WithFormat(types.OCIImageIndex),
				index.WithXDGRuntimePath(xdgPath),
			)
			h.AssertNil(t, err)
		})
	})
}
