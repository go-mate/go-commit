# go-commit
commit git project with golang cobra commands. format changed go source files.

# install

```bash
go install github.com/go-mate/go-commit/cmd/go-commit@latest
```

# command

```bash
cd project-path && go-commit -m 'message' --format-go
```

Can also add alias:
```bash
alias gcm='go-commit --username=yangyile --format-go'
alias gca='go-commit --username=yangyile --format-go --amend'
```

Can use `gcm` to commit all changes with a message and format Go source files:
```bash
cd project-path
gca -m message
```

Can use `gca` to amend the last commit with a message and format Go src files:
```bash
cd project-path
gca -m message
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
