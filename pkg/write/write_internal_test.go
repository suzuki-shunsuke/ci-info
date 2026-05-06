package write

import (
	"testing"

	"github.com/suzuki-shunsuke/ci-info/v2/pkg/github"
)

func Test_labelsTxt(t *testing.T) {
	t.Parallel()
	if labelsTxt(nil) != "" {
		t.Fatal("labelsTxt(nil) must be empty")
	}
	s := labelsTxt([]*github.Label{
		{
			Name: new("bug"),
		},
		{
			Name: new("foo"),
		},
	})
	exp := `bug
foo`
	if s != exp {
		t.Fatalf("wanted %s, got %s", exp, s)
	}
}
