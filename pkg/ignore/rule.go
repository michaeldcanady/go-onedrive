package ignore

// Rule represents a compiled ignore rule.
type Rule struct {
	Pattern  string // The compiled pattern string (glob or regex representation)
	Original string // The original rule text
	Negate   bool   // True if the rule starts with ! (re-include)
	DirOnly  bool   // True if the rule ends with /
	Line     int    // Source line number
}
