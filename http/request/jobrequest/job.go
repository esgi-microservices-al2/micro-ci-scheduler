package jobrequest

import "github.com/System-Glitch/goyave/v2/validation"

var (
	Store validation.RuleSet = validation.RuleSet{
		"cronexpression": {"required", "string", "max:255"},
		"name":           {"required", "string", "max:255"},
		"idproject":      {"required", "integer", "min:1"},
	}

)
