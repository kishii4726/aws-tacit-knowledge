package cmd

import (
	ec2Internal "awsselfrev/internal/aws/service/ec2"
	"awsselfrev/internal/color"
	"awsselfrev/internal/config"
	"awsselfrev/internal/table"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Check EC2 resources for best practices and configurations",
	Long: `This command checks various EC2 configurations and best practices such as:
- EBS default encryption
- Volume encryption
- Snapshot encryption`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		tbl := table.SetTable()
		client := ec2.NewFromConfig(cfg)
		_, levelWarning, levelAlert := color.SetLevelColor()

		ebsEncryptionEnabled, err := ec2Internal.IsEbsDefaultEncryptionEnabled(client)
		if err != nil {
			log.Fatalf("Failed to check EBS default encryption: %v", err)
		}
		if !ebsEncryptionEnabled {
			tbl.Append([]string{"EC2", levelWarning, "-", "Default encryption for EBS is not set"})
		}

		unencryptedVolumes, err := ec2Internal.IsVolumeEncrypted(client)
		if err != nil {
			log.Fatalf("Failed to check volume encryption: %v", err)
		}
		for _, v := range unencryptedVolumes {
			tbl.Append([]string{"EC2", levelAlert, v, "EBS encryption is not set"})
		}

		encryptedSnapshots, err := ec2Internal.IsSnapshotEncrypted(client)
		if err != nil {
			log.Fatalf("Failed to check snapshot encryption: %v", err)
		}
		for _, v := range encryptedSnapshots {
			tbl.Append([]string{"EC2", levelAlert, v, "EBS encryption is not set"})
		}

		table.Render("EC2", tbl)
	},
}

func init() {
	rootCmd.AddCommand(ec2Cmd)
}
