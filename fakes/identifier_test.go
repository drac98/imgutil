package fakes_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil/fakes"
	h "github.com/buildpacks/imgutil/testhelpers"
)

var fakeHash = "sha256:8ecc4820859d4006d17e8f4fd5248a81414c4e3b7ed5c34b623e23b3436fb1b2"

func TestFakeIdentifier(t *testing.T) {
	spec.Run(t, "FakeIdentifier", testFakeIdentifier, spec.Parallel(), spec.Report(report.Terminal{}))
}

func testFakeIdentifier(t *testing.T, when spec.G, it spec.S) {
	when("#Indentifier", func() {
		it("should create return expected result", func() {
			id := fakes.NewIdentifier(fakeHash)
			h.AssertEq(t, id.String(), fakeHash)
		})
	})
}
