# Autocompletion

This guide describes how you can enable command autocompletion for the STACKIT CLI by leveraging the functionality provided the [Cobra](https://github.com/spf13/cobra) framework.

The process may vary depending on the type of shell you are using and your operating system (OS).

You will need to start a new shell for the setup to take effect.

## bash

This process depends on the `bash-completion` package. If you don't have it installed already, you can install it via your OS's package manager.

#### Linux

```shell
stackit completion bash > ~/.bash_completion
```

#### macOS

```shell
stackit completion bash > $(brew --prefix)/etc/bash_completion.d/stackit
```

## zsh

If shell completion is not already enabled in your environment you will need to enable it by executing the following once:

```shell
echo "autoload -U compinit; compinit" >> ~/.zshrc
```

#### Linux

```shell
stackit completion zsh > "${fpath[1]}/_stackit"
```

#### macOS

```shell
stackit completion zsh > $(brew --prefix)/share/zsh/site-functions/_stackit
```

Additionally, you might also need to run:

```shell
source $(brew --prefix)/share/zsh/site-functions/_stackit >> ~/.zshrc
```

## PowerShell

You can load completions for your current shell session by running:

```shell
stackit completion powershell | Out-String | Invoke-Expression
```

To load completions for every new session, add the output of the above command to your PowerShell profile.

## fish

```shell
stackit completion fish > ~/.config/fish/completions/stackit.fish
```
