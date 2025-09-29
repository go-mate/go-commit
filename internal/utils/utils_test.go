// Package utils_test provides comprehensive testing for pattern matching utilities
// Tests cover exact matching, wildcard patterns, scoring systems, and edge cases
// Validates the sophisticated pattern matching algorithms used in Git remote URL matching
//
// utils_test 为模式匹配工具提供全面的测试
// 测试涵盖精确匹配、通配符模式、评分系统和边界情况
// 验证用于 Git 远程 URL 匹配的复杂模式匹配算法
package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestMatchRemotePattern_ExactMatch tests perfect pattern matching
// Validates that identical pattern and URL return correct score
//
// TestMatchRemotePattern_ExactMatch 测试完美模式匹配
// 验证相同的模式和 URL 返回正确的分数
func TestMatchRemotePattern_ExactMatch(t *testing.T) {
	pattern := "git@github.com:user/repo.git"
	remoteURL := "git@github.com:user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 28, score)
}

func TestMatchRemotePattern_WildcardMatch(t *testing.T) {
	pattern := "git@github.com:*"
	remoteURL := "git@github.com:user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 15, score)
}

func TestMatchRemotePattern_SingleWildcard(t *testing.T) {
	pattern := "*"
	remoteURL := "git@github.com:user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 0, score)
}

func TestMatchRemotePattern_NoMatch(t *testing.T) {
	pattern := "git@gitlab.com:*"
	remoteURL := "git@github.com:user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, -1, score)
}

func TestMatchRemotePattern_MultipleWildcards(t *testing.T) {
	pattern := "*@*.com:*"
	remoteURL := "git@github.com:user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 6, score)
}

func TestMatchRemotePattern_HTTPSMatch(t *testing.T) {
	pattern := "https://*.com/*"
	remoteURL := "https://github.com/user/repo.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 13, score)
}

func TestMatchRemotePattern_ComplexPattern(t *testing.T) {
	pattern := "git@github.*.com:go-mate/*"
	remoteURL := "git@github.mate.com:go-mate/go-commit.git"

	score := MatchRemotePattern(pattern, remoteURL)
	require.Equal(t, 24, score)
}

func TestMatchRemotePattern_EmptyStrings(t *testing.T) {
	// Empty pattern matches empty URL
	score := MatchRemotePattern("", "")
	require.Equal(t, 0, score)

	// Empty pattern doesn't match non-empty URL
	score = MatchRemotePattern("", "git@github.com")
	require.Equal(t, -1, score)

	// Non-empty pattern doesn't match empty URL
	score = MatchRemotePattern("git@*", "")
	require.Equal(t, -1, score)
}

// TestMatchRemotePattern_SpecificityRanking validates the scoring algorithm for pattern specificity
// Tests that more specific patterns receive higher scores for prioritized matching
//
// TestMatchRemotePattern_SpecificityRanking 验证模式特异性的评分算法
// 测试更具体的模式获得更高分数以实现优先匹配
func TestMatchRemotePattern_SpecificityRanking(t *testing.T) {
	remoteURL := "git@github.com:user/repo.git"

	// Test different patterns and their scores
	// 测试不同模式及其分数
	require.Equal(t, 28, MatchRemotePattern("git@github.com:user/repo.git", remoteURL))
	require.Equal(t, 20, MatchRemotePattern("git@github.com:user/*", remoteURL))
	require.Equal(t, 15, MatchRemotePattern("git@github.com:*", remoteURL))
	require.Equal(t, 9, MatchRemotePattern("git@*.com:*", remoteURL))
	require.Equal(t, 5, MatchRemotePattern("git@*:*", remoteURL))
	require.Equal(t, 0, MatchRemotePattern("*", remoteURL))
}

func TestMatchRemotePattern_RealWorldPatterns(t *testing.T) {
	// GitHub patterns
	require.Greater(t, MatchRemotePattern("git@github.com:*", "git@github.com:user/repo.git"), -1)
	require.Greater(t, MatchRemotePattern("https://github.com/*", "https://github.com/user/repo.git"), -1)
	require.Equal(t, -1, MatchRemotePattern("*://github.com/*", "git@github.com:user/repo.git"))

	// Company patterns
	require.Greater(t, MatchRemotePattern("git@*.company.com:*", "git@gitlab.company.com:team/project.git"), -1)
	require.Greater(t, MatchRemotePattern("git@*.company.com:*", "git@github.company.com:team/project.git"), -1)
	require.Equal(t, -1, MatchRemotePattern("git@*.company.com:*", "git@external.com:team/project.git"))

	// Flexible patterns
	require.Greater(t, MatchRemotePattern("*@github.com:*", "git@github.com:user/repo.git"), -1)
	require.Greater(t, MatchRemotePattern("*@github.com:*", "https@github.com:user/repo.git"), -1)

	// Edge cases
	require.Greater(t, MatchRemotePattern("", ""), -1)
	require.Greater(t, MatchRemotePattern("*", "anything"), -1)
	require.Greater(t, MatchRemotePattern("exact", "exact"), -1)
	require.Equal(t, -1, MatchRemotePattern("exact", "different"))
}
