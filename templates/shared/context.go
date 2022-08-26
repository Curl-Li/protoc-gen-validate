package shared

import (
	"bytes"
	"fmt"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/proto"

	"github.com/curl-li/protoc-gen-validate/validate"
)

type RuleContext struct {
	Field        pgs.Field
	Rules        proto.Message
	MessageRules *validate.MessageRules
	DefineErr    bool
	ErrBase      *validate.ErrorBase
	ErrIndex     int

	Typ        string
	WrapperTyp string

	OnKey            bool
	Index            string
	AccessorOverride string
}

func rulesContext(msg pgs.Message, f pgs.Field) (out RuleContext, err error) {
	out.Field = f
	out.DefineErr = true

	if msg != nil {
		if _, err = msg.Extension(validate.E_ErrorBase, &out.ErrBase); err != nil {
			return
		}
	}
	var rules validate.FieldRules
	if _, err = f.Extension(validate.E_Rules, &rules); err != nil {
		return
	}

	var wrapped bool
	if out.Typ, out.Rules, out.MessageRules, wrapped = resolveRules(f.Type(), &rules); wrapped {
		out.WrapperTyp = out.Typ
		out.Typ = "wrapper"
	}

	if out.Typ == "error" {
		err = fmt.Errorf("unknown rule type (%T)", rules.Type)
	}

	return
}

func (ctx RuleContext) Key(name, idx string) (out RuleContext, err error) {
	rules, ok := ctx.Rules.(*validate.MapRules)
	if !ok {
		err = fmt.Errorf("cannot get Key RuleContext from %T", ctx.Field)
		return
	}

	out.Field = ctx.Field
	out.AccessorOverride = name
	out.Index = idx

	var rule *validate.FieldRules
	for _, r := range rules.GetRules() {
		if r.Keys != nil {
			rule = r.GetKeys()
		}
	}

	out.Typ, out.Rules, out.MessageRules, _ = resolveRules(ctx.Field.Type().Key(), rule)

	if out.Typ == "error" {
		err = fmt.Errorf("unknown rule type (%T)", rules)
	}

	return
}

func (ctx RuleContext) KeyWithErrIndex(name, idx string, errIndex int) (RuleContext, error) {
	out, err := ctx.Key(name, idx)
	out.ErrIndex = errIndex
	return out, err
}

func (ctx RuleContext) Elem(name, idx string) (out RuleContext, err error) {
	out.Field = ctx.Field
	out.AccessorOverride = name
	out.Index = idx
	out.ErrBase = ctx.ErrBase

	var rules *validate.FieldRules
	switch r := ctx.Rules.(type) {
	case *validate.MapRules:
		for _, rr := range r.GetRules() {
			if rr.Values != nil {
				rules = rr.GetValues()
			}
		}
	case *validate.RepeatedRules:
		for _, rr := range r.GetRules() {
			if rr.Items != nil {
				rules = rr.GetItems()
			}
		}
	default:
		err = fmt.Errorf("cannot get Elem RuleContext from %T", ctx.Field)
		return
	}

	var wrapped bool
	if out.Typ, out.Rules, out.MessageRules, wrapped = resolveRules(ctx.Field.Type().Element(), rules); wrapped {
		out.WrapperTyp = out.Typ
		out.Typ = "wrapper"
	}

	if out.Typ == "error" {
		err = fmt.Errorf("unknown rule type (%T)", rules)
	}

	return
}

func (ctx RuleContext) ElemWithErrIndex(name, idx string, errIndex int) (RuleContext, error) {
	out, err := ctx.Elem(name, idx)
	out.ErrIndex = errIndex
	return out, err
}

func (ctx RuleContext) Unwrap(name string) (out RuleContext, err error) {
	if ctx.Typ != "wrapper" {
		err = fmt.Errorf("cannot unwrap non-wrapper type %q", ctx.Typ)
		return
	}

	return RuleContext{
		Field:            ctx.Field,
		Rules:            ctx.Rules,
		MessageRules:     ctx.MessageRules,
		Typ:              ctx.WrapperTyp,
		AccessorOverride: name,
		ErrBase:          ctx.ErrBase,
		DefineErr:        true,
	}, nil
}

