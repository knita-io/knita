package label

import executorv1 "github.com/knita-io/knita/api/executor/v1"

// MatchSelector returns true if the given labels satisfy both matchLabels and
// matchExpressions in the selector.
func MatchSelector(labels map[string]string, sel *executorv1.LabelSelector) bool {
	if sel == nil {
		return true
	}
	// 1) Simple map‐based matching (AND).
	for key, val := range sel.MatchLabels {
		v, ok := labels[key]
		if !ok {
			return false
		}
		if v != val {
			return false
		}
	}
	// 2) Expression‐based matching (each requirement is AND’d).
	for _, req := range sel.MatchExpressions {
		switch req.Operator {
		case executorv1.LabelSelectorRequirement_IN:
			// key must exist, and its value equal req.Values
			v, ok := labels[req.Key]
			if !ok {
				return false
			}
			found := false
			for _, allowed := range req.Values {
				if v == allowed {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case executorv1.LabelSelectorRequirement_NOT_IN:
			// key must exist, and its value equal req.Values
			v, ok := labels[req.Key]
			if !ok {
				return false
			}
			for _, banned := range req.Values {
				if v == banned {
					return false
				}
			}
		case executorv1.LabelSelectorRequirement_EXISTS:
			// key must exist (values is ignored)
			if _, ok := labels[req.Key]; !ok {
				return false
			}
		case executorv1.LabelSelectorRequirement_DOES_NOT_EXIST:
			// key must not exist (values is ignored)
			if _, ok := labels[req.Key]; ok {
				return false
			}
		default:
			return false
		}
	}
	return true
}
