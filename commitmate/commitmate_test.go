package commitmate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-xlan/gogit"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/must"
	"github.com/yyle88/osexec"
	"github.com/yyle88/rese"
)

// setupTestRepo creates a temporary git repository for testing
// Environment setup must succeed, so we use rese/must for all operations
func setupTestRepo() (string, func()) {
	// Create temporary directory - must succeed
	tempDIR := rese.V1(os.MkdirTemp("", "go-commit-test-*"))

	// Initialize git repository using osexec - must succeed
	execConfig := osexec.NewExecConfig().WithPath(tempDIR)

	// Initialize git repository - must succeed
	rese.V1(execConfig.Exec("git", "init"))

	// Configure git user (required for commits) - must succeed
	rese.V1(execConfig.Exec("git", "config", "user.name", "Test User"))
	rese.V1(execConfig.Exec("git", "config", "user.email", "test@example.com"))

	// Create initial commit to make it a valid repo - must succeed
	testFile := filepath.Join(tempDIR, "README.md")
	must.Done(os.WriteFile(testFile, []byte("# Test Repository\n"), 0644))

	// Add and commit initial file - must succeed
	rese.V1(execConfig.Exec("git", "add", "."))
	rese.V1(execConfig.Exec("git", "commit", "-m", "Initial commit"))

	// Return cleanup function
	cleanup := func() {
		_ = os.RemoveAll(tempDIR)
	}

	return tempDIR, cleanup
}

func TestGitCommit_NoChanges(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

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
	defer cleanup()

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
	require.Empty(t, status, "Working directory should be clean after commit")
}

func TestGitCommit_NoCommitFlag(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

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
	require.NotEmpty(t, status, "Working directory should have staged changes")
}

func TestGitCommit_WithGoFileFormatting(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

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
	defer cleanup()

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

	// Format files with custom filter that only allows main.go - must succeed for test setup
	allowOnlyMain := func(path string) bool {
		return filepath.Base(path) == "main.go"
	}
	must.Done(FormatChangedGoFiles(tempDIR, client, allowOnlyMain))

	// Check that only main.go was formatted (has proper spacing) - test verification uses require
	mainContent := rese.V1(os.ReadFile(filepath.Join(tempDIR, "main.go")))
	require.Contains(t, string(mainContent), "import \"fmt\"", "main.go should be formatted")

	// Check that generated files were NOT formatted (still have formatting issues) - test verification uses require
	pbContent := rese.V1(os.ReadFile(filepath.Join(tempDIR, "generated.pb.go")))
	require.Contains(t, string(pbContent), "import\"fmt\"", "pb.go should NOT be formatted")
}

func TestGitCommit_AmendCommit(t *testing.T) {
	tempDIR, cleanup := setupTestRepo()
	defer cleanup()

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
	require.Empty(t, status, "Working directory should be clean after amend")
}
