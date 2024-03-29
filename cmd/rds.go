/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"log"
	"os"

	"aws-tacit-knowledge/pkg/config"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

// rdsCmd represents the rds command
var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Service", "LEVEL", "MESSAGE"})
		client := rds.NewFromConfig(cfg)
		resp, err := client.DescribeDBClusters(context.TODO(), &rds.DescribeDBClustersInput{})
		if err != nil {
			log.Fatalf("%v", err)
		}
		var data [][]string

		for _, v := range resp.DBClusters {
			// Storageの暗号化確認
			if *&v.StorageEncrypted == false {
				data := append(data, []string{"RDS", "Alert", *v.DBClusterIdentifier + "のStorageが暗号化されていません"})
				for _, v := range data {
					table.Append(v)
				}
			}
			// 削除保護有効確認
			if *v.DeletionProtection == false {
				data := append(data, []string{"RDS", "Warning", *v.DBClusterIdentifier + "の削除保護が有効化されていません"})
				for _, v := range data {
					table.Append(v)
				}
			}
			// ログ出力確認 todo: ログ種類ごとに確認する
			if len(v.EnabledCloudwatchLogsExports) == 0 {
				data := append(data, []string{"RDS", "Warning", *v.DBClusterIdentifier + "でログ出力が設定されていません"})
				for _, v := range data {
					table.Append(v)
				}
			}
			for _, db_cluster_member := range v.DBClusterMembers {
				resp, err := client.DescribeDBInstances(context.TODO(), &rds.DescribeDBInstancesInput{
					DBInstanceIdentifier: db_cluster_member.DBInstanceIdentifier,
				})
				if err != nil {
					log.Fatalf("%v", err)
				}
				// 自動アップグレード
				for _, db_instance := range resp.DBInstances {
					if db_instance.AutoMinorVersionUpgrade == true {
						data := append(data, []string{"RDS", "Warning", *db_instance.DBInstanceIdentifier + "のマイナーバージョン自動アップグレードが有効化されています"})
						for _, v := range data {
							table.Append(v)
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(rdsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rdsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rdsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
