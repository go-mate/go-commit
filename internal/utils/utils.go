// Package utils provides advanced utility functions for go-commit pattern matching
// Contains sophisticated pattern matching utilities for Git remote URL matching
// Implements recursive glob matching algorithms with wildcard support and scoring systems
// Optimized for high-performance pattern matching in enterprise Git workflow automation
//
// utils 为 go-commit 模式匹配提供高级工具函数
// 包含用于 Git 远程 URL 匹配的复杂模式匹配工具
// 实现递归 glob 匹配算法，支持通配符和评分系统
// 为企业 Git 工作流自动化中的高性能模式匹配进行了优化
package utils

import "strings"

// MatchRemotePattern calculates intelligent match score for pattern against remote URL
// Implements advanced scoring system where score equals non-wildcard character count if matched
// Returns -1 for no match, 0+ for various specificity levels of successful matching
// Optimized for enterprise Git workflows with complex remote configurations
//
// Advanced Pattern Syntax Support:
// - "*" matches any number of characters (including zero) - universal wildcard
// - "git@github.com:*" matches "git@github.com:user/repo.git" - domain-specific
// - "*://*.com/*" matches "https://example.com/path" - protocol-agnostic
// - "git@*.company.com:team/*" - subdomain and path-specific matching
//
// Intelligent Scoring Algorithm:
// - Score equals total non-wildcard characters in pattern for specificity ranking
// - Exact matches naturally achieve highest scores (complete string length)
// - More specific patterns automatically receive higher priority scores
// - Enables automatic best-match selection in multi-pattern configurations
//
// MatchRemotePattern 为模式与远程 URL 计算智能匹配分数
// 实现高级评分系统，如果匹配则分数等于非通配符字符数
// 不匹配返回 -1，成功匹配的各种特异性级别返回 0+
// 为复杂远程配置的企业 Git 工作流进行了优化
//
// 高级模式语法支持：
// - "*" 匹配任意数量的字符（包括零个） - 通用通配符
// - "git@github.com:*" 匹配 "git@github.com:user/repo.git" - 域名特定
// - "*://*.com/*" 匹配 "https://example.com/path" - 协议无关
// - "git@*.company.com:team/*" - 子域名和路径特定匹配
//
// 智能评分算法：
// - 分数等于模式中非通配符字符总数，用于特异性排名
// - 精确匹配自然获得最高分数（完整字符串长度）
// - 更具体的模式自动获得更高的优先级分数
// - 在多模式配置中实现自动最佳匹配选择
func MatchRemotePattern(pattern, remoteURL string) int {
	// Use glob matching for all patterns
	// 对所有模式使用 glob 匹配
	if !matchGlob(remoteURL, pattern) {
		return -1
	}

	// Calculate score: count non-wildcard characters
	// Exact matches will naturally get highest scores (full length)
	// 计算分数：统计非通配符字符数量
	// 精确匹配会自然获得最高分数（完整长度）
	return countNonWildcardChars(pattern)
}

// matchGlob performs sophisticated glob pattern matching with advanced wildcard support
// Implements high-performance recursive matching algorithm for complex URL patterns
// Returns true if remoteURL successfully matches the specified wildcard pattern
//
// matchGlob 执行复杂的 glob 模式匹配，支持高级通配符
// 实现高性能递归匹配算法，用于复杂 URL 模式
// 如果 remoteURL 成功匹配指定的通配符模式则返回 true
func matchGlob(remoteURL, pattern string) bool {
	// Quick exact match check for performance
	// 快速精确匹配检查以提升性能
	if remoteURL == pattern {
		return true
	}
	// No wildcard in pattern, must be exact match
	// 模式中没有通配符，必须精确匹配
	if !strings.Contains(pattern, "*") {
		return false
	}
	return matchGlobRecursive([]rune(remoteURL), []rune(pattern))
}

// matchGlobRecursive recursively matches remoteURL against pattern
// Uses slice recursion to consume characters from both strings
//
// Algorithm:
// 1. If pattern is empty, check if remoteURL is also empty
// 2. If pattern starts with *, try matching 0 to N characters
// 3. If remoteURL is empty but pattern has non-wildcard, no match
// 4. Otherwise, match single character and recurse
//
// matchGlobRecursive 递归地将 remoteURL 与模式匹配
// 使用切片递归从两个字符串中消费字符
//
// 算法：
// 1. 如果模式为空，检查 remoteURL 是否也为空
// 2. 如果模式以 * 开头，尝试匹配 0 到 N 个字符
// 3. 如果 remoteURL 为空但模式有非通配符，则不匹配
// 4. 否则，匹配单个字符并递归
func matchGlobRecursive(remoteURL, pattern []rune) bool {
	// End of pattern
	// 模式结束
	if len(pattern) == 0 {
		return len(remoteURL) == 0
	}

	// Current character is wildcard
	// 当前字符是通配符
	if pattern[0] == '*' {
		// Try matching 0 or more characters
		// 尝试匹配 0 个或多个字符
		for i := 0; i <= len(remoteURL); i++ {
			if matchGlobRecursive(remoteURL[i:], pattern[1:]) {
				return true
			}
		}
		return false
	}

	// End of remoteURL but pattern continues (non-wildcard)
	// remoteURL 结束但模式继续（非通配符）
	if len(remoteURL) == 0 {
		return false
	}

	// Match single character
	// 匹配单个字符
	if remoteURL[0] == pattern[0] {
		return matchGlobRecursive(remoteURL[1:], pattern[1:])
	}

	// No match
	// 无匹配
	return false
}

// countNonWildcardChars performs precise character counting in patterns excluding wildcards
// Calculates specificity score by counting all non-asterisk characters in the pattern
// Essential component of the intelligent scoring system for pattern priority ranking
//
// countNonWildcardChars 对模式中除通配符外的字符进行精确计数
// 通过统计模式中所有非星号字符来计算特异性分数
// 是用于模式优先级排名的智能评分系统的关键组件
func countNonWildcardChars(pattern string) int {
	count := 0
	for _, char := range pattern {
		if char != '*' {
			count++
		}
	}
	return count
}
