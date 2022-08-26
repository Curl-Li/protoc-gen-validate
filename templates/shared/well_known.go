package shared

import (
	pgs "github.com/lyft/protoc-gen-star"

	"github.com/curl-li/protoc-gen-validate/validate"
)

type WellKnown string

const (
	Email    WellKnown = "email"
	Hostname WellKnown = "hostname"
	UUID     WellKnown = "uuid"
)

func FileNeeds(f pgs.File, wk WellKnown) bool {
	for _, msg := range f.AllMessages() {
		needed := Needs(msg, wk)
		if needed {
			return true
		}
	}

	return false
}

// Needs returns true if a well-known string validator is needed for this
// message.
func Needs(m pgs.Message, wk WellKnown) bool {

	for _, f := range m.Fields() {
		var rules validate.FieldRules
		if _, err := f.Extension(validate.E_Rules, &rules); err != nil {
			continue
		}

		switch {
		case f.Type().IsRepeated() && f.Type().Element().ProtoType() == pgs.StringT:
			var strRule *validate.StringRules
			for _, rr := range rules.GetRepeated().Rules {
				if rr.GetItems() != nil {
					strRule = rr.GetItems().GetString_()
				}
			}
			if strRulesNeeds(strRule, wk) {
				return true
			}
		case f.Type().IsMap():

			if f.Type().Key().ProtoType() == pgs.StringT {
				var strRule *validate.StringRules
				for _, rr := range rules.GetMap().Rules {
					if rr.Keys != nil {
						strRule = rr.GetKeys().GetString_()
					}
				}
				if strRulesNeeds(strRule, wk) {
					return true
				}
			}
			if f.Type().Element().ProtoType() == pgs.StringT {
				var strRule *validate.StringRules
				for _, rr := range rules.GetMap().Rules {
					if rr.Keys != nil {
						strRule = rr.GetValues().GetString_()
					}
				}
				if strRulesNeeds(strRule, wk) {
					return true
				}
			}
		case f.Type().ProtoType() == pgs.StringT:
			if strRulesNeeds(rules.GetString_(), wk) {
				return true
			}
		case f.Type().ProtoType() == pgs.MessageT && f.Type().IsEmbed() && f.Type().Embed().WellKnownType() == pgs.StringValueWKT:
			if strRulesNeeds(rules.GetString_(), wk) {
				return true
			}
		}
	}

	return false
}

func strRulesNeeds(rules *validate.StringRules, wk WellKnown) bool {
	for _, rule := range rules.Rules {
		switch wk {
		case Email:
			if rule.GetEmail() {
				return true
			}
		case Hostname:
			if rule.GetEmail() || rule.GetHostname() || rule.GetAddress() {
				return true
			}
		case UUID:
			if rule.GetUuid() {
				return true
			}
		}
	}
	return false
}