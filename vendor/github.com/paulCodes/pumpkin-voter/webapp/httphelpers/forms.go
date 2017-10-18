package httphelpers

import "github.com/bluele/gforms"

func CustomGformsNewFields(fields ...gforms.Field) *gforms.Fields {
	// TODO Move this into initCustomGformStyling.
	customTextClass := gforms.TextInputWidget(
		map[string]string{
			"class": "form-control uk-input",
		},
	)

	for _, field := range fields {
		switch field.(type) {
		case *gforms.TextField, *gforms.IntegerField:
			if field.GetWidget() == nil {
				field.SetWidget(customTextClass)
			}
		case *gforms.DateTimeField:
			if field.GetWidget() == nil {
				field.SetWidget(customTextClass)
			}
		}
	}
	return gforms.NewFields(fields...)
}
