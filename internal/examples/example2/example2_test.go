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

func TestLoadExample2Config(t *testing.T) {
	configPath := runpath.PARENT.Join("simple-personal-config.json")
	config := commitmate.LoadConfig(configPath)

	require.Len(t, config.Signatures, 3)
	require.Equal(t, "work-projects", config.Signatures[0].Name)
	require.Equal(t, "personal-github", config.Signatures[1].Name)
	require.Equal(t, "personal-gitlab", config.Signatures[2].Name)
}

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

func TestExample2GitCommitWorkProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example2-work-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "work-projects", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

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

func TestExample2GitCommitPersonalProject(t *testing.T) {
	tempDIR := rese.V1(os.MkdirTemp("", "example2-personal-test-*"))
	defer func() { must.Done(os.RemoveAll(tempDIR)) }()

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

	signature := rese.V1(commitmate.GetSignatureConfig(configPath, tempDIR))
	require.NotNil(t, signature)
	require.Equal(t, "personal-github", signature.Name)

	flags.Username = signature.Username
	flags.Eddress = signature.Eddress

	must.Done(commitmate.GitCommit(tempDIR, flags))

	authorName := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%an"))
	authorEmail := rese.V1(execConfig.Exec("git", "log", "-1", "--pretty=format:%ae"))
	require.Equal(t, "bob-dev", string(authorName))
	require.Equal(t, "bob.personal@gmail.com", string(authorEmail))
}
