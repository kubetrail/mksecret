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

	"github.com/kubetrail/mkphrase/pkg/flags"
	"github.com/kubetrail/mkphrase/pkg/run"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set passphrase as a new version",
	Long:  `This command writes passphrase as a new version to the named passphrase`,
	RunE:  run.Set,
	Args:  cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(setCmd)
	f := setCmd.Flags()
	b := filepath.Base

	f.String(b(flags.Name), "", "Name tag for the password (DNS1123 label format)")
}
