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

	if Code(id) != "T" {
		fmt.Println("Newly generated id is not of right code")
		t.FailNow()
	}
	fmt.Println(" - is correct code");

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

	if Code(id) != "T" {
		fmt.Println("Constant id is not of right code")
		t.FailNow()
	}
	fmt.Println(" - is correct code");

	tm := Date(id)
	date := tm.UTC().Format(time.RFC3339)

	if date != "2015-03-25T12:48:56Z" {
		fmt.Println("Constant id did not decode correct date")
		t.FailNow()
	}
	fmt.Println(" - date is correct", date)

	idBefore := Before(*tm)
	fmt.Println("Before id", idBefore)

	if idBefore >= id {
		fmt.Println("Before id is on the wrong side of constant id")
		t.FailNow()
	}
	fmt.Println(" - is less than constant id of same date")

	if idBefore[:4] != id[:4] {
		fmt.Println("Before id does not have same timefield as constant id")
		t.FailNow()
	}
	fmt.Println(" - matches constant id timefield")
	
}
