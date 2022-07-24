package matcher

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type matcher func(interface{}, interface{}) (bool, error)

//nolint:gochecknoglobals // an internal global here is more efficient than repeatedly creating the map in a hot path
var matchers map[reflect.Kind]matcher

//nolint:gochecknoinits // the assigned functions refer to `matchers` so we can't assign it directly: we need init()
func init() {
	matchers = map[reflect.Kind]matcher{
		reflect.Bool:    _matchPrimitive,
		reflect.String:  _matchPrimitive,
		reflect.Int:     _matchPrimitive,
		reflect.Int64:   _matchPrimitive,
		reflect.Float64: _matchPrimitive,
		reflect.Map:     _matchMap,
		reflect.Slice:   _matchSlice,
		reflect.Array:   _matchSlice,
	}
}

// JSONMatches checks if the JSON in `j` provided with the first argument
// satisfies the pattern in the second argument.
// Both `j` and `jSpec` are passed as byte slices.
// The pattern can be a valid literal value (in that case an exact match will
// be required), a special marker (a string starting with the hash character
// '#'), or any combination of these via arrays and objects.
func JSONMatches(j []byte, jSpec []byte) (bool, error) {
	var jAny interface{}
	err := json.Unmarshal(j, &jAny)
	if err != nil {
		return false, fmt.Errorf("can't unmarshal left argument: %w", err)
	}

	var specAny interface{}
	err = json.Unmarshal(jSpec, &specAny)
	if err != nil {
		return false, fmt.Errorf("can't unmarshal specifier argument: %w", err)
	}

	return _match(jAny, specAny)
}

// JSONStringMatches checks if the JSON string `j` provided with the first argument
// satisfies the pattern in the second argument.
// Both `j` and `jSpec` are passed as strings.
// The pattern can be a valid literal value (in that case an exact match will
// be required), a special marker (a string starting with the hash character
// '#'), or any combination of these via arrays and objects.
func JSONStringMatches(j string, jSpec string) (bool, error) {
	return JSONMatches([]byte(j), []byte(jSpec))
}

func _matchZero(x interface{}) (bool, error) {
	xV := reflect.ValueOf(x)
	if !xV.IsValid() {
		return true, nil
	}
	return false, nil
}

