package run

import (
	"bufio"
	"fmt"
	"path/filepath"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/google/uuid"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/kubetrail/mkphrase/pkg/flags"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/validation"
)

func Set(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	persistentFlags := getPersistentFlags(cmd)

	b := filepath.Base
	_ = viper.BindPFlag(flags.Name, cmd.Flags().Lookup(b(flags.Name)))
	name := viper.GetString(flags.Name)

	if err := setAppCredsEnvVar(persistentFlags.ApplicationCredentials); err != nil {
		err := fmt.Errorf("could not set Google Application credentials env. var: %w", err)
		return err
	}

	if len(name) == 0 {
		name = uuid.New().String()
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

	// Create the request to create the secret.
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", persistentFlags.Project),
		SecretId: name,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			Labels: map[string]string{
				LabelKey: AppName,
			},
		},
	}

	secret, err := client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		apiErr, ok := err.(*apierror.APIError)
		if ok {
			if apiErr.GRPCStatus().Code() == codes.AlreadyExists {
				secret, err = client.GetSecret(
					ctx,
					&secretmanagerpb.GetSecretRequest{
						Name: fmt.Sprintf("projects/%s/secrets/%s", persistentFlags.Project, name),
					},
				)
				if err != nil {
					return fmt.Errorf("failed to get secret: %w", err)
				}
			} else {
				return fmt.Errorf("failed to create secret: %w", err)
			}
		} else {
			return fmt.Errorf("failed to create a secret: %T, %w", err, err)
		}
	}

	labels := secret.GetLabels()
	if value, ok := labels[LabelKey]; !ok || value != AppName {
		return fmt.Errorf("secret is not being managed by this app")
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter passphrase: "); err != nil {
		return fmt.Errorf("failed to write to output: %w", err)
	}

	inputReader := bufio.NewReader(cmd.InOrStdin())
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read from input: %w", err)
	}

	if len(input) == 0 {
		return fmt.Errorf("invalid input of zero length")
	}

	// Build the request.
	data := []byte(input)
	dataCrc32C := int64(Crc32Sum(data))
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data:       data,
			DataCrc32C: &dataCrc32C,
		},
	}

	// Call the API.
	version, err := client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return fmt.Errorf("failed to add secret version: %w", err)
	}

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: version.Name,
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
			filepath.Base(version.GetName()),
			string(result.Payload.GetData()),
		},
	)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.Render() // Send output

	return nil
}
