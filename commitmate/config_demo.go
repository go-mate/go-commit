package commitmate

import (
	"strings"

	"github.com/yyle88/tern/zerotern"
	"github.com/yyle88/zaplog"
)

// GenerateConfigTemplate generates a configuration template based on current project
// Analyzes current Git remote URL and creates a suggested configuration template
// Provides starting configuration with placeholders in username and mailbox settings
// Outputs JSON template to stdout allowing simple copying and customization
//
// GenerateConfigTemplate 为当前项目生成配置模板
// 分析当前 Git 远程 URL 并创建建议的配置模板
// 提供带有用户名和邮箱设置占位符的启动配置
// 将 JSON 模板输出到 stdout 以便复制和自定义
func GenerateConfigTemplate(projectRoot string) *CommitConfig {
	zaplog.SUG.Debugln("generating config template based on project:", projectRoot)

	signatureConfig := generateSignatureTemplate(projectRoot)
	return &CommitConfig{
		Signatures: []*SignatureConfig{
			signatureConfig,
		},
	}
}

// generateSignatureTemplate creates a signature configuration template based on project Git remotes
// Extracts remote information and generates appropriate configuration patterns
//
// generateSignatureTemplate 基于项目 Git 远程创建签名配置模板
// 提取远程信息并生成合适的配置模式
func generateSignatureTemplate(projectRoot string) *SignatureConfig {
	// Get remote URL using shared function
	// 使用共享函数获取远程 URL
	remoteURL := getOriginRemoteURL(projectRoot)

	// Use default template when no remote URL is found
	// 假如没有找到远端的地址就给个默认的样例就行
	remoteURL = zerotern.VV(remoteURL, "git@github.com:username/repo.git")

	// Create template configuration based on detected remote
	// 基于检测到的远程创建模板配置
	return &SignatureConfig{
		Name:           generateConfigName(remoteURL),
		Username:       "your-username",
		Eddress:        "your-email@example.com",
		RemotePatterns: []string{generatePattern(remoteURL)},
	}
}

// generateConfigName creates a descriptive name of the configuration based on remote URL
// 基于远程 URL 为配置创建描述性名称
func generateConfigName(remoteURL string) string {
	slashIdx := strings.LastIndex(remoteURL, "/")
	if slashIdx == -1 {
		return "git-config"
	}
	repoName := remoteURL[slashIdx+1:]
	repoName = strings.TrimSuffix(repoName, ".git")
	return repoName + "-git-config"
}

// generatePattern creates a wildcard pattern based on the remote URL
// 基于远程 URL 创建通配符模式
func generatePattern(remoteURL string) string {
	slashIdx := strings.LastIndex(remoteURL, "/")
	if slashIdx == -1 {
		return remoteURL
	}
	return remoteURL[:slashIdx+1] + "*"
}
