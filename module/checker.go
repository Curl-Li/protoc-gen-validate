package module

import (
	"reflect"
	"regexp"
	"time"
	"unicode/utf8"

	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/curl-li/protoc-gen-validate/validate"
)

var unknown = ""

var httpHeaderName = "^:?[0-9a-zA-Z!#$%&'*+-.^_|~\x60]+$"

var httpHeaderValue = "^[^\u0000-\u0008\u000A-\u001F\u007F]*$"

var headerString = "^[^\u0000\u000A\u000D]*$" // For non-strict validation.

// Map from well known regex to regex pattern.
var regex_map = map[string]*string{
	"UNKNOWN":           &unknown,
	"HTTP_HEADER_NAME":  &httpHeaderName,
	"HTTP_HEADER_VALUE": &httpHeaderValue,
	"HEADER_STRING":     &headerString,
}

type FieldType interface {
	ProtoType() pgs.ProtoType
	Embed() pgs.Message
}

type Repeatable interface {
	IsRepeated() bool
}

func (m *Module) CheckRules(msg pgs.Message) {
	m.Push("msg: " + msg.Name().String())
	defer m.Pop()

	var disabled bool
	_, err := msg.Extension(validate.E_Disabled, &disabled)
	m.CheckErr(err, "unable to read validation extension from message")

	if disabled {
		m.Debug("validation disabled, skipping checks")
		return
	}

	for _, f := range msg.Fields() {
		m.Push(f.Name().String())

		var rules validate.FieldRules
		_, err = f.Extension(validate.E_Rules, &rules)
		m.CheckErr(err, "unable to read validation rules from field")

		if rules.GetMessage() != nil {
			m.MustType(f.Type(), pgs.MessageT, pgs.UnknownWKT)
			m.CheckMessage(f, &rules, false)
		}

		m.CheckFieldRules(f.Type(), &rules, false)

		m.Pop()
	}
}

func (m *Module) CheckFieldRules(typ FieldType, rules *validate.FieldRules, inject bool) {
	if rules == nil {
		return
	}

	switch r := rules.Type.(type) {
	case *validate.FieldRules_Float:
		m.MustType(typ, pgs.FloatT, pgs.FloatValueWKT)
		m.CheckFloat(r.Float, inject)
	case *validate.FieldRules_Double:
		m.MustType(typ, pgs.DoubleT, pgs.DoubleValueWKT)
		m.CheckDouble(r.Double, inject)
	case *validate.FieldRules_Int32:
		m.MustType(typ, pgs.Int32T, pgs.Int32ValueWKT)
		m.CheckInt32(r.Int32, inject)
	case *validate.FieldRules_Int64:
		m.MustType(typ, pgs.Int64T, pgs.Int64ValueWKT)
		m.CheckInt64(r.Int64, inject)
	case *validate.FieldRules_Uint32:
		m.MustType(typ, pgs.UInt32T, pgs.UInt32ValueWKT)
		m.CheckUInt32(r.Uint32, inject)
	case *validate.FieldRules_Uint64:
		m.MustType(typ, pgs.UInt64T, pgs.UInt64ValueWKT)
		m.CheckUInt64(r.Uint64, inject)
	case *validate.FieldRules_Sint32:
		m.MustType(typ, pgs.SInt32, pgs.UnknownWKT)
		m.CheckSInt32(r.Sint32, inject)
	case *validate.FieldRules_Sint64:
		m.MustType(typ, pgs.SInt64, pgs.UnknownWKT)
		m.CheckSInt64(r.Sint64, inject)
	case *validate.FieldRules_Fixed32:
		m.MustType(typ, pgs.Fixed32T, pgs.UnknownWKT)
		m.CheckFixed32(r.Fixed32, inject)
	case *validate.FieldRules_Fixed64:
		m.MustType(typ, pgs.Fixed64T, pgs.UnknownWKT)
		m.CheckFixed64(r.Fixed64, inject)
	case *validate.FieldRules_Sfixed32:
		m.MustType(typ, pgs.SFixed32, pgs.UnknownWKT)
		m.CheckSFixed32(r.Sfixed32, inject)
	case *validate.FieldRules_Sfixed64:
		m.MustType(typ, pgs.SFixed64, pgs.UnknownWKT)
		m.CheckSFixed64(r.Sfixed64, inject)
	case *validate.FieldRules_Bool:
		m.MustType(typ, pgs.BoolT, pgs.BoolValueWKT)
	case *validate.FieldRules_String_:
		m.MustType(typ, pgs.StringT, pgs.StringValueWKT)
		m.CheckStringRules(r.String_, inject)
	case *validate.FieldRules_Bytes:
		m.MustType(typ, pgs.BytesT, pgs.BytesValueWKT)
		m.CheckBytesRules(r.Bytes, inject)
	case *validate.FieldRules_Enum:
		m.MustType(typ, pgs.EnumT, pgs.UnknownWKT)
		m.CheckEnum(typ, r.Enum, inject)
	case *validate.FieldRules_Repeated:
		m.CheckRepeatedRules(typ, r.Repeated, inject)
	case *validate.FieldRules_Map:
		m.CheckMapRules(typ, r.Map, inject)
	case *validate.FieldRules_Any:
		m.CheckAnyRules(typ, r.Any, inject)
	case *validate.FieldRules_Duration:
		m.CheckDurationRules(typ, r.Duration, inject)
	case *validate.FieldRules_Timestamp:
		m.CheckTimestampRules(typ, r.Timestamp, inject)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", rules.Type)
	}
}

