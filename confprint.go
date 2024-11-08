package confprint

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

type configEntry struct {
	Key   string
	Value string
}

const (
	defaultMaskLength     = 8
	defaultLastCharacters = 3
	defaultMinSecretLen   = 20
)

// Option allows customizing the printer behavior
type Option func(*printer)

type printer struct {
	maskLength     int
	lastCharacters int
	minSecretLen   int
}

// WithMaskLength customizes the length of the masking (default: 8)
func WithMaskLength(length int) Option {
	return func(p *printer) {
		p.maskLength = length
	}
}

// WithVisibleSuffix customizes the number of visible characters at the end (default: 3)
func WithVisibleSuffix(length int) Option {
	return func(p *printer) {
		p.lastCharacters = length
	}
}

// WithMinSecretLen customizes minimum length for showing suffix (default: 20)
func WithMinSecretLen(length int) Option {
	return func(p *printer) {
		p.minSecretLen = length
	}
}

func newPrinter(opts ...Option) *printer {
	p := &printer{
		maskLength:     defaultMaskLength,
		lastCharacters: defaultLastCharacters,
		minSecretLen:   defaultMinSecretLen,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Print formats and prints configuration to the provided writer
func Print(w io.Writer, cfg interface{}, opts ...Option) error {
	p := newPrinter(opts...)

	val := reflect.ValueOf(cfg)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("cfg must be a struct or pointer to struct")
	}

	typ := val.Type()
	var entries []configEntry
	maxKeyWidth := 0

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		rawValue := fmt.Sprintf("%v", val.Field(i).Interface())

		formattedValue := rawValue
		if field.Tag.Get("safe") != "true" {
			formattedValue = p.maskValue(rawValue)
		}

		entries = append(entries, configEntry{
			Key:   field.Name,
			Value: formattedValue,
		})

		if len(field.Name) > maxKeyWidth {
			maxKeyWidth = len(field.Name)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	fmt.Fprintln(w, "=== Configuration ===")
	format := fmt.Sprintf("%%-%ds: %%s\n", maxKeyWidth)
	for _, entry := range entries {
		fmt.Fprintf(w, format, entry.Key, entry.Value)
	}

	return nil
}

func (p *printer) maskValue(value string) string {
	if len(value) >= p.minSecretLen {
		return fmt.Sprintf("%s%s",
			strings.Repeat("*", p.maskLength),
			value[len(value)-p.lastCharacters:])
	}
	return strings.Repeat("*", p.maskLength)
}
