package grammar

import (
	"strconv"
	"testing"
)

var ExprGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":   {ExpansionTuple{name: "<expr>"}},
	"<expr>":    {ExpansionTuple{name: "<term> + <expr>"}, ExpansionTuple{name: "<term> - <expr>"}, ExpansionTuple{name: "<term>"}},
	"<term>":    {ExpansionTuple{name: "<factor> * <term>"}, ExpansionTuple{name: "<factor> / <term>"}, ExpansionTuple{name: "<factor>"}},
	"<factor>":  {ExpansionTuple{name: "+<factor>"}, ExpansionTuple{name: "-<factor>"}, ExpansionTuple{name: "(<expr>)"}, ExpansionTuple{name: "<integer>.<integer>"}, ExpansionTuple{name: "<integer>"}},
	"<integer>": {ExpansionTuple{name: "<digit><integer>"}, ExpansionTuple{name: "<digit>"}},
	"<digit>":   {ExpansionTuple{name: "0"}, ExpansionTuple{name: "1"}, ExpansionTuple{name: "2"}, ExpansionTuple{name: "3"}, ExpansionTuple{name: "4"}, ExpansionTuple{name: "5"}, ExpansionTuple{name: "6"}, ExpansionTuple{name: "7"}, ExpansionTuple{name: "8"}, ExpansionTuple{name: "9"}},
}}
var CgiGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":   {ExpansionTuple{name: "<string>"}},
	"<string>":  {ExpansionTuple{name: "<letter>"}, ExpansionTuple{name: "<letter><string>"}},
	"<letter>":  {ExpansionTuple{name: "<plus>"}, ExpansionTuple{name: "<percent>"}, ExpansionTuple{name: "<other>"}},
	"<plus>":    {ExpansionTuple{name: "+"}},
	"<percent>": {ExpansionTuple{name: "%<hexdigit><hexdigit>"}},
	"<hexdigit>": {ExpansionTuple{name: "0"}, ExpansionTuple{name: "1"}, ExpansionTuple{name: "2"}, ExpansionTuple{name: "3"}, ExpansionTuple{name: "4"}, ExpansionTuple{name: "5"}, ExpansionTuple{name: "6"}, ExpansionTuple{name: "7"},
		ExpansionTuple{name: "8"}, ExpansionTuple{name: "9"}, ExpansionTuple{name: "a"}, ExpansionTuple{name: "b"}, ExpansionTuple{name: "c"}, ExpansionTuple{name: "d"}, ExpansionTuple{name: "e"}, ExpansionTuple{name: "f"}},
	"<other>": {ExpansionTuple{name: "0"}, ExpansionTuple{name: "1"}, ExpansionTuple{name: "2"}, ExpansionTuple{name: "3"}, ExpansionTuple{name: "4"}, ExpansionTuple{name: "5"}, ExpansionTuple{name: "a"}, ExpansionTuple{name: "b"}, ExpansionTuple{name: "c"}, ExpansionTuple{name: "d"}, ExpansionTuple{name: "e"}, ExpansionTuple{name: "-"}, ExpansionTuple{name: "_"}},
}}
var URLGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":     {ExpansionTuple{name: "<url>"}},
	"<url>":       {ExpansionTuple{name: "<scheme>://<authority><path><query>"}},
	"<scheme>":    {ExpansionTuple{name: "http"}, ExpansionTuple{name: "https"}, ExpansionTuple{name: "ftp"}, ExpansionTuple{name: "ftps"}},
	"<authority>": {ExpansionTuple{name: "<host>"}, ExpansionTuple{name: "<host>:<port>"}, ExpansionTuple{name: "<userinfo>@<host>"}, ExpansionTuple{name: "<userinfo>@<host>:<port>"}},
	"<host>":      {ExpansionTuple{name: "cispa.saarland"}, ExpansionTuple{name: "www.google.com"}, ExpansionTuple{name: "fuzzingbook.com"}},
	"<port>":      {ExpansionTuple{name: "80"}, ExpansionTuple{name: "8080"}, ExpansionTuple{name: "<nat>"}},
	"<nat>":       {ExpansionTuple{name: "<digit>"}, ExpansionTuple{name: "<digit><digit>"}},
	"<digit>":     {ExpansionTuple{name: "0"}, ExpansionTuple{name: "1"}, ExpansionTuple{name: "2"}, ExpansionTuple{name: "3"}, ExpansionTuple{name: "4"}, ExpansionTuple{name: "5"}, ExpansionTuple{name: "6"}, ExpansionTuple{name: "7"}, ExpansionTuple{name: "8"}, ExpansionTuple{name: "9"}},
	"<userinfo>":  {ExpansionTuple{name: "user:password"}},
	"<path>":      {ExpansionTuple{name: ""}, ExpansionTuple{name: "/"}, ExpansionTuple{name: "/<id>"}},
	"<id>":        {ExpansionTuple{name: "abc"}, ExpansionTuple{name: "def"}, ExpansionTuple{name: "x<digit><digit>"}},
	"<query>":     {ExpansionTuple{name: ""}, ExpansionTuple{name: "?<params>"}},
	"<params>":    {ExpansionTuple{name: "<param>"}, ExpansionTuple{name: "<param>&<params>"}},
	"<param>":     {ExpansionTuple{name: "<id>=<id>"}, ExpansionTuple{name: "<id>=<nat>"}},
}}

