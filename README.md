# go-commit

Smart Git commit tool with auto Go changed code formatting capabilities.


## Key Features

🎯 **Smart Commit Automation**: Intelligent staging, formatting, and committing with amend support  
⚡ **Auto Go Formatting**: Selective formatting of changed Go files with generated file exclusion  
🔄 **Signature-info Management**: Automatic Git signature selection based on remote URL patterns  
🌍 **Wildcard Patterns**: Sophisticated pattern matching for complex enterprise workflows  
📋 **Configuration-Driven**: JSON-based configuration with priority-based signature matching  

## Install

```bash
go install github.com/go-mate/go-commit/cmd/go-commit@latest
```

## Usage

```bash
# Basic commit with Go formatting
go-commit -m "some commit message" --format-go

# With signature info
go-commit -u "username" -e "example@example.com" -m "message" --format-go

# Use configuration file for auto choose signature info
go-commit -c "xx/xx/go-commit-config.json" -m "commit message" --format-go

# Amend previous commit
go-commit --amend -m "updated message" --format-go

# Force amend (even pushed to origin)
go-commit --amend --force -m "force amend message"
```

## Configuration

Create `go-commit-config.json` in system:

```json
{
  "signatures": [
    {
      "name": "work-github",
      "username": "work-user", 
      "eddress": "work@company.com",
      "remotePatterns": ["git@github.company.com:*"]
    },
    {
      "name": "play-github",
      "username": "play-user",
      "eddress": "play@gmail.com", 
      "remotePatterns": ["git@github.com:play-user/*"]
    }
  ]
}
```

See [configuration examples](internal/examples/)

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

---

## License

MIT License. See [LICENSE](LICENSE).

---

## Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

**Issues and Ideas:**
- 🐛 **Found a bug?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share your use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize by reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 🌟 **Success stories?** Share how this package improved your workflow
- 💬 **General feedback?** All suggestions and comments are welcome

**Code Contributions:**

1. **Fork**: Fork the repo on GitHub (using the webpage interface).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement your changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation for user-facing changes and use meaningful commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a pull request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## Support

Welcome to contribute to this project by submitting pull requests and reporting issues.

> If you find my projects valuable, please give some GitHub stars.
> Share it with (golang) programming teammates and friends who might benefit from it.
> 
> If you are writing tech blogs about development tools and workflows,
> we would be glad to provide content writing support to help promote this project.
> 
> We are committed to supporting open source and contributing to the (golang) ecosystem.
> Feedback helps us build enhanced tools that serve (golang) developers.
> 
> Working as a team we can make Git workflows more efficient and enjoyable.
> May the coding experience become more pleasant with this package!

**Happy Coding with this package!** 🎉

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-mate/go-commit.svg?variant=adaptive)](https://starchart.cc/go-mate/go-commit)
