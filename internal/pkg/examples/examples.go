package examples

import "fmt"

type Example struct {
	Description string
	Commands    []string
}

// Creates a new example
func NewExample(description string, commands ...string) Example {
	return Example{
		Description: description,
		Commands:    commands,
	}
}

// Returns the example formatted
func (e *Example) format() string {
	formatted := fmt.Sprintf("  %s\n", e.Description)
	for i, c := range e.Commands {
		formatted += fmt.Sprintf("  %s", c)
		if i != len(e.Commands)-1 {
			formatted += "\n"
		}
	}
	return formatted
}

// Builds a list of formatted examples
func Build(examples ...Example) string {
	formatted := ""
	for i, e := range examples {
		formatted += e.format()
		if i != len(examples)-1 {
			formatted += "\n\n"
		}
	}
	return formatted
}
