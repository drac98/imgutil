package imgutil_test

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/buildpacks/imgutil"
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
			h.AssertEq(t, trans.TLSClientConfig, (*tls.Config)(nil))
		})
	})
}
