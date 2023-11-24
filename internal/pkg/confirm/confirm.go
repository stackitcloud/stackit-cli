package confirm

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var errAborted = errors.New("operation aborted")

// Prompts the user for confirmation.
//
// Returns nil only if the user (explicitly) answers positive.
// Returns ErrAborted if the user answers negative.
func PromptForConfirmation(cmd *cobra.Command, prompt string) error {
	question := fmt.Sprintf("%s [y/N] ", prompt)
	reader := bufio.NewReader(cmd.InOrStdin())
	for i := 0; i < 3; i++ {
		cmd.Print(question)
		answer, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "y" || answer == "yes" {
			return nil
		}
		if answer == "" || answer == "n" || answer == "no" {
			return errAborted
		}
	}
	return fmt.Errorf("max number of wrong inputs")
}
