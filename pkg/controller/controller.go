package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func (ctrl *Controller) Run(ctx context.Context, params Params) error { //nolint:funlen,cyclop
	if err := ctrl.validateParams(params); err != nil {
		return fmt.Errorf("argument is invalid: %w", err)
	}

	isPR := params.PRNum > 0

	pr, err := ctrl.getPR(ctx, params)
	if err != nil {
		return err
	}

	if pr == nil {
		fmt.Fprintln(ctrl.Stdout, strings.Join([]string{
			"export " + params.Prefix + "HAS_ASSOCIATED_PR=false",
			"export " + params.Prefix + "IS_PR=false",
			"export " + params.Prefix + "REPO_OWNER=" + params.Owner,
			"export " + params.Prefix + "REPO_NAME=" + params.Repo,
		}, "\n"))
		return nil
	}

	files, _, err := ctrl.GitHub.GetPRFiles(ctx, github.ParamsGetPRFiles{
		Owner:    params.Owner,
		Repo:     params.Repo,
		PRNum:    pr.GetNumber(),
		FileSize: pr.GetChangedFiles(),
	})
	if err != nil {
		return fmt.Errorf("get pull request files: %w", err)
	}

	dir := params.Dir
	if dir == "" { //nolint:nestif
		d, err := os.MkdirTemp("", "ci-info")
		if err != nil {
			return fmt.Errorf("create a temporal directory: %w", err)
		}
		dir = d
	} else {
		if !filepath.IsAbs(dir) {
			d, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("convert -dir %s to absolute path: %w", dir, err)
			}
			dir = d
		}
		if err := os.MkdirAll(dir, 0o755); err != nil { //nolint:gomnd
			return fmt.Errorf("create a directory "+dir+": %w", err)
		}
	}

	ctrl.printEnvs(params.Prefix, dir, isPR, params.Owner, params.Repo, pr)

	if err := ctrl.writePRFilesJSON(filepath.Join(dir, "pr_files.json"), files); err != nil {
		return err
	}

	if err := ctrl.writePRJSON(filepath.Join(dir, "pr.json"), pr); err != nil {
		return err
	}

	if err := ctrl.writePRFilesTxt(filepath.Join(dir, "pr_files.txt"), files); err != nil {
		return err
	}

	if err := ctrl.writePRChangedFilesTxt(filepath.Join(dir, "pr_all_filenames.txt"), files); err != nil {
		return err
	}

	if err := ctrl.writeLabelsTxt(filepath.Join(dir, "labels.txt"), pr.Labels); err != nil {
		return fmt.Errorf("write labels.txt: %w", err)
	}
	return nil
}

func (ctrl *Controller) getPR(ctx context.Context, params Params) (*github.PullRequest, error) {
	prNum := params.PRNum
	if prNum <= 0 {
		logrus.WithFields(logrus.Fields{
			"owner": params.Owner,
			"repo":  params.Repo,
			"sha":   params.SHA,
		}).Debug("get pull request from SHA")
		prs, _, err := ctrl.GitHub.ListPRsWithCommit(ctx, github.ParamsListPRsWithCommit{
			Owner: params.Owner,
			Repo:  params.Repo,
			SHA:   params.SHA,
		})
		if err != nil {
			return nil, fmt.Errorf("list pull requests with a commit: %w", err)
		}
		logrus.WithFields(logrus.Fields{
			"size": len(prs),
		}).Debug("the number of pull requests assosicated with the commit")
		if len(prs) == 0 {
			return nil, nil //nolint:nilnil
		}
		prNum = prs[0].GetNumber()
	}
	pr, _, err := ctrl.GitHub.GetPR(ctx, github.ParamsGetPR{
		Owner: params.Owner,
		Repo:  params.Repo,
		PRNum: prNum,
	})
	if err != nil {
		return nil, fmt.Errorf("get a pull request: %w", err)
	}
	return pr, nil
}

