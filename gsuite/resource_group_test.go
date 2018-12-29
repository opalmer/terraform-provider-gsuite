package gsuite

import (
	"encoding/json"
	"fmt"
	directory "google.golang.org/api/admin/directory/v1"
	. "gopkg.in/check.v1"
	"net/http"
)

type ResourceGroupTest struct{}

var _ = Suite(&ResourceGroupTest{})

func (s *ResourceGroupTest) Test_resourceGroupRead(c *C) {
	group := &directory.Group{
		Id:                 "test-group",
		Name:               "test-group-name",
		AdminCreated:       true,
		DirectMembersCount: 42,
		Aliases:            []string{"alias-1", "alias-2"},
		NonEditableAliases: []string{"alias-3", "alias-4"},
	}

	cfg, d, stop := newTest(c, resourceGroup, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case fmt.Sprintf("/groups/%s", group.Id):
			c.Assert(r.Method, Equals, http.MethodGet)
			c.Assert(json.NewEncoder(w).Encode(group), IsNil)
		default:
			c.Fatalf("Unhandled route: %s (method: %s)", r.URL.Path, r.Method)
		}
	})
	defer stop()

	d.SetId(group.Id)
	c.Assert(resourceGroupRead(d, cfg), IsNil)
	c.Assert(d.Id(), Equals, group.Id)
	c.Assert(d.Get("direct_members_count"), Equals, int(group.DirectMembersCount))
	c.Assert(d.Get("admin_created"), Equals, group.AdminCreated)

	var aliases []string
	for _, alias := range d.Get("aliases").([]interface{}) {
		aliases = append(aliases, alias.(string))
	}
	c.Assert(aliases, DeepEquals, group.Aliases)

	var nonEditableAliases []string
	for _, alias := range d.Get("non_editable_aliases").([]interface{}) {
		nonEditableAliases = append(nonEditableAliases, alias.(string))
	}
	c.Assert(nonEditableAliases, DeepEquals, group.NonEditableAliases)
}
