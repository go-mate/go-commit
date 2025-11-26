// Package example3_test demonstrates multi-client freelance configuration patterns
// Tests signature management across multiple client projects and personal work
// Validates complex identity switching in freelance and consulting scenarios
//
// example3_test 演示多客户自由职业配置模式
// 测试跨多个客户项目和个人工作的签名管理
// 验证自由职业和咨询场景中的复杂身份切换
package example3_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-commit/commitmate"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
	"github.com/yyle88/runpath"
)

// TestLoadExample3Config validates multi-client configuration file loading
// Tests that configuration contains expected client and personal signature entries
//
// TestLoadExample3Config 验证多客户配置文件加载
// 测试配置包含预期的客户和个人签名条目
func TestLoadExample3Config(t *testing.T) {
	configPath := runpath.PARENT.Join("multi-client-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 4)
	require.Equal(t, "client-alpha", config.Signatures[0].Name)
	require.Equal(t, "client-beta", config.Signatures[1].Name)
	require.Equal(t, "personal-projects", config.Signatures[2].Name)
	require.Equal(t, "consulting-work", config.Signatures[3].Name)
}

// TestExample3ClientAlphaMatching validates pattern matching when working with client alpha projects
// Tests matching GitHub and GitLab repositories specific to client alpha
//
// TestExample3ClientAlphaMatching 验证与客户 alpha 项目工作时的模式匹配
// 测试特定于客户 alpha 的 GitHub 和 GitLab 仓库匹配
func TestExample3ClientAlphaMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("multi-client-config.json")
	config := commitmate.LoadConfig(configPath)

	// Test client-alpha GitHub matching
	signature := config.MatchSignature("git@github.com:client-alpha/webapp.git")
	require.NotNil(t, signature)
	require.Equal(t, "client-alpha", signature.Name)
	require.Equal(t, "charlie.contractor", signature.Username)
	require.Equal(t, "charlie@client-alpha.com", signature.Eddress)

	// Test client-alpha GitLab matching
	signature = config.MatchSignature("git@gitlab.client-alpha.com:frontend/react-app.git")
	require.NotNil(t, signature)
	require.Equal(t, "client-alpha", signature.Name)
}

// TestExample3ClientBetaMatching validates pattern matching when working with client beta projects
// Tests matching GitLab and custom Git repositories specific to client beta
//
// TestExample3ClientBetaMatching 验证与客户 beta 项目工作时的模式匹配
// 测试特定于客户 beta 的 GitLab 和自定义 Git 仓库匹配
func TestExample3ClientBetaMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("multi-client-config.json")
	config := commitmate.LoadConfig(configPath)

	// Test client-beta GitLab matching
	signature := config.MatchSignature("git@gitlab.com:client-beta/api-service.git")
	require.NotNil(t, signature)
	require.Equal(t, "client-beta", signature.Name)
	require.Equal(t, "charlie.freelancer", signature.Username)
	require.Equal(t, "charlie@client-beta.org", signature.Eddress)

	// Test client-beta internal git matching
	signature = config.MatchSignature("git@git.client-beta.internal:backend/database.git")
	require.NotNil(t, signature)
	require.Equal(t, "client-beta", signature.Name)
}

// TestExample3PersonalProjectsMatching validates pattern matching with personal side projects
// Tests matching personal GitHub and GitLab repositories outside client work
//
// TestExample3PersonalProjectsMatching 验证个人兴趣项目的模式匹配
// 测试客户工作之外的个人 GitHub 和 GitLab 仓库匹配
func TestExample3PersonalProjectsMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("multi-client-config.json")
	config := commitmate.LoadConfig(configPath)

	// Test personal GitHub matching
	signature := config.MatchSignature("git@github.com:charlie/my-tool.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-projects", signature.Name)
	require.Equal(t, "charlie", signature.Username)
	require.Equal(t, "charlie.freelancer@protonmail.com", signature.Eddress)

	// Test personal GitLab matching
	signature = config.MatchSignature("git@gitlab.com:charlie/side-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-projects", signature.Name)
}

// TestExample3GitCommitClientAlpha tests commit workflow with client alpha signature
// Creates client alpha repository, applies configuration, and verifies commit metadata
//
// TestExample3GitCommitClientAlpha 测试客户 alpha 签名的提交工作流程
// 创建客户 alpha 仓库，应用配置，并验证提交元数据
func TestExample3GitCommitClientAlpha(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-alpha-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:client-alpha/mobile-app.git"))

	testFile := filepath.Join(tempDIR, "auth.go")
	testContent := `package auth

import"net/http"

func ValidateToken(token string)bool{
	return len(token)>0
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("multi-client-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Implement token validation for mobile auth",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "charlie.contractor", flags.Username)
	require.Equal(t, "charlie@client-alpha.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "charlie.contractor", string(authorName))
	require.Equal(t, "charlie@client-alpha.com", string(authorEmail))
}

// TestExample3GitCommitPersonalProject tests commit workflow with personal project signature
// Creates personal repository, applies configuration, and verifies commit metadata
//
// TestExample3GitCommitPersonalProject 测试个人项目签名的提交工作流程
// 创建个人仓库，应用配置，并验证提交元数据
func TestExample3GitCommitPersonalProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-personal-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:charlie/cli-tool.git"))

	testFile := filepath.Join(tempDIR, "cmd.go")
	testContent := `package main

import"flag"

func main(){
	verbose:=flag.Bool("v",false,"verbose output")
	flag.Parse()
	if*verbose{
		println("Verbose mode enabled")
	}
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("multi-client-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add verbose flag support to CLI tool",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "charlie", flags.Username)
	require.Equal(t, "charlie.freelancer@protonmail.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "charlie", string(authorName))
	require.Equal(t, "charlie.freelancer@protonmail.com", string(authorEmail))
}

// TestExample3GitCommitConsultingWork tests commit workflow with consulting work signature
// Creates consulting repository, applies configuration, and verifies commit metadata
//
// TestExample3GitCommitConsultingWork 测试咨询工作签名的提交工作流程
// 创建咨询仓库，应用配置，并验证提交元数据
func TestExample3GitCommitConsultingWork(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-consulting-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@git.charlie-consulting.dev:client/architecture.git"))

	testFile := filepath.Join(tempDIR, "design.go")
	testContent := `package design

type Architecture struct{
	Components[]string
	Database   string
}

func NewArchitecture()Architecture{
	return Architecture{
		Components:[]string{"api","web","worker"},
		Database:"postgresql",
	}
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("multi-client-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Design initial system architecture",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "consulting-charlie", flags.Username)
	require.Equal(t, "hello@charlie-consulting.dev", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "consulting-charlie", string(authorName))
	require.Equal(t, "hello@charlie-consulting.dev", string(authorEmail))
}
