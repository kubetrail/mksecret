package run

import (
	"fmt"
	"path/filepath"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubetrail/mkphrase/pkg/flags"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"k8s.io/apimachinery/pkg/util/validation"
)

func Delete(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	_ = viper.BindPFlag(flags.Force, cmd.Flags().Lookup(filepath.Base(flags.Force)))

	name := args[0]
	force := viper.GetBool(flags.Force)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	if len(name) == 0 {
		return fmt.Errorf("please provide name of the password")
	}

	if !force {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Type secret name to delete: "); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
		var input string
		if _, err := fmt.Fscanln(cmd.InOrStdin(), &input); err != nil {
			return fmt.Errorf("failed to read from input: %w", err)
		}

		if input != name {
			return fmt.Errorf("input does not match secret name")
		}
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
	deleteRequest := &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", persistentFlags.Project, name),
	}

	// Call the API.
	if err := client.DeleteSecret(ctx, deleteRequest); err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}
