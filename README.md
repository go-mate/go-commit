[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-commit/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-commit/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-commit)](https://pkg.go.dev/github.com/go-mate/go-commit)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-commit/main.svg)](https://coveralls.io/github/go-mate/go-commit?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-commit.svg)](https://github.com/go-mate/go-commit/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-commit)](https://goreportcard.com/report/github.com/go-mate/go-commit)

# go-commit

Quick Git commit app with auto Go changed code formatting capabilities.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Quick Commit Automation**: Intelligent staging, formatting, and committing with amend support  
⚡ **Auto Go Formatting**: Selective formatting of changed Go files with generated file exclusion  
🔄 **Signature-info Management**: Automatic Git signature selection based on remote URL patterns  
🌍 **Wildcard Patterns**: Sophisticated pattern matching with complex enterprise workflows  
📋 **Configuration-Driven**: JSON-based configuration with score-based signature matching

## Installation

```bash
go install github.com/go-mate/go-commit/cmd/go-commit@latest
```

## Usage

```bash
# Quick commit with Go formatting
go-commit -m "some commit message" --format-go

# With signature info
go-commit -u "username" -e "example@example.com" -m "message" --format-go

# Use configuration file to auto choose signature info
go-commit -c "xx/xx/go-commit-config.json" -m "commit message" --format-go

# Amend previous commit
go-commit --amend -m "updated message" --format-go

# Force amend (even when pushed to origin)
go-commit --amend --force -m "force amend message"
```

## Configuration

Using a configuration file is adaptive but enables advanced features like automatic signature switching based on the project's remote URL.

To get started, you can generate a configuration template based on the current git remote:

```bash
# This creates a go-commit-config.json in current DIR
go-commit config example
```

This file allows you to define signatures with different git remotes:

```json
{
  "signatures": [
    {
      "name": "work-github", "username": "work-man", "eddress": "work@corp.com", "remotePatterns": [
      "git@github.corp.com:*"
    ]
    },
    {
      "name": "home-github", "username": "home-dev", "eddress": "home@example.com", "remotePatterns": [
      "git@github.com:home-dev/*"
    ]
    }
  ]
}
```

Examples:

- Project A with remote `git@github.corp.com:team/project-a` → auto commits as work-man(work@corp.com)
- Project B with remote `git@github.com:home-dev/project-b` → auto commits as home-dev(home@example.com)

This automatic switching makes multi-project workflow much more convenient.

**Validate Configuration:**

Once setting up the configuration, you can validate it:

```bash
# Check if config loads as expected and preview matched signature
go-commit config -c /path/to/go-commit-config.json
```

See the [configuration examples](internal/examples/) on advanced use cases.

## Recommended Aliases

```bash
# Quick commit with formatting
alias gcm='go-commit --username=yourname --format-go'

# Quick amend with formatting (extends gcm)
alias gca='gcm --amend'
```

### Usage Examples

```bash
# Commit with message and Go formatting
gcm -m "add new feature"

# Amend last commit
gca

# Change last commit message
gca -m "new commit message"

# Amend commit when signature info changes (even without code changes)
# Note: gca amends the commit if username/mailbox differs from previous commit
gca

# Force amend (dangerous - use with caution)
gca -m "force update pushed to remote" --force
```

### Advanced Usage Examples

```bash
# Stage changes without committing (during testing)
go-commit --no-commit --format-go

# Auto-format Go files and commit with auto-signature
go-commit -m "improve code format" --format-go --auto-sign

# Commit with specific username info (overrides config)
go-commit -u "John Doe" -e "john@corp.com" -m "hotfix" --format-go

# Use mailbox instead of eddress (improved semantics)
go-commit --mailbox "dev@team.com" -m "feature update" --format-go

# Configuration-driven commit (automatic signature based on remote)
go-commit -c ~/go-commit-config.json -m "automated commit" --format-go
```

---

## 🛠️ Development Notes

**Git Package Dependencies (Project-Specific Rule):**

In this project, we enforce a strict separation between production and test code:

- **Production code**: Use ONLY `gogit` (github.com/go-xlan/gogit) - provides object-based Git operations
- **Test code**: Use `gitgo` (github.com/go-xlan/gitgo) - provides chainable operations to set up test conditions

⚠️ **Important**: Do NOT use `gitgo` in production logic code. This rule applies to this project to maintain clean separation of concerns.

---

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
- 📖 **Documentation confusing?** Report it so we can enhance it
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize via reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Track project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package enhanced the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and adhere to Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/go-commit.svg?variant=adaptive)](https://starchart.cc/go-mate/go-commit)
