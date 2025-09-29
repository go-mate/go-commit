// Package commitmate: Advanced Git commit automation engine with intelligent Go formatting
// Features smart commit workflows with auto Go source code formatting and remote-based signature-info selection
// Provides seamless staging, formatting, committing, and amend operations with configuration-driven signatures
// Supports wildcard pattern matching for Git remote URLs to auto-select appropriate commit signatures
//
// commitmate: 高级 Git 提交自动化引擎，带有智能 Go 格式化功能
// 具有智能提交工作流程，包含自动 Go 源代码格式化和基于远程的身份选择
// 提供无缝的暂存、格式化、提交和 amend 操作，支持配置驱动的签名
// 支持 Git 远程 URL 的通配符模式匹配，自动选择合适的提交签名
package commitmate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-mate/go-commit/internal/utils"
	"github.com/go-xlan/gogit"
	"github.com/go-xlan/gogit/gogitassist"
	"github.com/go-xlan/gogit/gogitchange"
	"github.com/yyle88/erero"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/osexec"
	"github.com/yyle88/osexistpath/osmustexist"
	"github.com/yyle88/rese"
	"github.com/yyle88/tern/zerotern"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// CommitFlags represents the configuration for a commit operation
// Contains all custom options for customizing commit actions
// Supports amend mode, force operations, and selective Go formatting
//
// CommitFlags 代表提交操作的配置
// 包含所有用户指定的自定义提交行为选项
// 支持 amend 模式、强制操作和选择性 Go 格式化
type CommitFlags struct {
	Username string // Git account username // Git 账户用户名
	Message  string // Commit message content // 提交消息内容
	IsAmend  bool   // Enable amend mode on previous commit // 启用对上一次提交的 amend 模式
	IsForce  bool   // Force amend even if pushed to remote // 即使推送到远程也强制 amend
	Mailbox  string // Git account mailbox address (preferred) // Git 账户邮箱地址（优先）
	Eddress  string // Git account mailbox address (fallback) // Git 账户邮箱地址（备选）
	NoCommit bool   // Stage changes without committing // 仅暂存更改而不提交
	FormatGo bool   // Format changed Go files before commit // 提交前格式化已改变的 Go 文件
	AutoSign bool   // Use Git config as fallback // 使用 Git 配置作为备选
}

// ValidateFlags performs basic validation on commit flags and returns warnings
// Checks for logical conflicts and missing essential information
// Returns slice of warning messages for potential issues
//
// ValidateFlags 对提交标志执行基本验证并返回警告
// 检查逻辑冲突和缺失的基本信息
// 返回潜在问题的警告消息切片
func (f *CommitFlags) ValidateFlags() []string {
	var warnings []string

	// Check for force flag without amend
	// 检查没有 amend 的强制标志
	if f.IsForce && !f.IsAmend {
		warnings = append(warnings, "force flag set but amend is disabled - force has no effect")
	}

	// Check for commit message when not committing
	// 检查不提交时的提交消息
	if f.NoCommit && f.Message != "" {
		warnings = append(warnings, "commit message provided but no-commit flag is set")
	}

	// Check for missing authentication info (when not using AutoSign)
	// 检查缺失的身份验证信息（当不使用 AutoSign 时）
	if !f.AutoSign && f.Username == "" && f.Mailbox == "" && f.Eddress == "" {
		warnings = append(warnings, "no authentication info provided and auto-sign disabled")
	}

	return warnings
}

