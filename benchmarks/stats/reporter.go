package stats

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var printStats = flag.Bool("printStats", true, "Print stats to console")

// IncCounter increments a counter.
func IncCounter(name string, tags map[string]string, value int64) {
	name = addTagsToName(name, tags)
	if *printStats {
		fmt.Printf("IncCounter: %v = %v\n", name, value)
	}
}

func addTagsToName(name string, tags map[string]string) string {
	// The format we want is: host.endpoint.os.browser
	// if there's no host tag, then we don't use it.
	var keyOrder []string
	if _, ok := tags["host"]; ok {
		keyOrder = append(keyOrder, "host")
	}
	keyOrder = append(keyOrder, "endpoint")

	parts := []string{name}
	for _, k := range keyOrder {
		v, ok := tags[k]
		if !ok || v == "" {
			parts = append(parts, "no-"+k)
			continue
		}
		parts = append(parts, clean(v))
	}

	return strings.Join(parts, ".")
}

var specialChars = regexp.MustCompile(`[{}/\\:.]`)

// clean takes a string that may contain special characters, and replaces these
// characters with a '-'.
func clean(value string) string {
	return specialChars.ReplaceAllString(value, "-")
}
