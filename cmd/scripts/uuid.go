//go:build ignore

package main

import (
	"fmt"

	"github.com/canonflow/backend-starter/pkg/helpers"
)

func main() {
	uuid := helpers.GenerateUUIDV7()

	fmt.Printf("UUID V7: %s\n", uuid)
}
