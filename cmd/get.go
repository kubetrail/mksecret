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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a named phrase",
	Long:  `Retrieve a named phrase value`,
	RunE:  run.Get,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(getCmd)
	f := getCmd.Flags()
	b := filepath.Base

	f.String(b(flags.Version), "latest", "Get specific version")
	f.String(flags.Passphrase, "", "Encryption passphrase if required")
	f.Bool(flags.NoPrompt, false, "Hide all prompts")
}
