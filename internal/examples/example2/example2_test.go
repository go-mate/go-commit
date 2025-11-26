// Package example2_test demonstrates work vs personal project separation
// Tests signature matching between enterprise work projects and personal side projects
// Validates automatic identity switching based on Git remote patterns
//
// example2_test 演示工作项目与个人项目的分离
// 测试企业工作项目和个人兴趣项目之间的签名匹配
// 验证基于 Git 远程模式的自动身份切换
package example2_test

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

// TestLoadExample2Config validates configuration file loading with work and personal signatures
// Tests that configuration contains correct signature entries with proper names
//
// TestLoadExample2Config 验证配置文件加载和工作及个人签名
// 测试配置包含正确的签名条目和适当的名称
func TestLoadExample2Config(t *testing.T) {
	configPath := runpath.PARENT.Join("simple-personal-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 3)
	require.Equal(t, "work-projects", config.Signatures[0].Name)
	require.Equal(t, "personal-github", config.Signatures[1].Name)
	require.Equal(t, "personal-gitlab", config.Signatures[2].Name)
}

// TestExample2WorkPatternMatching validates pattern matching in work and personal contexts
// Tests matching enterprise GitHub, work GitLab, and personal GitHub/GitLab projects
//
// TestExample2WorkPatternMatching 验证工作和个人环境中的模式匹配
// 测试企业 GitHub、工作 GitLab 和个人 GitHub/GitLab 项目的匹配
func TestExample2WorkPatternMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("simple-personal-config.json")
	config := commitmate.LoadConfig(configPath)

	// Test work enterprise GitHub matching
	signature := config.MatchSignature("git@github.enterprise.company.com:team/project.git")
	require.NotNil(t, signature)
	require.Equal(t, "work-projects", signature.Name)
	require.Equal(t, "bob.smith", signature.Username)
	require.Equal(t, "bob.smith@company.com", signature.Eddress)

	// Test work GitLab matching
	signature = config.MatchSignature("git@gitlab.company.com:team/project.git")
	require.NotNil(t, signature)
	require.Equal(t, "work-projects", signature.Name)

	// Test personal GitHub matching
	signature = config.MatchSignature("git@github.com:bob-dev/my-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-github", signature.Name)
	require.Equal(t, "bob-dev", signature.Username)
	require.Equal(t, "bob.personal@gmail.com", signature.Eddress)

	// Test personal GitLab matching
	signature = config.MatchSignature("git@gitlab.com:bob-dev/my-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-gitlab", signature.Name)
	require.Equal(t, "bob-dev", signature.Username)
}

// TestExample2GitCommitWorkProject tests commit workflow with work project signature
// Creates enterprise repository, applies work configuration, and verifies commit metadata
//
// TestExample2GitCommitWorkProject 测试工作项目签名的提交工作流程
// 创建企业仓库，应用工作配置，并验证提交元数据
func TestExample2GitCommitWorkProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example2-work-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.enterprise.company.com:team/backend.git"))

	testFile := filepath.Join(tempDIR, "service.go")
	testContent := `package service

import"context"

func ProcessData(ctx context.Context)error{
	return nil
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("simple-personal-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add service layer for data processing",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "bob.smith", flags.Username)
	require.Equal(t, "bob.smith@company.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	client := rese.P1(gogit.New(tempDIR))
	status := rese.V1(client.Status())
	require.Empty(t, status)

	output := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%s"))
	require.Equal(t, "Add service layer for data processing", string(output))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "bob.smith", string(authorName))
	require.Equal(t, "bob.smith@company.com", string(authorEmail))
}

// TestExample2GitCommitPersonalProject tests commit workflow with personal project signature
// Creates personal repository, applies personal configuration, and verifies commit metadata
//
// TestExample2GitCommitPersonalProject 测试个人项目签名的提交工作流程
// 创建个人仓库，应用个人配置，并验证提交元数据
func TestExample2GitCommitPersonalProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example2-personal-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:bob-dev/hobby-project.git"))

	testFile := filepath.Join(tempDIR, "main.go")
	testContent := `package main

import"fmt"

func main(){
	fmt.Println("My personal hobby project!")
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("simple-personal-config.json")
	flags := &commitmate.CommitFlags{
		Message:  "Initial hobby project setup",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "bob-dev", flags.Username)
	require.Equal(t, "bob.personal@gmail.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "bob-dev", string(authorName))
	require.Equal(t, "bob.personal@gmail.com", string(authorEmail))
}
