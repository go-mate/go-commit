package commitmate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mate/go-commit/internal/utils"
	"github.com/go-xlan/gitgo"
	"github.com/go-xlan/gogit"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
)

// setupTestRepo creates a temporary git repository for testing purposes
// Environment setup must succeed, so we use rese/must for all operations
// Returns the temporary DIR path and cleanup function
//
// setupTestRepo 为测试目的创建临时 git 仓库
// 环境设置必须成功，因此我们对所有操作使用 rese/must
// 返回临时 DIR 路径和清理函数
func setupTestRepo() (string, func()) {
	// Create temp DIR - must succeed
	tempDIR := rese.V1(os.MkdirTemp("", "go-commit-test-*"))

	// Initialize git repository using gitgo - must succeed
	gcm := gitgo.New(tempDIR)
	gcm.Init().MustDone()

	// Configure git username and mailbox - must succeed
	rese.V1(osexec.NewExecConfig().WithPath(tempDIR).Exec("git", "config", "user.name", "Test Username"))
	rese.V1(osexec.NewExecConfig().WithPath(tempDIR).Exec("git", "config", "user.email", "test@example.com"))

	// Create initial commit to make it a valid repo - must succeed
	testFile := filepath.Join(tempDIR, "README.md")
	must.Done(os.WriteFile(testFile, []byte("# Test Repo\n"), 0644))

	// Add and commit initial file - must succeed
	gcm.Add().Commit("Initial commit").MustDone()

	// Return cleanup function
	cleanup := func() {
		must.Done(os.RemoveAll(tempDIR))
	}

	return tempDIR, cleanup
}

// TestGitCommit_NoChanges verifies behavior when repository has no uncommitted changes
// Should complete successfully without creating new commits
//
// TestGitCommit_NoChanges 验证仓库没有未提交更改时的行为
// 应该成功完成而不创建新提交
func TestGitCommit_NoChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "No changes commit",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))
}

