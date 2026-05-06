package output

import (
	"testing"

	"github.com/suzuki-shunsuke/ci-info/v2/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/github"
)

func Test_nonPREnv(t *testing.T) {
	t.Parallel()
	params := domain.Params{
		Prefix: "CI_INFO_",
		Owner:  "suzuki-shunsuke",
		Repo:   "foo",
	}
	s := NonPREnv(&params)
	exp := `export CI_INFO_HAS_ASSOCIATED_PR=false
export CI_INFO_IS_PR=false
export CI_INFO_REPO_OWNER=suzuki-shunsuke
export CI_INFO_REPO_NAME=foo`
	if s != exp {
		t.Fatalf("wanted %s, got %s", exp, s)
	}
}

func Test_prEnv(t *testing.T) {
	t.Parallel()
	s := PREnv("CI_INFO_", "/tmp/ci-info_yoo", false, "suzuki-shunsuke", "foo", &github.PullRequest{
		Number: new(10),
		Merged: new(true),
		Base: &github.PullRequestBranch{
			Ref: new("dev"),
		},
		Head: &github.PullRequestBranch{
			Ref: new("feature-branch"),
		},
		User: &github.User{
			Login: new("octocat"),
		},
	})
	exp := `export CI_INFO_IS_PR=false
export CI_INFO_HAS_ASSOCIATED_PR=true
export CI_INFO_PR_NUMBER=10
export CI_INFO_BASE_REF=dev
export CI_INFO_HEAD_REF=feature-branch
export CI_INFO_PR_AUTHOR=octocat
export CI_INFO_PR_MERGED=true
export CI_INFO_TEMP_DIR=/tmp/ci-info_yoo
export CI_INFO_REPO_OWNER=suzuki-shunsuke
export CI_INFO_REPO_NAME=foo`
	if s != exp {
		t.Fatalf("wanted %s, got %s", exp, s)
	}
}
