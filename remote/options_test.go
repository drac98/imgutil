package remote_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/fakes"
	"github.com/buildpacks/imgutil/remote"
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

				op := remote.WithKeychain(authn.DefaultKeychain)
				h.AssertNil(t, op(indexOptions))
				h.AssertNotNil(t, indexOptions.KeyChain)
			})
		})
		when("#WithXDGRuntimePath", func() {
			it("should create index from xdgPath", func() {
				var xdgPath = "xdg"
				op := remote.WithXDGRuntimePath(xdgPath)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.XdgPath, xdgPath)
			})
		})
		when("#PullInsecure", func() {
			it("should push to insecure index", func() {
				op := remote.PullInsecure()
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Insecure, true)
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
				op := remote.WithAll(true)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.All, true)
			})
		})
		when("#WithOS", func() {
			it("should add image with OS", func() {
				op := remote.WithOS(OS)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.OS, OS)
			})
		})
		when("#WithLocalImage", func() {
			it("should add local image", func() {
				h.AssertNotEq(t, indexAddOptions.Local, true)
				h.AssertNil(t, indexAddOptions.Image)

				img := fakes.NewImage("image", fakeHash.String(), fakes.NewIdentifier(fakeHash.String()))
				op := remote.WithLocalImage(img)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Local, true)
				h.AssertNotNil(t, indexAddOptions.Image)
			})
		})
		when("#WithArchitecture", func() {
			it("should add image with Architecture", func() {
				op := remote.WithArchitecture(arch)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.Architecture, arch)
			})
		})
		when("#WithVariant", func() {
			it("should add image with Variant", func() {
				op := remote.WithVariant(variant)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.Variant, variant)
			})
		})
		when("#WithOSVersion", func() {
			it("should add image with OSVersion", func() {
				op := remote.WithOSVersion(osVersion)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.OSVersion, osVersion)
			})
		})
		when("#WithFeatures", func() {
			it("should add image with Features", func() {
				op := remote.WithFeatures(features)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.Features, features), true)
			})
		})
		when("#WithOSFeatures", func() {
			it("should add image with OSFeatures", func() {
				op := remote.WithOSFeatures(osFeatures)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.Platform.OSFeatures, osFeatures), true)
			})
		})
		when("#WithAnnotations", func() {
			it("should add image with Annotations", func() {
				op := remote.UsingAnnotations(annotations)
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
				op := remote.WithPurge(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Purge, true)
			})
		})
		when("#WithTags", func() {
			it("should push with tags", func() {
				op := remote.WithTags(features...)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, imgutil.SliceContains(indexPushOptions.Tags, features), true)
			})
		})
		when("#WithInsecure", func() {
			it("should push to insecure", func() {
				op := remote.WithInsecure(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Insecure, true)
			})
		})
		when("#UsingFormat", func() {
			it("should push with format", func() {
				op := remote.EnsureFormat(types.OCIImageIndex)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.OCIImageIndex)
			})
			it("should return an error", func() {
				op := remote.EnsureFormat(types.OCIManifestSchema1)
				h.AssertNotNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.MediaType(""))
			})
		})
	})
}
