package main

import (
	"fmt"

	validator "gopkg.in/validator.v2"
)

func main() {
	validator.SetValidationFunc("valid_input", validate_inputs)

	b := generate_board()
	b.Vertical = example_input_vertical()
	b.Horizontal = example_input_horizontal()
	// b.Vertical = example_input_horizontal()
	// b.Horizontal = example_input_vertical()

	if errs := b.validate(); errs != nil {
		fmt.Println(errs)
	}

	// fmt.Println(b)

	fmt.Println(b.update())
	fmt.Println(b.RowStats)
	fmt.Println(b.ColStats)
	fmt.Println("")

	fmt.Println(b)

	fmt.Println(b.update())
	fmt.Println(b.RowStats)
	fmt.Println(b.ColStats)
	// fmt.Println("")

	// fmt.Println(b)
}