func TestGitCommit_WithNewFile(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file - must succeed for test setup
	newFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(newFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Add test file",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify file was committed - test verification uses require
	client := rese.P1(gogit.New(tempDIR))
	status := rese.V1(client.Status())
	require.Empty(t, status)
}

func TestGitCommit_NoCommitFlag(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file - must succeed for test setup
	newFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(newFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Should not commit",
		NoCommit: true, // Set NoCommit flag
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify file was NOT committed (still in staging) - test verification uses require
	client := rese.P1(gogit.New(tempDIR))
	status := rese.V1(client.Status())
	require.NotEmpty(t, status)
}

func TestGitCommit_WithGoFileFormatting(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a Go file with formatting issues - must succeed for test setup
	goFile := filepath.Join(tempDIR, "main.go")
	unformattedCode := "package main\n\nimport\"fmt\"\n\nfunc main(){\nfmt.Println(\"hello\")\n}\n"
	must.Done(os.WriteFile(goFile, []byte(unformattedCode), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Add Go file with formatting",
		NoCommit: false,
		FormatGo: true, // Enable Go formatting
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify file was formatted and committed - test verification uses require
	formattedContent := rese.V1(os.ReadFile(goFile))

	// Check that the file was formatted (spaces around import)
	require.Contains(t, string(formattedContent), "import \"fmt\"")
	require.Contains(t, string(formattedContent), "func main() {")
}

func TestFormatGoFiles_SkipGeneratedFiles(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create various Go files including generated ones - must succeed for test setup
	files := map[string]string{
		"main.go":                   "package main\nimport\"fmt\"\nfunc main(){fmt.Println(\"test\")}",
		"generated.pb.go":           "package main\nimport\"fmt\"\nfunc main(){fmt.Println(\"test\")}",
		"wire_gen.go":               "package main\nimport\"fmt\"\nfunc main(){fmt.Println(\"test\")}",
		"internal/data/ent/user.go": "package ent\nimport\"fmt\"\nfunc main(){fmt.Println(\"test\")}",
	}

	for filename, content := range files {
		fullPath := filepath.Join(tempDIR, filename)
		must.Done(os.MkdirAll(filepath.Dir(fullPath), 0755))
		must.Done(os.WriteFile(fullPath, []byte(content), 0644))
	}

	// Add files to git - must succeed for test setup
	client := rese.P1(gogit.New(tempDIR))
	must.Done(client.AddAll())

	// Format files with custom filter that just allows main.go - must succeed for test setup
	allowOnlyMain := func(path string) bool {
		return filepath.Base(path) == "main.go"
	}
	must.Done(FormatChangedGoFiles(tempDIR, client, allowOnlyMain))

	// Check that just main.go was formatted (has proper spacing) - test verification uses require
	mainContent := rese.V1(os.ReadFile(filepath.Join(tempDIR, "main.go")))
	require.Contains(t, string(mainContent), "import \"fmt\"")

	// Check that generated files were NOT formatted (still have formatting issues) - test verification uses require
	pbContent := rese.V1(os.ReadFile(filepath.Join(tempDIR, "generated.pb.go")))
	require.Contains(t, string(pbContent), "import\"fmt\"")
}

func TestGitCommit_AmendCommit(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "amend_test.txt")
	must.Done(os.WriteFile(testFile, []byte("original content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Original commit",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Modify the file - must succeed for test setup
	must.Done(os.WriteFile(testFile, []byte("amended content"), 0644))

	// Amend the commit
	flags.Message = "Amended commit message"
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the commit was amended - test verification uses require
	client := rese.P1(gogit.New(tempDIR))
	status := rese.V1(client.Status())
	require.Empty(t, status)
}

func TestGitCommit_AmendMessageWithoutCodeChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(testFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Original message",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the original commit
	client := rese.P1(gogit.New(tempDIR))
	headRef := rese.V1(client.Repo().Head())
	headCommit := rese.P1(client.Repo().CommitObject(headRef.Hash()))
	require.Equal(t, "Original message", headCommit.Message)

	// Amend just the message without any file changes
	flags.Message = "Updated message"
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the message was updated
	newHeadRef := rese.V1(client.Repo().Head())
	newHeadCommit := rese.P1(client.Repo().CommitObject(newHeadRef.Hash()))
	require.Equal(t, "Updated message", newHeadCommit.Message)

	// Verify status is clean
	status := rese.V1(client.Status())
	require.Empty(t, status)
}

func TestGitCommit_AmendSameMessageWithoutCodeChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(testFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Same message",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Get original commit hash
	client := rese.P1(gogit.New(tempDIR))
	headRef := rese.V1(client.Repo().Head())
	originalHash := headRef.Hash()

	// Try to amend with the same message - should be no-op
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the commit hash didn't change (no new commit created)
	newHeadRef := rese.V1(client.Repo().Head())
	newHash := newHeadRef.Hash()
	require.Equal(t, originalHash, newHash)
}

func TestGitCommit_AmendAuthorNameWithoutCodeChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(testFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Original User",
		Eddress:  "test@example.com",
		Message:  "Test message",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the original commit
	client := rese.P1(gogit.New(tempDIR))
	headRef := rese.V1(client.Repo().Head())
	headCommit := rese.P1(client.Repo().CommitObject(headRef.Hash()))
	require.Equal(t, "Original User", headCommit.Author.Name)

	// Amend just the author name without any file changes
	flags.Username = "Updated User"
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the author name was updated
	newHeadRef := rese.V1(client.Repo().Head())
	newHeadCommit := rese.P1(client.Repo().CommitObject(newHeadRef.Hash()))
	require.Equal(t, "Updated User", newHeadCommit.Author.Name)

	// Verify status is clean
	status := rese.V1(client.Status())
	require.Empty(t, status)
}

func TestGitCommit_AmendAuthorMailboxWithoutCodeChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(testFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "original@example.com",
		Message:  "Test message",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the original commit
	client := rese.P1(gogit.New(tempDIR))
	headRef := rese.V1(client.Repo().Head())
	headCommit := rese.P1(client.Repo().CommitObject(headRef.Hash()))
	require.Equal(t, "original@example.com", headCommit.Author.Email)

	// Amend just the mailbox without any file changes
	flags.Eddress = "updated@example.com"
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the mailbox was updated
	newHeadRef := rese.V1(client.Repo().Head())
	newHeadCommit := rese.P1(client.Repo().CommitObject(newHeadRef.Hash()))
	require.Equal(t, "updated@example.com", newHeadCommit.Author.Email)

	// Verify status is clean
	status := rese.V1(client.Status())
	require.Empty(t, status)
}

func TestGitCommit_AmendWithEmptyMessageWithoutCodeChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	// Create a new file and commit it first - must succeed for test setup
	testFile := filepath.Join(tempDIR, "test.txt")
	must.Done(os.WriteFile(testFile, []byte("test content"), 0644))

	flags := &CommitFlags{
		Username: "Test User",
		Eddress:  "test@example.com",
		Message:  "Original message",
		NoCommit: false,
		FormatGo: false,
		IsAmend:  false,
		IsForce:  false,
	}

	require.NoError(t, GitCommit(tempDIR, flags))

	// Get original commit hash
	client := rese.P1(gogit.New(tempDIR))
	headRef := rese.V1(client.Repo().Head())
	originalHash := headRef.Hash()

	// Try to amend with empty message and no file changes - should be no-op
	flags.Message = ""
	flags.IsAmend = true

	require.NoError(t, GitCommit(tempDIR, flags))

	// Verify the commit hash didn't change (no amend happened)
	newHeadRef := rese.V1(client.Repo().Head())
	newHash := newHeadRef.Hash()
	require.Equal(t, originalHash, newHash, "commit should not be amended when message is empty")
}

// setupTestRepoWithRemote creates a test repo with git remote configured
// Environment setup must succeed, so we use rese/must for all operations
func setupTestRepoWithRemote(remoteURL string) (string, func()) {
	tempDIR, cleanup := setupTestRepo()

	// Add remote to the repository using gitgo - must succeed
	gcm := gitgo.New(tempDIR)
	gcm.RemoteAdd("origin", remoteURL).MustDone()

	return tempDIR, cleanup
}

// createTestConfig creates a temporary config file for testing
// Environment setup must succeed, so we use rese/must for all operations
func createTestConfig(tempDIR string, config *CommitConfig) string {
	configPath := filepath.Join(tempDIR, "go-commit-config.json")
	configData := rese.V1(json.Marshal(config))
	must.Done(os.WriteFile(configPath, configData, 0644))
	return configPath
}

func TestMatchPattern_ExactMatch(t *testing.T) {
	pattern := "git@github.com:user/repo.git"
	remoteURL := "git@github.com:user/repo.git"

	score := utils.MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 28, score)
}

func TestMatchPattern_WildcardMatch(t *testing.T) {
	pattern := "git@github.com:*"
	remoteURL := "git@github.com:user/repo.git"

	score := utils.MatchRemotePattern(pattern, remoteURL)
	require.Greater(t, score, 0)
	require.Less(t, score, 24)
}

func TestMatchPattern_NoMatch(t *testing.T) {
	pattern := "git@gitlab.com:*"
	remoteURL := "git@github.com:user/repo.git"

	score := utils.MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, -1, score)
}

func TestMatchPattern_WildcardSpecificity(t *testing.T) {
	remoteURL := "git@github.com:user/repo.git"

	// More specific pattern should get higher score
	specificPattern := "git@github.com:user/*"
	generalPattern := "git@github.com:*"

	specificScore := utils.MatchRemotePattern(specificPattern, remoteURL)
	generalScore := utils.MatchRemotePattern(generalPattern, remoteURL)

	require.Greater(t, specificScore, generalScore)
}

func TestCommitConfig_MatchSignature_ExactMatch(t *testing.T) {
	config := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "exact-match",
				Username:       "exact-user",
				Eddress:        "exact@example.com",
				RemotePatterns: []string{"git@github.com:user/repo.git"},
			},
			&SignatureConfig{
				Name:           "wildcard-match",
				Username:       "wildcard-user",
				Eddress:        "wildcard@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	signature := config.MatchSignature("git@github.com:user/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "exact-match", signature.Name)
	require.Equal(t, "exact-user", signature.Username)
}

func TestCommitConfig_MatchSignature_WildcardMatch(t *testing.T) {
	config := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "github-match",
				Username:       "github-user",
				Eddress:        "github@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	signature := config.MatchSignature("git@github.com:different/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "github-match", signature.Name)
	require.Equal(t, "github-user", signature.Username)
}

func TestCommitConfig_MatchSignature_NoMatch(t *testing.T) {
	config := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "github-match",
				Username:       "github-user",
				Eddress:        "github@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	signature := config.MatchSignature("git@gitlab.com:user/repo.git")
	require.Nil(t, signature)
}

func TestCommitConfig_MatchSignature_Priority(t *testing.T) {
	config := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "general-github",
				Username:       "general-user",
				Eddress:        "general@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
			{
				Name:           "specific-user",
				Username:       "specific-user",
				Eddress:        "specific@example.com",
				RemotePatterns: []string{"git@github.com:specific/*"},
			},
		},
	}

	// More specific pattern should win
	signature := config.MatchSignature("git@github.com:specific/repo.git")
	require.NotNil(t, signature)
	require.Equal(t, "specific-user", signature.Name)
}