func Render(tpl *template.Template) func(ctx RuleContext) (string, error) {
	return func(ctx RuleContext) (string, error) {
		var b bytes.Buffer
		err := tpl.ExecuteTemplate(&b, ctx.Typ, ctx)
		return b.String(), err
	}
}

func resolveRules(typ interface{ IsEmbed() bool }, rules *validate.FieldRules) (ruleType string, rule proto.Message, messageRule *validate.MessageRules, wrapped bool) {
	switch r := rules.GetType().(type) {
	case *validate.FieldRules_Float:
		ruleType, rule, wrapped = "float", r.Float, typ.IsEmbed()
	case *validate.FieldRules_Double:
		ruleType, rule, wrapped = "double", r.Double, typ.IsEmbed()
	case *validate.FieldRules_Int32:
		ruleType, rule, wrapped = "int32", r.Int32, typ.IsEmbed()
	case *validate.FieldRules_Int64:
		ruleType, rule, wrapped = "int64", r.Int64, typ.IsEmbed()
	case *validate.FieldRules_Uint32:
		ruleType, rule, wrapped = "uint32", r.Uint32, typ.IsEmbed()
	case *validate.FieldRules_Uint64:
		ruleType, rule, wrapped = "uint64", r.Uint64, typ.IsEmbed()
	case *validate.FieldRules_Sint32:
		ruleType, rule, wrapped = "sint32", r.Sint32, false
	case *validate.FieldRules_Sint64:
		ruleType, rule, wrapped = "sint64", r.Sint64, false
	case *validate.FieldRules_Fixed32:
		ruleType, rule, wrapped = "fixed32", r.Fixed32, false
	case *validate.FieldRules_Fixed64:
		ruleType, rule, wrapped = "fixed64", r.Fixed64, false
	case *validate.FieldRules_Sfixed32:
		ruleType, rule, wrapped = "sfixed32", r.Sfixed32, false
	case *validate.FieldRules_Sfixed64:
		ruleType, rule, wrapped = "sfixed64", r.Sfixed64, false
	case *validate.FieldRules_Bool:
		ruleType, rule, wrapped = "bool", r.Bool, typ.IsEmbed()
	case *validate.FieldRules_String_:
		ruleType, rule, wrapped = "string", r.String_, typ.IsEmbed()
	case *validate.FieldRules_Bytes:
		ruleType, rule, wrapped = "bytes", r.Bytes, typ.IsEmbed()
	case *validate.FieldRules_Enum:
		ruleType, rule, wrapped = "enum", r.Enum, false
	case *validate.FieldRules_Repeated:
		ruleType, rule, wrapped = "repeated", r.Repeated, false
	case *validate.FieldRules_Map:
		ruleType, rule, wrapped = "map", r.Map, false
	case *validate.FieldRules_Any:
		ruleType, rule, wrapped = "any", r.Any, false
	case *validate.FieldRules_Duration:
		ruleType, rule, wrapped = "duration", r.Duration, false
	case *validate.FieldRules_Timestamp:
		ruleType, rule, wrapped = "timestamp", r.Timestamp, false
	case nil:
		if ft, ok := typ.(pgs.FieldType); ok && ft.IsRepeated() {
			return "repeated", &validate.RepeatedRules{}, rules.Message, false
		} else if ok && ft.IsMap() && ft.Element().IsEmbed() {
			return "map", &validate.MapRules{}, rules.Message, false
		} else if typ.IsEmbed() {
			return "message", rules.GetMessage(), rules.GetMessage(), false
		}
		return "none", nil, nil, false
	default:
		ruleType, rule, wrapped = "error", nil, false
	}

	return ruleType, rule, rules.Message, wrapped
}