func generateDigits() []ExpansionTuple {
	var digits []ExpansionTuple
	for i := 0; i < 10; i++ {
		digits = append(digits, ExpansionTuple{name: strconv.Itoa(i)})
	}
	return digits
}

var ExprEBNFGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":   {ExpansionTuple{name: "<expr>"}},
	"<expr>":    {ExpansionTuple{name: "<term> + <expr>"}, ExpansionTuple{name: "<term> - <expr>"}, ExpansionTuple{name: "<term>"}},
	"<term>":    {ExpansionTuple{name: "<factor> * <term>"}, ExpansionTuple{name: "<factor> / <term>"}, ExpansionTuple{name: "<factor>"}},
	"<factor>":  {ExpansionTuple{name: "<sign>?<factor>"}, ExpansionTuple{name: "(<expr>)"}, ExpansionTuple{name: "<integer>(.<integer>)?"}},
	"<sign>":    {ExpansionTuple{name: "+"}, ExpansionTuple{name: "-"}},
	"<integer>": {ExpansionTuple{name: "<digit>+"}},
	"<digit>":   generateDigits(),
}}

// Generating characters from '1' to '9'
func generateOneNine() []ExpansionTuple {
	var oneNine []ExpansionTuple
	for i := 1; i < 10; i++ {
		oneNine = append(oneNine, ExpansionTuple{name: strconv.Itoa(i)})
	}
	return oneNine
}

// Generating characters without quotes
func generateCharactersWithoutQuote() []ExpansionTuple {
	chars := "0123456789" + "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + "!#$%&'()*+,-./:;<=>?@[]^_`{|}~ "
	var characters []ExpansionTuple
	for _, char := range chars {
		characters = append(characters, ExpansionTuple{name: string(char)})
	}
	return characters
}

