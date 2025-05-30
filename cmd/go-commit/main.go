package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-xlan/gogit"
	"github.com/go-xlan/gogit/gogitassist"
	"github.com/go-xlan/gogit/gogitchange"
	"github.com/spf13/cobra"
	"github.com/yyle88/done"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/formatgo"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func main() {
	projectRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(projectRoot))

	commitFlags := &CommitFlags{}

	rootCmd := cobra.Command{
		Use:   "go-commit",
		Short: "go-commit",
		Long:  "go-commit",
		Run: func(cmd *cobra.Command, args []string) {
			GitCommit(projectRoot, commitFlags)
		},
	}
	//rootCmd.AddCommand(newLintRunCmd(projectRoot))
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Username, "username", "u", "", "git username")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Message, "message", "m", "", "commit message")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsAmend, "amend", "a", false, "is amend commit")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Eddress, "eddress", "e", "", "emails address")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.NoCommit, "no-commit", false, "not commit changes to git")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.FormatGo, "format-go", false, "format active go files")
	must.Done(rootCmd.Execute())
}

type CommitFlags struct {
	Username string
	Message  string
	IsAmend  bool
	Eddress  string
	NoCommit bool
	FormatGo bool
}

func GitCommit(projectRoot string, commitFlags *CommitFlags) {
	zaplog.SUG.Debugln(projectRoot, neatjsons.S(commitFlags))

	client := rese.P1(gogit.New(projectRoot))
	zaplog.SUG.Debugln(neatjsons.S(rese.V1(client.Status())))

	must.Done(client.AddAll())
	zaplog.SUG.Debugln(neatjsons.S(rese.V1(client.Status())))

	if commitFlags.FormatGo {
		zaplog.SUG.Debugln("format active go files")
		MustFormatGoFiles(projectRoot, client)

		must.Done(client.AddAll())
		zaplog.SUG.Debugln(neatjsons.S(rese.V1(client.Status())))
	}

	if len(rese.V1(client.Status())) == 0 {
		zaplog.SUG.Debugln("no change return")
		return
	}
	if commitFlags.NoCommit {
		zaplog.SUG.Debugln("no commit return")
		return
	}

	commitInfo := &gogit.CommitInfo{
		Name:    commitFlags.Username,
		Eddress: commitFlags.Eddress,
		Message: commitFlags.Message,
	}

	if commitFlags.IsAmend {
		done.VSE(client.AmendCommit(&gogit.AmendConfig{
			CommitInfo: commitInfo,
			ForceAmend: false,
		})).Done()
	} else {
		done.VSE(client.CommitAll(commitInfo)).Done()
	}

	gogitassist.DebugRepo(client.Repo())
}

func MustFormatGoFiles(projectRoot string, client *gogit.Client) {
	matchOptions := gogitchange.NewMatchOptions().MatchType(".go").MatchPath(func(path string) bool {
		zaplog.SUG.Debugln("path:", path)

		if strings.HasSuffix(path, ".pb.go") || //skip the pb code
			strings.HasSuffix(path, "/wire_gen.go") || //skip the wire code
			strings.Contains(path, "/internal/data/ent/") { //skip the auto gen code
			zaplog.SUG.Debugln("skip:", path)
			return false
		}

		zaplog.SUG.Debugln("pass:", path)
		return true
	})
	must.Done(gogitchange.NewChangedFileManager(projectRoot, client.Tree()).ForeachChangedGoFile(matchOptions, func(path string) error {
		if filepath.Ext(path) != ".go" {
			return nil
		}
		zaplog.ZAPS.Skip1.LOG.Info("golang-format-source", zap.String("path", path))

		must.Done(formatgo.FormatFile(path))
		return nil
	}))
}
