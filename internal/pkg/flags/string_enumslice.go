package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type stringEnumSliceFlag[T ~string] struct {
	ignoreCase bool
	options    []T
	value      []T
	valueSet   bool
	docs       string
	name       string
}

type StringEnumSliceFlagOption[T ~string] func(*stringEnumSliceFlag[T])

func IgnoreCase[T ~string]() StringEnumSliceFlagOption[T] {
	return func(f *stringEnumSliceFlag[T]) {
		f.ignoreCase = true
	}
}

func DefaultValues[T ~string](values ...T) StringEnumSliceFlagOption[T] {
	return func(f *stringEnumSliceFlag[T]) {
		f.value = append(f.value, values...)
	}
}

func StringEnumSliceFlag[T ~string](name string, possibleValues []T, docs string, opts ...StringEnumSliceFlagOption[T]) *stringEnumSliceFlag[T] {
	f := &stringEnumSliceFlag[T]{
		name: name,
		docs: docs,
	}
	for _, v := range possibleValues {
		if string(v) != "unknown_default_open_api" {
			f.options = append(f.options, v)
		}
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

var _ pflag.Value = &stringEnumSliceFlag[string]{}

func (s *stringEnumSliceFlag[T]) Register(cmd *cobra.Command) {
	cmd.Flags().Var(s, s.name, s.Usage())
}

func (s *stringEnumSliceFlag[T]) Usage() string {
	return s.docs + fmt.Sprintf(" (possible values: %s)", s.fmtValues(s.options))
}

func (s *stringEnumSliceFlag[T]) Get() []T {
	return s.value
}

func (s *stringEnumSliceFlag[T]) Name() string {
	return s.name
}

func (s *stringEnumSliceFlag[T]) String() string {
	return s.fmtValues(s.value)
}

func (s *stringEnumSliceFlag[T]) fmtValues(xs []T) string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range xs {
		sb.WriteString(string(v))
		if i != len(xs)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func (s *stringEnumSliceFlag[T]) Set(value string) error {
	// If the default value is still set, remove it
	// (Since we're going to append the incoming values to f.value)
	if !s.valueSet {
		s.value = []T{}
		s.valueSet = true
	}

	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}
	values := strings.Split(value, ",")
	return s.appendToValue(values)
}

func (s *stringEnumSliceFlag[T]) Type() string {
	return "stringSlice"
}

func (s *stringEnumSliceFlag[T]) appendToValue(values []string) error {
	for _, v := range values {
		v = strings.TrimSpace(v)

		foundValid := false
		for _, o := range s.options {
			if !s.ignoreCase && v == string(o) {
				s.value = append(s.value, T(v))
				foundValid = true
				break
			} else if s.ignoreCase && strings.EqualFold(v, string(o)) {
				s.value = append(s.value, T(strings.ToLower(v)))
				foundValid = true
				break
			}
		}

		if !foundValid {
			return fmt.Errorf("found value %q, expected one of %q", v, s.options)
		}
	}
	return nil
}
