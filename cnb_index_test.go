package imgutil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	"github.com/buildpacks/imgutil/layout"
	h "github.com/buildpacks/imgutil/testhelpers"

	cnbErrs "github.com/buildpacks/imgutil/errors"
)

func TestCNBIndex(t *testing.T) {
	dockerConfigDir, err := os.MkdirTemp("", "test.docker.config.dir")
	h.AssertNil(t, err)
	defer os.RemoveAll(dockerConfigDir)

	dockerRegistry = h.NewDockerRegistry(h.WithAuth(dockerConfigDir))
	dockerRegistry.Start(t)
	defer dockerRegistry.Stop(t)

	os.Setenv("DOCKER_CONFIG", dockerConfigDir)
	defer os.Unsetenv("DOCKER_CONFIG")

	spec.Run(t, "CNBIndex", testCnbIndex, spec.Parallel(), spec.Report(report.Terminal{}))
}

var (
	dockerRegistry *h.DockerRegistry

	repoName               = "busybox-multi-platform"
	NotFoundDigest, _      = name.NewDigest("some/not-found-digest@sha256:4f2fa90168e3ce7022d69a6d67f3f5ae5df1b92d801d2c51ffee341af635adb4")
	FoundDigest, _         = name.NewDigest(repoName + "@sha256:8a4415fb43600953cbdac6ec03c2d96d900bb21f8d78964837dad7f73b9afcdc")
	DockerFoundDigest, _   = name.NewDigest(repoName + "@sha256:742dbd9d350ccda2de4d9990d710b8a5f672a89eda0f2a4ee5403d72cbb02de0")
	FoundDigestOS          = "linux"
	FoundDigestArch        = "arm"
	FoundDigestVariant     = "v7"
	FoundDigestOSVersion   = "1.2.3"
	FoundDigestFeatures    = []string{"feature-3", "feature-4"}
	FoundDigestOSFeatures  = []string{"os-feature-3", "os-feature-4"}
	FoundDigestURLs        = []string{"https://bar.foo"}
	FoundDigestAnnotations = map[string]string{
		"com.docker.official-images.bashbrew.arch": "arm32v7",
		"org.opencontainers.image.base.name":       "scratch",
		"org.opencontainers.image.created":         "2024-02-28T00:44:18Z",
		"org.opencontainers.image.revision":        "185a3f7f21c307b15ef99b7088b228f004ff5f11",
		"org.opencontainers.image.source":          "https://github.com/docker-library/busybox.git",
		"org.opencontainers.image.url":             "https://hub.docker.com/_/busybox",
		"org.opencontainers.image.version":         "1.36.1-glibc",
	}

	// global directory and paths
	testDataDir = filepath.Join("layout", "testdata", "layout")
)

