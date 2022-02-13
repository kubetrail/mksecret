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
	"fmt"
	"strings"

	"github.com/kubetrail/mkphrase/pkg/run"

	"github.com/spf13/cobra"
)

var listCmdLong = `List all named phrases managed by this app
Internally it filters all secrets using label: labelKey=appName`

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List named phrases",
	Long: strings.ReplaceAll(
		strings.ReplaceAll(
			listCmdLong,
			"labelKey",
			run.LabelKey,
		),
		"appName",
		run.AppName,
	),
	RunE:    run.List,
	Args:    cobra.ExactArgs(0),
	Example: fmt.Sprintf("%s list", run.AppName),
}

func init() {
	rootCmd.AddCommand(listCmd)
}
