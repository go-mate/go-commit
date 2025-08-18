package commitmate

import (
	"path/filepath"
	"strings"

	"github.com/go-xlan/gogit"
	"github.com/go-xlan/gogit/gogitassist"
	"github.com/go-xlan/gogit/gogitchange"
	"github.com/yyle88/erero"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// CommitFlags represents the configuration for a commit operation
type CommitFlags struct {
	Username string
	Message  string
	IsAmend  bool
	IsForce  bool
	Eddress  string
	NoCommit bool
	FormatGo bool
}

// GitCommit performs the complete commit workflow with optional Go code formatting
func GitCommit(projectRoot string, commitFlags *CommitFlags) error {
	zaplog.SUG.Debugln(projectRoot, neatjsons.S(commitFlags))

	client, err := gogit.New(projectRoot)
	if err != nil {
		return erero.Wro(err)
	}

	status, err := client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(neatjsons.S(status))

	if err := client.AddAll(); err != nil {
		return erero.Wro(err)
	}

	status, err = client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	zaplog.SUG.Debugln(neatjsons.S(status))

	if commitFlags.FormatGo {
		zaplog.SUG.Debugln("format changed go files")
		if err := FormatChangedGoFiles(projectRoot, client, DefaultAllowFormat); err != nil {
			return erero.Wro(err)
		}

		if err := client.AddAll(); err != nil {
			return erero.Wro(err)
		}

		status, err = client.Status()
		if err != nil {
			return erero.Wro(err)
		}
		zaplog.SUG.Debugln(neatjsons.S(status))
	}

	status, err = client.Status()
	if err != nil {
		return erero.Wro(err)
	}
	if len(status) == 0 {
		zaplog.SUG.Debugln("no change return")
		return nil
	}
	if commitFlags.NoCommit {
		zaplog.SUG.Debugln("no commit return")
		return nil
	}

	commitInfo := &gogit.CommitInfo{
		Name:    commitFlags.Username,
		Eddress: commitFlags.Eddress,
		Message: commitFlags.Message,
	}

	if commitFlags.IsAmend {
		_, err = client.AmendCommit(&gogit.AmendConfig{
			CommitInfo: commitInfo,
			ForceAmend: commitFlags.IsForce,
		})
		if err != nil {
			return erero.Wro(err)
		}
	} else {
		_, err = client.CommitAll(commitInfo)
		if err != nil {
			return erero.Wro(err)
		}
	}

	gogitassist.DebugRepo(client.Repo())
	return nil
}

// FormatChangedGoFiles formats Go files that have been changed
// The allowFormat function determines which files should be formatted
func FormatChangedGoFiles(projectRoot string, client *gogit.Client, allowFormat func(path string) bool) error {
	matchOptions := gogitchange.NewMatchOptions().MatchType(".go").MatchPath(func(path string) bool {
		zaplog.SUG.Debugln("path:", path)

		pass := allowFormat(path)
		if pass {
			zaplog.SUG.Debugln("pass:", path)
		} else {
			zaplog.SUG.Debugln("skip:", path)
		}
		return pass
	})

	err := gogitchange.NewChangedFileManager(projectRoot, client.Tree()).ForeachChangedGoFile(matchOptions, func(path string) error {
		if filepath.Ext(path) != ".go" {
			return nil
		}
		zaplog.ZAPS.Skip1.LOG.Info("golang-format-source", zap.String("path", path))

		if err := formatgo.FormatFile(path); err != nil {
			return erero.Wro(err)
		}
		return nil
	})
	if err != nil {
		return erero.Wro(err)
	}
	return nil
}

// DefaultAllowFormat is the default allow function for Go files formatting
// It skips common generated files like .pb.go, wire_gen.go, and ent files
func DefaultAllowFormat(path string) bool {
	if strings.HasSuffix(path, ".pb.go") || // skip protobuf generated files
		strings.HasSuffix(path, "/wire_gen.go") || // skip wire generated files
		strings.Contains(path, "/internal/data/ent/") { // skip ent generated files
		return false
	}
	return true
}
