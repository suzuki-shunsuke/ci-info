package controller

import (
	"fmt"

	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func (c *Controller) outputNonPR(params Params) {
	fmt.Fprintf(c.stdout, `export %sHAS_ASSOCIATED_PR=false
export %sIS_PR=false
export %sREPO_OWNER=%s
export %sREPO_NAME=%s
`,
		params.Prefix,
		params.Prefix,
		params.Prefix, params.Owner,
		params.Prefix, params.Repo)
}

func (c *Controller) printEnvs(prefix, dir string, isPR bool, owner, repo string, pr *github.PullRequest) {
	fmt.Fprintf(c.stdout, `export %sIS_PR=%t
export %sHAS_ASSOCIATED_PR=true
export %sPR_NUMBER=%d
export %sBASE_REF=%s
export %sHEAD_REF=%s
export %sPR_AUTHOR=%s
export %sPR_MERGED=%t
export %sTEMP_DIR=%s
export %sREPO_OWNER=%s
export %sREPO_NAME=%s
`,
		prefix, isPR,
		prefix,
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
