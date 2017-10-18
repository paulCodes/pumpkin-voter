package pvtemplate

import (
	"bytes"
	"encoding/json"
	"github.com/flosch/pongo2"
	"regexp"
	"strings"
	"time"
)

func RegisterFilters() {
	pongo2.RegisterFilter("getValueByKey", FilterGetValueByKey)
	pongo2.RegisterFilter("makeTimeFromEpoch", FilterMakeTimeFromEpoch)
	pongo2.RegisterFilter("displayValue", FilterValidationDisplayValue)
	pongo2.RegisterFilter("json", FilterOutputJson)
}

/**
*	pongo2 filter that will get the value by key from a map[string]*pongo2.Value
 */
func FilterGetValueByKey(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	m := param.Interface().(map[string]string)
	return pongo2.AsValue(m[in.String()]), nil
}

func FilterMakeTimeFromEpoch(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	m := int64(in.Integer())
	t := time.Unix(m, 0)
	return pongo2.AsValue(t), nil
}

func FilterValidationDisplayValue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	r := strings.NewReplacer("_", " ", "[", " ", "]", "")
	t := r.Replace(in.String())
	return pongo2.AsValue(t), nil
}

func FilterOutputJson(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	out, _ := json.Marshal(in.Interface())
	return pongo2.AsValue(string(out)), nil
}

type tagTrimNode struct {
	wrapper *pongo2.NodeWrapper
}

var tagTrimRegexp = regexp.MustCompile(`^[\t\n\v\f\r ]+(.*)[\t\n\v\f\r ]+$`)

func (node *tagTrimNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	b := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB

	err := node.wrapper.Execute(ctx, b)
	if err != nil {
		return err
	}

	s := b.String()
	// Repeat this recursively
	changed := true
	for changed {
		s2 := tagTrimRegexp.ReplaceAllString(s, "$1")
		changed = s != s2
		s = s2
	}

	writer.WriteString(s)

	return nil
}

func tagTrimParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	trimNode := &tagTrimNode{}

	wrapper, _, err := doc.WrapUntilTag("endtrim")
	if err != nil {
		return nil, err
	}
	trimNode.wrapper = wrapper

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed trim-tag arguments.", nil)
	}

	return trimNode, nil
}

func init() {
	pongo2.RegisterTag("trim", tagTrimParser)
}