func TestLoadConfig_FileExists(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "config-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	testConfig := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "test-signature-info",
				Username:       "test-user",
				Eddress:        "test@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	configPath := createTestConfig(tempDIR, testConfig)

	config := LoadConfig(configPath)
	require.Len(t, config.Signatures, 1)
	require.Equal(t, "test-signature-info", config.Signatures[0].Name)
}

func TestLoadConfig_NoFileFound(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "config-test-*"))
	t.Cleanup(func() { must.Done(os.RemoveAll(tempDIR)) })

	nonExistentPath := filepath.Join(tempDIR, "non-existent-config.json")

	require.Panics(t, func() {
		config := LoadConfig(nonExistentPath)
		must.Full(config)
	})
}

func TestApplyProjectConfig_WithRemoteMatching(t *testing.T) {
	// 1. 立足项目 - 设置项目根目录
	tempDIR, cleanup := setupTestRepoWithRemote("git@github.com:user/repo.git")
	t.Cleanup(cleanup)

	// 2. 环境配置 - 创建测试配置
	testConfig := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "github-signature-info",
				Username:       "github-user",
				Eddress:        "github@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	// 3. 创建 flags - 用户的提交标志
	flags := &CommitFlags{
		Message: "test commit",
	}

	// 4. 读取配置并修饰 flags
	flags.ApplyProjectConfig(tempDIR, LoadConfig(createTestConfig(tempDIR, testConfig)))

	// 5. 检查结果 - 验证最终状态
	require.Equal(t, "github-user", flags.Username)
	require.Equal(t, "github@example.com", flags.Eddress)
}

