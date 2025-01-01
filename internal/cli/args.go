package cli

import (
	"flag"
	"fmt"
)

func Usage(flags *flag.FlagSet) string {
	longest := 0
	longestFn := func(flg *flag.Flag) {
		length := len(flg.Name)
		name, _ := flag.UnquoteUsage(flg)
		if len(name) > 0 {
			length += len(name) + 1
		}
		if length > longest {
			longest = length
		}
	}
	flags.VisitAll(longestFn)
	spacesStr := ""
	for i := 0; i < longest; i++ {
		spacesStr += " "
	}
	usage := ""
	usageFn := func(flg *flag.Flag) {
		name, rawUsage := flag.UnquoteUsage(flg)
		flgName := flg.Name
		if len(name) > 0 {
			flgName += " " + name
		}
		offsetLen := longest - len(flgName)
		offset := ""
		for i := 0; i < offsetLen; i++ {
			offset += " "
		}
		usage += fmt.Sprintf("  -%s%s    ", flgName, offset) + rawUsage + "\n"
	}
	flags.VisitAll(usageFn)
	return usage
}
