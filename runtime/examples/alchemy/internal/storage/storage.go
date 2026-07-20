package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/meowpow-png/mipe/runtime/examples/alchemy/internal/recipe"
)

func Load(path string) (recipe.Book, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return recipe.Book{}, fmt.Errorf("load recipes: %w", err)
	}

	var recipes []recipe.Recipe
	if err := json.Unmarshal(data, &recipes); err != nil {
		return recipe.Book{}, fmt.Errorf("parse recipes: %w", err)
	}
	return recipe.NewBook(recipes), nil
}
