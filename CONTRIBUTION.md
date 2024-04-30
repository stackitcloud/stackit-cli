# Contribute to the STACKIT CLI

Your contribution is welcome! Thank you for your interest in contributing to the STACKIT CLI. We greatly value your feedback, feature requests, additions to the code, bug reports or documentation extensions.

## Table of contents

- [Developer Guide](#developer-guide)
- [Useful Make commands](#useful-make-commands)
- [Repository structure](#repository-structure)
- [Implementing a new command](#implementing-a-new-command)
	- [Command file structure](#command-file-structure)
	- [Outputs, prints and debug logs](#outputs-prints-and-debug-logs)
- [Onboarding a new STACKIT service](#onboarding-a-new-stackit-service)
- [Local development](#local-development)
- [Code Contributions](#code-contributions)
- [Bug Reports](#bug-reports)

## Developer Guide

Prerequisites:

- [`Go`](https://go.dev/doc/install) 1.22+
- [`yamllint`](https://yamllint.readthedocs.io/en/stable/quickstart.html)

### Useful Make commands

These commands can be executed from the project root:

- `make project-tools`: install the required dependencies
- `make build`: compile the CLI and save the binary under _./bin/stackit_
- `make lint`: lint the code
- `make generate-docs`: generate Markdown documentation for every command
- `make test`: run unit tests

### Repository structure

The CLI commands are located under `internal/cmd`, where each folder includes the source code for each subcommand (including their own subcommands). Inside `pkg` you can find several useful packages that are shared by the commands and provide additional functionality such as `flags`, `globalflags`, `tables`, etc.

### Implementing a new command

Let's suppose you want to want to implement a new command `bar`, that would be the direct child of an existing command `stackit foo` (meaning it would be invoked as `stackit foo bar`):

1. You would start by creating a new folder `bar/` inside `internal/cmd/foo/`
2. Following with the creation of a file `bar.go` inside your new folder `internal/cmd/foo/bar/`
   1. The Go package should be similar to the command usage, in this case `package bar` would be an adequate name
   2. Please refer to the [Command file structure](./CONTRIBUTION.md/#command-file-structure) section for details on the strcutre of the file itself
3. To register the command `bar` as a child of the existing command `foo`, add `cmd.AddCommand(bar.NewCmd(p))` to the `addSubcommands` method of the constructor of the `foo` command
   1. In this case, `p` is the `printer` that is passed from the root command to all subcommands of the tree (refer to the [Outputs, prints and debug logs](./CONTRIBUTION.md/#outputs-prints-and-debug-logs) section for more details regarding the `printer`)

Please remeber to run `make generate-docs` after your changes to keep the commands' documentation updated.

#### Command file structure

Below is a typical structure of a CLI command:

```go
package bar

import (
	(...)
)

// Define consts for command flags
const (
   someArg = "MY_ARG"
   someFlag = "my-flag"
)

// Struct to model user input (arguments and/or flags)
type inputModel struct {
	*globalflags.GlobalFlagModel
	MyArg string
	MyFlag *string
}

// "bar" command constructor
func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bar",
		Short: "Short description of the command (is shown in the help of parent command)",
		Long:  "Long description of the command. Can contain some more information about the command usage. It is shown in the help of the current command.",
		Args:  args.SingleArg(someArg, utils.ValidateUUID), // Validate argument, with an optional validation function
		Example: examples.Build(
			examples.NewExample(
				`Do something with command "bar"`,
				"$ stackit foo bar arg-value --my-flag flag-value"),
			...
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p, cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("(...): %w", err)
			}

         projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}

         // Check API response "resp" and output accordingly
         if resp.Item == nil {
            p.Info("(...)", projectLabel)
				return nil
         }
			return outputResult(p, cmd, model.OutputFormat, instances)
		},
	}

	configureFlags(cmd)
	return cmd
}

// Configure command flags (type, default value, and description)
func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(myFlag, "defaultValue", "My flag description")
}

// Parse user input (arguments and/or flags)
func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
   myArg := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
      MyArg            myArg,
		MyFlag:          flags.FlagToStringPointer(cmd, myFlag),
	}, nil

	// Write the input model to the debug logs
	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

// Build request to the API
func buildRequest(ctx context.Context, model *inputModel, apiClient *foo.APIClient) foo.ApiListInstancesRequest {
	req := apiClient.GetBar(ctx, model.ProjectId, model.MyArg, someParam)
	return req
}

// Output result based on the configured output format
func outputResult(p *print.Printer, cmd *cobra.Command, outputFormat string, resources []foo.Resource) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resources, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal resource list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE")
		for i := range resources {
			resource := resources[i]
			table.AddRow(*resource.ResourceId, *resource.Name, *resource.State)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
```

Please remember to always add unit tests for `parseInput`, `buildRequest` (in `bar_test.go`), and any other util functions used.

If the new command `bar` is the first command in the CLI using a STACKIT service `foo`, please refer to [Onboarding a new STACKIT service](./CONTRIBUTION.md/#onboarding-a-new-stackit-service).

#### Outputs, prints and debug logs

The CLI has 4 different verbosity levels:

- `error`: For only displaying errors
- `warning`: For displaying user facing warnings _(and all of the above)_
- `info` (default): For displaying user facing info, such as operation success messages and spinners _(and all of the above)_
- `debug`: For displaying structured logs with different levels, including errors _(and all of the above)_

For prints that are specific to a certain log level, you can use the methods defined in the `print` package: `Error`, `Warn`, `Info`, and `Debug`.

For command outputs that should always be displayed, no matter the defined verbosity, you should use the `print` methods `Outputf` and `Outputln`. These should only be used for the actual output of the commands, which can usually be described by "I ran the command to see _this_".

### Onboarding a new STACKIT service

If you want to add a command that uses a STACKIT service `foo` that was not yet used by the CLI, you will first need to implement a few extra steps to configure the new service:

1.  Add a `FooCustomEndpointKey` key in `internal/pkg/config/config.go` (and add it to `ConfigKeys` and set the to default to `""` using `viper.SetDefault`)
2.  Update the `stackit config unset` and `stackit config unset` commands by adding flags to set and unset a custom endpoint for the `foo` service API, respectively, and update their unit tests
3.  Setup the SDK client configuration, using the authentication method configured in the CLI

    1.  This is done in `internal/pkg/services/foo/client/client.go`
    2.  Below is an example of a typical `client.go` file structure:

        ```go
        package client

        import (
           (...)
           "github.com/stackitcloud/stackit-sdk-go/services/foo"
        )

        func ConfigureClient(cmd *cobra.Command) (*foo.APIClient, error) {
           var err error
           var apiClient foo.APIClient
           var cfgOptions []sdkConfig.ConfigurationOption

           authCfgOption, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
           if err != nil {
              return nil, &errors.AuthError{}
           }
           cfgOptions = append(cfgOptions, authCfgOption, sdkConfig.WithRegion("eu01")) // Configuring region is needed if "foo" is a regional API

           customEndpoint := viper.GetString(config.fooCustomEndpointKey)

           if customEndpoint != "" {
              cfgOptions = append(cfgOptions, sdkConfig.WithEndpoint(customEndpoint))
           }

           apiClient, err = foo.NewAPIClient(cfgOptions...)
           if err != nil {
              return nil, &errors.AuthError{}
           }

           return apiClient, nil
        }
        ```

### Local development

To test your changes, you can either:

1. Build the application locally by running:

   ```bash
   $ go build -o ./bin/stackit
   ```

   To use the application from the root of the repository, you can run:

   ```bash
   $ ./bin/stackit [group] [subgroup] [command] [flags]
   ```

2. Skip building and run the Go application directly using:

   ```bash
   $ go run . [group] [subgroup] [command] [flags]
   ```

## Code Contributions

To make your contribution, follow these steps:

1. Check open or recently closed [Pull Requests](https://github.com/stackitcloud/stackit-cli/pulls) and [Issues](https://github.com/stackitcloud/stackit-cli/issues) to make sure the contribution you are making has not been already tackled by someone else.
2. Fork the repo.
3. Make your changes in a branch that is up-to-date with the original repo's `main` branch.
4. Commit your changes including a descriptive message
5. Create a pull request with your changes.
6. The pull request will be reviewed by the repo maintainers. If you need to make further changes, make additional commits to keep commit history. When the PR is merged, commits will be squashed.

## Bug Reports

If you would like to report a bug, please open a [GitHub issue](https://github.com/stackitcloud/stackit-cli/issues/new).

To ensure we can provide the best support to your issue, follow these guidelines:

1. Go through the existing issues to check if your issue has already been reported.
2. Make sure you are using the latest version of the provider, we will not provide bug fixes for older versions. Also, latest versions may have the fix for your bug.
3. Please provide as much information as you can about your environment, e.g. your version of Go, your version of the provider, which operating system you are using and the corresponding version.
4. Include in your issue the steps to reproduce it, along with code snippets and/or information about your specific use case. This will make the support process much easier and efficient.
