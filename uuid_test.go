package uuid

import (
	"fmt"
	"time"
	"testing"
)

func TestUUID(t *testing.T) {
	id := New("T")
	fmt.Println("New UUID", id)

	if !Valid(id) {
		fmt.Println("Newly generated id is not valid")
		t.FailNow()
	}
	fmt.Println(" - is valid")

	if Type(id) != "T" {
		fmt.Println("Newly generated id is not of right type")
		t.FailNow()
	}
	fmt.Println(" - is correct type");

	age := time.Now().Sub(*Date(id))
	if age.Seconds() > 5 {
		fmt.Println("Newly generated id is too old", age);
		t.FailNow()
	}
	if age.Seconds() < 0 {
		fmt.Println("Newly generated id is too young");
		t.FailNow()
	}
	fmt.Println(" - is only", age, "old");

	id = "8ZwLeTmHZDb89WknvTtq"
	fmt.Println("Constant UUID", id)

	if !Valid(id) {
		fmt.Println("Constant id is not valid")
		t.FailNow()
	}
	fmt.Println(" - is valid")

	if Type(id) != "T" {
		fmt.Println("Constant id is not of right type")
		t.FailNow()
	}
	fmt.Println(" - is correct type");

	date := Date(id).UTC().Format(time.RFC3339)

	if date != "2015-03-25T12:48:56Z" {
		fmt.Println("Constant id did not decode correct date")
		t.FailNow()
	}
	fmt.Println(" - date is correct", date)
}
