package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/meowpow-png/mipe/runtime/examples/alchemy/internal/recipe"
)

func List(w io.Writer, recipes []recipe.Recipe) {
	for _, item := range recipes {
		_, _ = fmt.Fprintln(w, item.Name)
	}
}

func Show(w io.Writer, item recipe.Recipe) {
	_, _ = fmt.Fprintf(
		w,
		"%s\n%s\nIngredients: %s\n",
		item.Name,
		item.Description,
		join(item.Ingredients),
	)
}

func Brew(w io.Writer, item recipe.Recipe) {
	_, _ = fmt.Fprintf(
		w,
		"Brewed %s: %s\n",
		item.Name,
		item.Description,
	)
}

func join(items []string) string {
	var result strings.Builder
	for i, item := range items {
		if i > 0 {
			result .WriteString(", ")
		}
		result .WriteString(item)
	}
	return result.String()
}
