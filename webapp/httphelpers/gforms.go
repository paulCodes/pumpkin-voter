package httphelpers

import (
	"github.com/bluele/gforms"
	"reflect"
	"strconv"
	"text/template"
	"time"
)

func initCustomGformStyling() {
	var defaultTemplates = `
{{define "TextTypeField"}}<input type="text" name="{{.Field.GetName | html}}" value="{{.Value | html}}"></input>{{end}}
{{define "BooleanTypeField"}}<input type="checkbox" name="{{.Field.GetName | html}}"{{if .Checked}} checked{{end}} class="uk-checkbox">{{end}}
{{define "SimpleWidget"}}<input type="{{.Type | html}}" name="{{.Field.GetName | html}}" value="{{.Value | html}}"{{range $attr, $val := .Attrs}} {{$attr | html}}="{{$val | html}}"{{end}}></input>{{end}}
{{define "SelectWidget"}}<select {{if .Multiple }}multiple {{end}}name="{{.Field.GetName | html}}"{{range $attr, $val := .Attrs}}{{$attr | html}}="{{$val | html}}"{{end}}>
{{range $idx, $val := .Options}}<option value="{{$val.Value | html}}"{{if $val.Selected }} selected{{end}}{{if $val.Disabled}} disabled{{end}}>{{$val.Label | html}}</option>
{{end}}</select>{{end}}
{{define "RadioWidget"}}{{$name := .Field.GetName}}{{range $idx, $val := .Options}}<input type="radio" name="{{$name | html}}" value="{{$val.Value | html}}"{{if or $val.Checked (eq $.Field.GetV.RawStr $val.Value) }} checked{{end}}{{if $val.Disabled}} disabled{{end}}>{{$val.Label | html}}
{{end}}{{end}}
{{define "CheckboxMultipleWidget"}}{{$name := .Field.GetName}}{{range $idx, $val := .Options}}<input type="checkbox" name="{{$name | html}}" value="{{$val.Value | html}}"{{if $val.Checked}} checked{{end}}{{if $val.Disabled}} disabled{{end}}>{{$val.Label | html}}
{{end}}{{end}}
{{define "FileTypeField"}}<input type="file" name="{{.Field.GetName | html}}"></input>{{end}}
{{define "TextAreaWidget"}}<textarea name="{{.Field.GetName | html}}" {{range $attr, $val := .Attrs}} {{$attr | html}}="{{$val | html}}"{{end}}>{{.Value | html}}</textarea>{{end}}
`

	// all templates of Field and Widget
	var replaceTemplate *template.Template

	var err error
	replaceTemplate, err = template.New("gforms").Parse(defaultTemplates)
	if err != nil {
		panic(err)
	}
	gforms.Template = replaceTemplate
}

// Create form instance from model data
func PopulateFormInstanceFromModel(fi *gforms.FormInstance, model interface{}) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		panic("Argument should be specified pointer type.")
	}

	mType := reflect.TypeOf(model).Elem()
	mValue := reflect.ValueOf(model).Elem()

	for i := 0; i < mValue.NumField(); i++ {
		typeField := mType.Field(i)
		tag := typeField.Tag.Get("gforms")
		if tag == "" {
			tag = typeField.Name
		} else if tag == "-" {
			continue
		}

		field, ok := fi.GetField(tag)
		if ok {
			valueField := mValue.Field(i)
			workField := valueField
			// Follow the indirection.
			for workField.Kind() == reflect.Ptr {
				workField = workField.Elem()
			}

			switch workField.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				field.SetInitial(strconv.FormatInt(workField.Int(), 10))
			case reflect.Float32, reflect.Float64:
				field.SetInitial(strconv.FormatFloat(workField.Float(), 'G', -1, 64))
			case reflect.String:
				field.SetInitial(workField.String())
			case reflect.Bool:
				if workField.Bool() {
					field.SetInitial("true")
				}
			case reflect.Struct:
				switch workField.Type().String() {
				case "time.Time":
					field.SetInitial(workField.Interface().(time.Time).Format(field.(*gforms.DateTimeFieldInstance).Format))
				}
			}
		}
	}
}
