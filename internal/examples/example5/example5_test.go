// Package example5_test demonstrates minimal universal configuration patterns
// Tests single signature matching across all Git remotes with wildcard patterns
// Validates simplified configuration approach suitable in single-identity scenarios
//
// example5_test 演示最小化通用配置模式
// 测试使用通配符模式在所有 Git 远程上的单一签名匹配
// 验证适用于单一身份场景的简化配置方式
package example5_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-commit/commitmate"
	"github.com/go-xlan/gogit"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestLoadExample5Config validates minimal configuration file loading
// Tests that configuration contains single universal signature with wildcard pattern
//
// TestLoadExample5Config 验证最小化配置文件加载
// 测试配置包含单个带通配符模式的通用签名
func TestLoadExample5Config(t *testing.T) {
	configPath := runpath.PARENT.Join("minimal-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 1)
	require.Equal(t, "main-identity", config.Signatures[0].Name)
	require.Equal(t, "developer", config.Signatures[0].Username)
	require.Equal(t, "developer@example.com", config.Signatures[0].Eddress)
	require.Equal(t, []string{"*"}, config.Signatures[0].RemotePatterns)
}

// TestExample5UniversalMatching validates universal pattern matching across all remotes
// Tests that single wildcard pattern matches GitHub, GitLab, Bitbucket, and custom domains
//
// TestExample5UniversalMatching 验证跨所有远程的通用模式匹配
// 测试单个通配符模式匹配 GitHub、GitLab、Bitbucket 和自定义域名
func TestExample5UniversalMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("minimal-config.json")
	config := commitmate.LoadConfig(configPath)

	// Test various remote URLs - all should match the universal pattern
	testCases := []string{
		"git@github.com:user/repo.git",
		"https://github.com/user/repo.git",
		"git@gitlab.com:user/repo.git",
		"https://gitlab.com/user/repo.git",
		"git@bitbucket.org:user/repo.git",
		"git@custom.domain.com:team/project.git",
		"https://custom.company.internal/project.git",
		"git@localhost:test/repo.git",
		"file:///local/path/repo.git",
	}

	for _, remoteURL := range testCases {
		signature := config.MatchSignature(remoteURL)
		require.NotNil(t, signature, "Should match remote: %s", remoteURL)
		require.Equal(t, "main-identity", signature.Name)
		require.Equal(t, "developer", signature.Username)
		require.Equal(t, "developer@example.com", signature.Eddress)
	}
}

// TestExample5GitCommitGitHub tests commit workflow with GitHub using universal configuration
// Creates GitHub repository, applies minimal configuration, and verifies commit metadata
//
// TestExample5GitCommitGitHub 测试使用通用配置的 GitHub 提交工作流程
// 创建 GitHub 仓库，应用最小化配置，并验证提交元数据
func TestExample5GitCommitGitHub(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-github-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:developer/simple-app.git"))

	testFile := filepath.Join(tempDIR, "app.go")
	testContent := `package main

import"fmt"

func main(){
	fmt.Println("Simple application with minimal config")
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("minimal-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Create simple application",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "developer", flags.Username)
	require.Equal(t, "developer@example.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	client := rese.P1(gogit.New(tempDIR))
	status := rese.V1(client.Status())
	require.Empty(t, status)

	output := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%s"))
	require.Equal(t, "Create simple application", string(output))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "developer", string(authorName))
	require.Equal(t, "developer@example.com", string(authorEmail))
}

// TestExample5GitCommitGitLab tests commit workflow with GitLab using universal configuration
// Creates GitLab repository, applies minimal configuration, and verifies commit metadata
//
// TestExample5GitCommitGitLab 测试使用通用配置的 GitLab 提交工作流程
// 创建 GitLab 仓库，应用最小化配置，并验证提交元数据
func TestExample5GitCommitGitLab(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-gitlab-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@gitlab.com:developer/utility-lib.git"))

	testFile := filepath.Join(tempDIR, "utils.go")
	testContent := `package utils

import"strings"

func CleanString(s string)string{
	return strings.TrimSpace(s)
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("minimal-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add string utility functions",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "developer", flags.Username)
	require.Equal(t, "developer@example.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "developer", string(authorName))
	require.Equal(t, "developer@example.com", string(authorEmail))
}

// TestExample5GitCommitCustomDomain tests commit workflow with custom domain using universal configuration
// Creates custom domain repository, applies minimal configuration, and verifies commit metadata
//
// TestExample5GitCommitCustomDomain 测试使用通用配置的自定义域名提交工作流程
// 创建自定义域名仓库，应用最小化配置，并验证提交元数据
func TestExample5GitCommitCustomDomain(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-custom-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@git.company.internal:team/project.git"))

	testFile := filepath.Join(tempDIR, "config.go")
	testContent := `package config

type Config struct{
	Host string
	Port int
}

func Default()Config{
	return Config{
		Host:"localhost",
		Port:8080,
	}
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("minimal-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add default configuration struct",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "developer", flags.Username)
	require.Equal(t, "developer@example.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "developer", string(authorName))
	require.Equal(t, "developer@example.com", string(authorEmail))
}

// TestExample5MinimalConfigAdvantages validates benefits of minimal universal configuration
// Tests simplicity, universal coverage, and ease of use with single signature approach
//
// TestExample5MinimalConfigAdvantages 验证最小化通用配置的优势
// 测试简洁性、通用覆盖和单签名方式的易用性
func TestExample5MinimalConfigAdvantages(t *testing.T) {
	configPath := runpath.PARENT.Join("minimal-config.json")
	config := commitmate.LoadConfig(configPath)

	// Verify the minimal config covers all scenarios
	testRemotes := []string{
		"git@github.com:opensource/project.git",
		"git@gitlab.company.com:internal/service.git",
		"https://bitbucket.org:personal/tool.git",
		"git@dev.example.org:experimental/prototype.git",
	}

	for _, remote := range testRemotes {
		signature := config.MatchSignature(remote)
		require.NotNil(t, signature)
		require.Equal(t, "main-identity", signature.Name)
		require.Equal(t, "developer", signature.Username)
		require.Equal(t, "developer@example.com", signature.Eddress)
	}

	// Verify simplicity - just one signature needed
	require.Len(t, config.Signatures, 1)

	// Verify universal coverage - single pattern covers everything
	require.Len(t, config.Signatures[0].RemotePatterns, 1)
	require.Equal(t, "*", config.Signatures[0].RemotePatterns[0])
}
