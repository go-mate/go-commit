// go-commit: Smart Git commit tool with auto Go code formatting
// Provides intelligent commit workflow with optional Go source formatting
// Supports amend operations, custom commit messages, and custom configuration
//
// go-commit: 智能 Git 提交工具，带有自动 Go 代码格式化
// 提供智能提交工作流程，可选的 Go 源代码格式化
// 还支持 amend 操作、自定义提交消息和用户配置
package main

import (
	"os"

	"github.com/go-mate/go-commit/commitmate"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	// Get current working DIR as project root
	// 获取当前工作 DIR 作为项目根目录
	projectRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(projectRoot))

	// Initialize commit configuration flags
	// 初始化提交配置标志
	commitFlags := &commitmate.CommitFlags{}

	// Configuration file path flag
	// 配置文件路径标志
	var configPath string

	// Define root command with comprehensive Git commit functionality
	// 定义具有全面 Git 提交功能的根命令
	rootCmd := cobra.Command{
		Use:   "go-commit",
		Short: "Smart Git commit tool with Go code formatting",
		Long:  "go-commit is a Git commit tool that auto formats changed Go code and provides flexible commit options",
		Run: func(cmd *cobra.Command, args []string) {
			// Try to load signature config if config file is provided
			// 如果提供了配置文件则尝试加载签名配置
			if configPath != "" {
				commitFlags.ApplyProjectConfig(projectRoot, commitmate.LoadConfig(configPath))
			}

			must.Done(commitmate.GitCommit(projectRoot, commitFlags))
		},
	}
	// Configure command line flags for commit customization
	// 配置用于提交自定义的命令行标志
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Username, "username", "u", "", "git username")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Message, "message", "m", "", "commit message")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsAmend, "amend", "a", false, "amend to the previous commit")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsForce, "force", "f", false, "force amend even pushed to remote")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Eddress, "eddress", "e", "", "email address")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.NoCommit, "no-commit", false, "stage changes without committing")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.FormatGo, "format-go", false, "format changed go files")

	// 配置其它信息比如用户信息
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to go-commit configuration file")

	// Add config subcommand for configuration management
	// 添加配置子命令用于配置管理
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management for go-commit",
		Long:  "Manage go-commit configurations, validate existing configs, or generate templates",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate config file path is provided
			// 验证提供了配置文件路径
			if configPath == "" {
				zaplog.SUG.Panicln("missing config path. use -c flag")
			}

			zaplog.SUG.Debugln("config path:", configPath)

			// Load and apply configuration
			// 加载并应用配置
			config := commitmate.LoadConfig(configPath)
			zaplog.SUG.Debugln("config items:", neatjsons.S(config))

			commitFlags.ApplyProjectConfig(projectRoot, config)
			zaplog.SUG.Debugln("commit flags:", neatjsons.S(commitFlags))
		},
	}

	// Add example subcommand to config for generating configuration template
	// 在 config 命令下添加 example 子命令用于生成配置模板
	configExampleCmd := &cobra.Command{
		Use:   "example",
		Short: "Generate configuration template for current project",
		Long:  "Generate a go-commit configuration template based on current project's Git remote URL",
		Run: func(cmd *cobra.Command, args []string) {
			configTemplate := commitmate.GenerateConfigTemplate(projectRoot)

			zaplog.SUG.Infoln("Generated configuration template:")
			// Output template as formatted JSON
			// 将模板输出为格式化的 JSON
			zaplog.SUG.Infoln(neatjsons.S(configTemplate))
			zaplog.SUG.Infoln("Save this template to a file (e.g., go-commit-config.json).")
		},
	}

	configCmd.AddCommand(configExampleCmd)
	rootCmd.AddCommand(configCmd)

	// Execute the CLI application
	// 执行 CLI 应用程序
	must.Done(rootCmd.Execute())
}
