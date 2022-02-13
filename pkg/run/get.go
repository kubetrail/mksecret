package run

import (
	"fmt"
	"path/filepath"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubetrail/mkphrase/pkg/flags"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

func Get(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	b := filepath.Base
	_ = viper.BindPFlag(flags.Version, cmd.Flags().Lookup(b(flags.Version)))

	name := args[0]
	version := viper.GetString(flags.Version)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	if len(name) == 0 {
		return fmt.Errorf("please provide name of the password")
	}

	if errs := validation.IsDNS1123Label(name); len(errs) > 0 {
		return fmt.Errorf("invalid name, need DNS1123Label format: %v", errs)
	}

	// Create the client.
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create secret manager client: %w", err)
	}
	defer client.Close()

	secret, err := client.GetSecret(
		ctx,
		&secretmanagerpb.GetSecretRequest{
			Name: fmt.Sprintf("projects/%s/secrets/%s", persistentFlags.Project, name),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to get secret: %w", err)
	}

	labels := secret.GetLabels()
	if value, ok := labels[LabelKey]; !ok || value != AppName {
		return fmt.Errorf("secret is not being managed by this app")
	}

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", persistentFlags.Project, name, version),
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return fmt.Errorf("failed to access secret version: %w", err)
	}

	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Version", "Phrase"})
	table.Append(
		[]string{
			name,
			filepath.Base(result.GetName()),
			string(result.Payload.GetData()),
		},
	)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.Render() // Send output

	return nil
}
