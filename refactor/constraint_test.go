package refactor

var invalidConstraint = &Constraint{IsUnique: NewBoolAttr("not bool")}
var invalidConstraintError = invalidConstraint.Validate()
