# go-import-checks

Golang source code linter to enforce some rules about inter-package imports inside you project. You can specify imports that are allowed and other that should be rejected.

This tool will report any violation of your rules.

The idea came from [go-cleanarch](https://github.com/roblaszczak/go-cleanarch), but I needed something more flexible.

## Install

For GO version 1.17 or later:

```bash
$ go install github.com/quantumcycle/go-import-checks@latest
```

For earlier GO versions:

```
GO111MODULE=on go get github.com/quantumcycle/go-import-checks
```

## Usage

You need to provide the linter with a config in yaml to define the rules you want to inforce.

### Rules syntax

The rule's syntax is a subset of glob since it's a custom implementation just for package names. It supports:

- '\*' to match a single path element
- '\*\*' to match multiple path elements
- '!xxx' to match anything different from 'xxx'

### Allow rules

Here is an example of some rules:

```yaml
imports-checks:
  - folder: "pkg/$lib"
    subpackages: true
    rules:
      - prefix: "github.com/quantumcycle/project1/"
        allow:
          - "pkg/**"
```

This translate to:

- relative to the root, for any `.go` file in `pkg/xxx`, or any subpackages, if an import starts with `github.com/quantumcycle/project1/`, it can only be `github.com/quantumcycle/project1/pkg/xxx`. In effect, it means that anything in `pkg` can only import other `pkg`, or any other packages outside of `github.com/quantumcycle/project1/`.

So the following imports would pass the checks:

- `github.com/quantumcycle/project1/pkg/crypto` (because it's in pkg)
- `github.com/quantumcycle/project1/pkg/convert` (because it's in pkg)
- `context` (because it's not in `github.com/quantumcycle/project1/`)
- `github.com/pkg/errors` (because it's not in `github.com/quantumcycle/project1/`)

And any of thoses would fail:

- `github.com/quantumcycle/project1/internal/api`
- `github.com/quantumcycle/project1/cmd/runner`

You can run this example:

```bash
cd validator/tests/allowed/restrict-same-package
go-import-checks --config=./config.yaml
```

You can use `subpackages: true` to indicate that the rules applied to this folder or anything below.

You can also use placeholders in the `folder` and reference them in the rules, like this:

```yaml
imports-checks:
  - folder: "internal/systems/$systemName"
    subpackages: true
    rules:
      - prefix: "github.com/quantumcycle/project1/"
        allow:
          - "internal/systems/$systemName/*"
          - "internal/systems/!$systemName/api"
          - "pkg/**"
```

You can run this example:

```bash
cd validator/tests/allowed/restrict-same-subpackage
go-import-checks --config=./config.yaml
```

In essence, it means that any package in `internal/systems/xxx` can import any of it's subpackages, or the `sdk` package of any other system package, or anything in `pkg`.

So in a file located in `github.com/quantumcycle/project1/internal/systems/catalog`, the following would pass:

- `github.com/quantumcycle/project1/internal/systems/catalog/data` (same system, sub-package)
- `github.com/quantumcycle/project1/internal/systems/catalog/data/model` (same system, sub-package)
- `github.com/quantumcycle/project1/internal/systems/pricing/api` (other system, api package)
- `github.com/quantumcycle/project1/pkg/calculator` (any pkg is allowed)

And the following would fail:

- `github.com/quantumcycle/project1/internal/systems/pricing/data` (only /api of other system is accepted)

Notice the `!` meaning that this part of the import must differ from the matching part in the folder path.

### Reject rules

Until now, we've seen rules to allow (whitelist) imports. You can also specify rejection rules (blacklisting) like this:

```yaml
imports-checks:
  - folder: "pkg/$lib"
    subpackages: true
    rules:
      - prefix: "github.com/quantumcycle/project1/"
        reject:
          - "internal/**"
```

For any lib package inside `pkg`, if a import starts with `github.com/quantumcycle/project1/`, it cannot be for `internal/xxx` packages, but anything else is fine. So when using `allow`, you specify what can be used, and with `reject`, you specify what cannot be used.

## Exclusion

There is some cases where it's ok to bypass the defined rules. You might need mocks to be able to bypass the rules for example. In these cases, use the `exclude` attribute to specify some regex for which the rules don't apply. Exclude are not case sensitive, but otherwise they are regex, so be sure to escape regex characters like dots.

```yaml
imports-checks:
  - folder: "internal/systems/$systemName"
    subpackages: true
    exclude:
      - ".*mock.*\\.go"
    rules:
      - prefix: "github.com/quantumcycle/project1/"
        allow:
          - "pkg/**"
```

## Running

You can run this example:

```bash
cd validator/tests/reject/reject-another-package
go-import-checks --config=./config.yaml
```
