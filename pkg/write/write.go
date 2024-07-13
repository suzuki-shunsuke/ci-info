package write

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

func Write(fs afero.Fs, dir string, pr *github.PullRequest, files []*github.CommitFile) error {
	if err := writeJSON(fs, filepath.Join(dir, "pr_files.json"), files); err != nil {
		return err
	}

	if err := writeJSON(fs, filepath.Join(dir, "pr.json"), pr); err != nil {
		return err
	}

	if err := writeFile(fs, filepath.Join(dir, "pr_files.txt"), []byte(prFilesTxt(files)+"\n")); err != nil {
		return err
	}

	if err := writeFile(fs, filepath.Join(dir, "pr_all_filenames.txt"), []byte(prChangedFilesTxt(files)+"\n")); err != nil {
		return err
	}

	if err := writeFile(fs, filepath.Join(dir, "labels.txt"), []byte(labelsTxt(pr.Labels)+"\n")); err != nil {
		return fmt.Errorf("write labels.txt: %w", err)
	}

	return nil
}

func writeFile(fs afero.Fs, p string, data []byte) error {
	return afero.WriteFile(fs, p, data, 0o644) //nolint:mnd,wrapcheck
}

func labelsTxt(labels []*github.Label) string {
	if len(labels) == 0 {
		return ""
	}
	labelNames := make([]string, len(labels))
	for i, label := range labels {
		labelNames[i] = label.GetName()
	}
	return strings.Join(labelNames, "\n")
}

func prFilesTxt(files []*github.CommitFile) string {
	if len(files) == 0 {
		return ""
	}
	prFileNames := make([]string, len(files))
	for i, file := range files {
		prFileNames[i] = file.GetFilename()
	}
	return strings.Join(prFileNames, "\n")
}

func prChangedFilesTxt(files []*github.CommitFile) string {
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
	return strings.Join(arr, "\n")
}

func writeJSON(fs afero.Fs, p string, data any) error {
	prJSON, err := fs.Create(p)
	if err != nil {
		return fmt.Errorf("create a file %s: %w", p, err)
	}
	defer prJSON.Close()
	if err := json.NewEncoder(prJSON).Encode(data); err != nil {
		return fmt.Errorf("encode data as JSON: %w", err)
	}
	return nil
}
