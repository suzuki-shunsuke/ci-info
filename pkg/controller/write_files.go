package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func (c *Controller) writeFiles(dir string, pr *github.PullRequest, files []*github.CommitFile) error {
	if err := c.writePRFilesJSON(filepath.Join(dir, "pr_files.json"), files); err != nil {
		return err
	}

	if err := c.writePRJSON(filepath.Join(dir, "pr.json"), pr); err != nil {
		return err
	}

	if err := c.writePRFilesTxt(filepath.Join(dir, "pr_files.txt"), files); err != nil {
		return err
	}

	if err := c.writePRChangedFilesTxt(filepath.Join(dir, "pr_all_filenames.txt"), files); err != nil {
		return err
	}

	if err := c.writeLabelsTxt(filepath.Join(dir, "labels.txt"), pr.Labels); err != nil {
		return fmt.Errorf("write labels.txt: %w", err)
	}

	return nil
}

func (c *Controller) writeLabelsTxt(p string, labels []*github.Label) error {
	labelNames := make([]string, len(labels))
	for i, label := range labels {
		labelNames[i] = label.GetName()
	}
	txt := ""
	if len(labelNames) != 0 {
		txt = strings.Join(labelNames, "\n") + "\n"
	}
	return c.writeFile(p, []byte(txt))
}

func (c *Controller) writePRFilesTxt(p string, files []*github.CommitFile) error {
	prFileNames := make([]string, len(files))
	for i, file := range files {
		prFileNames[i] = file.GetFilename()
	}
	txt := ""
	if len(prFileNames) != 0 {
		txt = strings.Join(prFileNames, "\n") + "\n"
	}
	return c.writeFile(p, []byte(txt))
}

func (c *Controller) writeFile(p string, data []byte) error {
	file, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("create a file %s: %w", p, err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write a file %s: %w", p, err)
	}
	return nil
}

func (c *Controller) writePRChangedFilesTxt(p string, files []*github.CommitFile) error {
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
	return c.writeFile(p, []byte(txt))
}

func (c *Controller) writePRJSON(p string, pr *github.PullRequest) error {
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

func (c *Controller) writePRFilesJSON(p string, files []*github.CommitFile) error {
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
