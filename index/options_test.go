package index_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/fakes"
	"github.com/buildpacks/imgutil/index"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func TestIndexOptions(t *testing.T) {
	spec.Run(t, "IndexOptions", testIndexOptions, spec.Parallel(), spec.Report(report.Terminal{}))
}

var (
	OS          = "some-os"
	arch        = "some-arch"
	variant     = "some-variant"
	osVersion   = "some-os-version"
	features    = []string{"feature1", "feature2"}
	osFeatures  = []string{"osFeature1", "osFeature2"}
	annotations = map[string]string{"key1": "value1", "key2": "value2"}
	fakeHash    = v1.Hash{Algorithm: "sha256", Hex: "8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"}
)

func testIndexOptions(t *testing.T, when spec.G, it spec.S) {
	when("#IndexOptions", func() {
		var indexOptions *imgutil.IndexOptions
		it.Before(func() {
			indexOptions = &imgutil.IndexOptions{}
		})
		when("#WithKeychain", func() {
			it("should auth with given keychain", func() {
				h.AssertNil(t, indexOptions.KeyChain)

				op := index.WithKeychain(authn.DefaultKeychain)
				h.AssertNil(t, op(indexOptions))
				h.AssertNotNil(t, indexOptions.KeyChain)
			})
		})
		when("#WithXDGRuntimePath", func() {
			it("should create index from xdgPath", func() {
				var xdgPath = "xdg"
				op := index.WithXDGRuntimePath(xdgPath)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.XdgPath, xdgPath)
			})
		})
		when("#PullInsecure", func() {
			it("should push to insecure index", func() {
				op := index.PullInsecure()
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Insecure, true)
			})
		})
		when("#WithFormat", func() {
			it("should support index type", func() {
				op := index.WithFormat(types.OCIImageIndex)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Format, types.OCIImageIndex)
			})
			it("should return an error", func() {
				op := index.WithFormat(types.OCIConfigJSON)
				h.AssertNotNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Format, types.MediaType(""))
			})
			it("should support image type", func() {
				op := index.WithFormat(types.DockerManifestSchema2)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Format, types.DockerManifestSchema2)
			})
		})
	})
	when("#IndexAddOptions", func() {
		var indexAddOptions *imgutil.IndexAddOptions
		it.Before(func() {
			indexAddOptions = &imgutil.IndexAddOptions{}
		})
		when("#WithAll", func() {
			it("should add all images", func() {
				op := index.WithAll(true)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.All, true)
			})
		})
		when("#WithOS", func() {
			it("should add image with OS", func() {
				op := index.WithOS(OS)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.OS, OS)
			})
		})
		when("#WithLocalImage", func() {
			it("should add local image", func() {
				h.AssertNotEq(t, indexAddOptions.Local, true)
				h.AssertNil(t, indexAddOptions.Image)

				img := fakes.NewImage("image", fakeHash.String(), fakes.NewIdentifier(fakeHash.String()))
				op := index.WithLocalImage(img)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Local, true)
				h.AssertNotNil(t, indexAddOptions.Image)
			})
		})
		when("#WithArchitecture", func() {
			it("should add image with Architecture", func() {
				op := index.WithArchitecture(arch)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Arch, arch)
			})
		})
		when("#WithVariant", func() {
			it("should add image with Variant", func() {
				op := index.WithVariant(variant)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Variant, variant)
			})
		})
		when("#WithOSVersion", func() {
			it("should add image with OSVersion", func() {
				op := index.WithOSVersion(osVersion)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.OSVersion, osVersion)
			})
		})
		when("#WithFeatures", func() {
			it("should add image with Features", func() {
				op := index.WithFeatures(features)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.Features, features), true)
			})
		})
		when("#WithOSFeatures", func() {
			it("should add image with OSFeatures", func() {
				op := index.WithOSFeatures(osFeatures)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.OSFeatures, osFeatures), true)
			})
		})
		when("#WithAnnotations", func() {
			it("should add image with Annotations", func() {
				op := index.WithAnnotations(annotations)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.MapContains(indexAddOptions.Annotations, annotations), true)
			})
		})
	})
	when("#IndexPushOptions", func() {
		var indexPushOptions *imgutil.IndexPushOptions
		it.Before(func() {
			indexPushOptions = &imgutil.IndexPushOptions{}
		})
		when("#WithPurge", func() {
			it("should purge index", func() {
				op := index.WithPurge(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Purge, true)
			})
		})
		when("#WithTags", func() {
			it("should push with tags", func() {
				op := index.WithTags(features...)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, imgutil.SliceContains(indexPushOptions.Tags, features), true)
			})
		})
		when("#WithInsecure", func() {
			it("should push to insecure", func() {
				op := index.WithInsecure(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Insecure, true)
			})
		})
		when("#UsingFormat", func() {
			it("should push with format", func() {
				op := index.UsingFormat(types.OCIImageIndex)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.OCIImageIndex)
			})
			it("should return an error", func() {
				op := index.UsingFormat(types.OCIManifestSchema1)
				h.AssertNotNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.MediaType(""))
			})
		})
	})
}
