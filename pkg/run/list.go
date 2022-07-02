package run

import (
	"encoding/json"
	"fmt"
	"path"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubetrail/mksecret/pkg/app"
	"github.com/kubetrail/mksecret/pkg/flags"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gopkg.in/yaml.v3"
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
		Filter:    fmt.Sprintf("labels.%s=%s", app.KeyManagedBy, app.Name),
	}

	// Call the API.
	secrets := client.ListSecrets(ctx, listRequest)
	if err != nil {
		return fmt.Errorf("failed to access secret version: %w", err)
	}

	switch persistentFlags.OutputFormat {
	case flags.OutputFormatNative:
		for {
			secret, err := secrets.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to list secrets: %w", err)
			}

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), path.Base(secret.GetName())); err != nil {
				return fmt.Errorf("failed to write to output: %w", err)
			}
		}
	case flags.OutputFormatJson:
		outputList := make([]string, 0, 128)
		for {
			secret, err := secrets.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to list secrets: %w", err)
			}

			outputList = append(outputList, path.Base(secret.GetName()))
		}
		jb, err := json.Marshal(outputList)
		if err != nil {
			return fmt.Errorf("failed to serialize output json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatYaml:
		outputList := make([]string, 0, 128)
		for {
			secret, err := secrets.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to list secrets: %w", err)
			}

			outputList = append(outputList, path.Base(secret.GetName()))
		}
		jb, err := yaml.Marshal(outputList)
		if err != nil {
			return fmt.Errorf("failed to serialize output yaml: %w", err)
		}

		if _, err := fmt.Fprint(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatTable:
		for {
			secret, err := secrets.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return fmt.Errorf("failed to list secrets: %w", err)
			}

			table.Append([]string{path.Base(secret.GetName())})
		}

		table.Render() // Send output
	}

	return nil
}
