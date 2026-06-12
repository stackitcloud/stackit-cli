package flags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type stringEnumFlag[T ~string] struct {
	ignoreCase bool
	options    []T
	value      T
	valueSet   bool
	docs       string
	name       string
}

type StringEnumFlagOption[T ~string] func(*stringEnumFlag[T])

func StringEnumIgnoreCase[T ~string]() StringEnumFlagOption[T] {
	return func(f *stringEnumFlag[T]) {
		f.ignoreCase = true
	}
}

func StringEnumDefaultValue[T ~string](value T) StringEnumFlagOption[T] {
	return func(f *stringEnumFlag[T]) {
		f.value = value
		f.valueSet = true
	}
}

func StringEnumFlag[T ~string](name string, possibleValues []T, docs string, opts ...StringEnumFlagOption[T]) *stringEnumFlag[T] {
	f := &stringEnumFlag[T]{
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

var _ pflag.Value = &stringEnumFlag[string]{}

func (s *stringEnumFlag[T]) Register(cmd *cobra.Command) {
	cmd.Flags().Var(s, s.name, s.Usage())
}

func (s *stringEnumFlag[T]) Usage() string {
	return s.docs + fmt.Sprintf(" (possible values: %s)", s.fmtValues(s.options))
}

func (s *stringEnumFlag[T]) Get() T {
	return s.value
}

func (s *stringEnumFlag[T]) Ptr() *T {
	if s.valueSet {
		return &s.value
	}
	return nil
}

func (s *stringEnumFlag[T]) Name() string {
	return s.name
}

func (s *stringEnumFlag[T]) String() string {
	return string(s.value)
}

func (s *stringEnumFlag[T]) fmtValues(xs []T) string {
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

func (s *stringEnumFlag[T]) Set(value string) error {
	v := strings.TrimSpace(value)

	if v == "" {
		return fmt.Errorf("value cannot be empty")
	}

	for _, o := range s.options {
		if !s.ignoreCase && v == string(o) {
			s.value = T(v)
			s.valueSet = true
			return nil
		} else if s.ignoreCase && strings.EqualFold(v, string(o)) {
			s.value = T(strings.ToLower(v))
			s.valueSet = true
			return nil
		}
	}

	return fmt.Errorf("found value %q, expected one of %q", v, s.options)
}

func (s *stringEnumFlag[T]) Type() string {
	return "string"
}
