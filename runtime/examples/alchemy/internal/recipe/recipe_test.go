package recipe

import "testing"

func TestBookFind(t *testing.T) {
	book := NewBook([]Recipe{{Name: "Healing Draught"}})
	got, err := book.Find("Healing Draught")
	if err != nil || got.Name != "Healing Draught" {
		t.Fatalf("Find() = %#v, %v", got, err)
	}
	if _, err := book.Find("Missing"); err == nil {
		t.Fatal("Find() should reject unknown recipes")
	}
}
