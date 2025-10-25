# ci-info

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/ci-info/main/LICENSE) | [INSTALL](INSTALL.md)

CLI tool to get CI related information.

## Motivation

We develop this tool to get some information in CI.

- Related Pull Request
  - PR Author
  - Pull Request Files
  - Labels
  - base and head branch
  - etc
- etc

## Requirements

GitHub Access Token with the pull requests read permission is required to get the information about the Pull Request.
In the public repository, GitHub Access Token is optional.

## Getting Started

Run the following command, which gets the information about https://github.com/suzuki-shunsuke/github-comment/pull/132 .

```console
$ ci-info run --owner suzuki-shunsuke --repo github-comment --pr 132
export CI_INFO_IS_PR=true
export CI_INFO_HAS_ASSOCIATED_PR=true
export CI_INFO_PR_NUMBER=132
export CI_INFO_BASE_REF=master
export CI_INFO_HEAD_REF=feat/add-silent-option
export CI_INFO_PR_AUTHOR=suzuki-shunsuke
export CI_INFO_PR_MERGED=true
export CI_INFO_REPO_OWNER=suzuki-shunsuke
export CI_INFO_REPO_NAME=github-comment
export CI_INFO_TEMP_DIR=/var/folders/w0/kzjzgvd52wg5s4jy5h0lcyqh0000gn/T/ci-info497729786
```

Then the shell script to export the environment variables are outputted and some files are created.
You can export them by `eval`.

```sh
eval "$(ci-info run --owner suzuki-shunsuke --repo github-comment --pr 132)"
```

Some files are created.

```sh
ls "$CI_INFO_TEMP_DIR"
```

- pr_files.txt: The list of pull request file paths which include a maximum of 3000 files
- pr_all_filenames.txt: The list of pull request file paths which include a maximum of 3000 files. In addition to `pr_files.txt`, the list of renamed file's `previous_filename` is included too.
- pr_files.json: [The response body of GitHub API List pull requests files](https://docs.github.com/en/free-pro-team@latest/rest/reference/pulls#list-pull-requests-files)
- pr.json: [The response body of GitHub API Get a pull request](https://docs.github.com/en/free-pro-team@latest/rest/reference/pulls#get-a-pull-request)
- labels.txt: The list of pull request label names

Note that the created directory and files aren't removed automatically.

## Usage

```console
$ ci-info help
NAME:
   ci-info - get CI information. https://github.com/suzuki-shunsuke/ci-info

USAGE:
   ci-info [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   run      get CI information
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

```console
$ ci-info run --help
NAME:
   ci-info run - get CI information

USAGE:
   ci-info run [command options] [arguments...]

OPTIONS:
   --owner value         repository owner
   --repo value          repository name
   --sha value           commit SHA
   --dir value           directory path where files are created. The directory is created by os.MkdirAll if it doesn't exist. By default the directory is created at Go's ioutil.TempDir
   --pr value            pull request number (default: 0)
   --github-token value  GitHub Access Token [$GITHUB_TOKEN, $GITHUB_ACCESS_TOKEN]
   --prefix value        The prefix of environment variable name (default: "CI_INFO_")
   --log-level value     log level
   --help, -h            show help (default: false)
```

## Complement options with Platform's built-in Environment variables

With [suzuki-shunske/go-ci-env](https://github.com/suzuki-shunsuke/go-ci-env) some parameters like `owner` and `repo` are gotten from Platform's built-in Environment variables.

## LICENSE

[MIT](LICENSE)
