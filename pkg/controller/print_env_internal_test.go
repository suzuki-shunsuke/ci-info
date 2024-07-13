package controller

import (
	"testing"

	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func Test_nonPREnv(t *testing.T) {
	t.Parallel()
	params := Params{
		Prefix: "CI_INFO_",
		Owner:  "suzuki-shunsuke",
		Repo:   "foo",
	}
	s := nonPREnv(params)
	exp := `export CI_INFO_HAS_ASSOCIATED_PR=false
export CI_INFO_IS_PR=false
export CI_INFO_REPO_OWNER=suzuki-shunsuke
export CI_INFO_REPO_NAME=foo`
	if s != exp {
		t.Fatalf("wanted %s, got %s", exp, s)
	}
}

func intP(i int) *int {
	return &i
}

func strP(i string) *string {
	return &i
}

func boolP(i bool) *bool {
	return &i
}

func Test_prEnv(t *testing.T) {
	t.Parallel()
	s := prEnv("CI_INFO_", "/tmp/ci-info_yoo", false, "suzuki-shunsuke", "foo", &github.PullRequest{
		Number: intP(10),
		Merged: boolP(true),
		Base: &github.PullRequestBranch{
			Ref: strP("dev"),
		},
		Head: &github.PullRequestBranch{
			Ref: strP("feature-branch"),
		},
		User: &github.User{
			Login: strP("octocat"),
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
