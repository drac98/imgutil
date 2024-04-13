package imgutil_test

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
	h "github.com/buildpacks/imgutil/testhelpers"
)

func TestNew(t *testing.T) {
	spec.Run(t, "New", testNew, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testNew(t *testing.T, when spec.G, it spec.S) {
	when("#NewAnnotate", func() {
		it("Annotate should not panic", func() {
			annotate := imgutil.NewAnnotate()
			h.AssertNotNil(t, annotate.Instance)
		})
	})
	when("#NewCNBIndex", func() {
		it("should return New Instance", func() {
			opts := imgutil.IndexOptions{
				IndexFormatOptions: imgutil.IndexFormatOptions{
					Format: types.OCIImageIndex,
				},
				IndexRemoteOptions: imgutil.IndexRemoteOptions{
					Insecure: true,
				},
				KeyChain:               authn.DefaultKeychain,
				XdgPath:                os,
				BaseImageIndexRepoName: arch,
			}

			idx := empty.Index
			h.AssertNotNil(t, idx)

			cnbIdx, err := imgutil.NewCNBIndex(arch, idx, opts)
			h.AssertNil(t, err)
			h.AssertEq(t, cnbIdx.ImageIndex, idx)

			h.AssertEq(t, cnbIdx.Format, opts.Format)
			h.AssertEq(t, cnbIdx.Insecure, opts.Insecure)
			h.AssertEq(t, cnbIdx.RepoName, opts.BaseImageIndexRepoName)
			h.AssertEq(t, cnbIdx.XdgPath, opts.XdgPath)
		})
	})
}
