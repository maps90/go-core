package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// ValidFormer valid interface
type ValidFormer interface {
	Valid(*Validation)
}

// Error show the error
type Error struct {
	Message, Key, Name, Field, Tmpl string
	Value                           interface{}
	LimitValue                      interface{}
}

// CustomErrorMessage for custom response message
type CustomErrorMessage struct {
	Message string `json:"errorMessage"`
	Key     string `json:"errorKey"`
	Field   string `json:"errorField"`
}

// String Returns the Message.
func (e *Error) String() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Result is returned from every validation method.
// It provides an indication of success, and a pointer to the Error (if any).
type Result struct {
	Error *Error
	Ok    bool
}

// Key Get Result by given key string.
func (r *Result) Key(key string) *Result {
	if r.Error != nil {
		r.Error.Key = key
	}
	return r
}

// Message Set Result message by string or format string with args
func (r *Result) Message(message string, args ...interface{}) *Result {
	if r.Error != nil {
		if len(args) == 0 {
			r.Error.Message = message
		} else {
			r.Error.Message = fmt.Sprintf(message, args...)
		}
	}
	return r
}

// A Validation context manages data validation and error messages.
type Validation struct {
	Errors    []*Error
	ErrorsMap map[string]*Error
}

// IsJson check whether string in json format or not
func (v *Validation) IsJson(s string) bool {
	var js map[string]interface{}

	return json.Unmarshal([]byte(s), &js) == nil
}

// SetCustomErrorMessage set error message with json on CustomErrorMessage format
func (v *Validation) SetCustomErrorMessage(msg string, key string, field string) string {
	errString := &CustomErrorMessage{msg, key, field}
	out, _ := json.Marshal(errString)

	return string(out)
}

// GetCustomErrorMessage extract data from json error message
func (v *Validation) GetCustomErrorMessage(msg string) (string, string, string) {
	if !v.IsJson(msg) {
		return msg, "", ""
	}

	res := CustomErrorMessage{}
	json.Unmarshal([]byte(msg), &res)

	return res.Message, res.Key, res.Field
}

// Clear Clean all ValidationError.
func (v *Validation) Clear() {
	v.Errors = []*Error{}
	v.ErrorsMap = nil
}

// HasErrors Has ValidationError nor not.
func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

// ErrorMap Return the errors mapped by key.
// If there are multiple validation errors associated with a single key, the
// first one "wins".  (Typically the first validation will be the more basic).
func (v *Validation) ErrorMap() map[string]*Error {
	return v.ErrorsMap
}

// Error Add an error to the validation context.
func (v *Validation) Error(message string, args ...interface{}) *Result {
	result := (&Result{
		Ok:    false,
		Error: &Error{},
	}).Message(message, args...)
	v.Errors = append(v.Errors, result.Error)
	return result
}

// Required Test that the argument is non-nil and non-empty (if string or list)
func (v *Validation) Required(obj interface{}, key string) *Result {
	return v.apply(Required{key}, obj)
}

func (v *Validation) PhoneNumber(obj interface{}, key string) *Result {
	return v.apply(IsPhone{Match{Regexp: phonePattern}, key}, obj)
}

// Max Test that the obj is less than max if obj's type is int
func (v *Validation) Max(obj interface{}, max int, key string) *Result {
	return v.apply(Max{max, key}, obj)
}

// Range Test that the obj is between mni and max if obj's type is int
func (v *Validation) Range(obj interface{}, min, max int, key string) *Result {
	return v.apply(Range{Min{Min: min}, Max{Max: max}, key}, obj)
}

// MinSize Test that the obj is longer than min size if type is string or slice
func (v *Validation) MinSize(obj interface{}, min int, key string) *Result {
	return v.apply(MinSize{min, key}, obj)
}

// MaxSize Test that the obj is shorter than max size if type is string or slice
func (v *Validation) MaxSize(obj interface{}, max int, key string) *Result {
	return v.apply(MaxSize{max, key}, obj)
}

// Length Test that the obj is same length to n if type is string or slice
func (v *Validation) Length(obj interface{}, n int, key string) *Result {
	return v.apply(Length{n, key}, obj)
}

// Alpha Test that the obj is [a-zA-Z] if type is string
func (v *Validation) Alpha(obj interface{}, key string) *Result {
	return v.apply(Alpha{key}, obj)
}

// Numeric Test that the obj is [0-9] if type is string
func (v *Validation) Numeric(obj interface{}, key string) *Result {
	return v.apply(Numeric{key}, obj)
}

// Float Test that the obj is float if type is string
func (v *Validation) Float(obj interface{}, key string) *Result {
	return v.apply(Float{key}, obj)
}

