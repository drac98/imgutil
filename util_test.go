package imgutil_test

import (
	"encoding/json"
	"testing"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func TestUtils(t *testing.T) {
	spec.Run(t, "Utils", testUtils, spec.Parallel(), spec.Report(report.Terminal{}))
}

var (
	emptyImage  = empty.Image
	OS          = "some-os"
	arch        = "some-arch"
	variant     = "some-variant"
	osVersion   = "some-os-version"
	features    = []string{"feature1", "feature2"}
	osFeatures  = []string{"osFeature1", "osFeature2"}
	urls        = []string{"url1", "url2"}
	annotations = map[string]string{"key1": "value1", "key2": "value2"}
	fakeHash    = v1.Hash{Algorithm: "sha256", Hex: "8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"}
)

func testUtils(t *testing.T, when spec.G, it spec.S) {
	when("#TaggableIndex", func() {
		var (
			taggableIndex *imgutil.TaggableIndex
			amd64Hash, _  = v1.NewHash("sha256:b9d056b83bb6446fee29e89a7fcf10203c562c1f59586a6e2f39c903597bda34")
			armv6Hash, _  = v1.NewHash("sha256:0bcc1b827b855c65eaf6e031e894e682b6170160b8a676e1df7527a19d51fb1a")
			indexManifest = v1.IndexManifest{
				SchemaVersion: 2,
				MediaType:     types.OCIImageIndex,
				Annotations: map[string]string{
					"test-key": "test-value",
				},
				Manifests: []v1.Descriptor{
					{
						MediaType: types.OCIManifestSchema1,
						Size:      832,
						Digest:    amd64Hash,
						Platform: &v1.Platform{
							OS:           "linux",
							Architecture: "amd64",
						},
					},
					{
						MediaType: types.OCIManifestSchema1,
						Size:      926,
						Digest:    armv6Hash,
						Platform: &v1.Platform{
							OS:           "linux",
							Architecture: "arm",
							OSVersion:    "v6",
						},
					},
				},
			}
		)
		it.Before(func() {
			taggableIndex = imgutil.NewTaggableIndex(&indexManifest)
		})
		it("should return RawManifest in expected format", func() {
			mfestBytes, err := taggableIndex.RawManifest()
			h.AssertNil(t, err)

			expectedMfestBytes, err := json.Marshal(indexManifest)
			h.AssertNil(t, err)

			h.AssertEq(t, mfestBytes, expectedMfestBytes)
		})
		it("should return expected digest", func() {
			digest, err := taggableIndex.Digest()
			h.AssertNil(t, err)
			h.AssertEq(t, digest.String(), "sha256:2375c0dfd06dd51b313fd97df5ecf3b175380e895287dd9eb2240b13eb0b5703")
		})
		it("should return expected size", func() {
			size, err := taggableIndex.Size()
			h.AssertNil(t, err)
			h.AssertEq(t, size, int64(547))
		})
		it("should return expected media type", func() {
			format, err := taggableIndex.MediaType()
			h.AssertNil(t, err)
			h.AssertEq(t, format, indexManifest.MediaType)
		})
	})
	when("#StringSet", func() {
		when("#NewStringSet", func() {
			it("should return not nil StringSet instance", func() {
				stringSet := imgutil.NewStringSet()
				h.AssertNotNil(t, stringSet)
				h.AssertEq(t, stringSet.StringSlice(), []string(nil))
			})
		})

		when("#Add", func() {
			var (
				stringSet *imgutil.StringSet
			)
			it.Before(func() {
				stringSet = imgutil.NewStringSet()
			})
			it("should add items", func() {
				item := "item1"
				stringSet.Add(item)

				h.AssertEq(t, stringSet.StringSlice(), []string{item})
			})
			it("should return added items", func() {
				items := []string{"item1", "item2", "item3"}
				for _, item := range items {
					stringSet.Add(item)
				}
				h.AssertEq(t, len(stringSet.StringSlice()), 3)
				h.AssertContains(t, stringSet.StringSlice(), items...)
			})
			it("should not support duplicates", func() {
				stringSet := imgutil.NewStringSet()
				item1 := "item1"
				item2 := "item2"
				items := []string{item1, item2, item1}
				for _, item := range items {
					stringSet.Add(item)
				}
				h.AssertEq(t, len(stringSet.StringSlice()), 2)
				h.AssertContains(t, stringSet.StringSlice(), []string{item1, item2}...)
			})
		})

		when("#Remove", func() {
			var (
				stringSet *imgutil.StringSet
				item      string
			)
			it.Before(func() {
				stringSet = imgutil.NewStringSet()
				item = "item1"
				stringSet.Add(item)
				h.AssertEq(t, stringSet.StringSlice(), []string{item})
			})
			it("should remove item", func() {
				stringSet.Remove(item)
				h.AssertEq(t, stringSet.StringSlice(), []string(nil))
			})
		})
	})
	when("#NewEmptyDockerIndex", func() {
		it("should return an empty docker index", func() {
			idx := imgutil.NewEmptyDockerIndex()
			h.AssertNotNil(t, idx)

			digest, err := idx.Digest()
			h.AssertNil(t, err)
			h.AssertNotEq(t, digest, v1.Hash{})

			format, err := idx.MediaType()
			h.AssertNil(t, err)
			h.AssertEq(t, format, types.DockerManifestList)
		})
	})
	when("#GetConfigFile", func() {
		it("should return ConfigFile", func() {
			config, err := imgutil.GetConfigFile(emptyImage)
			h.AssertNotNil(t, config)
			h.AssertNil(t, err)
		})
	})
	when("#GetManifest", func() {
		it("should return Manifest", func() {
			mfest, err := imgutil.GetManifest(emptyImage)
			h.AssertNotNil(t, mfest)
			h.AssertNil(t, err)
		})
	})
	when("#MapContains", func() {
		var (
			sampleMap    = map[string]string{"key1": "value1", "key2": "value2"}
			mapExists    = map[string]string{"key2": "value2"}
			mapNotExists = map[string]string{"key1": "value2"}
		)
		it("should return true", func() {
			h.AssertEq(t, imgutil.MapContains(sampleMap, mapExists), true)
		})
		it("should return false", func() {
			h.AssertEq(t, imgutil.MapContains(sampleMap, mapNotExists), false)
		})
	})
	when("#SliceContains", func() {
		var (
			sampleSlice    = []string{"item1", "item2", "item3"}
			sliceExists    = []string{"item2", "item1"}
			sliceNotExists = []string{"item1", "item2", "item3", "item4"}
		)
		it("should return true", func() {
			h.AssertEq(t, imgutil.SliceContains(sampleSlice, sliceExists), true)
		})
		it("should return false", func() {
			h.AssertEq(t, imgutil.SliceContains(sampleSlice, sliceNotExists), false)
		})
	})
	when("#MakeFileSafeName", func() {
		var (
			imageNameWithTag         = "cnbs/sample:image"
			expectedImageNameWithTag = "cnbs_sample-image"

			imageNameWithDigest         = "cnbs/sample@sha256:8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"
			expectedImageNameWithDigest = "cnbs_sample@sha256-8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"
		)
		it("should MakeFileNameSafe(Tag)", func() {
			h.AssertEq(t, imgutil.MakeFileSafeName(imageNameWithTag), expectedImageNameWithTag)
		})
		it("should MakeFileNameSafe(Digest)", func() {
			h.AssertEq(t, imgutil.MakeFileSafeName(imageNameWithDigest), expectedImageNameWithDigest)
		})
	})
	// when("#ValidateRepoName", func() {
	// 	it("should not return error when insecure", func() {
	// 		h.AssertNil(t, imgutil.ValidateRepoName(repoName, &imgutil.IndexOptions{IndexRemoteOptions: imgutil.IndexRemoteOptions{Insecure: true}}))
	// 	})
	// 	it("should not return an error when secure", func() {
	// 		h.AssertNil(t, imgutil.ValidateRepoName(repoName, &imgutil.IndexOptions{IndexRemoteOptions: imgutil.IndexRemoteOptions{Insecure: false}}))
	// 	})
	// })
}
