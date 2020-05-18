package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var statusCmd = &cobra.Command{
	Use:   "status [commandId]",
	Short: "get the status of a running bootstrap",
	Long:  `get the status of a running bootstrap`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := getBaseUrl(url)
		commandId := ""
		if len(args) > 0 && len(args[0]) > 0 {
			commandId = args[0]
		} else {
			log.Fatalln("invalid value for the commandId argument")
		}
		fameService := NewService(ctx)
		response, err := fameService.Status(commandId)
		if err != nil {
			log.Fatalf("failed to get status for command %s. exception %s\n", commandId, err)
		}
		fmt.Printf("RemainingMessagesToBeSent: %d, failedMessages: %d \n", response.RemainingMessagesSent, response.FailedMessages)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
