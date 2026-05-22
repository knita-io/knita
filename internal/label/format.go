package label

import (
	"fmt"
	"strings"

	executorv1 "github.com/knita-io/knita/api/executor/v1"
)

// FormatSelector returns a human-readable description of sel.
// e.g. "env=prod, tier in (frontend, backend), !deprecated"
func FormatSelector(sel *executorv1.LabelSelector) string {
	if sel == nil {
		return ""
	}
	var parts []string
	// 1) exact matches: key=value
	for k, v := range sel.MatchLabels {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	// 2) expressions
	for _, req := range sel.MatchExpressions {
		switch req.Operator {
		case executorv1.LabelSelectorRequirement_IN:
			parts = append(parts,
				fmt.Sprintf("%s in (%s)",
					req.Key,
					strings.Join(req.Values, ", "),
				),
			)
		case executorv1.LabelSelectorRequirement_NOT_IN:
			parts = append(parts,
				fmt.Sprintf("%s notin (%s)",
					req.Key,
					strings.Join(req.Values, ", "),
				),
			)
		case executorv1.LabelSelectorRequirement_EXISTS:
			parts = append(parts, req.Key)
		case executorv1.LabelSelectorRequirement_DOES_NOT_EXIST:
			parts = append(parts, "!"+req.Key)
		default:
			// fallback for unknown/unspecified operators
			parts = append(parts,
				fmt.Sprintf("%s<?>", req.Key),
			)
		}
	}
	if len(parts) == 0 {
		return "<none>"
	}
	return strings.Join(parts, ", ")
}
