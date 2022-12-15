package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AggregateValidations takes any number of schema.SchemaValidateDiagFunc and executes them one by one
// it returns a diagnostics object containing combined results of each validation function
func AggregateValidations(funcs ...schema.SchemaValidateDiagFunc) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		for _, f := range funcs {
			if err := f(i, path); err != nil {
				diags = append(diags, err...)
			}
		}
		return diags
	}
}

// IsNotBlank verifies whether given value is not blank and returns error if it is where "blank" means:
// - nil value
// - a collection with len == 0 in case the value is a map, array or slice
// - value equal to zero-value for given type (e.g. empty string)
func IsNotBlank(i interface{}, _ cty.Path) diag.Diagnostics {
	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		if val.Len() == 0 {
			return diag.Errorf("provided value cannot be blank")
		}
	default:
		if i == nil || reflect.DeepEqual(i, reflect.Zero(reflect.TypeOf(i)).Interface()) {
			return diag.Errorf("provided value cannot be blank")
		}
	}
	return nil
}

// ValidateJSON checks whether given value is a valid JSON object
func ValidateJSON(val interface{}, _ cty.Path) diag.Diagnostics {
	if str, ok := val.(string); ok {
		var target map[string]interface{}
		if err := json.Unmarshal([]byte(str), &target); err != nil {
			return diag.FromErr(fmt.Errorf("invalid JSON: %s", err))
		}
		return nil
	}
	return diag.Errorf("value is not a string: %s", val)
}

// ValidateNetwork defines network validation logic
func ValidateNetwork(i interface{}, _ cty.Path) diag.Diagnostics {
	val, ok := i.(string)
	if !ok {
		return diag.Errorf("'network' value is not a string: %v", i)
	}
	switch strings.ToLower(val) {
	case "production", "prod", "p", "staging", "stag", "s":
		return nil
	}
	return diag.Errorf("'%s' is an invalid network value: should be 'production', 'prod', 'p', 'staging', 'stag' or 's'", val)
}

// ValidateEmail checks if value is a valid email
func ValidateEmail(val interface{}, _ cty.Path) diag.Diagnostics {
	if str, ok := val.(string); ok {
		return diag.FromErr(validation.Validate(str, validation.Required, is.Email))
	}
	return diag.Errorf("value is not a string: %s", val)
}

// ValidateStringInSlice returns schema.SchemaValidateDiagFunc which tests if the value
// is a string and if it matches given slice of valid strings
func ValidateStringInSlice(valid []string) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		name := path[len(path)-1].(cty.GetAttrStep).Name
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %s to be string", name)
		}

		for _, s := range valid {
			if v == s {
				return nil
			}
		}

		return diag.Errorf("expected %s to be one of ['%s'], got %s", name, strings.Join(valid, "', '"), v)
	}
}

var (
	isRuleFormatValid = regexp.MustCompile(`^v[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString
)

// ValidateRuleFormat checks if value is a valid rule format
func ValidateRuleFormat(v interface{}, _ cty.Path) diag.Diagnostics {
	format, ok := v.(string)
	if !ok {
		return diag.Errorf("expected string, got %T", v)
	}

	if !isRuleFormatValid(format) {
		url := "https://techdocs.akamai.com/property-mgr/reference/latest-behaviors"
		return diag.Errorf(`"rule_format" must be of the form vYYYY-MM-DD (with a leading "v") see %s`, url)
	}

	return nil
}
