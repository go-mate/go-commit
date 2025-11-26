// Package example4_test demonstrates advanced pattern matching scenarios
// Tests exact matches, team-specific patterns, subdomain wildcards, and pattern priority
// Validates sophisticated matching algorithms in complex enterprise environments
//
// example4_test 演示高级模式匹配场景
// 测试精确匹配、团队特定模式、子域名通配符和模式优先级
// 验证复杂企业环境中的复杂匹配算法
package example4_test

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

// TestLoadExample4Config validates advanced pattern configuration file loading
// Tests that configuration contains expected pattern entries with correct structure
//
// TestLoadExample4Config 验证高级模式配置文件加载
// 测试配置包含预期的模式条目和正确结构
func TestLoadExample4Config(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 4)
	require.Equal(t, "exact-repo-match", config.Signatures[0].Name)
	require.Equal(t, "team-specific", config.Signatures[1].Name)
	require.Equal(t, "multi-subdomain", config.Signatures[2].Name)
	require.Equal(t, "protocol-flexible", config.Signatures[3].Name)
}

// TestExample4ExactRepoMatching validates exact repository pattern matching
// Tests highest priority matching with specific repository URLs
//
// TestExample4ExactRepoMatching 验证精确仓库模式匹配
// 测试特定仓库 URL 的最高优先级匹配
func TestExample4ExactRepoMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	// Test exact repository match - highest priority
	signature := config.MatchSignature("git@github.com:company/critical-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "exact-repo-match", signature.Name)
	require.Equal(t, "diana.lead", signature.Username)
	require.Equal(t, "diana.lead@critical-project.com", signature.Eddress)

	// Test HTTPS exact match
	signature = config.MatchSignature("https://github.com/company/critical-project.git")
	require.NotNil(t, signature)
	require.Equal(t, "exact-repo-match", signature.Name)
}

// TestExample4TeamSpecificMatching validates team-specific pattern matching
// Tests matching repositories within specific team namespaces
//
// TestExample4TeamSpecificMatching 验证团队特定模式匹配
// 测试特定团队命名空间内的仓库匹配
func TestExample4TeamSpecificMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	// Test team-alpha matching
	signature := config.MatchSignature("git@github.com:company/team-alpha/backend.git")
	require.NotNil(t, signature)
	require.Equal(t, "team-specific", signature.Name)
	require.Equal(t, "diana.dev", signature.Username)
	require.Equal(t, "diana@company.com", signature.Eddress)

	// Test team-beta matching
	signature = config.MatchSignature("https://github.com/company/team-beta/frontend.git")
	require.NotNil(t, signature)
	require.Equal(t, "team-specific", signature.Name)

	// Non-team project should not match team-specific
	signature = config.MatchSignature("git@github.com:company/general-project.git")
	require.NotNil(t, signature)
	require.NotEqual(t, "team-specific", signature.Name)
}

