// go-commit: Smart Git commit app with auto Go code formatting
// Provides intelligent commit workflow with adaptive Go source formatting
// Supports amend operations, custom commit messages, and custom configuration
//
// go-commit: 智能 Git 提交应用，带有自动 Go 代码格式化
// 提供智能提交工作流程，自适应的 Go 源代码格式化
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

// AppConfig holds application configuration options
// 应用配置保存应用程序配置选项
type AppConfig struct {
	ConfigPath string // Path to configuration file // 配置文件路径
}

func main() {
	// Get current working DIR as project root
	// 获取当前工作 DIR 作为项目根目录
	projectRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(projectRoot))

	// Initialize commit configuration flags
	// 初始化提交配置标志
	commitFlags := &commitmate.CommitFlags{}

	// Initialize app configuration
	// 初始化应用配置
	appConfig := &AppConfig{}

	// Create and configure root command
	// 创建并配置根命令
	rootCmd := createRootCommand(projectRoot, commitFlags, appConfig)

	// Add config command and its subcommands
	// 添加配置命令及其子命令
	configCmd := createConfigCommand(projectRoot, commitFlags, appConfig)
	configCmd.AddCommand(createConfigExampleCommand(projectRoot))

	rootCmd.AddCommand(configCmd)

	// Add independent config-example command (same features as config example)
	// 添加独立的 config-example 命令（与 config example 功能相同）
	configExampleIndependentCmd := createConfigExampleIndependentCommand(projectRoot)

	rootCmd.AddCommand(configExampleIndependentCmd)

	// Execute the CLI application
	// 执行 CLI 应用程序
	must.Done(rootCmd.Execute())
}

// createRootCommand creates the main root command with flags
// 创建主根命令和标志
func createRootCommand(projectRoot string, commitFlags *commitmate.CommitFlags, appConfig *AppConfig) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "go-commit",
		Short: "Smart Git commit app with Go code formatting",
		Long:  "go-commit is a Git commit app that auto formats changed Go code and provides flexible commit options",
		Run: func(cmd *cobra.Command, args []string) {
			// Load signature config if config file is provided
			// 如果提供了配置文件则加载签名配置
			if appConfig.ConfigPath != "" {
				commitFlags.ApplyProjectConfig(projectRoot, commitmate.LoadConfig(appConfig.ConfigPath))
			}

			// Validate commit flags and show warnings
			// 验证提交标志并显示警告
			if warnings := commitFlags.ValidateFlags(); len(warnings) > 0 {
				for _, warning := range warnings {
					zaplog.SUG.Warnf("⚠️  %s", warning)
				}
			}

			must.Done(commitmate.GitCommit(projectRoot, commitFlags))
		},
	}

	// Configure command line flags to enable commit customization
	// 配置用于提交自定义的命令行标志
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Username, "username", "u", "", "git username")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Message, "message", "m", "", "commit message")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsAmend, "amend", "a", false, "amend to the previous commit")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsForce, "force", "f", false, "force amend even pushed to remote")
	rootCmd.PersistentFlags().StringVar(&commitFlags.Mailbox, "mailbox", "", "mailbox address (preferred)")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Eddress, "eddress", "e", "", "mailbox address (fallback)")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.NoCommit, "no-commit", false, "stage changes without committing")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.FormatGo, "format-go", false, "format changed go files")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.AutoSign, "auto-sign", false, "auto-use git config signing info when unset")
	rootCmd.PersistentFlags().StringVarP(&appConfig.ConfigPath, "config", "c", "", "path to go-commit configuration file")

	return rootCmd
}

// createConfigCommand creates the config subcommand to manage configurations
// 创建用于配置管理的 config 子命令
func createConfigCommand(projectRoot string, commitFlags *commitmate.CommitFlags, appConfig *AppConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Configuration management of go-commit",
		Long:  "Manage go-commit configurations, validate existing configs, and generate templates",
		Run: func(cmd *cobra.Command, args []string) {
			// Validate config file path is provided
			// 验证提供了配置文件路径
			if appConfig.ConfigPath == "" {
				zaplog.SUG.Panicln("missing config path. use -c flag")
			}

			zaplog.SUG.Debugln("config path:", appConfig.ConfigPath)

			// Load and use configuration
			// 加载并使用配置
			config := commitmate.LoadConfig(appConfig.ConfigPath)
			zaplog.SUG.Debugln("config items:", neatjsons.S(config))

			commitFlags.ApplyProjectConfig(projectRoot, config)
			zaplog.SUG.Debugln("commit flags:", neatjsons.S(commitFlags))
		},
	}
}

// createConfigExampleCommand creates the config example subcommand
// 创建 config example 子命令
func createConfigExampleCommand(projectRoot string) *cobra.Command {
	return &cobra.Command{
		Use:   "example",
		Short: "Generate configuration template with current project",
		Long:  "Generate a go-commit configuration template based on current project's Git remote URL",
		Run: func(cmd *cobra.Command, args []string) {
			previewConfigTemplate(projectRoot)
		},
	}
}

// createConfigExampleIndependentCommand creates the independent config-example command
// 创建独立的 config-example 命令
func createConfigExampleIndependentCommand(projectRoot string) *cobra.Command {
	return &cobra.Command{
		Use:   "config-example",
		Short: "Generate configuration template with current project",
		Long:  "Generate a go-commit configuration template based on current project's Git remote URL",
		Run: func(cmd *cobra.Command, args []string) {
			previewConfigTemplate(projectRoot)
		},
	}
}

// previewConfigTemplate previews configuration template based on current project
// Generates and displays a sample configuration from project Git remotes
// Outputs formatted JSON template to enable simple copying and customization
//
// previewConfigTemplate 预览当前项目的配置模板
// 基于项目 Git 远程生成并显示示例配置
// 输出格式化的 JSON 模板以便复制和自定义
func previewConfigTemplate(projectRoot string) {
	configTemplate := commitmate.GenerateConfigTemplate(projectRoot)

	zaplog.SUG.Infoln("Generated configuration template:")
	// Output template as formatted JSON
	// 将模板输出为格式化的 JSON
	zaplog.SUG.Infoln(neatjsons.S(configTemplate))
	zaplog.SUG.Infoln("Save this template to a file (e.g., go-commit-config.json).")
}
