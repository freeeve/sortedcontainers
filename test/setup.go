// use test.sh

package main

import (
	"os"
	"strings"

	"github.com/clipperhouse/typewriter"
	_ "github.com/wfreeman/typewriters/sortedcontainer"
)

func main() {
	// don't let bad test or gen'd files get us stuck
	filter := func(f os.FileInfo) bool {
		return !strings.HasSuffix(f.Name(), "_test.go") && !strings.HasSuffix(f.Name(), "_container.go")
	}

	a, err := typewriter.NewAppFiltered("+test", filter)
	if err != nil {
		panic(err)
	}
	a.WriteAll()
}
