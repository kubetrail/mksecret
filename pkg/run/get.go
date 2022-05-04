package run

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"syscall"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kubetrail/mksecret/pkg/app"
	"github.com/kubetrail/mksecret/pkg/crypto"
	"github.com/kubetrail/mksecret/pkg/flags"
	"github.com/mr-tron/base58"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
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
	encrypted := false

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
	if value, ok := labels[app.KeyManagedBy]; !ok || value != app.Name {
		return fmt.Errorf("secret is not being managed by this app")
	}
	if value, ok := labels[app.KeyEncrypted]; ok && value == app.ValueTrue {
		encrypted = true
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

	payload := result.Payload.GetData()

	if encrypted {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Enter encryption password: "); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
		encryptionKey, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read encryption password from input: %w", err)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
		key, err := crypto.NewAesKeyFromPassphrase([]byte(encryptionKey))
		if err != nil {
			return fmt.Errorf("failed to generate new AES key: %w", err)
		}

		ciphertext, err := base58.Decode(string(result.Payload.GetData()))
		if err != nil {
			return fmt.Errorf("failed to base58 decode stored value: %w", err)
		}

		payload, err = crypto.DecryptWithAesKey(ciphertext, key)
		if err != nil {
			return fmt.Errorf("failed to decrypt data: %w", err)
		}
	}

	switch persistentFlags.OutputFormat {
	case flags.OutputFormatNative:
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(payload)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatJson:
		jb, err := json.Marshal(
			struct {
				Name    string `json:"name,omitempty"`
				Version string `json:"version,omitempty"`
				Payload string `json:"payload,omitempty"`
			}{
				Name:    name,
				Version: path.Base(result.GetName()),
				Payload: string(payload),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to serialize output json: %w", err)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), string(jb)); err != nil {
			return fmt.Errorf("failed to write to output: %w", err)
		}
	case flags.OutputFormatTable:
		table := tablewriter.NewWriter(cmd.OutOrStdout())
		table.SetHeader([]string{"Name", "Version", "Phrase"})
		table.Append(
			[]string{
				name,
				filepath.Base(result.GetName()),
				string(payload),
			},
		)
		table.SetBorder(false)
		table.SetColumnSeparator(" ")
		table.Render() // Send output
	}

	return nil
}
