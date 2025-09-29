// Configuration template testing for commitmate package
// Tests configuration template generation functionality with various Git remote scenarios
// Validates template creation for different repository configurations and edge cases
//
// commitmate 包的配置模板测试
// 测试各种 Git 远程场景下的配置模板生成功能
// 验证不同仓库配置和边界情况下的模板创建
package commitmate

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
)

func TestGenerateConfigTemplate(t *testing.T) {
	// Create temporary DIR for test repository
	// 为测试仓库创建临时 DIR
	tempDIR := rese.V1(os.MkdirTemp("", "test-config-template-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Initialize test git repository
	// 初始化测试 git 仓库
	repo := rese.V1(git.PlainInit(tempDIR, false))

	// Add origin remote
	// 添加 origin 远程
	remoteURL := "git@github.com:test-user/test-repo.git"
	rese.P1(repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	}))

	// Generate configuration template
	// 生成配置模板
	configTemplate := GenerateConfigTemplate(tempDIR)

	// Verify template structure
	// 验证模板结构
	require.NotNil(t, configTemplate)
	require.Len(t, configTemplate.Signatures, 1)

	signature := configTemplate.Signatures[0]
	require.Equal(t, "test-repo-git-config", signature.Name)
	require.Equal(t, "your-username", signature.Username)
	require.Equal(t, "your-email@example.com", signature.Eddress)
	require.Len(t, signature.RemotePatterns, 1)
	require.Equal(t, "git@github.com:test-user/*", signature.RemotePatterns[0])
}

func TestGenerateConfigTemplateWithoutRemote(t *testing.T) {
	// Create temporary DIR for test repository
	// 为测试仓库创建临时 DIR
	tempDIR := rese.V1(os.MkdirTemp("", "test-config-no-remote-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Initialize test git repository without remotes
	// 初始化没有远程的测试 git 仓库
	rese.V1(git.PlainInit(tempDIR, false))

	// Generate configuration template
	// 生成配置模板
	configTemplate := GenerateConfigTemplate(tempDIR)

	// Verify template uses default values
	// 验证模板使用默认值
	require.NotNil(t, configTemplate)
	require.Len(t, configTemplate.Signatures, 1)

	signature := configTemplate.Signatures[0]
	require.Equal(t, "repo-git-config", signature.Name)
	require.Equal(t, "your-username", signature.Username)
	require.Equal(t, "your-email@example.com", signature.Eddress)
	require.Len(t, signature.RemotePatterns, 1)
	require.Equal(t, "git@github.com:username/*", signature.RemotePatterns[0])
}

func TestGenerateConfigTemplateWithMultipleRemotes(t *testing.T) {
	// Create temporary DIR for test repository
	// 为测试仓库创建临时 DIR
	tempDIR := rese.V1(os.MkdirTemp("", "test-config-multiple-*"))
	defer func() {
		must.Done(os.RemoveAll(tempDIR))
	}()

	// Initialize test git repository
	// 初始化测试 git 仓库
	repo := rese.V1(git.PlainInit(tempDIR, false))

	// Add multiple remotes - origin should be preferred
	// 添加多个远程 - origin 应该被优先选择
	upstreamURL := "git@github.com:upstream/repo.git"
	originURL := "git@github.com:myuser/repo.git"

	rese.P1(repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{upstreamURL},
	}))

	rese.P1(repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{originURL},
	}))

	// Generate configuration template
	// 生成配置模板
	configTemplate := GenerateConfigTemplate(tempDIR)

	// Verify template uses origin remote
	// 验证模板使用 origin 远程
	require.NotNil(t, configTemplate)
	require.Len(t, configTemplate.Signatures, 1)

	signature := configTemplate.Signatures[0]
	require.Equal(t, "repo-git-config", signature.Name)
	require.Equal(t, "your-username", signature.Username)
	require.Equal(t, "your-email@example.com", signature.Eddress)
	require.Len(t, signature.RemotePatterns, 1)
	require.Equal(t, "git@github.com:myuser/*", signature.RemotePatterns[0])
}

// TestGenerateConfigName validates configuration name generation from various URL formats
// Tests different URL patterns and edge cases for consistent naming
//
// TestGenerateConfigName 验证从各种 URL 格式生成配置名称
// 测试不同的 URL 模式和边界情况以保证命名一致性
func TestGenerateConfigName(t *testing.T) {
	t.Run("GitHub SSH URL", func(t *testing.T) {
		result := generateConfigName("git@github.com:user/repo.git")
		require.Equal(t, "repo-git-config", result)
	})
	t.Run("GitHub HTTPS URL", func(t *testing.T) {
		result := generateConfigName("https://github.com/user/repo.git")
		require.Equal(t, "repo-git-config", result)
	})
	t.Run("URL without .git suffix", func(t *testing.T) {
		result := generateConfigName("git@github.com:user/repo")
		require.Equal(t, "repo-git-config", result)
	})
	t.Run("URL without slash", func(t *testing.T) {
		result := generateConfigName("invalid-url")
		require.Equal(t, "git-config", result)
	})
}

func TestGeneratePattern(t *testing.T) {
	t.Run("GitHub SSH URL", func(t *testing.T) {
		result := generatePattern("git@github.com:user/repo.git")
		require.Equal(t, "git@github.com:user/*", result)
	})
	t.Run("GitHub HTTPS URL", func(t *testing.T) {
		result := generatePattern("https://github.com/user/repo.git")
		require.Equal(t, "https://github.com/user/*", result)
	})
	t.Run("URL without slash", func(t *testing.T) {
		result := generatePattern("invalid-url")
		require.Equal(t, "invalid-url", result)
	})
}