var uuidRe = regexp.MustCompile(`(?i)^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)
var uuidV4Re = regexp.MustCompile(`(?i)^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}$`)

const (
	ignoreMarker  = "#ignore"
	nullMarker    = "#null"
	presentMarker = "#present"
)

//nolint:funlen,gocognit // reducing the number of statements would reduce legibility in this instance
func _matchWithSpecifier(x interface{}, spec string) (bool, error) {
	if x == nil && (spec == ignoreMarker || spec == nullMarker || spec == presentMarker) {
		return true, nil
	}
	xV := reflect.ValueOf(x)
	if !xV.IsValid() {
		return false, nil // here we now that ref is non-zero
	}

	//nolint:gomnd // the "magic" literal constant 2 here is clearer than a synthetic constant symbol
	specParts := strings.SplitN(spec, " ", 2)

	switch specParts[0] {
	case ignoreMarker:
		return true, nil
	case nullMarker:
		if xV.Kind() != reflect.Ptr {
			return false, nil
		}
		return xV.IsNil(), nil
	case "#notnull":
		if xV.Kind() != reflect.Ptr {
			return true, nil
		}
		return !xV.IsNil(), nil
	case presentMarker:
		return true, nil
	case "#notpresent":
		return false, nil
	case "#array":
		if (xV.Kind() != reflect.Array) && (xV.Kind() != reflect.Slice) {
			return false, nil
		}
		return true, nil
	case "#object":
		if xV.Kind() != reflect.Map {
			return false, nil
		}
		return true, nil
	case "#bool":
		fallthrough
	case "#boolean":
		if xV.Kind() != reflect.Bool {
			return false, nil
		}
		return true, nil
	case "#number":
		if (xV.Kind() != reflect.Int64) && (xV.Kind() != reflect.Float64) {
			return false, nil
		}
		return true, nil
	case "#string":
		if xV.Kind() != reflect.String {
			return false, nil
		}
		return true, nil
	case "#date":
		if xV.Kind() == reflect.String {
			xString, ok := x.(string)
			if !ok {
				return false, nil
			}
			_, err := time.Parse("2006-01-02", xString)
			if err == nil {
				return true, nil
			}
			return false, nil
		} else if xV.Kind() == reflect.Struct {
			_, ok := x.(time.Time)
			if ok {
				return true, nil
			}
		}
		return false, nil

	case "#datetime":
		if xV.Kind() == reflect.String {
			xString, ok := x.(string)
			if !ok {
				return false, nil
			}
			_, err := time.Parse(time.RFC3339, xString)
			if err == nil {
				return true, nil
			}
			return false, nil
		} else if xV.Kind() == reflect.Struct {
			_, ok := x.(time.Time)
			if ok {
				return true, nil
			}
		}
		return false, nil
	case "#uuid":
		if xV.Kind() != reflect.String {
			return false, nil
		}
		xString, ok := x.(string)
		if ok {
			return uuidRe.MatchString(xString), nil
		}
		return false, nil
	case "#uuid-v4":
		if xV.Kind() != reflect.String {
			return false, nil
		}
		xString, ok := x.(string)
		if ok {
			return uuidV4Re.MatchString(xString), nil
		}
		return false, nil
	case "#regex":
		//nolint:gomnd // the "magic" literal constant 2 here is clearer than a synthetic constant symbol
		if len(specParts) != 2 {
			return false, fmt.Errorf("expected exactly one argument for #regex")
		}
		r, err := regexp.Compile(specParts[1])
		if err != nil {
			return false, fmt.Errorf("invalid regex argument to #regex: %w", err)
		}
		if xV.Kind() != reflect.String {
			return false, nil
		}
		xString, ok := x.(string)
		if ok {
			return r.MatchString(xString), nil
		}
		return false, nil
		// TODO: "#[num] EXPR"
	}

	return false, fmt.Errorf("unsupported specifier '%s'", spec)
}

func _match(x interface{}, spec interface{}) (bool, error) {
	specV := reflect.ValueOf(spec)
	if !specV.IsValid() {
		return _matchZero(x)
	}

	if specV.Kind() == reflect.String {
		isMarker, specMarker := getMarker(spec)
		if isMarker {
			return _matchWithSpecifier(x, specMarker)
		}
	}

	xV := reflect.ValueOf(x)
	if !xV.IsValid() {
		return false, nil // here we now that spec is non-zero
	}

	if xV.Kind() != specV.Kind() {
		return false, nil
	}

	if m, ok := matchers[specV.Kind()]; ok {
		return m(x, spec)
	}
	tX := reflect.TypeOf(x)
	return false, fmt.Errorf("unable to compare %v (type: %v) - kind %v is not supported", x, tX, xV.Kind())
}

func _matchMap(x interface{}, y interface{}) (bool, error) {
	vX := reflect.ValueOf(x)
	if vX.Kind() != reflect.Map {
		return false, fmt.Errorf("wrong kind for left value, expected Map, got %v", vX.Kind())
	}
	if reflect.ValueOf(y).Kind() != reflect.Map {
		return false, fmt.Errorf("wrong kind for specifier value, expected Map, got %v", vX.Kind())
	}

	vY := reflect.ValueOf(y)

	matches, err := _matchMapCheckIteratingObject(vX, vY)
	if err != nil {
		return false, err
	}

	if !matches {
		return false, nil
	}

	matches, err = _matchMapCheckIteratingSpec(vX, vY)
	if err != nil {
		return false, err
	}

	return matches, nil
}

func _matchMapCheckIteratingObject(vX reflect.Value, vY reflect.Value) (bool, error) {
	matches := true
	iterX := vX.MapRange()
	for iterX.Next() {
		ySpecValue := vY.MapIndex(iterX.Key())
		if !ySpecValue.IsValid() {
			// missing spec for this key, skip...
			continue
		}
		itemMatches, err := _match(iterX.Value().Interface(), ySpecValue.Interface())
		if err != nil {
			return false, fmt.Errorf("can't compare map element %v: %w", iterX.Key().Interface(), err)
		}
		matches = matches && itemMatches
	}
	return matches, nil
}

func _matchMapCheckIteratingSpec(vX reflect.Value, vY reflect.Value) (bool, error) {
	matches := true
	iterY := vY.MapRange()
	for iterY.Next() {
		ySpecValue := iterY.Value()
		xValue := vX.MapIndex(iterY.Key())

		//nolint:wastedassign // defensive programming here...
		itemMatches := false
		if ySpecValue.Kind() == reflect.Interface {
			isMarker, marker := getMarker(ySpecValue.Interface())
			if isMarker {
				switch marker {
				case "#notpresent":
					itemMatches = !xValue.IsValid()
					matches = matches && itemMatches
					continue
				case presentMarker:
					itemMatches = xValue.IsValid()
					matches = matches && itemMatches
					continue
				}
			}
		}

		if !isMarker(ySpecValue.Interface(), ignoreMarker) {
			if !xValue.IsValid() {
				matches = false
			} else {
				var err error
				itemMatches, err = _match(xValue.Interface(), ySpecValue.Interface())
				if err != nil {
					return false, fmt.Errorf("can't compare map element %v: %w", iterY.Key().Interface(), err)
				}
				matches = matches && itemMatches
			}
		}
	}
	return matches, nil
}

func getMarker(y interface{}) (bool, string) {
	specString, ok := y.(string)
	if ok && strings.HasPrefix(specString, "#") {
		return true, specString
	}
	return false, ""
}

func isMarker(y interface{}, marker string) bool {
	isMarker, gotMarker := getMarker(y)
	return isMarker && marker == gotMarker
}

func _matchSlice(x interface{}, y interface{}) (bool, error) {
	vX := reflect.ValueOf(x)
	if vX.Kind() != reflect.Slice {
		return false, fmt.Errorf("wrong kind for left value, expected Slice, got %v", vX.Kind())
	}
	if reflect.ValueOf(y).Kind() != reflect.Slice {
		return false, fmt.Errorf("wrong kind for specifier value, expected Slice, got %v", vX.Kind())
	}

	vY := reflect.ValueOf(y)
	isArrayOf := false
	var arrayOf interface{}

	//nolint:gomnd // the "magic" literal constant 2 here is clearer than a synthetic constant symbol
	if vY.Len() == 2 {
		first := vY.Index(0).Interface()
		isMarker, firstMarker := getMarker(first)
		if isMarker && firstMarker == "#array-of" {
			isArrayOf = true
			arrayOf = vY.Index(1).Interface()
		}
	} else if vX.Len() != vY.Len() {
		return false, nil
	}

	matches := true
	sliceLen := vX.Len()
	for i := 0; i < sliceLen; i++ {
		var ySpecElem interface{}
		if isArrayOf {
			ySpecElem = arrayOf
		} else {
			ySpecElem = vY.Index(i).Interface()
		}
		itemMatches, err := _match(vX.Index(i).Interface(), ySpecElem)
		if err != nil {
			return false, fmt.Errorf("can't compare slice element %v: %w", i, err)
		}
		matches = matches && itemMatches
	}
	return matches, nil
}

func _matchPrimitive(x interface{}, y interface{}) (bool, error) {
	return reflect.DeepEqual(x, y), nil
}
