package cmd

import (
	"context"
	"log"

	"awsselfrev/internal/color"
	"awsselfrev/internal/config"
	"awsselfrev/internal/table"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var cloudwatchlogsCmd = &cobra.Command{
	Use:   "cloudwatchlogs",
	Short: "Checks CloudWatch Logs configurations for best practices",
	Long: `This command checks various CloudWatch Logs configurations and best practices such as:
- Log group retention settings`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		tbl := table.SetTable()
		client := cloudwatchlogs.NewFromConfig(cfg)
		_, _, levelAlert := color.SetLevelColor()

		checkLogGroupsRetention(client, tbl, levelAlert)

		table.Render("CloudWatchLogs", tbl)
	},
}

func checkLogGroupsRetention(client *cloudwatchlogs.Client, table *tablewriter.Table, levelAlert string) {
	resp, err := client.DescribeLogGroups(context.TODO(), &cloudwatchlogs.DescribeLogGroupsInput{})
	if err != nil {
		log.Fatalf("Failed to describe log groups: %v", err)
	}
	for _, logGroup := range resp.LogGroups {
		if logGroup.RetentionInDays == nil {
			table.Append([]string{"CloudWatchLogs", levelAlert, *logGroup.LogGroupName, "Retention is set to never expire"})
		}
	}
}

func init() {
	rootCmd.AddCommand(cloudwatchlogsCmd)
}
