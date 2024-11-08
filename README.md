# confprint

A tiny Go library for securely printing configuration values with smart masking of secrets.

## Installation

```bash
go get github.com/adyatlov/confprint
```

## Usage

Basic usage with [caarlos0/env](https://github.com/caarlos0/env):

```go
package main

import (
    "os"
    "github.com/adyatlov/confprint"
    "github.com/caarlos0/env/v10"
)

type Config struct {
    Port        int    `env:"PORT" envDefault:"8080" safe:"true"`
    Environment string `env:"ENV" envDefault:"development" safe:"true"`
    APIKey      string `env:"API_KEY"`                    // Will be masked
    SecretKey   string `env:"SECRET_KEY"`                 // Will be masked
}

func main() {
    cfg := Config{}
    if err := env.Parse(&cfg); err != nil {
        panic(err)
    }

    // Print to stdout with default settings
    confprint.Print(os.Stdout, &cfg)
}
```

Output:
```
=== Configuration ===
APIKey:      ********xyz
Environment: development
Port:        8080
SecretKey:   ********
```

### Customization

```go
confprint.Print(os.Stdout, &cfg,
    confprint.WithMaskLength(4),        // Show fewer asterisks
    confprint.WithVisibleSuffix(4),     // Show more characters at the end
    confprint.WithMinSecretLen(15),     // Adjust when to show ending chars
)
```

## Defaults

- Mask length: 8 characters (`********`)
- Visible suffix: 3 characters
- Minimum secret length for showing suffix: 20 characters

## License

MIT License - see LICENSE file
