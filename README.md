# go-commit

Smart Git commit tool with auto Go changed code formatting capabilities.

## Install

```bash
go install github.com/go-mate/go-commit/cmd/go-commit@latest
```

## Usage

```bash
# Basic usage
go-commit -m "your commit message" --format-go

# With user info
go-commit -u "username" -e "email@example.com" -m "message" --format-go

# Amend previous commit
go-commit --amend -m "updated message" --format-go

# Force amend (even after push)
go-commit --amend --force -m "force amend message"
```

## Recommended Aliases

```bash
# Quick commit with formatting
alias gcm='go-commit --username=yourname --format-go'

# Quick amend with formatting
alias gca='go-commit --username=yourname --format-go --amend'
```

### Usage Examples

```bash
# Commit with message and Go formatting
gcm -m "add new feature"

# Amend last commit
gca -m "fix typo in commit message"

# Force amend (dangerous - use with caution)
gca -m "force update pushed to remote" --force
```

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repo on GitHub (using the webpage interface).
2. Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. Navigate to the cloned project (`cd repo-name`)
4. Create a feature branch (`git checkout -b feature/xxx`).
5. Stage changes (`git add .`)
6. Commit changes (`git commit -m "Add feature xxx"`).
7. Push to the branch (`git push origin feature/xxx`).
8. Open a pull request on GitHub (on the GitHub webpage).

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

If you find this package valuable, give me some stars on GitHub! Thank you!!!

**Thank you for your support!**

**Happy Coding with this package!** 🎉

Give me stars. Thank you!!!

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/go-commit.svg?variant=adaptive)](https://starchart.cc/go-mate/go-commit)
