package confprint

import (
	"bytes"
	"testing"
)

type testConfig struct {
	APIKey      string `safe:"false"`
	Port        int    `safe:"true"`
	SecretKey   string
	Environment string `safe:"true"`
}

func TestPrint(t *testing.T) {
	tests := []struct {
		name    string
		cfg     interface{}
		opts    []Option
		wantErr bool
		want    string
	}{
		{
			name: "basic config with default options",
			cfg: testConfig{
				APIKey:      "very-long-secret-key-123",
				Port:        8080,
				SecretKey:   "short-secret",
				Environment: "development",
			},
			want: `=== Configuration ===
APIKey     : ********123
Environment: development
Port       : 8080
SecretKey  : ********
`,
		},
		{
			name: "pointer to struct",
			cfg: &testConfig{
				APIKey:      "very-long-secret-key-123",
				Port:        8080,
				SecretKey:   "short-secret",
				Environment: "development",
			},
			want: `=== Configuration ===
APIKey     : ********123
Environment: development
Port       : 8080
SecretKey  : ********
`,
		},
		{
			name: "custom mask length",
			cfg: testConfig{
				APIKey:      "very-long-secret-key-123",
				Port:        8080,
				SecretKey:   "short-secret",
				Environment: "development",
			},
			opts: []Option{WithMaskLength(4)},
			want: `=== Configuration ===
APIKey     : ****123
Environment: development
Port       : 8080
SecretKey  : ****
`,
		},
		{
			name: "custom visible suffix",
			cfg: testConfig{
				APIKey:      "very-long-secret-key-12345",
				Port:        8080,
				SecretKey:   "short-secret",
				Environment: "development",
			},
			opts: []Option{WithVisibleSuffix(5)},
			want: `=== Configuration ===
APIKey     : ********12345
Environment: development
Port       : 8080
SecretKey  : ********
`,
		},
		{
			name: "custom minimum secret length",
			cfg: testConfig{
				APIKey:      "short-key", // Less than default minSecretLen
				Port:        8080,
				SecretKey:   "also-short",
				Environment: "development",
			},
			opts: []Option{WithMinSecretLen(5)}, // Set to show suffix for shorter secrets
			want: `=== Configuration ===
APIKey     : ********key
Environment: development
Port       : 8080
SecretKey  : ********ort
`,
		},
		{
			name: "multiple options",
			cfg: testConfig{
				APIKey:      "very-long-secret-key-12345",
				Port:        8080,
				SecretKey:   "short-secret",
				Environment: "development",
			},
			opts: []Option{
				WithMaskLength(4),
				WithVisibleSuffix(5),
				WithMinSecretLen(15),
			},
			want: `=== Configuration ===
APIKey     : ****12345
Environment: development
Port       : 8080
SecretKey  : ****
`,
		},
		{
			name:    "non-struct input",
			cfg:     "not a struct",
			wantErr: true,
		},
		{
			name: "empty config",
			cfg:  testConfig{},
			want: `=== Configuration ===
APIKey     : ********
Environment: 
Port       : 0
SecretKey  : ********
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Print(&buf, tt.cfg, tt.opts...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("Print() output mismatch:\nwant:\n%s\ngot:\n%s", tt.want, got)
			}
		})
	}
}

func TestPrinterOptions(t *testing.T) {
	tests := []struct {
		name          string
		opts          []Option
		wantMaskLen   int
		wantSuffixLen int
		wantMinSecret int
	}{
		{
			name: "default options",
			opts: nil,
			wantMaskLen:   defaultMaskLength,
			wantSuffixLen: defaultLastCharacters,
			wantMinSecret: defaultMinSecretLen,
		},
		{
			name: "custom mask length",
			opts: []Option{WithMaskLength(4)},
			wantMaskLen:   4,
			wantSuffixLen: defaultLastCharacters,
			wantMinSecret: defaultMinSecretLen,
		},
		{
			name: "custom suffix length",
			opts: []Option{WithVisibleSuffix(5)},
			wantMaskLen:   defaultMaskLength,
			wantSuffixLen: 5,
			wantMinSecret: defaultMinSecretLen,
		},
		{
			name: "custom minimum secret length",
			opts: []Option{WithMinSecretLen(15)},
			wantMaskLen:   defaultMaskLength,
			wantSuffixLen: defaultLastCharacters,
			wantMinSecret: 15,
		},
		{
			name: "all options combined",
			opts: []Option{
				WithMaskLength(4),
				WithVisibleSuffix(5),
				WithMinSecretLen(15),
			},
			wantMaskLen:   4,
			wantSuffixLen: 5,
			wantMinSecret: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newPrinter(tt.opts...)

			if p.maskLength != tt.wantMaskLen {
				t.Errorf("maskLength = %d, want %d", p.maskLength, tt.wantMaskLen)
			}
			if p.lastCharacters != tt.wantSuffixLen {
				t.Errorf("lastCharacters = %d, want %d", p.lastCharacters, tt.wantSuffixLen)
			}
			if p.minSecretLen != tt.wantMinSecret {
				t.Errorf("minSecretLen = %d, want %d", p.minSecretLen, tt.wantMinSecret)
			}
		})
	}
}
