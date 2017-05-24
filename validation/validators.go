package validation

import (
	"fmt"
	"github.com/pmylund/sortutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// MessageTmpls store commond validate template
var MessageTmpls = map[string]string{
	"Required":      "is required",
	"Min":           "cannot be less than %d",
	"Max":           "must be less than %d",
	"Range":         "range is between %d to %d",
	"MinSize":       "minimum size is %d",
	"MaxSize":       "maximum size is %d",
	"Length":        "required length is %d",
	"Alpha":         "must be valid alpha characters",
	"Numeric":       "must be valid numeric characters",
	"AlphaNumeric":  "must be valid alpha or numeric characters",
	"Match":         "must match %s",
	"NoMatch":       "must not match %s",
	"AlphaDash":     "must be valid alpha or numeric or dash(-_) characters",
	"Base64":        "must be valid base64 characters",
	"IsDate":        "must be valid date format. eg: '%s'",
	"DateBefore":    "must be set after or equal to %s",
	"SliceMatch":    "only valid for (%s)",
	"Float":         "must be valid decimal/integer value",
	"Duplicate":     "duplicate value detected",
	"Incremental":   "must be in incremental value, start from 1",
	"Phone":         "must be valid phone number",
	"Email":         "must be valid email address",
	"PositiveFloat": "must be positive decimal number (> 0.00)",
	"Name":          "must be valid name.",
}

func SetDefaultMessage(msg map[string]string) {
	if len(msg) == 0 {
		return
	}

	for name := range msg {
		MessageTmpls[name] = msg[name]
	}
}

// Validator interface
type Validator interface {
	IsSatisfied(interface{}) bool
	DefaultMessage() string
	GetKey() string
	GetLimitValue() interface{}
}

// Required struct
type Required struct {
	Key string
}

// IsSatisfied judge whether obj has value
func (r Required) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if _, ok := obj.(bool); ok {
		return true
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if i, ok := obj.(uint); ok {
		return i != 0
	}
	if i, ok := obj.(int8); ok {
		return i != 0
	}
	if i, ok := obj.(uint8); ok {
		return i != 0
	}
	if i, ok := obj.(int16); ok {
		return i != 0
	}
	if i, ok := obj.(uint16); ok {
		return i != 0
	}
	if i, ok := obj.(uint32); ok {
		return i != 0
	}
	if i, ok := obj.(int32); ok {
		return i != 0
	}
	if i, ok := obj.(int64); ok {
		return i != 0
	}
	if i, ok := obj.(uint64); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

// DefaultMessage return the default error message
func (r Required) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Required"])
}

// GetKey return the r.Key
func (r Required) GetKey() string {
	return r.Key
}

// GetLimitValue return nil now
func (r Required) GetLimitValue() interface{} {
	return nil
}

// Min check struct
type Min struct {
	Min int
	Key string
}

// IsSatisfied judge whether obj is valid
func (m Min) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num >= m.Min
	}
	return false
}

// DefaultMessage return the default min error message
func (m Min) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Min"], m.Min)
}

// GetKey return the m.Key
func (m Min) GetKey() string {
	return m.Key
}

// GetLimitValue return the limit value, Min
func (m Min) GetLimitValue() interface{} {
	return m.Min
}

// Max validate struct
type Max struct {
	Max int
	Key string
}

// IsSatisfied judge whether obj is valid
func (m Max) IsSatisfied(obj interface{}) bool {
	num, ok := obj.(int)
	if ok {
		return num <= m.Max
	}
	return false
}

// DefaultMessage return the default max error message
func (m Max) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Max"], m.Max)
}

// GetKey return the m.Key
func (m Max) GetKey() string {
	return m.Key
}

// GetLimitValue return the limit value, Max
func (m Max) GetLimitValue() interface{} {
	return m.Max
}

// Range Requires an integer to be within Min, Max inclusive.
type Range struct {
	Min
	Max
	Key string
}

// IsSatisfied judge whether obj is valid
func (r Range) IsSatisfied(obj interface{}) bool {
	return r.Min.IsSatisfied(obj) && r.Max.IsSatisfied(obj)
}

// DefaultMessage return the default Range error message
func (r Range) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Range"], r.Min.Min, r.Max.Max)
}

// GetKey return the m.Key
func (r Range) GetKey() string {
	return r.Key
}

// GetLimitValue return the limit value, Max
func (r Range) GetLimitValue() interface{} {
	return []int{r.Min.Min, r.Max.Max}
}

// MinSize Requires an array or string to be at least a given length.
type MinSize struct {
	Min int
	Key string
}

// IsSatisfied judge whether obj is valid
func (m MinSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) >= m.Min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

// DefaultMessage return the default MinSize error message
func (m MinSize) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["MinSize"], m.Min)
}

// GetKey return the m.Key
func (m MinSize) GetKey() string {
	return m.Key
}

// GetLimitValue return the limit value
func (m MinSize) GetLimitValue() interface{} {
	return m.Min
}

// MaxSize Requires an array or string to be at most a given length.
type MaxSize struct {
	Max int
	Key string
}

