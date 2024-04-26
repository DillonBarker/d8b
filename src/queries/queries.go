package queries

import (
	"os"

	"github.com/BurntSushi/toml"
)

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

func AddQuery(choice Choice) error {
	choices, err := LoadQueries()

	if err != nil {
		return err
	}

	file, err := os.Create("queries.toml")

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)

	choices.Choices = append(choices.Choices, choice)

	if err := encoder.Encode(choices); err != nil {
		return err
	}

	return nil
}

func RemoveQuery(choice Choice) error {
	choices, err := LoadQueries()

	if err != nil {
		return err
	}

	file, err := os.Create("queries.toml")

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := toml.NewEncoder(file)

	choices.Choices = append(choices.Choices, choice)

	newChoices := Choices{}

	for _, item := range choices.Choices {
		if item.Name != choice.Name {
			newChoices.Choices = append(newChoices.Choices, item)
		}
	}

	if err := encoder.Encode(newChoices); err != nil {
		return err
	}

	return nil
}