// TestExample4MultiSubdomainMatching validates wildcard subdomain pattern matching
// Tests matching across multiple subdomains with wildcard patterns
//
// TestExample4MultiSubdomainMatching 验证通配符子域名模式匹配
// 测试使用通配符模式跨多个子域名的匹配
func TestExample4MultiSubdomainMatching(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	// Test wildcard subdomain matching
	signature := config.MatchSignature("git@dev.multi-env.com:project/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "multi-subdomain", signature.Name)
	require.Equal(t, "diana.ops", signature.Username)
	require.Equal(t, "diana.ops@multi-env.com", signature.Eddress)

	// Test HTTPS subdomain matching
	signature = config.MatchSignature("https://staging.multi-env.com/project/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "multi-subdomain", signature.Name)

	// Test internal subdomain matching
	signature = config.MatchSignature("git@dev.project.internal:service/api.git")
	require.NotNil(t, signature)
	require.Equal(t, "multi-subdomain", signature.Name)
}

// TestExample4PatternPriority validates pattern matching priority rules
// Tests that more specific patterns override less specific ones
//
// TestExample4PatternPriority 验证模式匹配优先级规则
// 测试更具体的模式覆盖不太具体的模式
func TestExample4PatternPriority(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	// Exact match should win over team-specific
	signature := config.MatchSignature("git@github.com:company/critical-project.git")
	require.Equal(t, "exact-repo-match", signature.Name)

	// Team-specific should win over protocol-flexible for team projects
	signature = config.MatchSignature("git@github.com:company/team-alpha/service.git")
	require.Equal(t, "team-specific", signature.Name)

	// Unknown patterns should fall back to protocol-flexible
	signature = config.MatchSignature("git@unknown.domain.com:any/repo.git")
	require.Equal(t, "protocol-flexible", signature.Name)
	require.Equal(t, "diana", signature.Username)
}

// TestExample4GitCommitExactRepo tests commit workflow with exact repository matching
// Creates exact match repository, applies configuration, and verifies commit metadata
//
// TestExample4GitCommitExactRepo 测试精确仓库匹配的提交工作流程
// 创建精确匹配仓库，应用配置，并验证提交元数据
func TestExample4GitCommitExactRepo(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-exact-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:company/critical-project.git"))

	testFile := filepath.Join(tempDIR, "security.go")
	testContent := `package security

import"crypto/rand"

func GenerateToken()[]byte{
	token:=make([]byte,32)
	rand.Read(token)
	return token
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("advanced-patterns.json")
	flags := &commitmate.CommitFlags{
		Message:  "Implement secure token generation",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "diana.lead", flags.Username)
	require.Equal(t, "diana.lead@critical-project.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.lead", string(authorName))
	require.Equal(t, "diana.lead@critical-project.com", string(authorEmail))
}

// TestExample4GitCommitTeamSpecific tests commit workflow with team-specific matching
// Creates team repository, applies configuration, and verifies commit metadata
//
// TestExample4GitCommitTeamSpecific 测试团队特定匹配的提交工作流程
// 创建团队仓库，应用配置，并验证提交元数据
func TestExample4GitCommitTeamSpecific(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-team-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@github.com:company/team-alpha/microservice.git"))

	testFile := filepath.Join(tempDIR, "handler.go")
	testContent := `package handler

import"net/http"

func HealthCheck(w http.ResponseWriter,r*http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("advanced-patterns.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add health check endpoint for microservice",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "diana.dev", flags.Username)
	require.Equal(t, "diana@company.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.dev", string(authorName))
	require.Equal(t, "diana@company.com", string(authorEmail))
}

// TestExample4GitCommitMultiSubdomain tests commit workflow with subdomain wildcard matching
// Creates subdomain repository, applies configuration, and verifies commit metadata
//
// TestExample4GitCommitMultiSubdomain 测试子域名通配符匹配的提交工作流程
// 创建子域名仓库，应用配置，并验证提交元数据
func TestExample4GitCommitMultiSubdomain(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-subdomain-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	execConfig := osexec.NewExecConfig().WithPath(tempDIR)
	rese.V1(execConfig.Exec("git", "init"))
	rese.V1(execConfig.Exec("git", "remote", "add", "origin", "git@staging.multi-env.com:infrastructure/deployment.git"))

	testFile := filepath.Join(tempDIR, "deploy.go")
	testContent := `package deploy

import"context"

func DeployService(ctx context.Context,name string)error{
	// Deployment logic here
	return nil
}
`
	must.Done(os.WriteFile(testFile, []byte(testContent), 0644))

	configPath := runpath.PARENT.Join("advanced-patterns.json")
	flags := &commitmate.CommitFlags{
		Message:  "Add automated deployment service",
		FormatGo: true,
	}

	config := commitmate.LoadConfig(configPath)
	flags.ApplyProjectConfig(tempDIR, config)

	// Verify project config was applied correctly
	require.Equal(t, "diana.ops", flags.Username)
	require.Equal(t, "diana.ops@multi-env.com", flags.Eddress)

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.ops", string(authorName))
	require.Equal(t, "diana.ops@multi-env.com", string(authorEmail))
}
