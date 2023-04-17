/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:  "comment",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment := args[0]
		issueID := args[1]
		runnerID := os.Getenv("GITHUB_RUN_ID")
		safeRepo := os.Getenv("GITHUB_REPOSITORY")
		path := os.Getenv("GITHUB_EVENT_PATH")
		logger.Debug("path: %s", path)
		runnerURL := fmt.Sprintf("https://github.com/%s/actions/runs/%s", safeRepo, runnerID)
		switch {
		case strings.HasPrefix(comment, "/sealos_bot_release"):
			data := strings.Split(comment, " ")
			if len(data) == 2 && utils.ValidateVersion(data[1]) {

				msg := config.GlobalsConfig.GetMessage("release_success", "release action finished successfully!")
				return gh.SendMsgToIssue(issueID, msg, runnerURL, safeRepo)
			} else {
				msg := config.GlobalsConfig.GetMessage("release_format_error", "release action failed!")
				logger.Error("command format is error: %s ex. /sealos_bot_release {tag}", comment)
				return gh.SendMsgToIssue(issueID, msg, runnerURL, safeRepo)
			}
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		printEnvs()
		if err := checkPermission(); err != nil {
			return err
		}
		if err := checkGithubEnv(); err != nil {
			return err
		}
		return nil
	},
}

func checkGithubEnv() error {
	if _, ok := os.LookupEnv("GITHUB_RUN_ID"); !ok {
		return fmt.Errorf("error: GITHUB_RUN_ID is not set. Please set the GITHUB_RUN_ID environment variable")
	}
	if _, ok := os.LookupEnv("GITHUB_REPOSITORY"); !ok {
		return fmt.Errorf("error: GITHUB_REPOSITORY is not set. Please set the GITHUB_REPOSITORY environment variable")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(commentCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
