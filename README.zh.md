# go-commit

快捷的 Git 提交工具，具备自动 Go 代码格式化功能。

---

## 英文文档

[ENGLISH README](README.md)

## 核心特性

🎯 **智能提交自动化**: 智能暂存、格式化和提交，支持 amend 模式  
⚡ **自动 Go 格式化**: 选择性格式化修改的 Go 文件，排除生成文件  
🔄 **签名信息管理**: 基于远程 URL 模式的自动 Git 签名选择  
🌍 **通配符模式**: 复杂企业工作流的高级模式匹配  
📋 **配置驱动**: 基于 JSON 的配置，支持优先级签名匹配

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

使用配置文件是可选的，但它能让您使用更多高级功能，例如根据项目的远程URL自动切换签名。

您可以根据当前项目的 Git 远程仓库来快速生成一份配置模板，以此开始：

```bash
# 这会在当前目录下创建一个 go-commit-config.json 文件
go-commit config example
```

该文件允许您为不同的远程仓库定义签名，格式如下：

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

# 快速修改并格式化
alias gca='go-commit --username=yourname --format-go --amend'
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

---

## 许可证

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 贡献代码

欢迎贡献代码！报告 bug、建议功能、贡献代码：

**问题与想法:**

- 🐛 **发现缺陷?** 在 GitHub 上提交 issue，包含重现步骤
- 💡 **功能构想?** 创建 issue 讨论建议
- 📖 **文档疑惑?** 报告问题帮助我们改进
- 🚀 **功能需求?** 分享你的使用场景帮助我们理解需求
- ⚡ **性能瓶颈?** 报告慢操作帮助我们优化
- 🔧 **配置困扰?** 询问复杂设置相关问题
- 🌟 **成功故事?** 分享这个包如何改善你的工作流程
- 💬 **意见建议?** 欢迎所有建议和评论

**代码贡献:**

1. **Fork**: 在 GitHub 上 Fork 仓库 (使用网页界面)。
2. **克隆**: 克隆 Fork 的项目 (`git clone https://github.com/yourname/repo-name.git`)。
3. **导航**: 进入克隆的项目 (`cd repo-name`)
4. **分支**: 创建功能分支 (`git checkout -b feature/xxx`)。
5. **编码**: 实现你的更改并编写全面的测试
6. **测试**: (Golang 项目) 确保测试通过 (`go test ./...`) 并遵循 Go 代码风格约定
7. **文档**: 更新面向用户的更改文档并使用有意义的提交消息
8. **暂存**: 暂存更改 (`git add .`)
9. **提交**: 提交更改 (`git commit -m "Add feature xxx"`) 确保向后兼容
10. **推送**: 推送到分支 (`git push origin feature/xxx`)。
11. **PR**: 在 GitHub 上创建 Pull Request (在 GitHub 网页上) 并提供详细描述。

请确保测试通过并包含相关文档更新。

---

## 支持

非常欢迎通过提交 pull request 和报告问题来为此项目做出贡献。

> 如果你觉得我的项目有价值，请给一些 GitHub 星标。
> 把它分享给 (golang) 编程同事和朋友。
>
> 如果你正在写关于开发工具和工作流程的技术博客，
> 我们很乐意提供内容写作支持来帮助推广这个项目。
>
> 我们致力于支持开源并为 (golang) 生态系统做出贡献。
> 反馈帮助我们构建服务 (golang) 开发者的增强工具。
>
> 作为一个团队协作，我们可以让 Git 工作流更高效、更愉快。
> 愿这个包让编码体验变得更加愉快！

**使用这个包快乐编程!** 🎉

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-mate/go-commit.svg?variant=adaptive)](https://starchart.cc/go-mate/go-commit)