var JsonEbnfGrammar = Grammar{G: map[string][]ExpansionTuple{
	"<start>":      {ExpansionTuple{name: "<json>"}},
	"<json>":       {ExpansionTuple{name: "<element>"}},
	"<element>":    {ExpansionTuple{name: "<ws><value><ws>"}},
	"<value>":      {ExpansionTuple{name: "<object>"}, ExpansionTuple{name: "<array>"}, ExpansionTuple{name: "<string>"}, ExpansionTuple{name: "<number>"}, ExpansionTuple{name: "true"}, ExpansionTuple{name: "false"}, ExpansionTuple{name: "null"}, ExpansionTuple{name: "'; DROP TABLE STUDENTS"}},
	"<object>":     {ExpansionTuple{name: "{<ws>"}, ExpansionTuple{name: "{<members>"}},
	"<members>":    {ExpansionTuple{name: "<member>(,<members>)*"}},
	"<member>":     {ExpansionTuple{name: "<ws><string><ws>:<element>"}},
	"<array>":      {ExpansionTuple{name: "[<ws>"}, ExpansionTuple{name: "[<elements>"}},
	"<elements>":   {ExpansionTuple{name: "<element>(,<elements>)*"}},
	"<string>":     {ExpansionTuple{name: "\"<characters>\""}},
	"<characters>": {ExpansionTuple{name: "<character>*"}},
	"<character>":  generateCharactersWithoutQuote(),
	"<number>":     {ExpansionTuple{name: "<int><frac><exp>"}},
	"<int>":        {ExpansionTuple{name: "<digit>"}, ExpansionTuple{name: "<onenine><digits>"}, ExpansionTuple{name: "-<digit>"}, ExpansionTuple{name: "-<onenine><digits>"}},
	"<digits>":     {ExpansionTuple{name: "<digit>+"}},
	"<digit>":      append(generateDigits(), ExpansionTuple{name: "<onenine>"}),
	"<onenine>":    generateOneNine(),
	"<frac>":       {ExpansionTuple{name: ""}, ExpansionTuple{name: ".<digits>"}},
	"<exp>":        {ExpansionTuple{name: ""}, ExpansionTuple{name: "E<sign><digits>"}, ExpansionTuple{name: "e<sign><digits>"}},
	"<sign>":       {ExpansionTuple{name: ""}, ExpansionTuple{name: "+"}, ExpansionTuple{name: "-"}},
	"<ws>":         {ExpansionTuple{name: " "}},
}}

func TestNonTerminals(t *testing.T) {
	var tests = []struct {
		input    string
		expected []string
	}{
		{"<term> * <factor>", []string{"<term>", "<factor>"}},
		{"<digit><integer>", []string{"<digit>", "<integer>"}},
		{"1 < 3 > 2", []string{}},
		{"1 <3> 2", []string{"<3>"}},
		{"1 + 2", []string{}},
	}

	for _, test := range tests {
		ans := NonTerminals(test.input)
		if len(ans) != len(test.expected) {

			t.Errorf("NonTerminals(%s) = %v; want %v", test.input, ans, test.expected)
		} else {
			for i, v := range ans {
				if v != test.expected[i] {
					t.Errorf("NonTerminals(%s) = %v; want %v", test.input, ans, test.expected)
					break
				}
			}
		}
	}
}

func TestIsNonTerminals(t *testing.T) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{"<abc>", true},
		{"<symbol-1>", true},
		{"+", false},
	}

	for _, test := range tests {
		ans := IsNonTerminals(test.input)
		if ans != test.expected {
			t.Errorf("NonTerminals(%s) = %v; want %v", test.input, ans, test.expected)
		}
	}
}

func TestExtendGrammar(t *testing.T) {
	var tests = []struct {
		key      string
		expected interface{}
	}{
		{"<identifier>", []string{"<idchar>", "<identifier><idchar>"}},
		{"<idchar>", []string{"a", "b", "c", "d"}},
	}

	SimpleNonterminalGrammar.Extend(Grammar{G: map[string][]ExpansionTuple{
		"<identifier>": []ExpansionTuple{{name: "<idchar>"}, {name: "<identifier><idchar>"}},
		"<idchar>":     []ExpansionTuple{{name: "a"}, {name: "b"}, {name: "c"}, {name: "d"}},
	}})

	for _, test := range tests {
		for i, v := range SimpleNonterminalGrammar.G[test.key] {
			if v.GetName() != test.expected.([]string)[i] {
				t.Errorf("Extend() for key %s = %v; want %v", test.key, SimpleNonterminalGrammar.G[test.key], test.expected.([]string)[i])
				break
			}
		}
	}
}

