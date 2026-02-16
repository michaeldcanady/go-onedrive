package main

import (
	"fmt"

	stduritemplate "github.com/std-uritemplate/std-uritemplate/go/v2"
)

const (
	pathTest = "/me/drive/root:/{item_path}:/content"
)

func main() {
	subs := make(stduritemplate.Substitutions)

	subs["item_path"] = "this/is/a/very/log/and/weird/path.js"

	fmt.Println(stduritemplate.Expand(pathTest, subs))
}
