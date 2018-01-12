package main

import (
	"fmt"

	nocomment "github.com/maja42/no-comment"
)

func main() {
	input := `
	Line comment //gets stripped away!
	Block /*gets stripped away!*/ comment
	"Quoted /*text*/ \" stays 'til the end! \\" /* comments don't \*/ really!
	"And quoted text \
	//can also span multiple lines!"
	`

	out := nocomment.StripCStyleComments(input)

	fmt.Println(out)
	// Prints:
	// 	Reservoir Dogs
	// 	Airplane!

}