// AlphaNumeric Test that the obj is [0-9a-zA-Z] if type is string
func (v *Validation) AlphaNumeric(obj interface{}, key string) *Result {
	return v.apply(AlphaNumeric{key}, obj)
}

// Match Test that the obj matches regexp if type is string
func (v *Validation) Match(obj interface{}, regex *regexp.Regexp, key string) *Result {
	return v.apply(Match{regex, key}, obj)
}

// NoMatch Test that the obj doesn't match regexp if type is string
func (v *Validation) NoMatch(obj interface{}, regex *regexp.Regexp, key string) *Result {
	return v.apply(NoMatch{Match{Regexp: regex}, key}, obj)
}

// AlphaDash Test that the obj is [0-9a-zA-Z_-] if type is string
func (v *Validation) AlphaDash(obj interface{}, key string) *Result {
	return v.apply(AlphaDash{NoMatch{Match: Match{Regexp: alphaDashPattern}}, key}, obj)
}

func (v *Validation) IsDate(obj interface{}, format, key string) *Result {
	return v.apply(IsDate{format, key}, obj)
}

func (v *Validation) DateBefore(obj interface{}, reference, format, key string) *Result {
	return v.apply(DateBefore{reference, format, key}, obj)
}

func (v *Validation) SliceMatch(obj interface{}, haystack interface{}, key string) *Result {
	return v.apply(SliceMatch{haystack, key}, obj)
}

func (v *Validation) Duplicate(obj interface{}, monitor, key string) *Result {
	return v.apply(Duplicate{monitor, key}, obj)
}

func (v *Validation) Incremental(obj interface{}, key string) *Result {
	return v.apply(Incremental{key}, obj)
}

func (v *Validation) Email(obj interface{}, key string) *Result {
	return v.apply(Email{Match{Regexp: emailPattern}, key}, obj)
}

func (v *Validation) PositiveFloat(obj interface{}, key string) *Result {
	return v.apply(PositiveFloat{Match{Regexp: positiveFloatPattern}, key}, obj)
}


func (v *Validation) apply(chk Validator, obj interface{}) *Result {
	result := &Result{Ok: true}
	_, ok := chk.(Required)
	if !ok {
		if o := checkEmpty(obj); o {
			return &Result{Ok: true}
		}
	}
	if chk.IsSatisfied(obj) {
		return result
	}

	// Add the error to the validation context.
	key := chk.GetKey()
	Name := key
	Field := ""
	msg := chk.DefaultMessage()

	if v.IsJson(key) {
		msg, key, Field = v.GetCustomErrorMessage(key)
		msg = v.SetCustomErrorMessage(msg+" "+chk.DefaultMessage(), key, Field)
	}

	parts := strings.Split(key, ".")
	if len(parts) == 2 {
		Field = parts[0]
		Name = parts[1]
	}

	err := &Error{
		Message:    msg,
		Key:        key,
		Name:       Name,
		Field:      Field,
		Value:      obj,
		Tmpl:       MessageTmpls[Name],
		LimitValue: chk.GetLimitValue(),
	}
	v.setError(err)

	// Also return it in the result.
	return &Result{
		Ok:    false,
		Error: err,
	}
}

func (v *Validation) setError(err *Error) {
	v.Errors = append(v.Errors, err)
	if v.ErrorsMap == nil {
		v.ErrorsMap = make(map[string]*Error)
	}
	if _, ok := v.ErrorsMap[err.Field]; !ok {
		v.ErrorsMap[err.Field] = err
	}
}

// SetError Set error message for one field in ValidationError
func (v *Validation) SetError(fieldName string, errMsg string) *Error {
	err := &Error{Key: fieldName, Field: fieldName, Tmpl: errMsg, Message: errMsg}
	v.setError(err)
	return err
}

// Check Apply a group of validators to a field, in order, and return the
// ValidationResult from the first one that fails, or the last one that
// succeeds.
func (v *Validation) Check(obj interface{}, checks ...Validator) *Result {
	var result *Result
	for _, check := range checks {
		result = v.apply(check, obj)
		if !result.Ok {
			return result
		}
	}
	return result
}

func (v *Validation) Validate(obj interface{}) *Validation {
	v.Valid(obj)
	return v
}

