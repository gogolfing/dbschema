package strval

import "fmt"

const (
	True  = "true"
	False = "false"
)

func Value(v *string, def string) string {
	if v == nil {
		return def
	}
	return *v
}

func Bool(v *string, def bool) bool {
	if v == nil {
		return def
	}
	if *v == True {
		return true
	}
	return false
}

func ValidateBool(value *string) error {
	if value == nil {
		return nil
	}
	if *value != True && *value != False {
		return fmt.Errorf("must be %q, %q, or not present", True, False)
	}
	return nil
}
