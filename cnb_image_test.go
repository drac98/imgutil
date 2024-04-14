package imgutil_test

import (
	"bytes"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/fakes"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func TestCNBImage(t *testing.T) {
	spec.Run(t, "CNBImage", testCnbImage, spec.Parallel(), spec.Report(report.Terminal{}))
}

var digest = "sha256:8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"

func testCnbImage(t *testing.T, when spec.G, it spec.S) {
	when("#MutateConfigFile", func() {
		it("should mutate both configFile and manifest", func() {
			fakeImg := fakes.NewImage("cnbs/some:image", digest, fakes.NewIdentifier(digest))
			cnbImg, err := imgutil.NewCNBImage(imgutil.ImageOptions{
				BaseImage: fakeImg.UnderlyingImage(),
			})
			h.AssertNil(t, err)
			h.AssertNotNil(t, cnbImg)

			err = cnbImg.MutateConfigFile(func(c *v1.ConfigFile) {
				c.OS = OS
				c.Architecture = arch
				c.Variant = variant
				c.OSVersion = osVersion
				c.OSFeatures = osFeatures
			})
			h.AssertNil(t, err)

			config, err := cnbImg.ConfigFile()
			h.AssertNil(t, err)
			h.AssertNotNil(t, config)

			mfest, err := cnbImg.Manifest()
			h.AssertNil(t, err)
			h.AssertNotNil(t, mfest)

			h.AssertEq(t, config.OS, OS)
			h.AssertEq(t, config.Architecture, arch)
			h.AssertEq(t, config.Variant, variant)
			h.AssertEq(t, config.OSVersion, osVersion)
			h.AssertEq(t, imgutil.SliceContains(config.OSFeatures, osFeatures), true)
		})
	})
	when("#V1Image", func() {
		var (
			config = v1.Descriptor{
				MediaType:   types.DockerConfigJSON,
				Size:        256,
				Digest:      v1.Hash{Hex: "8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2", Algorithm: "sha256"},
				URLs:        urls,
				Annotations: annotations,
				Platform: &v1.Platform{
					OS:           OS,
					Architecture: arch,
					Variant:      variant,
					OSVersion:    osVersion,
					Features:     features,
					OSFeatures:   osFeatures,
				},
			}
			v1Image = imgutil.NewV1Image(emptyImage, config)
		)
		it("should return expected Manifest", func() {
			mfest, err := v1Image.Manifest()
			h.AssertNil(t, err)
			h.AssertNotNil(t, mfest)

			h.AssertEq(t, mfest.Config, config)
		})
		it("should return expected RawManifest", func() {
			rawMfest, err := v1Image.RawManifest()
			h.AssertNil(t, err)
			h.AssertNotEq(t, len(rawMfest), 0)

			mfest, err := v1.ParseManifest(bytes.NewReader(rawMfest))
			h.AssertNil(t, err)
			h.AssertEq(t, mfest.Config, config)
		})
	})
	when("Getters and Setters", func() {
		var (
			cnbImage, err = imgutil.NewCNBImage(imgutil.ImageOptions{
				BaseImage: emptyImage,
				Platform: imgutil.Platform{
					OS:           "linux",
					Architecture: "arm",
				},
			})
		)
		it.Before(func() {
			h.AssertNil(t, err)
			h.AssertNotNil(t, cnbImage)
		})
		it("should set and get Features", func() {
			h.AssertNil(t, cnbImage.SetFeatures(features))
			feats, err := cnbImage.Features()
			h.AssertNil(t, err)
			h.AssertEq(t, imgutil.SliceContains(feats, features), true)
		})
		it("should set and get OSFeatures", func() {
			h.AssertNil(t, cnbImage.SetOSFeatures(osFeatures))
			osFeats, err := cnbImage.OSFeatures()
			h.AssertNil(t, err)
			h.AssertEq(t, imgutil.SliceContains(osFeats, osFeatures), true)
		})
		it("should set and get URLs", func() {
			h.AssertNil(t, cnbImage.SetURLs(urls))
			u, err := cnbImage.URLs()
			h.AssertNil(t, err)
			h.AssertEq(t, imgutil.SliceContains(u, urls), true)
		})
		it("should set and get Annotations", func() {
			annos, err := cnbImage.Annotations()
			h.AssertNotNil(t, err)
			h.AssertEq(t, annos, map[string]string(nil))

			h.AssertNil(t, cnbImage.SetAnnotations(annotations))
			annos, err = cnbImage.Annotations()
			h.AssertNil(t, err)
			h.AssertEq(t, imgutil.MapContains(annos, annotations), true)
		})
	})
}
