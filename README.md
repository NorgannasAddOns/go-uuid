# go-uuid

20 character quasi-sequential complex unique id generator with type character, embedded timestamp, and check digit.

## Example:

    package main
    
    import (
      "fmt"
      "github.com/NorgannasAddOns/go-uuid"
    )

    func main() {
      id := uuid.New("H")
      fmt.Println("UUID", id)
      fmt.Println("Is valid?", uuid.Valid(id))
      fmt.Println("Generated around", uuid.Date(id))
    }

## License:

Project code is released under CC0 license:

<a rel="license" href="http://creativecommons.org/publicdomain/zero/1.0/">
<img src="http://i.creativecommons.org/p/zero/1.0/88x31.png" style="border-style: none;" alt="CC0" />
</a>
    
  
