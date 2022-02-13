package run

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/kubetrail/mkphrase/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type persistentFlagValues struct {
	ApplicationCredentials string `json:"applicationCredentials,omitempty"`
	Project                string `json:"project,omitempty"`
}

func getPersistentFlags(cmd *cobra.Command) persistentFlagValues {
	rootCmd := cmd.Root().PersistentFlags()
	b := filepath.Base

	_ = viper.BindPFlag(flags.GoogleProjectID, rootCmd.Lookup(b(flags.GoogleProjectID)))
	_ = viper.BindPFlag(flags.GoogleApplicationCredentials, rootCmd.Lookup(b(flags.GoogleApplicationCredentials)))

	_ = viper.BindEnv(flags.GoogleProjectID, "GOOGLE_PROJECT_ID")

	applicationCredentials := viper.GetString(flags.GoogleApplicationCredentials)
	project := viper.GetString(flags.GoogleProjectID)

	return persistentFlagValues{
		ApplicationCredentials: applicationCredentials,
		Project:                project,
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
