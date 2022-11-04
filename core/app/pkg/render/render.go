package render

import (
	"github.com/unrolled/render"
	"net/http"
)

var rnd *render.Render

func init() {
	rnd = render.New()
}

func JSON(w http.ResponseWriter, v interface{}) error {
	return rnd.JSON(w, http.StatusOK, v)
}
