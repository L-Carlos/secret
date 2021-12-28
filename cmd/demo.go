package main

import (
	"fmt"
	"log"

	"github.com/L-Carlos/secret"
)

func main() {
	v := secret.MemoryVault("fake-key")

	err := v.Set("demo_key", "demo_value")
	must(err)

	plain, err := v.Get("demo_key")
	must(err)

	fmt.Println("plain:", plain)

}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
