// Package commitmate provides Git utility functions
// Contains functions to retrieve Git remote and config information
//
// commitmate 包提供 Git 工具函数
// 包含获取 Git 远程和配置信息的函数
package commitmate

import (
	"github.com/go-xlan/gitgo"
	"github.com/go-xlan/gogit"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

// getOriginRemoteURL extracts the origin remote URL from a Git repo
// Prioritizes 'origin' remote, falls back to first available remote
// Returns blank string when no remotes exist
//
// getOriginRemoteURL 从 Git 仓库提取 origin 远程 URL
// 优先使用 'origin' 远程，回退到第一个可用远程
// 当没有远程时返回空字符串
func getOriginRemoteURL(projectRoot string) string {
	client := rese.P1(gogit.New(projectRoot))

	// Try origin remote first
	// 优先尝试 origin 远程
	if url, err := client.GetRemoteURL("origin"); err != nil {
		zaplog.SUG.Debugln("cannot get origin remote:", err)
	} else {
		return url
	}

	// Fallback to first available remote
	// 回退到第一个可用远程
	if url, err := client.GetFirstRemoteURL(); err != nil {
		zaplog.SUG.Debugln("cannot get first remote:", err)
	} else {
		return url
	}

	return ""
}

// getGitConfigValue retrieves a configuration value from Git config in the specified project DIR
// Returns blank string if command fails or config item doesn't exist
//
// getGitConfigValue 从指定项目 DIR 的 Git 配置获取配置值
// 如果命令失败或配置键不存在则返回空字符串
func getGitConfigValue(projectRoot, key string) string {
	value, err := gitgo.New(projectRoot).ConfigGet(key)
	if err != nil {
		zaplog.SUG.Debugln("cannot get git config", key, ":", err)
		return ""
	}
	return value
}
