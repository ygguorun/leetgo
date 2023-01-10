package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
	"github.com/j178/leetgo/lang"
	"github.com/j178/leetgo/leetcode"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:    "git",
	Hidden: true,
	Short:  "Git related commands",
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Add, commit and push your code to remote repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c := leetcode.NewClient()
		qs, err := leetcode.ParseQID(args[0], c)
		if err != nil {
			return err
		}
		if len(qs) > 1 {
			return fmt.Errorf("multiple questions found")
		}
		result, err := lang.GeneratePathsOnly(qs[0])
		if err != nil {
			return err
		}
		err = gitAddCommitPush(result)
		return err
	},
}

func init() {
	gitCmd.AddCommand(gitPushCmd)
}

func gitAddCommitPush(genResult *lang.GenerateResult) error {
	files := make([]string, 0, len(genResult.Files))
	for _, f := range genResult.Files {
		files = append(files, f.Path)
	}
	err := runCmd("git", "add", files...)
	if err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	var msg string
	prompt := &survey.Editor{
		Message: "Commit message",
		Default: fmt.Sprintf(
			"Add solution for %s. %s",
			genResult.Question.QuestionFrontendId,
			genResult.Question.GetTitle(),
		),
	}
	err = survey.AskOne(prompt, &msg)
	if err != nil {
		return fmt.Errorf("git commit message: %w", err)
	}
	if msg == "" {
		return errors.New("git commit message: empty message")
	}
	err = runCmd("git", "commit", "-m", msg)
	if err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	err = runCmd("git", "push")
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil
}

func runCmd(command string, subcommand string, args ...string) error {
	cmd := exec.Command(command, subcommand)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
