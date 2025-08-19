// Package commitmate: Core Git commit automation engine with Go formatting capabilities
// Provides intelligent commit workflow with auto Go source code formatting
// Handles staging, formatting, committing, and amend operations seamlessly
//
// commitmate: 带有 Go 格式化功能的核心 Git 提交自动化引擎
// 提供智能提交工作流程，带有自动 Go 源代码格式化
// 无缝处理暂存、格式化、提交和 amend 操作
package commitmate

import (
	"path/filepath"
	"strings"

	"github.com/go-xlan/gogit"
	"github.com/go-xlan/gogit/gogitassist"
	"github.com/go-xlan/gogit/gogitchange"
	"github.com/yyle88/erero"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// CommitFlags represents the configuration for a commit operation
// Contains all custom options for customizing commit behavior
// Supports amend mode, force operations, and selective Go formatting
//
// CommitFlags 代表提交操作的配置
// 包含所有用户指定的自定义提交行为选项
// 支持 amend 模式、强制操作和选择性 Go 格式化
type CommitFlags struct {
	Username string // Git author username // Git 作者用户名
	Message  string // Commit message content // 提交消息内容
	IsAmend  bool   // Whether to amend previous commit // 是否 amend 上一次提交
	IsForce  bool   // Force amend even if pushed to remote // 即使推送到远程也强制 amend
	Eddress  string // Git author email address // Git 作者邮箱地址
	NoCommit bool   // Stage changes without committing // 仅暂存更改而不提交
	FormatGo bool   // Format changed Go files before commit // 提交前格式化已改变的 Go 文件
}

// GitCommit performs the complete commit workflow with optional Go code formatting
// Stages all changes, optionally formats Go files, and creates or amends commits
// Returns error if any step in the commit process fails
//
// 执行完整的提交工作流程，可选的 Go 代码格式化
// 暂存所有更改，可选格式化 Go 文件，并创建或 amend 提交
// 如果提交过程中的任何步骤失败则返回错误
func GitCommit(projectRoot string, commitFlags *CommitFlags) error {
	// Log project context and commit configuration
	// 记录项目上下文和提交配置
	zaplog.SUG.Debugln(projectRoot, neatjsons.S(commitFlags))

	// Initialize Git client for the project
	// 为项目初始化 Git 客户端
	client, err := gogit.New(projectRoot)
	if err != nil {
		return erero.Wro(err)
	}

	// Check initial repository status
	// 检查初始仓库状态
	status, err := client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(neatjsons.S(status))

	// Stage all changes for commit
	// 为提交暂存所有更改
	if err := client.AddAll(); err != nil {
		return erero.Wro(err)
	}

	// Verify staged changes
	// 验证已暂存的更改
	status, err = client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(neatjsons.S(status))

	// Format Go files if requested
	// 如果请求则格式化 Go 文件
	if commitFlags.FormatGo {
		zaplog.SUG.Debugln("format changed go files")

		// Apply Go formatting to changed files
		// 对已改变的文件应用 Go 格式化
		if err := FormatChangedGoFiles(projectRoot, client, DefaultAllowFormat); err != nil {
			return erero.Wro(err)
		}

		// Re-stage files after formatting
		// 格式化后重新暂存文件
		if err := client.AddAll(); err != nil {
			return erero.Wro(err)
		}

		// Check status after formatting
		// 格式化后检查状态
		status, err = client.Status()
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(neatjsons.S(status))
	}

	// Final status check before commit
	// 提交前的最终状态检查
	status, err = client.Status()
	if err != nil {
		return erero.Wro(err)
	}

	// Exit immediately if no changes to commit
	// 如果没有更改要提交则提前退出
	if len(status) == 0 {
		zaplog.SUG.Debugln("no change return")
		return nil
	}

	// Exit if staging without commit was requested
	// 如果请求仅暂存而不提交则退出
	if commitFlags.NoCommit {
		zaplog.SUG.Debugln("no commit return")
		return nil
	}

	// Prepare commit information from flags
	// 从标志准备提交信息
	commitInfo := &gogit.CommitInfo{
		Name:    commitFlags.Username,
		Eddress: commitFlags.Eddress,
		Message: commitFlags.Message,
	}

	// Execute commit or amend based on flags
	// 根据标志执行提交或 amend
	if commitFlags.IsAmend {
		// Amend the previous commit
		// Amend 上一次提交
		_, err = client.AmendCommit(&gogit.AmendConfig{
			CommitInfo: commitInfo,
			ForceAmend: commitFlags.IsForce,
		})
		if err != nil {
			return erero.Wro(err)
		}
	} else {
		// Create new commit
		// 创建新提交
		_, err = client.CommitAll(commitInfo)
		if err != nil {
			return erero.Wro(err)
		}
	}

	// Debug repository state after commit
	// 提交后调试仓库状态
	gogitassist.DebugRepo(client.Repo())
	return nil
}

// FormatChangedGoFiles formats Go files that have been changed
// Uses allowFormat function to determine which files should be formatted
// Applies Go formatting to eligible files and logs the process
//
// 格式化已改变的 Go 文件
// 使用 allowFormat 函数确定哪些文件应该被格式化
// 对符合条件的文件应用 Go 格式化并记录过程
func FormatChangedGoFiles(projectRoot string, client *gogit.Client, allowFormat func(path string) bool) error {
	// Configure matching options for Go files with custom filter
	// 为 Go 文件配置带有自定义过滤器的匹配选项
	matchOptions := gogitchange.NewMatchOptions().MatchType(".go").MatchPath(func(path string) bool {
		zaplog.SUG.Debugln("path:", path)

		// Apply custom format filter
		// 应用用户定义的格式过滤器
		pass := allowFormat(path)
		if pass {
			zaplog.SUG.Debugln("pass:", path)
		} else {
			zaplog.SUG.Debugln("skip:", path)
		}
		return pass
	})

	// Process each changed Go file with formatting
	// 处理每个已改变的 Go 文件进行格式化
	err := gogitchange.NewChangedFileManager(projectRoot, client.Tree()).ForeachChangedGoFile(matchOptions, func(path string) error {
		// Double-check file extension to ensure correctness
		// 为安全起见双重检查文件扩展名
		if filepath.Ext(path) != ".go" {
			return nil
		}

		// Log formatting operation
		// 记录格式化操作
		zaplog.ZAPS.Skip1.LOG.Info("golang-format-source", zap.String("path", path))

		// Apply Go formatting to the file
		// 对文件应用 Go 格式化
		if err := formatgo.FormatFile(path); err != nil {
			return erero.Wro(err)
		}
		return nil
	})
	if err != nil {
		return erero.Wro(err)
	}
	return nil
}

// DefaultAllowFormat is the default filter function for Go files formatting
// Skips common generated files like .pb.go, wire_gen.go, and ent files
// Returns true if the file should be formatted, false to skip
//
// DefaultAllowFormat 是 Go 文件格式化的默认过滤函数
// 跳过常见的生成文件，如 .pb.go、wire_gen.go 和 ent 文件
// 如果文件应该被格式化则返回 true，跳过则返回 false
func DefaultAllowFormat(path string) bool {
	// Skip various types of generated files
	// 跳过各种类型的生成文件
	if strings.HasSuffix(path, ".pb.go") || // skip protobuf generated files // 跳过 protobuf 生成文件
		strings.HasSuffix(path, "/wire_gen.go") || // skip wire generated files // 跳过 wire 生成文件
		strings.Contains(path, "/internal/data/ent/") { // skip ent generated files // 跳过 ent 生成文件
		return false
	}
	// Allow formatting for all other Go files
	// 允许格式化所有其他 Go 文件
	return true
}
