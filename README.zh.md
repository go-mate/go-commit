[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-mate/go-commit/release.yml?branch=main&label=BUILD)](https://github.com/go-mate/go-commit/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-mate/go-commit)](https://pkg.go.dev/github.com/go-mate/go-commit)
[![Coverage Status](https://img.shields.io/coveralls/github/go-mate/go-commit/main.svg)](https://coveralls.io/github/go-mate/go-commit?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-mate/go-commit.svg)](https://github.com/go-mate/go-commit/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-mate/go-commit)](https://goreportcard.com/report/github.com/go-mate/go-commit)

# go-commit

快捷的 Git 提交应用，具备自动 Go 代码格式化功能。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🎯 **智能提交自动化**: 智能暂存、格式化和提交，支持 amend 模式  
⚡ **自动 Go 格式化**: 选择性格式化修改的 Go 文件，排除生成文件  
🔄 **签名信息管理**: 基于远程 URL 模式的自动 Git 签名选择  
🌍 **通配符模式**: 复杂企业工作流的高级模式匹配  
📋 **配置驱动**: 基于 JSON 的配置，支持评分式签名匹配

## 安装

```bash
go install github.com/go-mate/go-commit/cmd/go-commit@latest
```

## 使用方法

```bash
# 基本提交并格式化 Go 代码
go-commit -m "some commit message" --format-go

# 使用签名信息
go-commit -u "username" -e "example@example.com" -m "message" --format-go

# 使用配置文件自动选择签名信息
go-commit -c "xx/xx/go-commit-config.json" -m "commit message" --format-go

# 修改上一次提交
go-commit --amend -m "updated message" --format-go

# 强制修改 (即使已推送到远程)
go-commit --amend --force -m "force amend message"
```

## 配置

使用配置文件是自适应的，但它能让您使用更多高级功能，例如根据项目的远程URL自动切换签名。

您可以根据当前项目的 Git 远程代码库来快速生成一份配置模板，以此开始：

```bash
# 这会在当前文件夹下创建一个 go-commit-config.json 文件
go-commit config example
```

该文件允许您为不同的远程代码库定义签名，格式如下：

```json
{
  "signatures": [
    {
      "name": "work-github", "username": "work-user", "eddress": "work@company.com", "remotePatterns": [
      "git@github.company.com:*"
    ]
    },
    {
      "name": "play-github", "username": "play-user", "eddress": "play@example.com", "remotePatterns": [
      "git@github.com:play-user/*"
    ]
    }
  ]
}
```

示例:

- 项目 A 的远程地址为 `git@github.company.com:team/project-a` → 自动使用 work-user(work@company.com) 提交
- 项目 B 的远程地址为 `git@github.com:play-user/project-b` → 自动使用 play-user(play@example.com) 提交

这种自动切换功能让多项目工作流变得更加便捷。

**验证配置:**

设置好配置文件后，您可以验证其是否正确：

```bash
# 检查配置是否正确加载并预览匹配的签名
go-commit config -c /path/to/go-commit-config.json
```

如果希望了解更多高级用法，请参阅[配置示例](internal/examples/)。

## 推荐别名

```bash
# 快速提交并格式化
alias gcm='go-commit --username=yourname --format-go'

# 快速追加提交并格式化（扩展 gcm）
alias gca='gcm --amend'
```

### 使用示例

```bash
# 提交消息并格式化 Go 代码
gcm -m "添加个新功能"

# 追加最后一次提交
gca

# 修改最后一次提交
gca -m "新的提交信息"

# 强制追加 (危险 - 谨慎使用)
gca -m "修改提交信息" --force
```

### 高级使用示例

```bash
# 仅暂存更改而不提交（用于测试）
go-commit --no-commit --format-go

# 自动格式化 Go 文件并使用自动签名提交
go-commit -m "改进代码格式" --format-go --auto-sign

# 使用特定用户信息提交（覆盖配置）
go-commit -u "张三" -e "zhangsan@company.com" -m "紧急修复" --format-go

# 使用 mailbox 而非 eddress，语义更清晰
go-commit --mailbox "developer@team.com" -m "功能更新" --format-go

# 配置驱动的提交（基于远程自动选择签名）
go-commit -c ~/go-commit-config.json -m "自动化提交" --format-go
```

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-mate/go-commit.svg?variant=adaptive)](https://starchart.cc/go-mate/go-commit)