// IsSatisfied judge whether obj is valid
func (m MaxSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) <= m.Max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}

	return false
}

// DefaultMessage return the default MaxSize error message
func (m MaxSize) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["MaxSize"], m.Max)
}

// GetKey return the m.Key
func (m MaxSize) GetKey() string {
	return m.Key
}

// GetLimitValue return the limit value
func (m MaxSize) GetLimitValue() interface{} {
	return m.Max
}

// Length Requires an array or string to be exactly a given length.
type Length struct {
	N   int
	Key string
}

// IsSatisfied judge whether obj is valid
func (l Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) == l.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == l.N
	}
	return false
}

// DefaultMessage return the default Length error message
func (l Length) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Length"], l.N)
}

// GetKey return the m.Key
func (l Length) GetKey() string {
	return l.Key
}

// GetLimitValue return the limit value
func (l Length) GetLimitValue() interface{} {
	return l.N
}

// Alpha check the alpha
type Alpha struct {
	Key string
}

// IsSatisfied judge whether obj is valid
func (a Alpha) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
		return true
	}
	return false
}

// DefaultMessage return the default Length error message
func (a Alpha) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Alpha"])
}

// GetKey return the m.Key
func (a Alpha) GetKey() string {
	return a.Key
}

// GetLimitValue return the limit value
func (a Alpha) GetLimitValue() interface{} {
	return nil
}

// Float check number
type Float struct {
	Key string
}

func (f Float) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		if _, err := strconv.ParseFloat(str, 64); err != nil {
			return false
		}
		return true
	}
	return false
}

func (f Float) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Float"])
}

func (f Float) GetKey() string {
	return f.Key
}

// GetLimitValue return the limit value
func (f Float) GetLimitValue() interface{} {
	return nil
}

// Numeric check number
type Numeric struct {
	Key string
}

// IsSatisfied judge whether obj is valid
func (n Numeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		if nok := strings.HasPrefix(str, "0"); nok {
			if len(str) > 1 {
				return false
			}
		}
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
		return true
	}
	return false
}

// DefaultMessage return the default Length error message
func (n Numeric) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Numeric"])
}

// GetKey return the n.Key
func (n Numeric) GetKey() string {
	return n.Key
}

// GetLimitValue return the limit value
func (n Numeric) GetLimitValue() interface{} {
	return nil
}

// AlphaNumeric check alpha and number
type AlphaNumeric struct {
	Key string
}

// IsSatisfied judge whether obj is valid
func (a AlphaNumeric) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') && ('9' < v || v < '0') {
				return false
			}
		}
		return true
	}
	return false
}

// DefaultMessage return the default Length error message
func (a AlphaNumeric) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["AlphaNumeric"])
}

// GetKey return the a.Key
func (a AlphaNumeric) GetKey() string {
	return a.Key
}

// GetLimitValue return the limit value
func (a AlphaNumeric) GetLimitValue() interface{} {
	return nil
}

// Match Requires a string to match a given regex.
type Match struct {
	Regexp *regexp.Regexp
	Key    string
}

// IsSatisfied judge whether obj is valid
func (m Match) IsSatisfied(obj interface{}) bool {
	return m.Regexp.MatchString(fmt.Sprintf("%v", obj))
}

// DefaultMessage return the default Match error message
func (m Match) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Match"], m.Regexp.String())
}

// GetKey return the m.Key
func (m Match) GetKey() string {
	return m.Key
}

// GetLimitValue return the limit value
func (m Match) GetLimitValue() interface{} {
	return m.Regexp.String()
}

// NoMatch Requires a string to not match a given regex.
type NoMatch struct {
	Match
	Key string
}

// IsSatisfied judge whether obj is valid
func (n NoMatch) IsSatisfied(obj interface{}) bool {
	return !n.Match.IsSatisfied(obj)
}

// DefaultMessage return the default NoMatch error message
func (n NoMatch) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["NoMatch"], n.Regexp.String())
}

// GetKey return the n.Key
func (n NoMatch) GetKey() string {
	return n.Key
}

// GetLimitValue return the limit value
func (n NoMatch) GetLimitValue() interface{} {
	return n.Regexp.String()
}

var alphaDashPattern = regexp.MustCompile("[^\\d\\w-_]")

// AlphaDash check not Alpha
type AlphaDash struct {
	NoMatch
	Key string
}

// DefaultMessage return the default AlphaDash error message
func (a AlphaDash) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["AlphaDash"])
}

// GetKey return the n.Key
func (a AlphaDash) GetKey() string {
	return a.Key
}

// GetLimitValue return the limit value
func (a AlphaDash) GetLimitValue() interface{} {
	return nil
}

var base64Pattern = regexp.MustCompile("^(?:[A-Za-z0-99+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")

// Base64 check struct
type Base64 struct {
	Match
	Key string
}

// DefaultMessage return the default Base64 error message
func (b Base64) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Base64"])
}

// GetKey return the b.Key
func (b Base64) GetKey() string {
	return b.Key
}

// GetLimitValue return the limit value
func (b Base64) GetLimitValue() interface{} {
	return nil
}