func testCnbIndex(t *testing.T, when spec.G, it spec.S) {
	// Getters

	var (
		idx    *imgutil.CNBIndex
		err    error
		tmpDir string
		// localPath     string
		baseIndexPath string
	)

	it.Before(func() {
		// creates the directory to save all the OCI images on disk
		tmpDir, err = os.MkdirTemp("", "layout-image-indexes")
		h.AssertNil(t, err)

		// image index directory on disk
		baseIndexPath = filepath.Join(testDataDir, repoName)
		// global directory and paths
		testDataDir = filepath.Join("layout", "testdata", "layout")

		index, err := layout.NewIndex(repoName, tmpDir, imgutil.FromBaseImageIndex(baseIndexPath))
		h.AssertNil(t, err)

		idx, err = imgutil.NewCNBIndex(repoName, index.ImageIndex, imgutil.IndexOptions{BaseImageIndexRepoName: baseIndexPath})
		h.AssertNil(t, err)

		// localPath = filepath.Join(tmpDir, repoName)
	})

	it.After(func() {
		err := os.RemoveAll(tmpDir)
		h.AssertNil(t, err)
	})

	when("#OS", func() {
		it("should return OS of the image", func() {
			os, err := idx.OS(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, os, FoundDigestOS)
		})
		it("should return digest not found error", func() {
			os, err := idx.OS(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, os, "")
		})
	})

	when("#Architecture", func() {
		it("should return Architecture of the image", func() {
			arch, err := idx.Architecture(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, arch, FoundDigestArch)
		})
		it("should return digest not found error", func() {
			arch, err := idx.Architecture(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, arch, "")
		})
	})

	when("#Variant", func() {
		it("should return variant of the image", func() {
			variant, err := idx.Variant(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, variant, FoundDigestVariant)
		})
		it("should return digest not found error", func() {
			variant, err := idx.Variant(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, variant, "")
		})
	})

	when("#OSVersion", func() {
		it("should return OSVersion of the image", func() {
			osVersion, err := idx.OSVersion(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, osVersion, FoundDigestOSVersion)
		})
		it("should return digest not found error", func() {
			osVersion, err := idx.OSVersion(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, osVersion, "")
		})
	})

	when("#Features", func() {
		it("should return Features of the image", func() {
			feats, err := idx.Features(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, feats, FoundDigestFeatures)
		})
		it("should return digest not found error", func() {
			feats, err := idx.Features(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, feats, []string(nil))
		})
	})

	when("#OSFeatures", func() {
		it("should return OSFeatures of the image", func() {
			osFeats, err := idx.OSFeatures(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, osFeats, FoundDigestOSFeatures)
		})
		it("should return digest not found error", func() {
			osFeats, err := idx.OSFeatures(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, osFeats, []string(nil))
		})
	})

	when("#URLs", func() {
		it("should return URLs of the image", func() {
			urls, err := idx.URLs(FoundDigest)
			h.AssertNil(t, err)
			h.AssertEq(t, urls, FoundDigestURLs)
		})
		it("should return digest not found error", func() {
			urls, err := idx.URLs(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, urls, []string(nil))
		})
	})

	when("#Annotations", func() {
		when("OCI", func() {
			it("should return Annotations of the image", func() {
				annos, err := idx.Annotations(FoundDigest)
				h.AssertNil(t, err)
				h.AssertEq(t, annos, FoundDigestAnnotations)
			})
		})
		when("Docker", func() {
			it("should not return Annotations of the image", func() {
				annos, err := idx.Annotations(DockerFoundDigest)
				h.AssertNotNil(t, err)
				h.AssertEq(t, annos, map[string]string(nil))
			})
		})
		it("should return digest not found error", func() {
			annos, err := idx.Annotations(NotFoundDigest)
			h.AssertEq(t, err.Error(), cnbErrs.NewDigestNotFoundError(NotFoundDigest.DigestStr()).Error())
			h.AssertEq(t, annos, map[string]string(nil))
		})
	})

	// Setters

	when("#SetOS", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetOS for ImageIndex with digest", func() {})
		it("should #SetOS for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
	})
	when("#SetArchitecture", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetArchitecture for ImageIndex with digest", func() {})
		it("should #SetArchitecture for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
	})
	when("#SetVariant", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetVariant for ImageIndex with digest", func() {})
		it("should #SetVariant for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
	})
	when("#SetOSVersion", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetOSVersion for ImageIndex with digest", func() {})
		it("should #SetOSVersion for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
	})
	when("#SetFeatures", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetFeatures for ImageIndex with digest", func() {})
		it("should #SetFeatures for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
		when("duplicates", func() {
			it("should not add duplicate #Features", func() {})
		})
	})
	when("#SetOSFeatures", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetOSFeatures for ImageIndex with digest", func() {})
		it("should #SetOSFeatures for Image with digest", func() {})
		it("should return an error when no image or index found with given digest", func() {})
		when("duplicates", func() {
			it("should not add duplicate #OSFeatures", func() {})
		})
	})
	when("#SetURLs", func() {
		it("should return an error when image is removed", func() {})
		it("should #SetURLs for ImageIndex with digest", func() {})
		it("should #SetURLs for Image with digest", func() {})
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
		})
		when("docker", func() {
			it("should not #SetAnnotations for ImageIndex with digest", func() {})
			it("should not #SetAnnotations for Image with digest", func() {})
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
		when("#Setters", func() {
			it("should #SetOS for newly added Image with digest", func() {})
			it("should #SetArchitecture for newly added Image with digest", func() {})
			it("should #SetVariant for newly added Image with digest", func() {})
			it("should #SetOSVersion for newly added Image with digest", func() {})
			it("should #SetFeatures for newly added Image with digest", func() {})
			it("should #SetOSFeatures for newly added Image with digest", func() {})
			it("should #SetURLs for newly added Image with digest", func() {})
			it("should #SetAnnotations for newly added Image with digest", func() {})
			when("Docker", func() {
				it("should not #SetAnnotations for newly added Image with digest", func() {})
			})
		})
		when("#Getters", func() {
			it("should get #OS for newly added Image with digest", func() {})
			it("should get #Architecture for newly added Image with digest", func() {})
			it("should get #Variant for newly added Image with digest", func() {})
			it("should get #OSVersion for newly added Image with digest", func() {})
			it("should get #Features for newly added Image with digest", func() {})
			it("should get #OSFeatures for newly added Image with digest", func() {})
			it("should get #URLs for newly added Image with digest", func() {})
			it("should get #Annotations for newly added Image with digest", func() {})
			when("Docker", func() {
				it("should not get #Annotations for newly added Image with digest", func() {})
			})
		})
	})
	when("#Save", func() {
		it("should save ImageIndex", func() {})
		it("should save Annotated ImageIndex", func() {})
		it("should not add images with duplicate digest", func() {})
		it("should save changes in expected manner", func() {})
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
}
