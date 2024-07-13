package write

import (
	"testing"

	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func strP(i string) *string {
	return &i
}

func Test_labelsTxt(t *testing.T) {
	t.Parallel()
	if labelsTxt(nil) != "" {
		t.Fatal("labelsTxt(nil) must be empty")
	}
	s := labelsTxt([]*github.Label{
		{
			Name: strP("bug"),
		},
		{
			Name: strP("foo"),
		},
	})
	exp := `bug
foo`
	if s != exp {
		t.Fatalf("wanted %s, got %s", exp, s)
	}
}
