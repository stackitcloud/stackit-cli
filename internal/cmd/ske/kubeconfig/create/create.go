package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	skeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameArg = "CLUSTER_NAME"

	disableWritingFlag = "disable-writing"
	expirationFlag     = "expiration"
	filepathFlag       = "filepath"
	loginFlag          = "login"
	overwriteFlag      = "overwrite"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName    string
	DisableWriting bool
	ExpirationTime *string
	Filepath       *string
	Login          bool
	Overwrite      bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", clusterNameArg),
		Short: "Creates or update a kubeconfig for a SKE cluster",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
			"Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster, if the config exists in the kubeconfig file the information will be updated.",
			"By default, the kubeconfig information of the SKE cluster is merged into the default kubeconfig file of the current user. If the kubeconfig file doesn't exist, a new one will be created.",
			"You can override this behavior by specifying a custom filepath using the --filepath flag or by setting the KUBECONFIG env variable (fallback).\n",
			"An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is 1h.\n",
			"Note that the format is <value><unit>, e.g. 30d for 30 days and you can't combine units."),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create or update a kubeconfig for the SKE cluster with name "my-cluster. If the config exits in the kubeconfig file the information will be updated."`,
				"$ stackit ske kubeconfig create my-cluster"),
			examples.NewExample(
				`Get a login kubeconfig for the SKE cluster with name "my-cluster". `+
					"This kubeconfig does not contain any credentials and instead obtains valid credentials via the `stackit ske kubeconfig login` command.",
				"$ stackit ske kubeconfig create my-cluster --login"),
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 30 days. If the config exits in the kubeconfig file the information will be updated.`,
				"$ stackit ske kubeconfig create my-cluster --expiration 30d"),
			examples.NewExample(
				`Create or update a kubeconfig for the SKE cluster with name "my-cluster" and set the expiration time to 2 months. If the config exits in the kubeconfig file the information will be updated.`,
				"$ stackit ske kubeconfig create my-cluster --expiration 2M"),
			examples.NewExample(
				`Create or update a kubeconfig for the SKE cluster with name "my-cluster" in a custom filepath. If the config exits in the kubeconfig file the information will be updated.`,
				"$ stackit ske kubeconfig create my-cluster --filepath /path/to/config"),
			examples.NewExample(
				`Get a kubeconfig for the SKE cluster with name "my-cluster" without writing it to a file and format the output as json`,
				"$ stackit ske kubeconfig create my-cluster --disable-writing --output-format json"),
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster. It will OVERWRITE your current kubeconfig file."`,
				"$ stackit ske kubeconfig create my-cluster --overwrite true"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			if !model.DisableWriting {
				var prompt string
				if model.Overwrite {
					prompt = fmt.Sprintf("Are you sure you want to create a kubeconfig for SKE cluster %q? This will OVERWRITE your current kubeconfig file, if it exists.", model.ClusterName)
				} else {
					prompt = fmt.Sprintf("Are you sure you want to update your kubeconfig for SKE cluster %q? This will update your kubeconfig file. \nIf it the kubeconfig file doesn't exists, it will create a new one.", model.ClusterName)
				}
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			var (
				kubeconfig     string
				respKubeconfig *ske.Kubeconfig
				respLogin      *ske.LoginKubeconfig
			)

			if !model.Login {
				req, err := buildRequestCreate(ctx, model, apiClient)
				if err != nil {
					return fmt.Errorf("build kubeconfig create request: %w", err)
				}
				respKubeconfig, err = req.Execute()
				if err != nil {
					return fmt.Errorf("create kubeconfig for SKE cluster: %w", err)
				}
				if respKubeconfig.Kubeconfig == nil {
					return fmt.Errorf("no kubeconfig returned from the API")
				}
				kubeconfig = *respKubeconfig.Kubeconfig
			} else {
				req, err := buildRequestLogin(ctx, model, apiClient)
				if err != nil {
					return fmt.Errorf("build login kubeconfig create request: %w", err)
				}
				respLogin, err = req.Execute()
				if err != nil {
					return fmt.Errorf("create login kubeconfig for SKE cluster: %w", err)
				}
				if respLogin.Kubeconfig == nil {
					return fmt.Errorf("no login kubeconfig returned from the API")
				}
				kubeconfig = *respLogin.Kubeconfig
			}

			// Create the config file
			var kubeconfigPath string
			if model.Filepath == nil {
				kubeconfigPath, err = skeUtils.GetDefaultKubeconfigPath()
				if err != nil {
					return fmt.Errorf("get default kubeconfig path: %w", err)
				}
			} else {
				kubeconfigPath = *model.Filepath
			}

			if !model.DisableWriting {
				if model.Overwrite {
					err = skeUtils.WriteConfigFile(kubeconfigPath, kubeconfig)
				} else {
					err = skeUtils.MergeKubeConfig(kubeconfigPath, kubeconfig)
				}
				if err != nil {
					return fmt.Errorf("write kubeconfig file: %w", err)
				}
				params.Printer.Outputf("\nSet kubectl context to %s with: kubectl config use-context %s\n", model.ClusterName, model.ClusterName)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ClusterName, kubeconfigPath, respKubeconfig, respLogin)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(disableWritingFlag, false, fmt.Sprintf("Disable the writing of kubeconfig. Set the output format to json or yaml using the --%s flag to display the kubeconfig.", globalflags.OutputFormatFlag))
	cmd.Flags().BoolP(loginFlag, "l", false, "Create a login kubeconfig that obtains valid credentials via the STACKIT CLI. This flag is mutually exclusive with the expiration flag.")
	cmd.Flags().String(filepathFlag, "", "Path to create the kubeconfig file. Will fall back to KUBECONFIG env variable if not set. In case both aren't set, the kubeconfig is created as file named 'config' in the .kube folder in the user's home directory.")
	cmd.Flags().StringP(expirationFlag, "e", "", "Expiration time for the kubeconfig in seconds(s), minutes(m), hours(h), days(d) or months(M). Example: 30d. By default, expiration time is 1h")
	cmd.Flags().Bool(overwriteFlag, false, "Overwrite the kubeconfig file.")
	cmd.MarkFlagsMutuallyExclusive(loginFlag, expirationFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	expTime := flags.FlagToStringPointer(p, cmd, expirationFlag)

	if expTime != nil {
		var err error
		expTime, err = skeUtils.ConvertToSeconds(*expTime)
		if err != nil {
			return nil, &errors.FlagValidationError{
				Flag:    expirationFlag,
				Details: err.Error(),
			}
		}
	}

	disableWriting := flags.FlagToBoolValue(p, cmd, disableWritingFlag)

	isInvalidOutputFormat := globalFlags.OutputFormat == "" || globalFlags.OutputFormat == print.NoneOutputFormat || globalFlags.OutputFormat == print.PrettyOutputFormat
	if disableWriting && isInvalidOutputFormat {
		return nil, fmt.Errorf("when setting the flag --%s, you must specify --%s as one of the values: %s",
			disableWritingFlag, globalflags.OutputFormatFlag, fmt.Sprintf("%s, %s", print.JSONOutputFormat, print.YAMLOutputFormat))
	}

	model := inputModel{
		ClusterName:     clusterName,
		DisableWriting:  disableWriting,
		ExpirationTime:  expTime,
		Filepath:        flags.FlagToStringPointer(p, cmd, filepathFlag),
		GlobalFlagModel: globalFlags,
		Login:           flags.FlagToBoolValue(p, cmd, loginFlag),
		Overwrite:       flags.FlagToBoolValue(p, cmd, overwriteFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequestCreate(ctx context.Context, model *inputModel, apiClient *ske.APIClient) (ske.ApiCreateKubeconfigRequest, error) {
	req := apiClient.CreateKubeconfig(ctx, model.ProjectId, model.Region, model.ClusterName)

	payload := ske.CreateKubeconfigPayload{}

	if model.ExpirationTime != nil {
		payload.ExpirationSeconds = model.ExpirationTime
	}

	return req.CreateKubeconfigPayload(payload), nil
}

func buildRequestLogin(ctx context.Context, model *inputModel, apiClient *ske.APIClient) (ske.ApiGetLoginKubeconfigRequest, error) {
	return apiClient.GetLoginKubeconfig(ctx, model.ProjectId, model.Region, model.ClusterName), nil
}

func outputResult(p *print.Printer, outputFormat, clusterName, kubeconfigPath string, respKubeconfig *ske.Kubeconfig, respLogin *ske.LoginKubeconfig) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		var err error
		var details []byte
		if respKubeconfig != nil {
			details, err = json.MarshalIndent(respKubeconfig, "", "  ")
		} else if respLogin != nil {
			details, err = json.MarshalIndent(respLogin, "", "  ")
		}
		if err != nil {
			return fmt.Errorf("marshal SKE Kubeconfig: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		var err error
		var details []byte
		if respKubeconfig != nil {
			details, err = yaml.MarshalWithOptions(respKubeconfig, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		} else if respLogin != nil {
			details, err = yaml.MarshalWithOptions(respLogin, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		}
		if err != nil {
			return fmt.Errorf("marshal SKE Kubeconfig: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var expiration string
		if respKubeconfig != nil {
			expiration = fmt.Sprintf(", with expiration date %v (UTC)", utils.ConvertTimePToDateTimeString(respKubeconfig.ExpirationTimestamp))
		}
		p.Outputf("Updated kubeconfig file for cluster %s in %q%s\n", clusterName, kubeconfigPath, expiration)

		return nil
	}
}
