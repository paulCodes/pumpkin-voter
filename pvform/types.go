package pvform

type ThreeCol struct {
	GroupName  string
	LabelWidth int
	InputWidth int
	Fields     []FormField
}

type FormField struct {
	Title          string
	Type           string
	ClarifyingText string
	IsRequired     bool
	Initial        string
}