type IsDate struct {
	Format string
	Key    string
}

func (d IsDate) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["IsDate"], d.Format)
}

func (d IsDate) GetKey() string {
	return d.Key
}

func (d IsDate) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		_, err := time.Parse(d.Format, str)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

// GetLimitValue return the limit value
func (d IsDate) GetLimitValue() interface{} {
	return nil
}

type DateBefore struct {
	Reference string
	Format    string
	Key       string
}

func (d DateBefore) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["DateBefore"], d.Reference)
}

func (d DateBefore) GetKey() string {
	return d.Key
}

func (d DateBefore) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		s, err := time.Parse(d.Format, d.Reference)
		if err != nil {
			return false
		}
		if e, err := time.Parse(d.Format, str); err != nil {
			return false
		} else if e.Before(s) {
			return false
		}
		return true
	}
	return false
}

func (d DateBefore) GetLimitValue() interface{} {
	return nil
}

type SliceMatch struct {
	Haystack interface{}
	Key      string
}

func (i SliceMatch) DefaultMessage() string {
	var hs []string
	switch reflect.TypeOf(i.Haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i.Haystack)
		for i := 0; i < s.Len(); i++ {
			hs = append(hs, fmt.Sprintf("%v", s.Index(i).Interface()))
		}
	}

	return fmt.Sprintf(MessageTmpls["SliceMatch"], strings.Join(hs, " | "))
}

func (i SliceMatch) GetKey() string {
	return i.Key
}

func (i SliceMatch) IsSatisfied(obj interface{}) bool {
	switch reflect.TypeOf(i.Haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i.Haystack)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(obj, s.Index(i).Interface()) == true {
				return true
			}
		}
	}
	return false
}

func (i SliceMatch) GetLimitValue() interface{} {
	return nil
}

type Duplicate struct {
	Monitor string
	Key     string
}

func (d Duplicate) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Duplicate"])
}

func (d Duplicate) GetKey() string {
	return d.Key
}

func (d Duplicate) IsSatisfied(obj interface{}) bool {
	checkDupe := make(map[interface{}]bool)
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			if !s.Index(i).FieldByName(d.Monitor).IsValid() {
				continue
			}
			ss := s.Index(i).Interface()
			rf := reflect.ValueOf(ss).FieldByName(d.Monitor).Interface()
			if str, ok := rf.(string); ok {
				if len(str) == 0 {
					continue
				}
			}

			if !checkDupe[rf] {
				checkDupe[rf] = true
			} else {
				return false
			}
		}
		return true
	}
	checkDupe = nil
	return false

}

func (d Duplicate) GetLimitValue() interface{} {
	return nil
}

type Incremental struct {
	Key string
}

func (i Incremental) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Incremental"])
}

func (i Incremental) GetKey() string {
	return i.Key
}

func (i Incremental) IsSatisfied(obj interface{}) bool {
	sortutil.AscByField(obj, i.Key)
	t := reflect.TypeOf(obj)
	if t.Kind() != reflect.Slice {
		return false
	}
	d := reflect.ValueOf(obj)
	previousSeq := 0
	for w := 0; w < d.Len(); w++ {
		ss := d.Index(w).Interface()
		rf := reflect.ValueOf(ss).FieldByName(i.Key).Interface()
		k := reflect.TypeOf(rf)
		z2 := 0

		switch k.Kind() {
		case reflect.String:
			z2, _ = strconv.Atoi(rf.(string))
		case reflect.Int:
			z2 = rf.(int)
		}

		z1 := previousSeq
		if z1+1 != z2 {
			return false
		}
		previousSeq += 1
	}
	return true
}

func (i Incremental) GetLimitValue() interface{} {
	return nil
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

type Email struct {
	Match
	Key string
}

func (e Email) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["Email"])
}

func (e Email) GetKey() string {
	return e.Key
}

func (e Email) GetLimitValue() interface{} {
	return nil
}

var positiveFloatPattern = regexp.MustCompile("^[0-9.]")

type PositiveFloat struct {
	Match
	Key string
}

func (b PositiveFloat) DefaultMessage() string {
	return fmt.Sprint(MessageTmpls["PositiveFloat"])
}

func (b PositiveFloat) GetKey() string {
	return b.Key
}

func (b PositiveFloat) GetLimitValue() interface{} {
	return nil
}


type IsPhone struct {
	Match
	Key string
}

var phonePattern = regexp.MustCompile(`^(\+?\d[1-9]\s*-?|\d{2}[1-9])?\s*\d([- ]?\d){4,}$`)

func (p IsPhone) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Phone"])
}

func (p IsPhone) GetKey() string {
	return p.Key
}

func (p IsPhone) GetLimitValue() interface{} {
	return nil
}

type IsName struct {
	Match
	Key string
}

var namePattern = regexp.MustCompile(`^[a-zA-Z ".-]+$`)

func (p IsName) DefaultMessage() string {
	return fmt.Sprintf(MessageTmpls["Name"])
}

func (p IsName) GetKey() string {
	return p.Key
}

func (p IsName) GetLimitValue() interface{} {
	return nil
}