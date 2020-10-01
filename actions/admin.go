package actions

import (
	"github.com/gobuffalo/buffalo"
)

// AdminShow default implementation.
func AdminShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("admin/show.html"))
}