func TestApplyProjectConfig_ConfigFlagOverride(t *testing.T) {
	// 1. 立足项目 - 设置项目根目录
	tempDIR, cleanup := setupTestRepoWithRemote("git@github.com:user/repo.git")
	t.Cleanup(cleanup)

	// 2. 环境配置 - 创建测试配置
	testConfig := &CommitConfig{
		Signatures: []*SignatureConfig{
			&SignatureConfig{
				Name:           "github-signature-info",
				Username:       "github-user",
				Eddress:        "github@example.com",
				RemotePatterns: []string{"git@github.com:*"},
			},
		},
	}

	// 3. 创建 flags - 用户的提交标志
	flags := &CommitFlags{
		Message:  "test commit",
		Username: "override-user", // Pre-set username will be overridden by config
	}

	// 4. 读取配置并修饰 flags
	flags.ApplyProjectConfig(tempDIR, LoadConfig(createTestConfig(tempDIR, testConfig)))

	// 5. 检查结果 - 验证最终状态
	require.Equal(t, "github-user", flags.Username)       // Config value overrides existing flag
	require.Equal(t, "github@example.com", flags.Eddress) // Config value applied to empty field
}

func TestAutoSignFlag(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	t.Cleanup(cleanup)

	t.Run("AutoSign enabled fills empty username and eddress", func(t *testing.T) {
		// Create a new file for commit - must succeed for test setup
		// 创建用于提交的新文件 - 测试设置必须成功
		newFile := filepath.Join(tempDIR, "test1.txt")
		must.Done(os.WriteFile(newFile, []byte("test content 1"), 0644))

		flags := &CommitFlags{
			Username: "", // Empty username
			Eddress:  "", // Empty eddress
			Message:  "Test commit with AutoSign",
			AutoSign: true, // Enable AutoSign
		}

		// This should work and fill username/eddress from git config (global or local)
		// 这应该工作并从 git 配置填充用户名/邮箱（全局或本地）
		require.NoError(t, GitCommit(tempDIR, flags))

		// Verify that git log shows the correct author from git config
		// 验证 git log 显示从 git 配置读取的正确作者
		execConfig := osexec.NewExecConfig().WithPath(tempDIR)
		authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
		authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))

		// Should not be empty since AutoSign should have filled from git config
		// 不应该为空，因为 AutoSign 应该从 git 配置填充
		require.NotEmpty(t, string(authorName))
		require.NotEmpty(t, string(authorEmail))
		require.NotEqual(t, "gogit", string(authorName)) // Should not be default gogit user
	})

	t.Run("AutoSign disabled keeps empty username and eddress", func(t *testing.T) {
		// Create a new file for commit - must succeed for test setup
		// 创建用于提交的新文件 - 测试设置必须成功
		newFile := filepath.Join(tempDIR, "test2.txt")
		must.Done(os.WriteFile(newFile, []byte("test content 2"), 0644))

		flags := &CommitFlags{
			Username: "", // Empty username
			Eddress:  "", // Empty eddress
			Message:  "Test commit without AutoSign",
			AutoSign: false, // Disable AutoSign
		}

		// This should still work but use empty username/eddress for commit
		// 这应该仍然工作但使用空的用户名/邮箱进行提交
		require.NoError(t, GitCommit(tempDIR, flags))

		// Verify that git log shows default gogit user (not from git config)
		// 验证 git log 显示默认 gogit 用户（不是从 git 配置读取）
		execConfig := osexec.NewExecConfig().WithPath(tempDIR)
		authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
		authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))

		// Should be gogit default user since AutoSign is disabled
		// 应该是 gogit 默认用户，因为 AutoSign 被禁用
		require.Equal(t, "gogit", string(authorName))
		require.Equal(t, "gogit@github.com/go-xlan/gogit", string(authorEmail))
	})

	t.Run("AutoSign enabled but username already provided", func(t *testing.T) {
		// Create a new file for commit - must succeed for test setup
		// 创建用于提交的新文件 - 测试设置必须成功
		newFile := filepath.Join(tempDIR, "test3.txt")
		must.Done(os.WriteFile(newFile, []byte("test content 3"), 0644))

		flags := &CommitFlags{
			Username: "Manual User",        // Already provided username
			Eddress:  "manual@example.com", // Already provided eddress
			Message:  "Test commit with existing info",
			AutoSign: true, // Enable AutoSign but shouldn't override existing info
		}

		// Should use provided username/eddress instead of reading from git config
		// 应该使用提供的用户名/邮箱而不是从 git 配置读取
		require.NoError(t, GitCommit(tempDIR, flags))

		// Verify that git log shows the manually provided user info (not from git config)
		// 验证 git log 显示手动提供的用户信息（不是从 git 配置读取）
		execConfig := osexec.NewExecConfig().WithPath(tempDIR)
		authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
		authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))

		// Should use the manually provided info, not git config
		// 应该使用手动提供的信息，而不是 git 配置
		require.Equal(t, "Manual User", string(authorName))
		require.Equal(t, "manual@example.com", string(authorEmail))
	})
}

