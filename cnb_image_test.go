package imgutil_test

import (
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/fakes"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func Test_CNB_Image(t *testing.T) {
	spec.Run(t, "CNB_Image", test_cnb_image, spec.Parallel(), spec.Report(report.Terminal{}))
}

var digest = "sha256:8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"

func test_cnb_image(t *testing.T, when spec.G, it spec.S) {
	when("#MutateConfigFile", func() {
		it("should mutate both configFile and manifest", func() {
			fakeImg := fakes.NewImage("cnbs/some:image", digest, fakes.NewIdentifier(digest))
			cnbImg, err := imgutil.NewCNBImage(imgutil.ImageOptions{
				BaseImage: fakeImg.UnderlyingImage(),
			})
			h.AssertNil(t, err)
			h.AssertNotNil(t, cnbImg)

			err = cnbImg.MutateConfigFile(func(c *v1.ConfigFile) {
				c.OS = os
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

			h.AssertEq(t, config.OS, os)
			h.AssertEq(t, config.Architecture, arch)
			h.AssertEq(t, config.Variant, variant)
			h.AssertEq(t, config.OSVersion, osVersion)
			h.AssertEq(t, imgutil.SliceContains(config.OSFeatures, osFeatures), true)

			h.AssertEq(t, mfest.Config.Platform.OS, os)
			h.AssertEq(t, mfest.Config.Platform.Architecture, arch)
			h.AssertEq(t, mfest.Config.Platform.Variant, variant)
			h.AssertEq(t, mfest.Config.Platform.OSVersion, osVersion)
			h.AssertEq(t, imgutil.SliceContains(mfest.Config.Platform.OSFeatures, osFeatures), true)
		})
	})
}
