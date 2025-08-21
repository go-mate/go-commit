package example1_test

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

func TestLoadExample1Config(t *testing.T) {
	// Get path to the config file in same DIR
	// 获取同一 DIR 中配置文件的路径
	configPath := runpath.PARENT.Join("go-commit-config.example.json")

	// Load the configuration
	// 加载配置
	config := commitmate.LoadConfig(configPath)

	// Verify basic structure
	// 验证基本结构
	require.Len(t, config.Signatures, 4)
	require.Equal(t, "personal-github", config.Signatures[0].Name)
	require.Equal(t, "personal-gitlab", config.Signatures[1].Name)
	require.Equal(t, "github-contributions", config.Signatures[2].Name)
	require.Equal(t, "fallback-default", config.Signatures[3].Name)
}

func TestExample1PatternMatching(t *testing.T) {
	// Load config for testing
	// 加载配置进行测试
	configPath := runpath.PARENT.Join("go-commit-config.example.json")
	config := commitmate.LoadConfig(configPath)

	// Test GitHub personal projects matching
	// 测试 GitHub 个人项目匹配
	signature := config.MatchSignature("git@github.com:alice/my-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-github", signature.Name)
	require.Equal(t, "alice", signature.Username)

	// Test GitLab personal projects matching
	// 测试 GitLab 个人项目匹配
	signature = config.MatchSignature("git@gitlab.com:alice/my-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "personal-gitlab", signature.Name)
	require.Equal(t, "alice", signature.Username)

	// Test open source contributions matching
	// 测试开源贡献匹配
	signature = config.MatchSignature("git@github.com:golang/go.git")
	require.NotNil(t, signature)
	require.Equal(t, "github-contributions", signature.Name)
	require.Equal(t, "alice-oss", signature.Username)

	// Test fallback matching for unknown remotes
	// 测试未知远程的兜底匹配
	signature = config.MatchSignature("git@unknown.com:some/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "fallback-default", signature.Name)
	require.Equal(t, "alice", signature.Username)
}

func TestExample1GitCommitWithSignature(t *testing.T) {
	// Create temp DIR for test repository
	// 为测试仓库创建临时 DIR
	tempDIR := rese.V1(os.MkdirTemp("", "example1-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

	// Initialize git repository using osexec
	// 使用 osexec 初始化 git 仓库
	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))

	// Set up git remote for GitHub personal project
	// 为 GitHub 个人项目设置 git 远程
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:alice/test-project.git"))

	// Create test file
	// 创建测试文件
	testFile := filepath.Join(tempDIR, "main.go")
	testContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello from alice's personal project!")
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	// Get config path and perform commit
	// 获取配置路径并执行提交
	configPath := runpath.PARENT.Join("go-commit-config.example.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add main.go with hello message",
		FormatGo: true,
	}

	// Apply project config to commit flags
	// 将项目配置应用到提交标志
	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify config was applied correctly
	// 验证配置被正确应用
	require.Equal(t, "alice", flags.Username)
	require.Equal(t, "alice.dev@gmail.com", flags.Eddress)

	// Perform the commit
	// 执行提交
	must.Done(commitmate.GitCommit(tempDIR, flags))

	// Use gogit client to verify commit
	// 使用 gogit 客户端验证提交
	client := rese.P1(gogit.New(tempDIR))

	// Verify working DIR is clean after commit
	// 验证提交后工作 DIR 是干净的
	status := rese.V1(client.Status())
	require.Empty(t, status)

	// Verify commit was created with correct message
	// 验证提交已创建且包含正确消息
	output := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%s"))
	require.Equal(t, "Add main.go with hello message", string(output))

	// Verify commit author
	// 验证提交作者
	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "alice", string(authorName))
	require.Equal(t, "alice.dev@gmail.com", string(authorEmail))
}
