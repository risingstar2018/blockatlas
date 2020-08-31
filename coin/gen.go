// +build coins

//go:generate rm -f coins.go
//go:generate go run gen.go

package main

import (
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	coinFile     = "coins.yml"
	filename     = "coins.go"
	templateFile = `// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
// using data from coins.yml
package coin

import (
	"fmt"
)

// Coin is the native currency of a blockchain
type Coin struct {
	ID               uint
	Handle           string
	Symbol           string
	PreferedSymbol   string
	Name             string
	Decimals         uint
	BlockTime        int
	MinConfirmations int64
	SampleAddr       string
}

func (c *Coin) String() string {
	return fmt.Sprintf("[%s] %s (#%d)", c.Symbol, c.Name, c.ID)
}

const (
{{- range .Coins }}
{{- if .PreferedSymbol}}
	{{ .PreferedSymbol }} = {{ .ID }}
{{- else}}
	{{ .Symbol }} = {{ .ID }}
{{- end}}
{{- end }}
)

var Coins = map[uint]Coin{
{{- range .Coins }}
{{- if .PreferedSymbol }}
	{{ .PreferedSymbol }}: {
{{- else }}
	{{ .Symbol }}: {
{{- end }}
		ID:               {{.ID}},
		Handle:           "{{.Handle}}",
		Symbol:           "{{.Symbol}}",
{{- if .PreferedSymbol }}
		PreferedSymbol:   "{{.PreferedSymbol}}",
{{- end }}
		Name:             "{{.Name}}",
		Decimals:         {{.Decimals}},
		BlockTime:        {{.BlockTime}},
		MinConfirmations: {{.MinConfirmations}},
		SampleAddr:       "{{.SampleAddr}}",
	},
{{- end }}
}

{{- range .Coins }}
func {{ .Handle.Capitalize }}() Coin {
{{- if .PreferedSymbol }}
	return Coins[{{ .PreferedSymbol }}]
{{- else }}
	return Coins[{{ .Symbol }}]
{{- end}}
}

{{- end }}

`
)

type Handle string

func (h Handle) Capitalize() string {
	return strings.Title(string(h))
}

type Coin struct {
	ID               uint   `yaml:"id"`
	Handle           Handle `yaml:"handle"`
	Symbol           string `yaml:"symbol"`
	PreferedSymbol   string `yaml:"preferedSymbol,omitempty"`
	Name             string `yaml:"name"`
	Decimals         uint   `yaml:"decimals"`
	BlockTime        int    `yaml:"blockTime"`
	MinConfirmations int64  `yaml:"minConfirmations"`
	SampleAddr       string `yaml:"sampleAddress"`
}

func main() {
	coinFile := getValidParameter("COIN_FILE", coinFile)
	var coinList []Coin
	coin, err := os.Open(coinFile)
	dec := yaml.NewDecoder(coin)
	err = dec.Decode(&coinList)
	if err != nil {
		log.Panic(err)
	}

	goFile := getValidParameter("COIN_GO_FILE", filename)
	f, err := os.Create(goFile)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	coinsTemplate := template.Must(template.New("").Parse(templateFile))
	err = coinsTemplate.Execute(f, map[string]interface{}{
		"Timestamp": time.Now(),
		"Coins":     coinList,
	})
	if err != nil {
		log.Panic(err)
	}
}

func getValidParameter(env, variable string) string {
	e, ok := os.LookupEnv(env)
	if ok {
		return e
	}
	return variable
}