func TestCommitFlags_ValidateFlags(t *testing.T) {
	t.Run("Valid configuration should have no warnings", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "test-user",
			Mailbox:  "test@example.com",
			Message:  "test message",
			IsAmend:  false,
			IsForce:  false,
			NoCommit: false,
			AutoSign: false,
		}

		warnings := flags.ValidateFlags()
		require.Empty(t, warnings)
	})

	t.Run("Force flag without amend should generate warning", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "test-user",
			Mailbox:  "test@example.com",
			Message:  "test message",
			IsAmend:  false,
			IsForce:  true, // Force without amend
		}

		warnings := flags.ValidateFlags()
		require.Len(t, warnings, 1)
		require.Contains(t, warnings[0], "force flag set but amend is disabled")
	})

	t.Run("Commit message with no-commit should generate warning", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "test-user",
			Mailbox:  "test@example.com",
			Message:  "test message", // Message with NoCommit
			NoCommit: true,
		}

		warnings := flags.ValidateFlags()
		require.Len(t, warnings, 1)
		require.Contains(t, warnings[0], "commit message provided but no-commit flag is set")
	})

	t.Run("Missing authentication info without auto-sign should generate warning", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "", // Empty
			Mailbox:  "", // Empty
			Eddress:  "", // Empty
			Message:  "test message",
			AutoSign: false, // Auto-sign disabled
		}

		warnings := flags.ValidateFlags()
		require.Len(t, warnings, 1)
		require.Contains(t, warnings[0], "no authentication info provided and auto-sign disabled")
	})

	t.Run("Missing authentication info with auto-sign should have no warnings", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "", // Empty
			Mailbox:  "", // Empty
			Eddress:  "", // Empty
			Message:  "test message",
			AutoSign: true, // Auto-sign enabled
		}

		warnings := flags.ValidateFlags()
		require.Empty(t, warnings)
	})

	t.Run("Multiple issues should generate multiple warnings", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "",             // Empty
			Mailbox:  "",             // Empty
			Eddress:  "",             // Empty
			Message:  "test message", // Message with NoCommit
			IsAmend:  false,
			IsForce:  true,  // Force without amend
			NoCommit: true,  // NoCommit with message
			AutoSign: false, // Auto-sign disabled
		}

		warnings := flags.ValidateFlags()
		require.Len(t, warnings, 3)
		require.Contains(t, warnings[0], "force flag set but amend is disabled")
		require.Contains(t, warnings[1], "commit message provided but no-commit flag is set")
		require.Contains(t, warnings[2], "no authentication info provided and auto-sign disabled")
	})

	t.Run("Mailbox as fallback for eddress should work", func(t *testing.T) {
		flags := &CommitFlags{
			Username: "test-user",
			Mailbox:  "",                 // Empty mailbox
			Eddress:  "test@example.com", // But eddress provided
			Message:  "test message",
			AutoSign: false,
		}

		warnings := flags.ValidateFlags()
		require.Empty(t, warnings) // Should have no warnings since eddress is provided
	})
}