func TestSRange(t *testing.T) {

	var tests = []struct {
		input    string
		expected []string
	}{
		{"abc", []string{"a", "b", "c"}},
		{"a", []string{"a"}},
		{"", []string{}},
	}

	for _, test := range tests {
		ans := SRange(test.input)
		for i, v := range ans {
			if v != test.expected[i] {
				t.Errorf("SRange(%s) = %v; want %v", test.input, ans, test.expected)
				break
			}
		}
	}
}

func TestCRange(t *testing.T) {
	var tests = []struct {
		start    string
		end      string
		expected []string
	}{
		{"a", "c", []string{"a", "b", "c"}},
		{"a", "a", []string{"a"}},
		{"a", "b", []string{"a", "b"}},
		{"0", "9", []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}},
	}

	for _, test := range tests {
		ans := CRange(test.start, test.end)
		for i, v := range ans {
			if v != test.expected[i] {
				t.Errorf("CRange(%s, %s) = %v; want %v", test.start, test.end, ans, test.expected)
				break
			}
		}
	}
}

func TestIsValidGrammar(t *testing.T) {

	if ExprGrammar.IsValidGrammar("", nil) == false {
		t.Errorf("IsValidGrammar() = %v; want %v", ExprGrammar.IsValidGrammar("", nil), true)
	}

	if CgiGrammar.IsValidGrammar("", nil) == false {
		t.Errorf("IsValidGrammar() = %v; want %v", CgiGrammar.IsValidGrammar("", nil), true)
	}
	if URLGrammar.IsValidGrammar("", nil) == false {
		t.Errorf("IsValidGrammar() = %v; want %v", URLGrammar.IsValidGrammar("", nil), true)
	}
	if ExprEBNFGrammar.IsValidGrammar("", nil) == false {
		t.Errorf("IsValidGrammar() = %v; want %v", ExprEBNFGrammar.IsValidGrammar("", nil), true)
	}
	if JsonEbnfGrammar.IsValidGrammar("", nil) == false {
		t.Errorf("IsValidGrammar() = %v; want %v", JsonEbnfGrammar.IsValidGrammar("", nil), true)
	}
}

func TestGrammar_Visualize(t *testing.T) {
	ExprGrammar.Visualize("expr.png")
	CgiGrammar.Visualize("cgi.png")
	URLGrammar.Visualize("url.png")
	ExprEBNFGrammar.Visualize("expr_ebnf.png")
	JsonEbnfGrammar.Visualize("json_ebnf.png")

}

func TestExpansionTuple_Expand(t *testing.T) {
	e := ExpansionTuple{name: "<digit>*<integer>"}
	res := e.Expand()
	if len(res) != 3 {
		t.Errorf("Expand() = %v; want %v", res, 3)
	}
	if res[0].GetName() != "<digit>" {
		t.Errorf("Expand() = %v; want %v", res[0].GetName(), "<digit>")
	}
	if res[1].GetName() != "*" {
		t.Errorf("Expand() = %v; want %v", res[1].GetName(), "*")
	}
	if res[2].GetName() != "<integer>" {
		t.Errorf("Expand() = %v; want %v", res[2].GetName(), "<integer>")
	}
}

func TestExpansionTuple_Expand2(t *testing.T) {
	e := ExpansionTuple{name: ""}
	res := e.Expand()
	if len(res) != 1 {
		t.Errorf("Expand() = %v; want %v", res, 3)
	}
	if res[0].GetName() != "" {
		t.Errorf("Expand() = %v; want %v", res[0].GetName(), "<digit>")
	}
}

func TestGrammar_SymbolCost(t *testing.T) {
	testCases := []struct {
		name   string
		symbol string
		cost   float64
	}{
		{"Test case 1", "<digit>", 1},
		{"Test case 2", "<expr>", 5},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			m := make(map[string]struct{})
			cost := ExprGrammar.SymbolCost(ExpansionTuple{name: tt.symbol}, m)
			if cost != tt.cost {
				t.Errorf("SymbolCost() = %v; want %v", cost, tt.cost)
			}
		})
	}
}
