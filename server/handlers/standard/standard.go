package standard

import (
	"net/http"

	"github.com/arcade55/nzflights_webui/webui/pages"
)

// homeHandler is a simple handler that checks for the visitor ID.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// The middleware ensures this cookie will always exist.

	page := pages.HomePage()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	page.RenderStream(w)
}
