package run

import (
	"fmt"
	"hash/crc32"
	"os"

	"github.com/kubetrail/mksecret/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type persistentFlagValues struct {
	ApplicationCredentials string `json:"applicationCredentials,omitempty"`
	Project                string `json:"project,omitempty"`
	OutputFormat           string `json:"outputFormat,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()

	_ = viper.BindPFlag(flags.GoogleProjectID, rootCmd.Lookup(flags.GoogleProjectID))
	_ = viper.BindPFlag(flags.GoogleApplicationCredentials, rootCmd.Lookup(flags.GoogleApplicationCredentials))
	_ = viper.BindPFlag(flags.OutputFormat, rootCmd.Lookup(flags.OutputFormat))

	_ = viper.BindEnv(flags.GoogleProjectID, "GOOGLE_PROJECT_ID")

	applicationCredentials := viper.GetString(flags.GoogleApplicationCredentials)
	project := viper.GetString(flags.GoogleProjectID)
	outputFormat := viper.GetString(flags.OutputFormat)

	return persistentFlagValues{
		ApplicationCredentials: applicationCredentials,
		Project:                project,
		OutputFormat:           outputFormat,
	}
}

func setAppCredsEnvVar(applicationCredentials string) error {
	if len(applicationCredentials) > 0 {
		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", applicationCredentials); err != nil {
			err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
			return err
		}
	}

	return nil
}

// Crc32Sum produces crc32 sum
func Crc32Sum(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}
