package queries

import "github.com/BurntSushi/toml"

type Choice struct {
	Name  string `toml:"name"`
	Query string `toml:"query"`
}

type Choices struct {
	Choices []Choice `toml:"choice"`
}

func LoadQueries() (Choices, error) {
	var choices Choices
	if _, err := toml.DecodeFile("queries.toml", &choices); err != nil {
		return choices, err
	}
	return choices, nil
}