var (
	errOwnerRequired      = errors.New("owner is required")
	errRepoRequired       = errors.New("repo is required")
	errSHAOrPRNumRequired = errors.New("sha or pr number is required")
)

func (ctrl *Controller) validateParams(params Params) error {
	if params.Owner == "" {
		return errOwnerRequired
	}
	if params.Repo == "" {
		return errRepoRequired
	}
	if params.PRNum <= 0 && params.SHA == "" {
		return errSHAOrPRNumRequired
	}
	return nil
}

func (ctrl *Controller) writeLabelsTxt(p string, labels []*github.Label) error {
	labelNames := make([]string, len(labels))
	for i, label := range labels {
		labelNames[i] = label.GetName()
	}
	txt := ""
	if len(labelNames) != 0 {
		txt = strings.Join(labelNames, "\n") + "\n"
	}
	return ctrl.writeFile(p, []byte(txt))
}

func (ctrl *Controller) writePRFilesTxt(p string, files []*github.CommitFile) error {
	prFileNames := make([]string, len(files))
	for i, file := range files {
		prFileNames[i] = file.GetFilename()
	}
	txt := ""
	if len(prFileNames) != 0 {
		txt = strings.Join(prFileNames, "\n") + "\n"
	}
	return ctrl.writeFile(p, []byte(txt))
}

func (ctrl *Controller) writeFile(p string, data []byte) error {
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create a file "+p+": %w", err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write a file "+p+": %w", err)
	}
	return nil
}

func (ctrl *Controller) writePRChangedFilesTxt(p string, files []*github.CommitFile) error {
	prFileNames := make(map[string]struct{}, len(files))
	for _, file := range files {
		prFileNames[file.GetFilename()] = struct{}{}
		prevFileName := file.GetPreviousFilename()
		if prevFileName != "" {
			prFileNames[prevFileName] = struct{}{}
		}
	}
	arr := make([]string, 0, len(prFileNames))
	for k := range prFileNames {
		arr = append(arr, k)
	}
	txt := ""
	if len(prFileNames) != 0 {
		txt = strings.Join(arr, "\n") + "\n"
	}
	return ctrl.writeFile(p, []byte(txt))
}

func (ctrl *Controller) writePRJSON(p string, pr *github.PullRequest) error {
	prJSON, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create a file "+p+": %w", err)
	}
	defer prJSON.Close()
	if err := json.NewEncoder(prJSON).Encode(pr); err != nil {
		return fmt.Errorf("encode a pull request as JSON: %w", err)
	}
	return nil
}

func (ctrl *Controller) writePRFilesJSON(p string, files []*github.CommitFile) error {
	prFilesJSON, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create a file "+p+": %w", err)
	}
	defer prFilesJSON.Close()
	if err := json.NewEncoder(prFilesJSON).Encode(files); err != nil {
		return fmt.Errorf("encode a pull request files as JSON: %w", err)
	}
	return nil
}

func (ctrl *Controller) printEnvs(prefix, dir string, isPR bool, owner, repo string, pr *github.PullRequest) {
	fmt.Fprintln(ctrl.Stdout, strings.Join([]string{
		"export " + prefix + "IS_PR=" + strconv.FormatBool(isPR),
		"export " + prefix + "HAS_ASSOCIATED_PR=true",
		"export " + prefix + "PR_NUMBER=" + strconv.Itoa(pr.GetNumber()),
		"export " + prefix + "BASE_REF=" + pr.GetBase().GetRef(),
		"export " + prefix + "HEAD_REF=" + pr.GetHead().GetRef(),
		"export " + prefix + "PR_AUTHOR=" + pr.GetUser().GetLogin(),
		"export " + prefix + "PR_MERGED=" + strconv.FormatBool(pr.GetMerged()),
		"export " + prefix + "TEMP_DIR=" + dir,
		"export " + prefix + "REPO_OWNER=" + owner,
		"export " + prefix + "REPO_NAME=" + repo,
	}, "\n"))
}
