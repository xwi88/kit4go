package json

import (
	"testing"

	"github.com/xwi88/kit4go/json"
)

type user struct {
	Blog string
	Name string
}

func TestJson(t *testing.T) {
	u := user{"xwi88", "http://github.com/xwi88"}
	t.Logf("pkg name: %v", json.PKG)

	if b, err := json.MarshalIndent(u, "", "  "); err != nil {
		t.Error(err)
	} else {
		t.Logf("data: %#v, MarshalIndent: %v", u, string(b))
	}
}
