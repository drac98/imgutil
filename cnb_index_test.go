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
}
