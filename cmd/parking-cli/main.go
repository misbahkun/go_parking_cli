
package main

import (
	"os"
	"strconv"

	parking "github.com/misbahkun/go_parking_cli/v2"
)

func main() {
	// Default parking capacity
	capacity := 10

	// If argument is provided, use it as capacity
	if len(os.Args) > 1 {
		if c, err := strconv.Atoi(os.Args[1]); err == nil {
			capacity = c
		}
	}

	// Start the parking CLI with the specified capacity
	parking.ParkingCLI(capacity)
}