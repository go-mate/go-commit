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

func TestLoadExample3Config(t *testing.T) {
	configPath := runpath.PARENT.Join("multi-client-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 4)
	require.Equal(t, "client-alpha", config.Signatures[0].Name)
	require.Equal(t, "client-beta", config.Signatures[1].Name)
	require.Equal(t, "personal-projects", config.Signatures[2].Name)
	require.Equal(t, "consulting-work", config.Signatures[3].Name)
}

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

func TestExample3GitCommitClientAlpha(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-alpha-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

func TestExample3GitCommitPersonalProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-personal-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

func TestExample3GitCommitConsultingWork(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example3-consulting-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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
