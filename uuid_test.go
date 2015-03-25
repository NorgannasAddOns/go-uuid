package uuid

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := New("T")
	fmt.Println("UUID", uuid)
	fmt.Println("Date", Date(uuid))
	fmt.Println("Valid", Valid(uuid))

	fmt.Println("Date", Date("L8ZEmwuQu2WoXJQhZuz7"))
}
