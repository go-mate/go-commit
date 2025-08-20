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

func TestLoadExample4Config(t *testing.T) {
	configPath := runpath.PARENT.Join("advanced-patterns.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 4)
	require.Equal(t, "exact-repo-match", config.Signatures[0].Name)
	require.Equal(t, "team-specific", config.Signatures[1].Name)
	require.Equal(t, "multi-subdomain", config.Signatures[2].Name)
	require.Equal(t, "protocol-flexible", config.Signatures[3].Name)
}

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

func TestExample4GitCommitExactRepo(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-exact-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "exact-repo-match", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.lead", string(authorName))
	require.Equal(t, "diana.lead@critical-project.com", string(authorEmail))
}

func TestExample4GitCommitTeamSpecific(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-team-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "team-specific", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.dev", string(authorName))
	require.Equal(t, "diana@company.com", string(authorEmail))
}

func TestExample4GitCommitMultiSubdomain(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example4-subdomain-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "multi-subdomain", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "diana.ops", string(authorName))
	require.Equal(t, "diana.ops@multi-env.com", string(authorEmail))
}
