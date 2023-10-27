package output

import "strings"

func JoinMetricName(name ...string) string {
	return strings.Join(name, "_")
}
