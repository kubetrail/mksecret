package run

import (
	"fmt"
	"path/filepath"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func List(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name"})
	table.SetBorder(false)
	table.SetColumnSeparator(" ")

	// Create the client.
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	listRequest := &secretmanagerpb.ListSecretsRequest{
		Parent:    fmt.Sprintf("projects/%s", persistentFlags.Project),
		PageSize:  0,
		PageToken: "",
		Filter:    fmt.Sprintf("labels.%s=%s", LabelKey, AppName),
	}

	// Call the API.
	secrets := client.ListSecrets(ctx, listRequest)
	if err != nil {
		return fmt.Errorf("failed to access secret version: %w", err)
	}

	for {
		secret, err := secrets.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		table.Append([]string{filepath.Base(secret.GetName())})
	}

	table.Render() // Send output

	return nil
}