// Valid Validate a struct.
// the obj parameter must be a struct or a struct pointer
func (v *Validation) Valid(obj interface{}) (b bool, err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
		return
	}

	for i := 0; i < objT.NumField(); i++ {
		var vfs []ValidFunc
		if vfs, err = getValidFuncs(objT.Field(i)); err != nil {
			return
		}
		for _, vf := range vfs {
			if _, err = funcs.Call(vf.Name,
				mergeParam(v, objV.Field(i).Interface(), vf.Params)...); err != nil {
				return
			}
		}
	}

	if !v.HasErrors() {
		if form, ok := obj.(ValidFormer); ok {
			form.Valid(v)
		}
	}

	return !v.HasErrors(), nil
}

// RecursiveValid Recursively validate a struct.
// Step1: Validate by v.Valid
// Step2: If pass on step1, then reflect obj's fields
// Step3: Do the Recursively validation to all struct, struct pointer or slice struct fields
func (v *Validation) RecursiveValid(objc interface{}) (bool, error) {
	//Step 1: validate obj itself firstly
	// fails if objc is not struct
	pass, err := v.Valid(objc)
	if err != nil || false == pass {
		return pass, err // Stop recursive validation
	}
	// Step 2: Validate struct's struct fields
	objT := reflect.TypeOf(objc)
	objV := reflect.ValueOf(objc)

	if isStructPtr(objT) {
		objT = objT.Elem()
		objV = objV.Elem()
	}

	for i := 0; i < objT.NumField(); i++ {

		t := objT.Field(i).Type
		if _, ok := objV.Field(i).Interface().(*time.Time); ok {
			continue
		}

		// Recursive applies to struct or pointer to struct fields
		if isStruct(t) || isStructPtr(t) {
			// Step 3: do the recursive validation
			// Only valid the Public field recursively
			if objV.Field(i).CanInterface() {
				pass, err = v.RecursiveValid(objV.Field(i).Interface())
			}
		}

		//range over slices struct
		if t.Kind() == reflect.Slice {
			o := objV.Field(i).Interface()
			d := reflect.ValueOf(o)
			for k := 0; k < d.Len(); k++ {
				pass, err = v.Valid(d.Index(k).Interface())
			}
		}
	}
	return pass, err
}

// Valid Validate a struct (except Required validation).
// the obj parameter must be a struct or a struct pointer
func (v *Validation) ValidWithException(obj interface{}, exceptions []string) (b bool, err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
		return
	}

	for i := 0; i < objT.NumField(); i++ {
		var vfs []ValidFunc
		if vfs, err = getValidFuncs(objT.Field(i)); err != nil {
			return
		}

	LoopTagValid:
		for _, vf := range vfs {
			if len(exceptions) > 0 {
				for _, vex := range exceptions {
					if strings.ToLower(vf.Name) == vex {
						continue LoopTagValid
					}
				}
			}

			if _, err = funcs.Call(vf.Name,
				mergeParam(v, objV.Field(i).Interface(), vf.Params)...); err != nil {
				return
			}
		}
	}

	if !v.HasErrors() {
		if form, ok := obj.(ValidFormer); ok {
			form.Valid(v)
		}
	}

	return !v.HasErrors(), nil
}

// RecursiveValid Recursively validate a struct.
// Step1: Validate by v.Valid
// Step2: If pass on step1, then reflect obj's fields
// Step3: Do the Recursively validation to all struct, struct pointer or slice struct fields
func (v *Validation) RecursiveValidWithException(objc interface{}, exceptions []string) (bool, error) {
	//Step 1: validate obj itself firstly
	// fails if objc is not struct
	pass, err := v.ValidWithException(objc, exceptions)

	if err != nil || false == pass {
		return pass, err // Stop recursive validation
	}

	// Step 2: Validate struct's struct fields
	objT := reflect.TypeOf(objc)
	objV := reflect.ValueOf(objc)

	if isStructPtr(objT) {
		objT = objT.Elem()
		objV = objV.Elem()
	}

	for i := 0; i < objT.NumField(); i++ {
		if objT.Field(i).Tag.Get("valid") != "" {
			t := objT.Field(i).Type
			if _, ok := objV.Field(i).Interface().(*time.Time); ok {
				continue
			}

			// Recursive applies to struct or pointer to struct fields
			if isStruct(t) || isStructPtr(t) {
				// Step 3: do the recursive validation
				// Only valid the Public field recursively
				if objV.Field(i).CanInterface() {
					pass, err = v.RecursiveValidWithException(objV.Field(i).Interface(), exceptions)
				}
			}

			//range over slices struct
			if t.Kind() == reflect.Slice {
				o := objV.Field(i).Interface()
				d := reflect.ValueOf(o)
				for k := 0; k < d.Len(); k++ {
					pass, err = v.ValidWithException(d.Index(k).Interface(), exceptions)
				}
			}

		}
	}
	return pass, err
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
