package output

import (
	"fmt"

	"github.com/suzuki-shunsuke/ci-info/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func NonPREnv(params domain.Params) string {
	return fmt.Sprintf(`export %sHAS_ASSOCIATED_PR=false
export %sIS_PR=false
export %sIS_ISSUE=false
export %sREPO_OWNER=%s
export %sREPO_NAME=%s`,
		params.Prefix,
		params.Prefix,
		params.Prefix,
		params.Prefix, params.Owner,
		params.Prefix, params.Repo)
}

func IssueEnv(params domain.Params) string {
	return fmt.Sprintf(`export %sHAS_ASSOCIATED_PR=false
export %sIS_PR=false
export %sIS_ISSUE=true
export %sNUMBER=%d
export %sISSUE_NUMBER=%d
export %sREPO_OWNER=%s
export %sREPO_NAME=%s`,
		params.Prefix,
		params.Prefix,
		params.Prefix,
		params.Prefix, params.IssueNum,
		params.Prefix, params.IssueNum,
		params.Prefix, params.Owner,
		params.Prefix, params.Repo)
}

func PREnv(prefix, dir string, isPR bool, owner, repo string, pr *github.PullRequest) string {
	return fmt.Sprintf(`export %sIS_PR=%t
export %sIS_ISSUE=false
export %sHAS_ASSOCIATED_PR=true
export %sPR_NUMBER=%d
export %sNUMBER=%d
export %sBASE_REF=%s
export %sHEAD_REF=%s
export %sPR_AUTHOR=%s
export %sPR_MERGED=%t
export %sTEMP_DIR=%s
export %sREPO_OWNER=%s
export %sREPO_NAME=%s`,
		prefix, isPR,
		prefix,
		prefix,
		prefix, pr.GetNumber(),
		prefix, pr.GetNumber(),
		prefix, pr.GetBase().GetRef(),
		prefix, pr.GetHead().GetRef(),
		prefix, pr.GetUser().GetLogin(),
		prefix, pr.GetMerged(),
		prefix, dir,
		prefix, owner,
		prefix, repo,
	)
}