// GitCommit performs the complete commit workflow with selective Go code formatting
// Stages all changes, formats Go files when needed, and creates commits as requested
// Returns error if some step in the commit process fails
//
// 执行完整的提交工作流程，可选的 Go 代码格式化
// 暂存所有更改，可选格式化 Go 文件，并创建或 amend 提交
// 如果提交过程中的某个步骤失败则返回错误
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

	// Check initial repo status
	// 检查初始代码库状态
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

	// Check staged changes
	// 检查已暂存的更改
	status, err = client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(neatjsons.S(status))

	// Format Go files if requested
	// 如果请求则格式化 Go 文件
	if commitFlags.FormatGo {
		zaplog.SUG.Debugln("format changed go files")

		// Format changed Go files
		// 对已改变的文件应用 Go 格式化
		if err := FormatChangedGoFiles(projectRoot, client, DefaultAllowFormat); err != nil {
			return erero.Wro(err)
		}

		// Re-stage files when formatting done
		// 格式化完成后重新暂存文件
		if err := client.AddAll(); err != nil {
			return erero.Wro(err)
		}

		// Check status when formatting done
		// 格式化完成后检查状态
		status, err = client.Status()
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(neatjsons.S(status))
	}

	// Exit when no changes to commit
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
	// Prioritize Mailbox instead of Eddress for mailbox address
	// 优先使用 Mailbox 而非 Eddress 作为邮箱地址
	mailbox := zerotern.VV(commitFlags.Mailbox, commitFlags.Eddress)

	commitInfo := &gogit.CommitInfo{
		Name:    commitFlags.Username,
		Mailbox: mailbox,
		Message: commitFlags.Message,
	}

	// Use empty username/eddress from Git config as fallback (when allowed)
	// 从 Git 配置使用空的用户名/邮箱作为备选（当允许时）
	if commitFlags.AutoSign {
		useConfigSignInfo(projectRoot, commitInfo)
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

	// Debug repo state when commit done
	// 提交完成后调试代码库状态
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
	// Configure matching options for Go files with custom function
	// 为 Go 文件配置带有自定义过滤器的匹配选项
	matchOptions := gogitchange.NewMatchOptions().MatchType(".go").MatchPath(func(path string) bool {
		zaplog.SUG.Debugln("path:", path)

		// Use custom format function
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

		// Format the Go file
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

// ApplyProjectConfig applies project-specific configuration to commit flags
// Resolves appropriate signature from config based on project remote URLs
// Auto-selects and applies the best matching signature configuration
//
// ApplyProjectConfig 将项目特定配置应用到提交标志
// 基于项目远程 URL 从配置中解析合适的签名
// 自动选择并使用最佳匹配的签名配置
func (f *CommitFlags) ApplyProjectConfig(projectRoot string, config *CommitConfig) {
	zaplog.SUG.Debugln("applying project config to commit flags")
	f.ApplySignature(rese.V1(config.ResolveSignature(projectRoot)))
}

// ApplySignature applies signature configuration to flags
// Signature config values override existing flag values when signature fields are not empty
//
// ApplySignature 将签名配置应用到标志
// 当签名字段非空时，签名配置值会覆盖现有标志值
func (f *CommitFlags) ApplySignature(signature *SignatureConfig) {
	if signature != nil {
		zaplog.SUG.Debugln("applying signature config:", neatjsons.S(signature))
		// Use signature config value if available, otherwise keep existing flag value
		// Uses zerotern.VV to use config values or keep existing flag values as fallback
		// 如果配置中有值就使用配置值，否则保留现有标志值
		// 使用 zerotern.VV 优先使用配置值，其次使用现有标志值作为备选
		f.Username = zerotern.VV(signature.Username, f.Username)

		// Prioritize mailbox instead of eddress from signature config
		// 从签名配置中优先使用 mailbox 而非 eddress
		mailbox := zerotern.VV(signature.Mailbox, signature.Eddress)

		// Set mailbox to both fields to maintain compatibility
		// 将 mailbox 设置到两个字段以保持兼容性
		if mailbox != "" {
			f.Mailbox = zerotern.VV(mailbox, f.Mailbox)
			f.Eddress = zerotern.VV(mailbox, f.Eddress)
		}
	}
}

// SignatureConfig represents a Git signature configuration with advanced pattern matching
// Maps Git remote URL patterns to corresponding account username and mailbox settings
// Supports sophisticated wildcard matching for flexible remote pattern definitions
// Enables automatic signature-info switching based on repo remote configurations
//
// SignatureConfig 代表具有高级模式匹配的 Git 签名配置
// 将 Git 远程 URL 模式映射到相应的账户用户名和邮箱设置
// 支持复杂的通配符匹配以实现灵活的远程模式定义
// 基于代码库远程配置实现自动身份切换
type SignatureConfig struct {
	Name           string   `json:"name"`           // Config name for reference // 配置名称用于引用
	Username       string   `json:"username"`       // Git username for commits // 用于提交的 Git 用户名
	Mailbox        string   `json:"mailbox"`        // Git mailbox for commits (preferred) // 用于提交的 Git 邮箱（优先）
	Eddress        string   `json:"eddress"`        // Git mailbox for commits (fallback) // 用于提交的 Git 邮箱（备选）
	RemotePatterns []string `json:"remotePatterns"` // Remote URL patterns (supports wildcards) // 远程 URL 模式（支持通配符）
}

// CommitConfig represents the comprehensive configuration system for go-commit
// Contains intelligent signature mappings for automated Git commit operations
// Enables automatic signature selection based on Git remote URL pattern matching
// Supports scoring-based matching with wildcard patterns for enterprise and custom workflows
//
// CommitConfig 代表 go-commit 应用的全面配置系统
// 包含用于自动化 Git 提交操作的智能签名映射
// 基于 Git 远程 URL 模式匹配实现自动签名选择
// 支持基于评分的通配符模式匹配，适用于企业和自定义工作流程
type CommitConfig struct {
	Signatures []*SignatureConfig `json:"signatures"` // List of configured signatures // 配置的签名列表
}

// LoadConfig loads the go-commit configuration from the specified file path
// Reads, validates, and parses the JSON configuration file with signature mappings
// Utilizes osmustexist for file validation and rese/must for robust error handling
// Returns complete loaded configuration suited for signature matching operations
//
// LoadConfig 从指定文件路径加载 go-commit 配置
// 读取、验证并解析包含签名映射的 JSON 配置文件
// 使用 osmustexist 进行文件验证和 rese/must 进行强健的错误处理
// 返回完全加载的配置，准备进行签名匹配操作
func LoadConfig(configPath string) *CommitConfig {
	data := rese.A1(os.ReadFile(osmustexist.FILE(configPath)))

	// Parse JSON configuration
	// 解析 JSON 配置
	var config CommitConfig
	must.Done(json.Unmarshal(data, &config))

	// Validate loaded configuration
	// 验证加载的配置
	validateConfig(&config)

	zaplog.SUG.Debugln("loaded config from:", configPath)
	return &config
}

// validateConfig performs basic validation on the loaded configuration
// Ensures signature configurations have required fields and valid patterns
//
// validateConfig 对加载的配置执行基本验证
// 确保签名配置具有必需字段和有效模式
func validateConfig(config *CommitConfig) {
	for idx, signature := range config.Signatures {
		// Check core signature fields for completeness
		// 检查关键签名字段的完整性
		if signature.Username == "" {
			zaplog.SUG.Warnf("signature[%d] missing username", idx)
		}
		// Check for mailbox address (mailbox preferred instead of eddress)
		// 检查邮箱地址（优先 mailbox 而非 eddress）
		if signature.Mailbox == "" && signature.Eddress == "" {
			zaplog.SUG.Warnf("signature[%d] missing mailbox (mailbox or eddress)", idx)
		}
		if len(signature.RemotePatterns) == 0 {
			zaplog.SUG.Warnf("signature[%d] missing remote patterns", idx)
		}
	}
}

// ResolveSignature resolves Git signature based on project repo remotes
// Extracts Git remote URLs and performs pattern-based signature matching
// Prioritizes 'origin' remote but falls back to first available remote for signature resolution
// Returns the best matched signature or nil if no suitable patterns match the remote configuration
//
// ResolveSignature 基于项目仓库远程解析 Git 签名
// 提取 Git 远程 URL 并执行基于模式的签名匹配
// 优先使用 'origin' 远程，但在签名解析时回退到第一个可用远程
// 返回最佳匹配的签名，如果没有合适的模式匹配远程配置则返回 nil
func (config *CommitConfig) ResolveSignature(projectRoot string) (*SignatureConfig, error) {
	// Get Git remote URL
	// 获取 Git 远程 URL
	client := rese.P1(gogit.New(projectRoot))

	// 获取远程仓库地址，而且这里允许为空
	remotes := rese.V1(client.Repo().Remotes())

	// Use origin remote if available
	// 如果可用则使用 origin 远程
	var remoteURL string
	if len(remotes) > 0 {
		for _, remote := range remotes {
			remoteConfig := remote.Config()
			if remoteConfig.Name == "origin" && len(remoteConfig.URLs) > 0 {
				remoteURL = remoteConfig.URLs[0]
				break
			}
		}
		// Use first remote if no origin
		// 如果没有 origin 则回退到第一个远程
		if remoteURL == "" && len(remotes[0].Config().URLs) > 0 {
			remoteURL = remotes[0].Config().URLs[0]
		}
	}

	zaplog.SUG.Debugln("remote URL:", remoteURL)

	// Find and return matching signature
	// 查找并返回匹配的签名
	if remoteURL != "" {
		signature := config.MatchSignature(remoteURL)
		if signature != nil {
			zaplog.SUG.Debugln("matched signature:", signature.Name)
			return signature, nil
		}
		return nil, nil
	}

	// No matching signature found
	// 没有找到匹配的签名
	return nil, nil
}

// MatchSignature finds the best signature configuration for the specified remote URL
// Employs sophisticated pattern matching with wildcards and score-based ranking selection
// Evaluates all configured signature patterns and returns the highest-scoring match
// Returns the best matched signature or nil if no patterns match the remote URL
//
// MatchSignature 为指定的远程 URL 找到最佳的签名配置
// 采用复杂的通配符模式匹配和基于评分的优先级选择
// 评估所有配置的签名模式并返回得分最高的匹配
// 返回最佳匹配的签名，如果没有模式匹配远程 URL 则返回 nil
func (config *CommitConfig) MatchSignature(remoteURL string) *SignatureConfig {
	var bestMatch *SignatureConfig
	bestMatchScore := -1

	// Iterate through all configured signatures
	// 遍历所有配置的签名
	for _, signature := range config.Signatures {
		// Check each remote pattern for this signature
		// 检查此签名的每个远程模式
		for _, pattern := range signature.RemotePatterns {
			score := utils.MatchRemotePattern(pattern, remoteURL)
			if score > bestMatchScore {
				bestMatchScore = score
				bestMatch = signature
			}
		}
	}

	return bestMatch
}

// DefaultAllowFormat is the default check function for Go files formatting
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
	// Enable formatting for all remaining Go files
	// 允许格式化所有其他 Go 文件
	return true
}

// useConfigSignInfo sets empty username/eddress fields from Git configuration
// Uses "git config user.name" and "git config user.email" as fallback
//
// useConfigSignInfo 从 Git 配置填充空的用户名/邮箱字段
// 使用 "git config user.name" 和 "git config user.email" 作为备选
func useConfigSignInfo(projectRoot string, commitInfo *gogit.CommitInfo) {
	// Just get Git config if username is empty
	// 仅在用户名为空时才尝试获取 Git 配置
	if commitInfo.Name == "" {
		if gitUsername := getGitConfigValue(projectRoot, "user.name"); gitUsername != "" {
			commitInfo.Name = gitUsername
			zaplog.SUG.Debugln("using git config user.name:", gitUsername)
		}
	}

	// Just get Git config if eddress is empty
	// 仅在邮箱为空时才尝试获取 Git 配置
	if commitInfo.Mailbox == "" {
		if gitEddress := getGitConfigValue(projectRoot, "user.email"); gitEddress != "" {
			commitInfo.Mailbox = gitEddress
			zaplog.SUG.Debugln("using git config user.email:", gitEddress)
		}
	}
}

// getGitConfigValue retrieves a configuration value from Git config in the specified project DIR
// Returns blank string if command fails or config item doesn't exist
//
// getGitConfigValue 从指定项目 DIR 的 Git 配置获取配置值
// 如果命令失败或配置键不存在则返回空字符串
func getGitConfigValue(projectRoot, key string) string {
	execConfig := osexec.NewExecConfig().WithPath(projectRoot)
	output, err := execConfig.Exec("git", "config", key)
	if err != nil {
		zaplog.SUG.Debugln("cannot get git config", key, ":", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}
