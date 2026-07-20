package recipe

import "fmt"

type Recipe struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Ingredients []string `json:"ingredients"`
}

type Book struct {
	recipes []Recipe
}

func NewBook(recipes []Recipe) Book {
	return Book{recipes: recipes}
}

func (b Book) All() []Recipe {
	return b.recipes
}

func (b Book) Find(name string) (Recipe, error) {
	for _, recipe := range b.recipes {
		if recipe.Name == name {
			return recipe, nil
		}
	}
	return Recipe{}, fmt.Errorf("recipe %q not found", name)
}
