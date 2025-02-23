package view

import (
	g "maragu.dev/gomponents"
)

// HTMX attributes
func HxPost(url string) g.Node {
	return g.Attr("hx-post", url)
}
