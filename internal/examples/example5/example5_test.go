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

func TestLoadExample5Config(t *testing.T) {
	configPath := runpath.PARENT.Join("minimal-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 1)
	require.Equal(t, "main-identity", config.Signatures[0].Name)
	require.Equal(t, "developer", config.Signatures[0].Username)
	require.Equal(t, "developer@example.com", config.Signatures[0].Eddress)
	require.Equal(t, []string{"*"}, config.Signatures[0].RemotePatterns)
}

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

func TestExample5GitCommitGitHub(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-github-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "main-identity", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

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

func TestExample5GitCommitGitLab(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-gitlab-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "main-identity", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "developer", string(authorName))
	require.Equal(t, "developer@example.com", string(authorEmail))
}

func TestExample5GitCommitCustomDomain(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example5-custom-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "main-identity", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "developer", string(authorName))
	require.Equal(t, "developer@example.com", string(authorEmail))
}

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

	// Verify simplicity - only one signature needed
	require.Len(t, config.Signatures, 1)

	// Verify universal coverage - single pattern covers everything
	require.Len(t, config.Signatures[0].RemotePatterns, 1)
	require.Equal(t, "*", config.Signatures[0].RemotePatterns[0])
}
