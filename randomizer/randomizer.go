package randomizer

import (
	"fmt"
	"net/http"

	"github.com/sethvargo/go-password/password"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	res, _ := password.Generate(10, 5, 0, false, false)
	fmt.Fprint(w, res)
	return
}
