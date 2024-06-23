# gh-git-describe

`gh` extension to execute `git describe` on a remote GitHub repository.


## Installation

```console
gh extension install thombashi/gh-git-describe
```


## Upgrade

```console
gh extension upgrade git-describe
```


## Usage

### Command help

```
      --cache-dir string   cache directory path. If not specified, use the system's temporary directory.
      --log-level string   log level (debug, info, warn, error) (default "info")
      --no-cache           disable cache
  -R, --repo string        [required] GitHub repository ID
```

### Examples

```console
$ gh git-describe -R actions/checkout -- --tags a5ac7e51b41094c92402da3b24376905380afc29
v4.1.6
$ gh git-describe -R actions/checkout -- --tags b80ff79f1755d06ba70441c368a6fe801f5f3a62
v4.1.6-2-gb80ff79
```
