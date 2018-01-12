package nocomment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripCStyleComments_multiline(t *testing.T) {
	res := StripCStyleComments(`line1 // comment // comment
		line 2
		line 3 // comment
		line 4 /// comment
		line 5 //// comment
		line 6 ///// comment`)
	assert.Equal(t, `line1 
		line 2
		line 3 
		line 4 
		line 5 
		line 6 `, res)
}

func TestStripCStyleComments_lineComments(t *testing.T) {
	var testCases = []struct {
		In       string
		Expected string
	}{
		// simple
		{
			"text // text\n",
			"text \n"},
		{
			`text // text`,
			`text `},
		{
			"te/xt /",
			"te/xt /"},

		// with some quotes
		{
			`text "txt // text`,
			`text "txt // text`}, // If the quote is not closed, comments are not sripped!
		{
			"text \"txt // text\ntext \"txt // text",
			"text \"txt // text\ntext \"txt // text"}, // If the quote is not closed before a new-line, comments are not sripped!
		{
			`text ""txt // text """`,
			`text ""txt `},
		{
			`text """"txt // text """"`,
			`text """"txt `},
		{
			`text """""txt // text`,
			`text """""txt // text`},
		// comments within quotes are ignored
		{
			`text "//" text`,
			`text "//" text`},
		{
			`text """//""" text`,
			`text """//""" text`},
		// quote spans multiple lines
		{
			"text \"line1 //\\\nline // 2\\\nline 3\"//comment",
			"text \"line1 //\\\nline // 2\\\nline 3\""},
		// quote spans multiple lines and does not close
		{
			"text \"line1 //\\\nline // 2\\\nline 3//comment\nnext\"line//com",
			"text \"line1 //\\\nline // 2\\\nline 3//comment\nnext\"line//com"},
		// comments cannot be escaped
		{
			`text \// text`,
			`text \`},
		{
			"text // text\\\ntext",
			"text \ntext"},
		{
			`text \/* text */`,
			`text \`},
		{
			`text /* text \*/ text2 */`,
			`text  text2 */`},
		// escapes before comment-similar strings work as usual
		{
			`text \/text`,
			`text \/text`},
		{
			`text \/"text//text"`,
			`text \/"text//text"`},
		// escaped quotes
		{
			`text \"txt // text`,
			`text \"txt `},
		{
			`text "txt \"// text"// text`,
			`text "txt \"// text"`},
		{
			`text "\"txt // text \"""`,
			`text "\"txt // text \"""`},
		{
			`text "\"\""txt // text "\"\""`,
			`text "\"\""txt `},
		{
			`text """"\"txt// text`,
			`text """"\"txt`},

		// escaped escapes before quotes
		{
			`text \\"txt // text`,
			`text \\"txt // text`},
		{
			`text "txt \\"// text"// text`,
			`text "txt \\"`},
		{
			`text "\\"\\"\z"txt // text "\"\""`,
			`text "\\"\\"\z"txt `},
		{
			`text \\\\"\\\\"\\"\\\\\\"// text`,
			`text \\\\"\\\\"\\"\\\\\\"`},
	}

	for i, testCase := range testCases {
		out := StripCStyleComments(testCase.In)
		assert.Equal(t, testCase.Expected, out, "Testcase %d. input: %s", i, testCase.In)
	}
}

func TestStripCStyleComments_multiLineComments_startOnly(t *testing.T) {
	var testCases = []struct {
		In       string
		Expected string
	}{
		// simple
		{
			"text /* text\n",
			"text "},
		{
			`text /* text`,
			`text `},
		{
			"text /",
			"text /"},

		// with some quotes
		{
			`text "txt /* text`,
			`text "txt /* text`}, // If the quote is not closed, comments are not sripped!
		{
			"text \"txt /* text\ntext \"txt /* text",
			"text \"txt /* text\ntext \"txt /* text"}, // If the quote is not closed before a new-line, comments are not sripped!
		{
			`text ""txt /* text """`,
			`text ""txt `},
		{
			`text """"txt /* text """"`,
			`text """"txt `},
		{
			`text """""txt /* text`,
			`text """""txt /* text`},
		// comments within quotes are ignored
		{
			`text "/*" text`,
			`text "/*" text`},
		{
			`text """/*""" text`,
			`text """/*""" text`},
		// quote spans multiple lines
		{
			"text \"line1 /*\\\nline /* 2\\\nline 3\"/*comment",
			"text \"line1 /*\\\nline /* 2\\\nline 3\""},
		// quote spans multiple lines and does not close
		{
			"text \"line1 /*\\\nline /* 2\\\nline 3/*comment\nnext\"line/*com",
			"text \"line1 /*\\\nline /* 2\\\nline 3/*comment\nnext\"line/*com"},

		// escaped quotes
		{
			`text \"txt /* text`,
			`text \"txt `},
		{
			`text "txt \"/* text"/* text`,
			`text "txt \"/* text"`},
		{
			`text "\"txt /* text \"""`,
			`text "\"txt /* text \"""`},
		{
			`text "\"\""txt /* text "\"\""`,
			`text "\"\""txt `},
		{
			`text """"\"txt/* text`,
			`text """"\"txt`},

		// escaped escapes before quotes
		{
			`text \\"txt /* text`,
			`text \\"txt /* text`},
		{
			`text "txt \\"/* text"/* text`,
			`text "txt \\"`},
		{
			`text "\\"\\"\z"txt /* text "\"\""`,
			`text "\\"\\"\z"txt `},
		{
			`text \\\\"\\\\"\\"\\\\\\"/* text`,
			`text \\\\"\\\\"\\"\\\\\\"`},
	}

	for i, testCase := range testCases {
		out := StripCStyleComments(testCase.In)
		assert.Equal(t, testCase.Expected, out, "Testcase %d. input: %s", i, testCase.In)
	}
}

func TestStripCStyleComments_multiLineComments_startAndEnd(t *testing.T) {
	var testCases = []struct {
		In       string
		Expected string
	}{
		{
			"text/*text*/text",
			"texttext"},
		{
			"text/**text**/text",
			"texttext"},
		{
			"text/***text/***/text",
			"texttext"},
		{
			"text/*te\n\nxt*/text",
			"texttext"},
		{
			"text/*text*\ntext*/text",
			"texttext"},
		{
			"000/**/111/*/222/**/333",
			"000111333"},
		{
			"000/** /111/*/222/**/333",
			"000222333"},
	}

	for i, testCase := range testCases {
		out := StripCStyleComments(testCase.In)
		assert.Equal(t, testCase.Expected, out, "Testcase %d. input: %s", i, testCase.In)
	}
}

func TestStripCStyleComments_mixedComments(t *testing.T) {
	var testCases = []struct {
		In       string
		Expected string
	}{
		{
			"text/*te//xt*/text",
			"texttext"},
		{
			"text//*text*/text",
			"text"},
		{
			"text///*text*/text",
			"text"},
		{
			"text/*te\n//\nxt*/text",
			"texttext"},
	}

	for i, testCase := range testCases {
		out := StripCStyleComments(testCase.In)
		assert.Equal(t, testCase.Expected, out, "Testcase %d. input: %s", i, testCase.In)
	}
}