func (m *Module) MustType(typ FieldType, pt pgs.ProtoType, wrapper pgs.WellKnownType) {
	if emb := typ.Embed(); emb != nil && emb.IsWellKnown() && emb.WellKnownType() == wrapper {
		m.MustType(emb.Fields()[0].Type(), pt, pgs.UnknownWKT)
		return
	}

	if typ, ok := typ.(Repeatable); ok {
		m.Assert(!typ.IsRepeated(),
			"repeated rule should be used for repeated fields")
	}

	m.Assert(typ.ProtoType() == pt,
		" expected rules for ",
		typ.ProtoType().Proto(),
		" but got ",
		pt.Proto(),
	)
}

func (m *Module) CheckFloat(rules *validate.FloatRules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *float32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckDouble(rules *validate.DoubleRules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *float64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckInt32(rules *validate.Int32Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckInt64(rules *validate.Int64Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckUInt32(rules *validate.UInt32Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *uint32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckUInt64(rules *validate.UInt64Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *uint64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckSInt32(rules *validate.SInt32Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckSInt64(rules *validate.SInt64Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckFixed32(rules *validate.Fixed32Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *uint32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckFixed64(rules *validate.Fixed64Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *uint64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckSFixed32(rules *validate.SFixed32Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int32
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckSFixed64(rules *validate.SFixed64Rules, inject bool) {
	var (
		in, notIn            bool
		ct, lt, lte, gt, gte *int64
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			ct = r.Const
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lt` and `const` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			lt = r.Lt
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `lte` and `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			lte = r.Lte
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gt` and `const` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			gt = r.Gt
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `gte` and `const` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			gte = gte
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkNums(len(r.In), len(r.NotIn), r.Const, r.Lt, r.Lte, r.Gt, r.Gte)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckStringRules(rules *validate.StringRules, inject bool) {
	var (
		pattern, wk                                          bool
		length, minLen, maxLen, lenBs, minBs, maxBs          *uint64
		in, notIn, ct, prefix, suffix, contains, notContains bool
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(!ct, "cannot have multi `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `const` and `len` rules on the same field")
			m.Assert(minLen == nil, "cannot have both `const` and `min_len` rules on the same field")
			m.Assert(maxLen == nil, "cannot have both `const` and `max_len` rules on the same field")
			m.Assert(lenBs == nil, "cannot have both `const` and `len_bytes` rules on the same field")
			m.Assert(minBs == nil, "cannot have both `const` and `min_bytes` rules on the same field")
			m.Assert(maxBs == nil, "cannot have both `const` and `max_bytes` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(!prefix, "cannot have both `const` and `prefix` rules on the same field")
			m.Assert(!suffix, "cannot have both `const` and `suffix` rules on the same field")
			m.Assert(!contains, "cannot have both `const` and `contains` rules on the same field")
			m.Assert(!notContains, "cannot have both `const` and `not_contains` rules on the same field")
			m.Assert(!pattern, "cannot have both `const` and `pattern` rules on the same field")
			m.Assert(!wk, "cannot have both `const` and `well_known` rules on the same field")
		}
		if r.Len != nil {
			m.Assert(length == nil, "cannot have multi `len` rules on the same field")
			m.Assert(!ct, "cannot have both `len` and `const` rules on the same field")
			m.Assert(minLen == nil, "cannot have both `len` and `min_len` rules on the same field")
			m.Assert(maxLen == nil, "cannot have both `len` and `max_len` rules on the same field")
			length = r.Len
		}
		if r.MinLen != nil {
			m.Assert(minLen == nil, "cannot have multi `min_len` rules on the same field")
			m.Assert(!ct, "cannot have both `min_len` and `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `len` and `min_len` rules on the same field")
			minLen = r.MinLen
		}
		if r.MaxLen != nil {
			m.Assert(maxLen == nil, "cannot have multi `max_len` rules on the same field")
			m.Assert(!ct, "cannot have both `max_len` and `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `len` and `max_len` rules on the same field")
			maxLen = r.MaxLen
		}
		if r.LenBytes != nil {
			m.Assert(lenBs == nil, "cannot have multi `len_bytes` rules on the same field")
			m.Assert(!ct, "cannot have both `len_bytes` and `const` rules on the same field")
			m.Assert(minBs == nil, "cannot have both `len_bytes` and `min_bytes` rules on the same field")
			m.Assert(maxBs == nil, "cannot have both `len_bytes` and `max_bytes` rules on the same field")
			lenBs = r.LenBytes
		}
		if r.MinBytes != nil {
			m.Assert(minBs == nil, "cannot have multi `min_bytes` rules on the same field")
			m.Assert(!ct, "cannot have both `min_bytes` and `const` rules on the same field")
			m.Assert(lenBs == nil, "cannot have both `len_bytes` and `min_bytes` rules on the same field")
			minBs = r.MinBytes
		}
		if r.MaxLen != nil {
			m.Assert(maxBs == nil, "cannot have multi `max_bytes` rules on the same field")
			m.Assert(!ct, "cannot have both `max_bytes` and `const` rules on the same field")
			m.Assert(lenBs == nil, "cannot have both `len_bytes` and `max_bytes` rules on the same field")
			maxBs = r.MaxBytes
		}
		if r.Pattern != nil {
			m.Assert(!pattern, "cannot have multi `pattern` rules on the same field")
			m.Assert(!ct, "cannot have both `pattern` and `const` rules on the same field")
			m.Assert(!wk, "cannot have both `pattern` and `well_known` rules on the same field")
			pattern = true
		}
		if r.Prefix != nil {
			m.Assert(!prefix, "cannot have multi `prefix` rules on the same field")
			prefix = true
		}
		if r.Suffix != nil {
			m.Assert(!suffix, "cannot have multi `suffix` rules on the same field")
			suffix = true
		}
		if r.Contains != nil {
			m.Assert(!contains, "cannot have multi `contains` rules on the same field")
			contains = true
		}
		if r.NotContains != nil {
			m.Assert(!notContains, "cannot have multi `not_contains` rules on the same field")
			notContains = true
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(!ct, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(!ct, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.WellKnown != nil {
			m.Assert(!wk, "cannot have multi `well_known` rules on the same field")
			m.Assert(!ct, "cannot have both `well_known` and `const` rules on the same field")
			m.Assert(!pattern, "cannot have both `pattern` and `well_known` rules on the same field")
			wk = true
		}

		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckString(r)
	}

	m.checkLen(length, minLen, maxLen)
	m.checkLen(lenBs, minBs, maxBs)
	m.checkMinMax(minLen, maxLen)
	m.checkMinMax(minBs, maxBs)
}

func (m *Module) CheckString(r *validate.StringRule) {
	m.checkLen(r.Len, r.MinLen, r.MaxLen)
	m.checkLen(r.LenBytes, r.MinBytes, r.MaxBytes)
	m.checkMinMax(r.MinLen, r.MaxLen)
	m.checkMinMax(r.MinBytes, r.MaxBytes)
	m.checkIns(len(r.In), len(r.NotIn))
	m.checkWellKnownRegex(r.GetWellKnownRegex(), r)
	m.checkPattern(r.Pattern, len(r.In))

	if r.MaxLen != nil {
		max := int(r.GetMaxLen())
		m.Assert(utf8.RuneCountInString(r.GetPrefix()) <= max, "`prefix` length exceeds the `max_len`")
		m.Assert(utf8.RuneCountInString(r.GetSuffix()) <= max, "`suffix` length exceeds the `max_len`")
		m.Assert(utf8.RuneCountInString(r.GetContains()) <= max, "`contains` length exceeds the `max_len`")

		m.Assert(
			r.MaxBytes == nil || r.GetMaxBytes() >= r.GetMaxLen(),
			"`max_len` cannot exceed `max_bytes`")
	}

	if r.MaxBytes != nil {
		max := int(r.GetMaxBytes())
		m.Assert(len(r.GetPrefix()) <= max, "`prefix` length exceeds the `max_bytes`")
		m.Assert(len(r.GetSuffix()) <= max, "`suffix` length exceeds the `max_bytes`")
		m.Assert(len(r.GetContains()) <= max, "`contains` length exceeds the `max_bytes`")
	}
	m.Assert(len(r.Error.GetMethod()) > 0, "method to create error instance can not be empty")
}

func (m *Module) CheckBytesRules(rules *validate.BytesRules, inject bool) {
	var (
		pattern, wk                             bool
		length, minLen, maxLen                  *uint64
		in, notIn, ct, prefix, suffix, contains bool
	)
	for _, r := range rules.Rules {
		if r.Const != nil {
			m.Assert(!ct, "cannot have multi `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `const` and `len` rules on the same field")
			m.Assert(minLen == nil, "cannot have both `const` and `min_len` rules on the same field")
			m.Assert(maxLen == nil, "cannot have both `const` and `max_len` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(!prefix, "cannot have both `const` and `prefix` rules on the same field")
			m.Assert(!suffix, "cannot have both `const` and `suffix` rules on the same field")
			m.Assert(!contains, "cannot have both `const` and `contains` rules on the same field")
			m.Assert(!pattern, "cannot have both `const` and `pattern` rules on the same field")
			m.Assert(!wk, "cannot have both `const` and `well_known` rules on the same field")
		}
		if r.Len != nil {
			m.Assert(length == nil, "cannot have multi `len` rules on the same field")
			m.Assert(!ct, "cannot have both `len` and `const` rules on the same field")
			m.Assert(minLen == nil, "cannot have both `len` and `min_len` rules on the same field")
			m.Assert(maxLen == nil, "cannot have both `len` and `max_len` rules on the same field")
			length = r.Len
		}
		if r.MinLen != nil {
			m.Assert(minLen == nil, "cannot have multi `min_len` rules on the same field")
			m.Assert(!ct, "cannot have both `min_len` and `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `len` and `min_len` rules on the same field")
			minLen = r.MinLen
		}
		if r.MaxLen != nil {
			m.Assert(maxLen == nil, "cannot have multi `max_len` rules on the same field")
			m.Assert(!ct, "cannot have both `max_len` and `const` rules on the same field")
			m.Assert(length == nil, "cannot have both `len` and `max_len` rules on the same field")
			maxLen = r.MaxLen
		}
		if r.Pattern != nil {
			m.Assert(!pattern, "cannot have multi `pattern` rules on the same field")
			m.Assert(!ct, "cannot have both `pattern` and `const` rules on the same field")
			m.Assert(!wk, "cannot have both `pattern` and `well_known` rules on the same field")
			pattern = true
		}
		if r.Prefix != nil {
			m.Assert(!prefix, "cannot have multi `prefix` rules on the same field")
			prefix = true
		}
		if r.Suffix != nil {
			m.Assert(!suffix, "cannot have multi `suffix` rules on the same field")
			suffix = true
		}
		if r.Contains != nil {
			m.Assert(!contains, "cannot have multi `contains` rules on the same field")
			contains = true
		}
		if len(r.In) > 0 {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(!ct, "cannot have both `in` and `const` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if len(r.NotIn) > 0 {
			m.Assert(!notIn, "cannot have multi `in` rules on the same field")
			m.Assert(!ct, "cannot have both `not_in` and `const` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.WellKnown != nil {
			m.Assert(!wk, "cannot have multi `well_known` rules on the same field")
			m.Assert(!ct, "cannot have both `well_known` and `const` rules on the same field")
			m.Assert(!pattern, "cannot have both `pattern` and `well_known` rules on the same field")
			wk = true
		}

		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckBytes(r)
	}

	m.checkLen(length, minLen, maxLen)
	m.checkMinMax(minLen, maxLen)
}

func (m *Module) CheckBytes(r *validate.BytesRule) {
	m.checkMinMax(r.MinLen, r.MaxLen)
	m.checkIns(len(r.In), len(r.NotIn))
	m.checkPattern(r.Pattern, len(r.In))

	if r.MaxLen != nil {
		max := int(r.GetMaxLen())
		m.Assert(len(r.GetPrefix()) <= max, "`prefix` length exceeds the `max_len`")
		m.Assert(len(r.GetSuffix()) <= max, "`suffix` length exceeds the `max_len`")
		m.Assert(len(r.GetContains()) <= max, "`contains` length exceeds the `max_len`")
	}
}

func (m *Module) CheckEnum(ft FieldType, r *validate.EnumRules, inject bool) {
	m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
	m.checkIns(len(r.In), len(r.NotIn))

	if r.GetDefinedOnly() && len(r.In) > 0 {
		typ, ok := ft.(interface {
			Enum() pgs.Enum
		})

		if !ok {
			m.Failf("unexpected field type (%T)", ft)
		}

		defined := typ.Enum().Values()
		vals := make(map[int32]struct{}, len(defined))

		for _, val := range defined {
			vals[val.Value()] = struct{}{}
		}

		for _, in := range r.In {
			if _, ok = vals[in]; !ok {
				m.Failf("undefined `in` value (%d) conflicts with `defined_only` rule")
			}
		}
	}
}

func (m *Module) CheckMessage(f pgs.Field, rules *validate.FieldRules, inject bool) {
	m.Assert(f.Type().IsEmbed(), "field is not embedded but got message rules")
	emb := f.Type().Embed()
	if emb != nil && emb.IsWellKnown() {
		switch emb.WellKnownType() {
		case pgs.AnyWKT:
			m.Failf("Any rules should be used for Any fields")
		case pgs.DurationWKT:
			m.Failf("Duration rules should be used for Duration fields")
		case pgs.TimestampWKT:
			m.Failf("Timestamp rules should be used for Timestamp fields")
		}
	}

	if rules.Type != nil && rules.GetMessage().GetSkip() {
		m.Failf("Skip should not be used with WKT scalar rules")
	}

	if rules.GetMessage().GetRequired() {
		m.Assert(inject || rules.GetMessage().GetError() != nil, "error should be defined when message is required")
	}
}

func (m *Module) CheckRepeatedRules(ft FieldType, rules *validate.RepeatedRules, inject bool) {
	typ := m.mustFieldType(ft)

	m.Assert(typ.IsRepeated(), "field is not repeated but got repeated rules")
	var (
		unique, itemsRule  bool
		minItems, maxItems *uint64
	)
	for _, r := range rules.Rules {
		if r.MinItems != nil {
			m.Assert(minItems == nil, "cannot have multi `min_items` rules on the same field")
			minItems = r.MinItems
		}
		if r.MaxItems != nil {
			m.Assert(maxItems == nil, "cannot have multi `max_items` rules on the same field")
			maxItems = r.MaxItems
		}
		if r.Unique != nil {
			m.Assert(!unique, "cannot have multi `unique` rules on the same field")
		}
		if r.Items != nil {
			m.Assert(!itemsRule, "cannot have multi `items` rules on the same field")
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckRepeated(typ, r)
	}
	m.checkMinMax(minItems, maxItems)
}

func (m *Module) CheckRepeated(typ pgs.FieldType, r *validate.RepeatedRule) {

	m.checkMinMax(r.MinItems, r.MaxItems)

	if r.GetUnique() {
		m.Assert(
			!typ.Element().IsEmbed(),
			"unique rule is only applicable for scalar types")
	}

	m.Push("items")
	m.CheckFieldRules(typ.Element(), r.Items, true)
	m.Pop()
}

func (m *Module) CheckMapRules(ft FieldType, rules *validate.MapRules, inject bool) {
	typ := m.mustFieldType(ft)
	m.Assert(typ.IsMap(), "field is not a map but got map rules")

	var (
		noSparse             bool
		keysRule, valuesRule bool
		minPairs, maxPairs   *uint64
	)

	for _, r := range rules.Rules {
		if r.MinPairs != nil {
			m.Assert(minPairs == nil, "cannot have multi `min_pairs` rules on the same field")
			minPairs = r.MinPairs
		}
		if r.MaxPairs != nil {
			m.Assert(maxPairs == nil, "cannot have multi `max_pairs` rules on the same field")
			maxPairs = r.MaxPairs
		}
		if r.NoSparse != nil {
			m.Assert(!noSparse, "cannot have multi `no_sparse` rules on the same field")
			noSparse = true
		}
		if r.Keys != nil {
			m.Assert(!keysRule, "cannot have multi `keys` rules on the same field")
			keysRule = true
		}
		if r.Values != nil {
			m.Assert(!valuesRule, "cannot have multi `values` rules on the same field")
			valuesRule = true
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckMap(typ, r)
	}

	m.checkMinMax(minPairs, maxPairs)
}

func (m *Module) CheckMap(typ pgs.FieldType, r *validate.MapRule) {
	m.checkMinMax(r.MinPairs, r.MaxPairs)

	if r.GetNoSparse() {
		m.Assert(
			typ.Element().IsEmbed(),
			"no_sparse rule is only applicable for embedded message types",
		)
	}

	m.Push("keys")
	m.CheckFieldRules(typ.Key(), r.Keys, true)
	m.Pop()

	m.Push("values")
	m.CheckFieldRules(typ.Element(), r.Values, true)
	m.Pop()
}

func (m *Module) CheckAnyRules(ft FieldType, rules *validate.AnyRules, inject bool) {
	var in, notIn, required bool
	for _, r := range rules.GetRules() {
		if r.In != nil {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			in = true
		}
		if r.NotIn != nil {
			m.Assert(!notIn, "cannot have multi `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Required != nil {
			m.Assert(!required, "cannot have multi `required` rules on the same field")
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.checkIns(len(r.In), len(r.NotIn))
	}
}

func (m *Module) CheckDurationRules(ft FieldType, rules *validate.DurationRules, inject bool) {
	var (
		lt, lte, gt, gte, ct *time.Duration
		required, in, notIn  bool
	)
	for _, r := range rules.GetRules() {
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			m.Assert(!in, "cannot have both `const` and `in` rules on the same field")
			m.Assert(!notIn, "cannot have both `const` and `not_in` rules on the same field")
			ct = m.checkDur(r.GetConst())
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `lt` rules on the same field")
			lt = m.checkDur(r.GetLt())
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `lte` rules on the same field")
			lte = m.checkDur(r.GetLte())
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `gt` rules on the same field")
			gt = m.checkDur(r.GetGt())
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `gte` rules on the same field")
			gte = m.checkDur(r.GetGte())
		}
		if r.In != nil {
			m.Assert(!in, "cannot have multi `in` rules on the same field")
			m.Assert(!notIn, "cannot have both `in` and `not_in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `in` rules on the same field")
			in = true
		}
		if r.NotIn != nil {
			m.Assert(!notIn, "cannot have multi `not_in` rules on the same field")
			m.Assert(!in, "cannot have both `in` and `not_in` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `not_in` rules on the same field")
			notIn = true
		}
		if r.Required != nil {
			m.Assert(!required, "cannot have multi `required` rules on the same field")
			required = true
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckDuration(ft, r)
	}
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckDuration(ft FieldType, r *validate.DurationRule) {
	m.checkNums(
		len(r.GetIn()),
		len(r.GetNotIn()),
		m.checkDur(r.GetConst()),
		m.checkDur(r.GetLt()),
		m.checkDur(r.GetLte()),
		m.checkDur(r.GetGt()),
		m.checkDur(r.GetGte()))

	for _, v := range r.GetIn() {
		m.Assert(v != nil, "cannot have nil values in `in`")
		m.checkDur(v)
	}

	for _, v := range r.GetNotIn() {
		m.Assert(v != nil, "cannot have nil values in `not_in`")
		m.checkDur(v)
	}
}

func (m *Module) CheckTimestampRules(ft FieldType, rules *validate.TimestampRules, inject bool) {
	var (
		withIn                 *time.Duration
		lt, lte, gt, gte, ct   *int64
		required, ltNow, gtNow bool
	)
	for _, r := range rules.GetRules() {
		if r.Required != nil {
			m.Assert(!required, "cannot have multi `required` rules on the same field")
		}
		if r.Const != nil {
			m.Assert(ct == nil, "cannot have multi `const` rules on the same field")
			m.Assert(lt == nil, "cannot have both `const` and `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `const` and `lte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `const` and `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `const` and `gte` rules on the same field")
			m.Assert(!ltNow, "cannot have both `const` and `lt_now` rules on the same field")
			m.Assert(!gtNow, "cannot have both `const` and `gt_now` rules on the same field")
			m.Assert(withIn == nil, "cannot have both `const` and `within` rules on the same field")
			ct = m.checkTS(r.GetConst())
		}
		if r.Lt != nil {
			m.Assert(lt == nil, "cannot have multi `lt` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lt` and `lte` rules on the same field")
			m.Assert(!ltNow, "cannot have both `lt` and `lt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `lt` rules on the same field")
			lt = m.checkTS(r.GetLt())
		}
		if r.Lte != nil {
			m.Assert(lte == nil, "cannot have multi `lte` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lte` rules on the same field")
			m.Assert(!ltNow, "cannot have both `lte` and `lt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `lte` rules on the same field")
			lte = m.checkTS(r.GetLte())
		}
		if r.Gt != nil {
			m.Assert(gt == nil, "cannot have multi `gt` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gt` and `gte` rules on the same field")
			m.Assert(!gtNow, "cannot have both `gt` and `gt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `gt` rules on the same field")
			gt = m.checkTS(r.GetGt())
		}
		if r.Gte != nil {
			m.Assert(gte == nil, "cannot have multi `gte` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gte` rules on the same field")
			m.Assert(!gtNow, "cannot have both `gte` and `gt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `gte` rules on the same field")
			gte = m.checkTS(r.GetGte())
		}
		if r.LtNow != nil {
			m.Assert(!ltNow, "cannot have multi `lt_now` rules on the same field")
			m.Assert(lt == nil, "cannot have both `lt` and `lt_now` rules on the same field")
			m.Assert(lte == nil, "cannot have both `lte` and `lt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `lt_now` rules on the same field")
			ltNow = true
		}
		if r.GtNow != nil {
			m.Assert(!gtNow, "cannot have multi `gt_now` rules on the same field")
			m.Assert(gt == nil, "cannot have both `gt` and `gt_now` rules on the same field")
			m.Assert(gte == nil, "cannot have both `gte` and `gt_now` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `gt_now` rules on the same field")
			gtNow = true
		}
		if r.Within != nil {
			m.Assert(withIn == nil, "cannot have multi `within` rules on the same field")
			m.Assert(ct == nil, "cannot have both `const` and `within` rules on the same field")
			withIn = m.checkDur(r.GetWithin())
		}
		m.Assert(inject || r.Error != nil, "cannot have nil error on rule")
		m.CheckTimestamp(ft, r)
	}
	m.Assert(withIn == nil || (ltNow || gtNow), "within` rule cannot be used with absolute `lt/gt` rules")
	m.checkNums(0, 0, ct, lt, lte, gt, gte)
}

func (m *Module) CheckTimestamp(ft FieldType, r *validate.TimestampRule) {
	m.checkNums(0, 0,
		m.checkTS(r.GetConst()),
		m.checkTS(r.GetLt()),
		m.checkTS(r.GetLte()),
		m.checkTS(r.GetGt()),
		m.checkTS(r.GetGte()))

	m.Assert(
		(r.LtNow == nil && r.GtNow == nil) || (r.Lt == nil && r.Lte == nil && r.Gt == nil && r.Gte == nil),
		"`now` rules cannot be mixed with absolute `lt/gt` rules")

	m.Assert(
		r.Within == nil || (r.Lt == nil && r.Lte == nil && r.Gt == nil && r.Gte == nil),
		"`within` rule cannot be used with absolute `lt/gt` rules")

	m.Assert(
		r.LtNow == nil || r.GtNow == nil,
		"both `now` rules cannot be used together")

	dur := m.checkDur(r.Within)
	m.Assert(
		dur == nil || *dur > 0,
		"`within` rule must be positive and non-zero")
}

func (m *Module) mustFieldType(ft FieldType) pgs.FieldType {
	typ, ok := ft.(pgs.FieldType)
	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}

	return typ
}

func (m *Module) checkNums(in, notIn int, ci, lti, ltei, gti, gtei interface{}) {
	m.checkIns(in, notIn)

	c := reflect.ValueOf(ci)
	lt, lte := reflect.ValueOf(lti), reflect.ValueOf(ltei)
	gt, gte := reflect.ValueOf(gti), reflect.ValueOf(gtei)

	m.Assert(
		c.IsNil() ||
			in == 0 && notIn == 0 &&
				lt.IsNil() && lte.IsNil() &&
				gt.IsNil() && gte.IsNil(),
		"`const` can be the only rule on a field",
	)

	m.Assert(
		in == 0 ||
			lt.IsNil() && lte.IsNil() &&
				gt.IsNil() && gte.IsNil(),
		"cannot have both `in` and range constraint rules on the same field",
	)

	m.Assert(
		lt.IsNil() || lte.IsNil(),
		"cannot have both `lt` and `lte` rules on the same field",
	)

	m.Assert(
		gt.IsNil() || gte.IsNil(),
		"cannot have both `gt` and `gte` rules on the same field",
	)

	if !lt.IsNil() {
		m.Assert(gt.IsNil() || !reflect.DeepEqual(lti, gti),
			"cannot have equal `gt` and `lt` rules on the same field")
		m.Assert(gte.IsNil() || !reflect.DeepEqual(lti, gtei),
			"cannot have equal `gte` and `lt` rules on the same field")
	} else if !lte.IsNil() {
		m.Assert(gt.IsNil() || !reflect.DeepEqual(ltei, gti),
			"cannot have equal `gt` and `lte` rules on the same field")
		m.Assert(gte.IsNil() || !reflect.DeepEqual(ltei, gtei),
			"use `const` instead of equal `lte` and `gte` rules")
	}
}

func (m *Module) checkIns(in, notIn int) {
	m.Assert(
		in == 0 || notIn == 0,
		"cannot have both `in` and `not_in` rules on the same field")
}

func (m *Module) checkMinMax(min, max *uint64) {
	if min == nil || max == nil {
		return
	}

	m.Assert(
		*min <= *max,
		"`min` value is greater than `max` value")
}

func (m *Module) checkLen(len, min, max *uint64) {
	if len == nil {
		return
	}

	m.Assert(
		min == nil,
		"cannot have both `len` and `min_len` rules on the same field")

	m.Assert(
		max == nil,
		"cannot have both `len` and `max_len` rules on the same field")
}

func (m *Module) checkWellKnownRegex(wk validate.KnownRegex, r *validate.StringRule) {
	if wk != 0 {
		m.Assert(r.Pattern == nil, "regex `well_known_regex` and regex `pattern` are incompatible")
		var non_strict = r.Strict != nil && *r.Strict == false
		if (wk.String() == "HTTP_HEADER_NAME" || wk.String() == "HTTP_HEADER_VALUE") && non_strict {
			// Use non-strict header validation.
			r.Pattern = regex_map["HEADER_STRING"]
		} else {
			r.Pattern = regex_map[wk.String()]
		}
	}
}

func (m *Module) checkPattern(p *string, in int) {
	if p != nil {
		m.Assert(in == 0, "regex `pattern` and `in` rules are incompatible")
		_, err := regexp.Compile(*p)
		m.CheckErr(err, "unable to parse regex `pattern`")
	}
}

func (m *Module) checkDur(d *durationpb.Duration) *time.Duration {
	if d == nil {
		return nil
	}

	dur, err := d.AsDuration(), d.CheckValid()
	m.CheckErr(err, "could not resolve duration")
	return &dur
}

func (m *Module) checkTS(ts *timestamppb.Timestamp) *int64 {
	if ts == nil {
		return nil
	}

	t, err := ts.AsTime(), ts.CheckValid()
	m.CheckErr(err, "could not resolve timestamp")
	return proto.Int64(t.UnixNano())
}
