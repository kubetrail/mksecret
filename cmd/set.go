/*
Copyright © 2022 kubetrail.io authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"path/filepath"

	"github.com/kubetrail/mksecret/pkg/flags"
	"github.com/kubetrail/mksecret/pkg/run"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a secret as a new version",
	Long:  `This command writes a secret as a new version to the named secret`,
	RunE:  run.Set,
}

func init() {
	rootCmd.AddCommand(setCmd)
	f := setCmd.Flags()
	b := filepath.Base

	f.String(b(flags.Name), "", "Name tag for the secret (DNS1123 label format)")
	f.Bool(b(flags.Encrypt), false, "Turn on encryption (true when passphrase is provided)")
	f.String(flags.Passphrase, "", "Encryption passphrase")
	f.Bool(flags.NoPrompt, false, "Hide all prompts")

	_ = setCmd.RegisterFlagCompletionFunc(
		flags.Name,
		func(
			cmd *cobra.Command,
			args []string,
			toComplete string,
		) (
			[]string,
			cobra.ShellCompDirective,
		) {
			return []string{
					"example-names-in-dns1123-format",
					"my-mnemonic-1",
					"eth-hex-seed-key",
					"my-super-secret",
				},
				cobra.ShellCompDirectiveDefault
		},
	)
}
