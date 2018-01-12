package main

import (
	"fmt"

	"github.com/maja42/no-comment"
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
	// 		Line comment
	// 		Block  comment
	// 		"Quoted /*text*/ \" stays 'til the end! \\"  really!
	// 		"And quoted text \
	// 		//can also span multiple lines!"
}
