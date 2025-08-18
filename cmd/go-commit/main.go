package main

import (
	"os"

	"github.com/go-mate/go-commit/commitmate"
	"github.com/spf13/cobra"
	"github.com/yyle88/eroticgo"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	projectRoot := rese.C1(os.Getwd())
	zaplog.SUG.Debugln(eroticgo.GREEN.Sprint(projectRoot))

	commitFlags := &commitmate.CommitFlags{}

	rootCmd := cobra.Command{
		Use:   "go-commit",
		Short: "Smart Git commit tool with Go code formatting",
		Long:  "go-commit is a Git commit tool that auto formats changed Go code and provides flexible commit options",
		Run: func(cmd *cobra.Command, args []string) {
			must.Done(commitmate.GitCommit(projectRoot, commitFlags))
		},
	}
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Username, "username", "u", "", "git username")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Message, "message", "m", "", "commit message")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsAmend, "amend", "a", false, "amend to the previous commit")
	rootCmd.PersistentFlags().BoolVarP(&commitFlags.IsForce, "force", "f", false, "force amend even pushed to remote")
	rootCmd.PersistentFlags().StringVarP(&commitFlags.Eddress, "eddress", "e", "", "email address")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.NoCommit, "no-commit", false, "stage changes without committing")
	rootCmd.PersistentFlags().BoolVar(&commitFlags.FormatGo, "format-go", false, "format changed go files")
	must.Done(rootCmd.Execute())
}
