# no comment!
[![Build Status](https://travis-ci.org/maja42/no-comment.svg?branch=master)](https://travis-ci.org/maja42/no-comment)
[![GoDoc](https://godoc.org/github.com/maja42/no-comment?status.svg)](https://godoc.org/github.com/maja42/no-comment)

No-comment is a simple, light-weight library that can remove comments of various styles and dialects.

The library avoids expensive regular-expression matching and uses minimalstic state machines for optimal performance.

## Supported comment styles

Currently, only C-Style comments can be stripped. Future dialects will be added.

If you require a specific style, create a feature request in the issue tracker.
Pull-requests for new dialects are welcome, as long as they are well-tested and cover all special cases correctly.


## Special cases

All special cases are handled correctly and the behaviour is validated with unit-tests.
Although different comment styles have different properties, many share the same special cases. These can include:


- Block comments can span multiple lines

- Comments cannot start within quoted regions 

	```text "quoted // text"```

- Quoted regions can be escaped

	```text \"not quoted "quoted \" quoted"```

- Escape characters can be escaped

	``` text \\\\"quoted text\\\\" not quoted```

- Quoted regions can span multiple lines, if new-line characters are escaped

	```text "quoted text \\n quoted text on second line"```

- Comments are not allowed to start, if a quote does not have a matching end-quote

	``` text "nearly quoted text //not a comment```

- Comment starts and ends cannot be escaped
	``` text \/*comment\*/ not a comment```

The above rules are some common examples that can be encountered and which are correctly treated by the library.
Additional special cases might apply for certain  dialects.


## Usage

```go
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
```

## Installation

`go get github.com/maja42/no-comment`

## Demo

The _example folder contains a small example for stripping c-style comments.

```
cd _example/
go run main.go
```