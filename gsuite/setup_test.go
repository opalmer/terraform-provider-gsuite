package gsuite

import (
	"github.com/hashicorp/terraform/helper/schema"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

func newTest(c *C, resourceFunc func() *schema.Resource, handler http.HandlerFunc) (*Config, *schema.ResourceData, func()) {
	ts := httptest.NewServer(handler)
	cfg := &Config{basePath: ts.URL}
	c.Assert(cfg.loadAndValidate(), IsNil)
	return cfg, resourceFunc().TestResourceData(), ts.Close
}
