package imgutil_test

import (
	// "bytes"
	"testing"

	// v1 "github.com/google/go-containerregistry/pkg/v1"
	// "github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	// "github.com/buildpacks/imgutil"
	// "github.com/buildpacks/imgutil/fakes"
	// h "github.com/buildpacks/imgutil/testhelpers"
)

func TestCNBIndex(t *testing.T) {
	spec.Run(t, "CNB_Index", testCnbIndex, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testCnbIndex(t *testing.T, when spec.G, it spec.S) {
	when("#CNBIndex", func() {
		// Getters

		when("#OS", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated OS", func() {})
			it("should return OS of newly added image", func() {})
			it("should return OS of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#Architecture", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated Architecture", func() {})
			it("should return Architecture of newly added image", func() {})
			it("should return Architecture of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#Variant", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated Variant", func() {})
			it("should return variant of newly added image", func() {})
			it("should return variant of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#OSVersion", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated OSVersion", func() {})
			it("should return OSVersion of newly added image", func() {})
			it("should return OSVersion of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#Features", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated Features", func() {})
			it("should return Features of newly added image", func() {})
			it("should return Features of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#OSFeatures", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated OSFeatures", func() {})
			it("should return OSFeatures of newly added image", func() {})
			it("should return OSFeatures of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#URLs", func() {
			it("should return an error when image is removed", func() {})
			it("should return latest annotated URLs", func() {})
			it("should return URLs of newly added image", func() {})
			it("should return URLs of the image", func() {})
			it("shoudl return digest not found error", func() {})
		})
		when("#Annotations", func() {
			it("should return an error when image is removed", func() {})
			when("OCI", func() {
				it("should return latest annotated Annotations", func() {})
				it("should return Annotations of newly added image", func() {})
				it("should return Annotations of the image", func() {})
			})
			when("Docker", func() {
				it("should not return latest annotated Annotations", func() {})
				it("should not return Annotations of newly added image", func() {})
				it("should not return Annotations of the image", func() {})
			})
			it("shoudl return digest not found error", func() {})
		})

		// Setters

		when("#SetOS", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetOS for ImageIndex with digest", func() {})
			it("should #SetOS for Image with digest", func() {})
			it("should #SetOS for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
		})
		when("#SetArchitecture", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetArchitecture for ImageIndex with digest", func() {})
			it("should #SetArchitecture for Image with digest", func() {})
			it("should #SetArchitecture for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
		})
		when("#SetVariant", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetVariant for ImageIndex with digest", func() {})
			it("should #SetVariant for Image with digest", func() {})
			it("should #SetVariant for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
		})
		when("#SetOSVersion", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetOSVersion for ImageIndex with digest", func() {})
			it("should #SetOSVersion for Image with digest", func() {})
			it("should #SetOSVersion for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
		})
		when("#SetFeatures", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetFeatures for ImageIndex with digest", func() {})
			it("should #SetFeatures for Image with digest", func() {})
			it("should #SetFeatures for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
			when("duplicates", func() {
				it("should not add duplicate #Features", func() {})
			})
		})
		when("#SetOSFeatures", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetOSFeatures for ImageIndex with digest", func() {})
			it("should #SetOSFeatures for Image with digest", func() {})
			it("should #SetOSFeatures for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
			when("duplicates", func() {
				it("should not add duplicate #OSFeatures", func() {})
			})
		})
		when("#SetURLs", func() {
			it("should return an error when image is removed", func() {})
			it("should #SetURLs for ImageIndex with digest", func() {})
			it("should #SetURLs for Image with digest", func() {})
			it("should #SetURLs for newly added Image with digest", func() {})
			it("should return an error when no image or index found with given digest", func() {})
			when("duplicates", func() {
				it("should not add duplicate #URLs", func() {})
			})
		})
		when("#SetAnnotations", func() {
			it("should return an error when image is removed", func() {})
			when("OCI", func() {
				it("should #SetAnnotations for ImageIndex with digest", func() {})
				it("should #SetAnnotations for Image with digest", func() {})
				it("should #SetAnnotations for newly added Image with digest", func() {})
			})
			when("docker", func() {
				it("should not #SetAnnotations for ImageIndex with digest", func() {})
				it("should not #SetAnnotations for Image with digest", func() {})
				it("should not #SetAnnotations for newly added Image with digest", func() {})
			})
			it("should return an error when no image or index found with given digest", func() {})
			when("duplicates", func() {
				it("should not add duplicate #Annotations", func() {})
			})
		})

		// MISC

		when("#Add", func() {
			when("Local", func() {
				it("should return an error if Index not saved", func() {})
				it("should append image", func() {})
			})
			when("Single Image", func() {
				it("should append given annotations", func() {})
				it("should add image to index", func() {})
			})
			when("Index", func() {
				when("All images", func() {
					it("should add all images", func() {})
					it("should add annotations to all images", func() {})
				})
				when("Target Specific Image", func() {
					it("should add an image with given platform", func() {})
					it("should return an error when no image is found with given platform", func() {})
					it("should add annotations to Target Specific Image", func() {})
				})
				when("No Target Specific Options provided", func() {
					it("should add Platform Specific Image", func() {})
					it("should return an error when Platform Specific Image not found", func() {})
					it("should add Annotations to Platform Specific Image", func() {})
				})
			})
		})
		when("#Save", func() {
			it("should save ImageIndex", func() {})
			it("should save Annotated ImageIndex", func() {})
			it("should not add images with duplicate digest", func() {})
		})
		when("#Push", func() {
			it("should return an error when annotated changes not saved", func() {})
			it("should return an error when Invalid Push Format provided", func() {})
			it("should push with current format when Format not specified", func() {})
			it("should pruge index afer push when requested", func() {})
		})
		when("#Inspect", func() {
			it("should return an error when annotated changes not saved", func() {})
			it("should output index in expected format", func() {})
		})
		when("#Remove", func() {
			it("should remove specified image", func() {})
			it("should remove all images when index specified", func() {})
			it("should return an error when image not found", func() {})
		})
		when("#Delete", func() {
			it("should delete local index", func() {})
			it("should return an error when index not exists", func() {})
		})
	})
}
