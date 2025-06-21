# ci-info

[![Build Status](https://github.com/suzuki-shunsuke/ci-info/workflows/CI/badge.svg)](https://github.com/suzuki-shunsuke/ci-info/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/ci-info)](https://goreportcard.com/report/github.com/suzuki-shunsuke/ci-info)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/ci-info.svg)](https://github.com/suzuki-shunsuke/ci-info)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/ci-info/main/LICENSE)

CLI tool to get CI related information.

## Motivation

We develop this tool to get some information in CI.

* Related Pull Request
  * PR Author
  * Pull Request Files
  * Labels
  * base and head branch
  * etc
* etc

## Install

* [Homebrew](#homebrew)
* [aqua](#aqua)
* [GitHub Releases](#github-releases)

### Homebrew

You can install ci-info with [Homebrew](https://brew.sh/).

```console
$ brew install --cask suzuki-shunsuke/ci-info/ci-info
```

## aqua

You can install ci-info with [aqua](https://aquaproj.github.io/).

```console
$ aqua g -i suzuki-shunsuke/ci-info
```

## GitHub Releases

Please download a binary from the [release page](https://github.com/suzuki-shunsuke/ci-info/releases).

<details>
<summary>Verify downloaded binaries from GitHub Releases</summary>

You can verify downloaded binaries using some tools.

1. [Cosign](https://github.com/sigstore/cosign)
1. [slsa-verifier](https://github.com/slsa-framework/slsa-verifier)
1. [GitHub CLI](https://cli.github.com/)

#### 1. Cosign

You can install Cosign by aqua.

```sh
aqua g -i sigstore/cosign
```

```sh
gh release download -R suzuki-shunsuke/ci-info v2.3.1
cosign verify-blob \
  --signature ci-info_2.3.1_checksums.txt.sig \
  --certificate ci-info_2.3.1_checksums.txt.pem \
  --certificate-identity-regexp 'https://github\.com/suzuki-shunsuke/go-release-workflow/\.github/workflows/release\.yaml@.*' \
  --certificate-oidc-issuer "https://token.actions.githubusercontent.com" \
  ci-info_2.3.1_checksums.txt
```

Output:

```
Verified OK
```

After verifying the checksum, verify the artifact.

```sh
cat ci-info_2.3.1_checksums.txt | sha256sum -c --ignore-missing
```

#### 2. slsa-verifier

You can install slsa-verifier by aqua.

```sh
aqua g -i slsa-framework/slsa-verifier
```

```sh
gh release download -R suzuki-shunsuke/ci-info v2.3.1
slsa-verifier verify-artifact ci-info_2.3.1_darwin_arm64.tar.gz \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/suzuki-shunsuke/ci-info \
  --source-tag v2.3.1
```

Output:

```
Verified signature against tlog entry index 136878875 at URL: https://rekor.sigstore.dev/api/v1/log/entries/108e9186e8c5677a7ac053c11af84554df024d7c465abc4ae459493bd09be4875df26f45c1ffda32
Verified build using builder "https://github.com/slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@refs/tags/v2.0.0" at commit 69950dff0ec546640c90cbcaf23df344d0b612cd
Verifying artifact ci-info_2.3.1_darwin_arm64.tar.gz: PASSED
```

#### 3. GitHub CLI

ci-info >= v2.3.1

You can install GitHub CLI by aqua.

```sh
aqua g -i cli/cli
```

```sh
gh release download -R suzuki-shunsuke/ci-info v2.3.1 -p ci-info_2.3.1_darwin_arm64.tar.gz
gh attestation verify ci-info_2.3.1_darwin_arm64.tar.gz \
  -R suzuki-shunsuke/ci-info \
  --signer-workflow suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml
```

Output:

```
Loaded digest sha256:7fec0b88d213986b16605dd8e64f6230e4b4fc605a0ce4c2fd9698fdc40d3e2d for file://ci-info_2.3.1_darwin_arm64.tar.gz
Loaded 1 attestation from GitHub API
✓ Verification succeeded!

sha256:7fec0b88d213986b16605dd8e64f6230e4b4fc605a0ce4c2fd9698fdc40d3e2d was attested by:
REPO                                 PREDICATE_TYPE                  WORKFLOW
suzuki-shunsuke/go-release-workflow  https://slsa.dev/provenance/v1  .github/workflows/release.yaml@7f97a226912ee2978126019b1e95311d7d15c97a
```

</details>

## Requirements

GitHub Access Token is required to get the information about the Pull Request.
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

```
$ eval "$(ci-info run --owner suzuki-shunsuke --repo github-comment --pr 132)"
```

Some files are created.

```
$ ls "$CI_INFO_TEMP_DIR"
```

* pr_files.txt: The list of pull request file paths which include a maximum of 3000 files
* pr_all_filenames.txt: The list of pull request file paths which include a maximum of 3000 files. In addition to `pr_files.txt`, the list of renamed file's `previous_filename` is included too.
* pr_files.json: [The response body of GitHub API List pull requests files](https://docs.github.com/en/free-pro-team@latest/rest/reference/pulls#list-pull-requests-files)
* pr.json: [The response body of GitHub API Get a pull request](https://docs.github.com/en/free-pro-team@latest/rest/reference/pulls#get-a-pull-request)
* labels.txt: The list of pull request label names

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
