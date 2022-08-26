package shared

import (
	pgs "github.com/lyft/protoc-gen-star"

	"github.com/curl-li/protoc-gen-validate/validate"
)

// Disabled returns true if validations are disabled for msg
func Disabled(msg pgs.Message) (disabled bool, err error) {
	_, err = msg.Extension(validate.E_Disabled, &disabled)
	return
}

// Ignore returns true if validations aren't to be generated for msg
func Ignored(msg pgs.Message) (ignored bool, err error) {
	_, err = msg.Extension(validate.E_Ignored, &ignored)
	return
}

// RequiredOneOf returns true if the oneof field requires a field to be set
func RequiredOneOf(oo pgs.OneOf) (required bool, err error) {
	var oneof validate.OneOf
	if _, err = oo.Extension(validate.E_Oneof, &oneof); err != nil {
		return
	}
	required = oneof.GetRequired()
	return
}

// OneOfRule returns oneof rule
func OneOfRule(oo pgs.OneOf) (rule *validate.OneOf, err error) {
	_, err = oo.Extension(validate.E_Oneof, &rule)
	return
}
