package imgutil_test

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/fakes"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func TestOptions(t *testing.T) {
	spec.Run(t, "Options", testOptions, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testOptions(t *testing.T, when spec.G, it spec.S) {
	when("#GetTransport", func() {
		it("should return InsecureTransport", func() {
			roundTripper := imgutil.GetTransport(true)
			trans := roundTripper.(*http.Transport)
			h.AssertEq(t, trans.TLSClientConfig.InsecureSkipVerify, true)
		})
		it("should return SecureTransport", func() {
			roundTripper := imgutil.GetTransport(false)
			trans := roundTripper.(*http.Transport)
			if trans.TLSClientConfig == nil {
				h.AssertEq(t, trans.TLSClientConfig, (*tls.Config)(nil))
			} else {
				h.AssertEq(t, trans.TLSClientConfig.InsecureSkipVerify, false)
			}
		})
	})
	when("#IndexOptions", func() {
		var indexOptions *imgutil.IndexOptions
		it.Before(func() {
			indexOptions = &imgutil.IndexOptions{}
		})
		when("#FromBaseImageIndex", func() {
			it("should set base image index", func() {
				var baseImageIndex = "image-index"
				op := imgutil.FromBaseImageIndex(baseImageIndex)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.BaseImageIndexRepoName, baseImageIndex)
			})
		})
		when("#FromBaseImageIndexInstance", func() {
			it("should set base image index from instance", func() {
				op := imgutil.FromBaseImageIndexInstance(empty.Index)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.BaseIndex, empty.Index)
			})
		})
		when("#WithKeychain", func() {
			it("should auth with given keychain", func() {
				h.AssertNil(t, indexOptions.KeyChain)

				op := imgutil.WithKeychain(authn.DefaultKeychain)
				h.AssertNil(t, op(indexOptions))
				h.AssertNotNil(t, indexOptions.KeyChain)
			})
		})
		when("#WithXDGRuntimePath", func() {
			it("should create index from xdgPath", func() {
				var xdgPath = "xdg"
				op := imgutil.WithXDGRuntimePath(xdgPath)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.XdgPath, xdgPath)
			})
		})
		when("#PullInsecure", func() {
			it("should push to insecure index", func() {
				op := imgutil.PullInsecure()
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Insecure, true)
			})
		})
		when("#WithFormat", func() {
			it("should create index with format", func() {
				op := imgutil.WithFormat(types.OCIImageIndex)
				h.AssertNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Format, types.OCIImageIndex)
			})
			it("should return an error", func() {
				op := imgutil.WithFormat(types.DockerConfigJSON)
				h.AssertNotNil(t, op(indexOptions))
				h.AssertEq(t, indexOptions.Format, types.MediaType(""))
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
				op := imgutil.WithAll(true)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.All, true)
			})
		})
		when("#WithOS", func() {
			it("should add image with OS", func() {
				op := imgutil.WithOS(OS)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.OS, OS)
			})
		})
		when("#WithLocalImage", func() {
			it("should add local image", func() {
				h.AssertNotEq(t, indexAddOptions.Local, true)
				h.AssertNil(t, indexAddOptions.Image)

				img := fakes.NewImage("image", fakeHash.String(), fakes.NewIdentifier(fakeHash.String()))
				op := imgutil.WithLocalImage(img)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Local, true)
				h.AssertNotNil(t, indexAddOptions.Image)
			})
		})
		when("#WithArchitecture", func() {
			it("should add image with Architecture", func() {
				op := imgutil.WithArchitecture(arch)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.Architecture, arch)
			})
		})
		when("#WithVariant", func() {
			it("should add image with Variant", func() {
				op := imgutil.WithVariant(variant)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.Variant, variant)
			})
		})
		when("#WithOSVersion", func() {
			it("should add image with OSVersion", func() {
				op := imgutil.WithOSVersion(osVersion)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, indexAddOptions.Platform.OSVersion, osVersion)
			})
		})
		when("#WithFeatures", func() {
			it("should add image with Features", func() {
				op := imgutil.WithFeatures(features)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.Features, features), true)
			})
		})
		when("#WithOSFeatures", func() {
			it("should add image with OSFeatures", func() {
				op := imgutil.WithOSFeatures(osFeatures)
				h.AssertNil(t, op(indexAddOptions))
				h.AssertEq(t, imgutil.SliceContains(indexAddOptions.Platform.OSFeatures, osFeatures), true)
			})
		})
		when("#WithAnnotations", func() {
			it("should add image with Annotations", func() {
				op := imgutil.WithAnnotations(annotations)
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
				op := imgutil.WithPurge(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Purge, true)
			})
		})
		when("#WithTags", func() {
			it("should push with tags", func() {
				op := imgutil.WithTags(features...)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, imgutil.SliceContains(indexPushOptions.Tags, features), true)
			})
		})
		when("#WithInsecure", func() {
			it("should push to insecure", func() {
				op := imgutil.WithInsecure(true)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Insecure, true)
			})
		})
		when("#UsingFormat", func() {
			it("should push with format", func() {
				op := imgutil.UsingFormat(types.OCIImageIndex)
				h.AssertNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.OCIImageIndex)
			})
			it("should return an error", func() {
				op := imgutil.UsingFormat(types.OCIManifestSchema1)
				h.AssertNotNil(t, op(indexPushOptions))
				h.AssertEq(t, indexPushOptions.Format, types.MediaType(""))
			})
		})
	})
}
