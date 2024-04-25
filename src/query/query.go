package query

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/DillonBarker/d8b/src/model"
)

func LoadQueries() (map[string]string, model.Choices, error) {
	var queries map[string]string
	var choices model.Choices

	if _, err := toml.DecodeFile("queries.toml", &choices); err != nil {
		return queries, choices, err
	}

	queryMap := make(map[string]string)
	for _, choice := range choices.Choice {
		queryMap[choice.Name] = choice.Query
	}

	return queryMap, choices, nil
}

func SaveQueries(choices model.Choices) error {
	file, err := os.OpenFile("queries.toml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := toml.NewEncoder(file).Encode(choices); err != nil {
		return err
	}
	return nil
}
