package util

import (
	"math"
	"strings"
	"time"

	"github.com/infobloxopen/atlas-app-toolkit/query"
)

const (
	fieldNamedListType = "type"
	fieldItems         = "items"
	filtDelim          = "."
)

type FilteringIteratorCallback func(f interface{}) interface{}

func IterateFiltering(f *query.Filtering, callback FilteringIteratorCallback) {
	doCallback := func(fi interface{}) interface{} {
		return callback(fi)
	}

	var getOperator func(interface{}) interface{}

	getOperator = func(f interface{}) interface{} {
		val := f.(*query.LogicalOperator)

		left := val.GetLeft()
		switch leftVal := left.(type) {
		case *query.LogicalOperator_LeftOperator:
			val.SetLeft(getOperator(leftVal.LeftOperator))

		case *query.LogicalOperator_LeftStringCondition:
			val.SetLeft(doCallback(leftVal.LeftStringCondition))

		case *query.LogicalOperator_LeftNumberCondition:
			val.SetLeft(doCallback(leftVal.LeftNumberCondition))

		case *query.LogicalOperator_LeftNullCondition:
			val.SetLeft(doCallback(leftVal.LeftNullCondition))

		case *query.LogicalOperator_LeftStringArrayCondition:
			val.SetLeft(doCallback(leftVal.LeftStringArrayCondition))
		}

		right := val.GetRight()
		switch rightVal := right.(type) {
		case *query.LogicalOperator_RightOperator:
			val.SetRight(getOperator(rightVal.RightOperator))

		case *query.LogicalOperator_RightStringCondition:
			val.SetRight(doCallback(rightVal.RightStringCondition))

		case *query.LogicalOperator_RightNumberCondition:
			val.SetRight(doCallback(rightVal.RightNumberCondition))

		case *query.LogicalOperator_RightNullCondition:
			val.SetRight(doCallback(rightVal.RightNullCondition))

		case *query.LogicalOperator_RightStringArrayCondition:
			val.SetLeft(doCallback(rightVal.RightStringArrayCondition))
		}
		return val
	}

	root := f.GetRoot()
	switch val := root.(type) {
	case *query.Filtering_Operator:
		f.SetRoot(getOperator(val.Operator))

	case *query.Filtering_StringCondition:
		f.SetRoot(doCallback(val.StringCondition))

	case *query.Filtering_NumberCondition:
		f.SetRoot(doCallback(val.NumberCondition))

	case *query.Filtering_NullCondition:
		f.SetRoot(doCallback(val.NullCondition))
	case *query.Filtering_StringArrayCondition:
		f.SetRoot(doCallback(val.StringArrayCondition))
	}
}

type ValidateCallback func(fieldName string, cond interface{}) error

type FilterDescription struct {
	IsRequired bool
	Validate   ValidateCallback
}

func FilterValidateInt(field string, cond interface{}) error {
	cnd, ok := cond.(*query.NumberCondition)
	if !ok || math.Trunc(cnd.Value) != cnd.Value {
		return InvalidArgumentErr("Invalid filter `%s` value, suported only integers", field)
	}
	return nil
}

func FilterValidateBool(field string, cond interface{}) error {
	cnd, ok := cond.(*query.StringCondition)

	if !ok || cnd.Type != query.StringCondition_EQ || (cnd.Value != "true" && cnd.Value != "false") {
		return InvalidArgumentErr("Invalid filter `%s` value, suported only 'true' and 'false' values", field)
	}
	return nil
}

func FilterValidateTimestamp(field string, cond interface{}) error {
	cnd, ok := cond.(*query.StringCondition)

	if !ok || (cnd.Type != query.StringCondition_EQ &&
		cnd.Type != query.StringCondition_GE &&
		cnd.Type != query.StringCondition_GT &&
		cnd.Type != query.StringCondition_LE &&
		cnd.Type != query.StringCondition_LT) {
		return InvalidArgumentErr("Invalid condition %s for filter field `%s`", cnd.Type, field)
	}
	if !isISOTimestamp(cnd.Value) {
		return InvalidArgumentErr("Invalid filter `%s` value, supported only iso timestamp", field)
	}
	return nil
}

func isISOTimestamp(ts string) bool {
	_, err := time.Parse(time.RFC3339, ts)
	return err == nil
}

// VerifyOrderFields verifies if sorting fields specified exists in set of valid fields
func VerifyOrderFields(f *query.Sorting, validFields map[string]struct{}) error {
	for _, criteria := range f.GetCriterias() {
		field := criteria.GetTag()
		if _, ok := validFields[field]; !ok {
			return InvalidArgumentErr("Invalid order_by field: '%s'", field)
		}
	}
	return nil
}

func RenameOrderFields(f *query.Sorting, alias map[string]string) {
	for _, criteria := range f.GetCriterias() {
		field := criteria.GetTag()
		if name, ok := alias[field]; ok {
			criteria.Tag = name
		}
	}
}

func VerifyFilterFields(f *query.Filtering, m map[string]*FilterDescription) (err error) {
	receivedFields := map[string]struct{}{}

	checkExpectedFields := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}

		var field string

		switch cond.(type) {
		case *query.StringCondition:
			sc, _ := cond.(*query.StringCondition)
			field = strings.Join(sc.FieldPath, filtDelim)
		case *query.NumberCondition:
			nc, _ := cond.(*query.NumberCondition)
			field = strings.Join(nc.FieldPath, filtDelim)
		case *query.NullCondition:
			nc, _ := cond.(*query.NullCondition)
			field = strings.Join(nc.FieldPath, filtDelim)
		case *query.StringArrayCondition:
			nc, _ := cond.(*query.StringArrayCondition)
			field = strings.Join(nc.FieldPath, filtDelim)
		}

		params, ok := m[field]
		if !ok {
			err = InvalidArgumentErr("Invalid filter field: %s", field)
		}
		if params != nil && params.Validate != nil {
			err = params.Validate(field, cond)
		}
		receivedFields[field] = struct{}{}
		return cond
	}

	IterateFiltering(f, checkExpectedFields)

	for field, params := range m {
		if params != nil && params.IsRequired {
			if _, ok := receivedFields[field]; !ok {
				err = InvalidArgumentErr("Filtering by field '%s' is required", field)
				break
			}
		}
	}

	return err
}

// RenameFilterFields .. remaps API-level fields to SQL-level fields
// This function is useful to avoid exposing internal names to the API
func RenameFilterFields(f *query.Filtering, alias map[string]string) {
	if nil == alias {
		return
	}
	IterateFiltering(f, func(cond interface{}) interface{} {
		switch cond.(type) {
		case *query.StringCondition:
			sc, _ := cond.(*query.StringCondition)
			field := strings.Join(sc.FieldPath, filtDelim)
			if name, ok := alias[field]; ok {
				sc.FieldPath = strings.Split(name, filtDelim)
			}
		case *query.NumberCondition:
			nc, _ := cond.(*query.NumberCondition)
			field := strings.Join(nc.FieldPath, filtDelim)
			if name, ok := alias[field]; ok {
				nc.FieldPath = strings.Split(name, filtDelim)
			}
		case *query.NullCondition:
			nc, _ := cond.(*query.NullCondition)
			field := strings.Join(nc.FieldPath, filtDelim)
			if name, ok := alias[field]; ok {
				nc.FieldPath = strings.Split(name, filtDelim)
			}
		case *query.StringArrayCondition:
			nc, _ := cond.(*query.StringArrayCondition)
			field := strings.Join(nc.FieldPath, filtDelim)
			if name, ok := alias[field]; ok {
				nc.FieldPath = strings.Split(name, filtDelim)
			}
		}
		return cond
	})
}

func FilterExactNumberCondition(f *query.Filtering) (field string, value int, err error) {
	checkNumberCondition := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}

		if len(field) > 0 {
			err = InvalidArgumentErr("Found more than one condition for filtering")
			return cond
		}

		nc, ok := cond.(*query.NumberCondition)
		if ok && !nc.IsNegative && nc.Type == query.NumberCondition_EQ {
			field = strings.Join(nc.FieldPath, ".")
			value = int(nc.Value)
		} else if ok {
			err = InvalidArgumentErr("Only '==' operation is supported for filtering")
		} else {
			err = InvalidArgumentErr("Found non-number condition for filtering")
		}
		return cond
	}

	IterateFiltering(f, checkNumberCondition)
	return field, value, err
}

func FilterExactNullCondition(f *query.Filtering) (field string, err error) {
	checkNullCondition := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}

		if len(field) > 0 {
			err = InvalidArgumentErr("Found more than one condition for filtering")
			return cond
		}

		nc, ok := cond.(*query.NullCondition)
		if ok && !nc.IsNegative {
			field = strings.Join(nc.FieldPath, ".")
		} else if ok {
			err = InvalidArgumentErr("Only '==' operation is supported for filtering")
		} else {
			err = InvalidArgumentErr("Found non-null condition for filtering")
		}
		return cond
	}
	IterateFiltering(f, checkNullCondition)
	return field, err
}

func FilterExactStringCondition(f *query.Filtering) (field string, value string, err error) {
	checkStringCondition := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}

		if len(field) > 0 {
			err = InvalidArgumentErr("Found more than one condition for filtering")
			return cond
		}

		nc, ok := cond.(*query.StringCondition)
		if ok && !nc.IsNegative && nc.Type == query.StringCondition_EQ {
			field = strings.Join(nc.FieldPath, ".")
			value = string(nc.Value)
		} else if ok {
			err = InvalidArgumentErr("Only '==' operation is supported for filtering")
		} else {
			err = InvalidArgumentErr("Found non-string condition for filtering")
		}
		return cond
	}
	IterateFiltering(f, checkStringCondition)
	return field, value, err
}

func FilterStringArrayCondition(f *query.Filtering) (field string, values []string, err error) {
	checkStringArrayCondition := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}

		if len(field) > 0 {
			err = InvalidArgumentErr("Found more than one condition for filtering")
			return cond
		}

		nc, ok := cond.(*query.StringArrayCondition)
		if ok && !nc.IsNegative {
			field = strings.Join(nc.FieldPath, ".")
			values = nc.Values
		} else if !ok {
			err = InvalidArgumentErr("Found non-array string condition for filtering")
		} else {
			err = InvalidArgumentErr("Found negative array string condition for filtering")
		}
		return cond
	}
	IterateFiltering(f, checkStringArrayCondition)
	return field, values, err
}

type ListTypeConversion func(string) (int, error)

func FilterNamedListType(f *query.Filtering, listNameToValue ListTypeConversion) (err error) {
	checkConditions := func(cond interface{}) interface{} {

		if err != nil {
			return cond
		}
		field := ""

		switch cnd := cond.(type) {

		case *query.StringCondition:
			field = strings.Join(cnd.FieldPath, ".")

			if field == fieldNamedListType {
				if cnd.Type != query.StringCondition_EQ {
					err = InvalidArgumentErr("Only '!=' and '==' condition operations  available for `%s` filter", fieldNamedListType)
					return cond
				}
				var typeVal int
				typeVal, err = listNameToValue(cnd.Value)
				if err != nil {
					err = InvalidArgumentErr("Invalid filter `%s` value", fieldNamedListType)
					return cond
				}

				return &query.NumberCondition{
					FieldPath:  []string{"list_type"},
					Value:      float64(typeVal),
					Type:       query.NumberCondition_Type(cnd.Type),
					IsNegative: cnd.IsNegative,
				}
			}
		case *query.NumberCondition:
			field = strings.Join(cnd.FieldPath, ".")
		case *query.StringArrayCondition:
			field = strings.Join(cnd.FieldPath, ".")
		case *query.NumberArrayCondition:
			field = strings.Join(cnd.FieldPath, ".")
		case *query.NullCondition:
			field = strings.Join(cnd.FieldPath, ".")
		}
		if field == fieldNamedListType {
			err = InvalidArgumentErr("Invalid filtering `%s` value, supported only string type", fieldNamedListType)
			return cond
		}

		return cond
	}

	IterateFiltering(f, checkConditions)
	return err
}

func GetItemsCondition(f *query.Filtering, operation query.StringCondition_Type) (field string, value string, isNegative bool, err error) {
	isNonItemsConditionFound := false

	checkConditions := func(cond interface{}) interface{} {
		if err != nil {
			return cond
		}
		fName := ""
		switch cnd := cond.(type) {
		case *query.StringCondition:
			f := strings.Join(cnd.FieldPath, ".")
			if f == fieldItems {
				field = f
				value = cnd.Value
				isNegative = cnd.IsNegative
				if cnd.Type != operation {
					err = InvalidArgumentErr("Filtering by `%s` only %s operation support", fieldItems, query.StringCondition_Type_name[int32(operation)])
				}
				return cond
			} else {
				isNonItemsConditionFound = true
			}
		case *query.NumberCondition:
			fName = strings.Join(cnd.FieldPath, ".")
			isNonItemsConditionFound = true
		case *query.StringArrayCondition:
			fName = strings.Join(cnd.FieldPath, ".")
			isNonItemsConditionFound = true
		case *query.NumberArrayCondition:
			fName = strings.Join(cnd.FieldPath, ".")
			isNonItemsConditionFound = true
		case *query.NullCondition:
			fName = strings.Join(cnd.FieldPath, ".")
			isNonItemsConditionFound = true
		default:
			isNonItemsConditionFound = true
		}
		if fName == fieldItems {
			err = InvalidArgumentErr("Invalid filtering `%s` value, supported only string type", fieldItems)
			return cond
		}

		return cond
	}

	IterateFiltering(f, checkConditions)
	if isNonItemsConditionFound && field != "" {
		err = InvalidArgumentErr("Filtering by `items` can't be used with additional filtering options")
	}
	return
}

func FilterPolicyIdSupport(f *query.Filtering) (err error) {

	checkConditions := func(cond interface{}) interface{} {

		if err != nil {
			return cond
		}

		switch cnd := cond.(type) {
		// handle policy_id filter
		case *query.NumberCondition:
			field := strings.Join(cnd.FieldPath, ".")
			if field == "policy_id" {
				cnd.FieldPath[0] = "security_policy_id"
			}

		}
		return cond
	}

	IterateFiltering(f, checkConditions)
	return err
}

type DefaultIdCallback func() (int, error)

func FilterBoolToDefaultId(f *query.Filtering, getIdFunc DefaultIdCallback, boolFieldName string, defaultIdFieldName string) (err error) {
	var defaultId int

	checkConditions := func(cond interface{}) interface{} {

		if err != nil {
			return cond
		}

		switch cnd := cond.(type) {
		// handle default_security_policy filter
		case *query.StringCondition:

			field := strings.Join(cnd.FieldPath, ".")
			if field == boolFieldName {

				if defaultId == 0 {
					defaultId, err = getIdFunc()
					if err != nil {
						return cond
					}
				}

				if cnd.Value == "true" || cnd.Value == "false" {
					nct := query.NumberCondition_EQ
					if cnd.Type != query.StringCondition_EQ {
						err = InvalidArgumentErr("Only '!=' and '==' condition operations  available for `%s` filter", boolFieldName)
						return cond
					}

					if cnd.Value == "false" {
						cnd.IsNegative = !cnd.IsNegative
					}
					return &query.NumberCondition{
						FieldPath:  []string{defaultIdFieldName},
						Value:      float64(defaultId),
						Type:       nct,
						IsNegative: cnd.IsNegative,
					}
				}
				err = InvalidArgumentErr("Only 'true' or 'false' boolean condition values available for `%s` filter", boolFieldName)
			}
		}
		return cond
	}

	IterateFiltering(f, checkConditions)
	return err
}

func FilterValidateMdmType(field string, cond interface{}) error {
	var validMdmTypes = []string{"generic", "android"}
	cnd, ok := cond.(*query.StringCondition)

	if !ok || cnd.Type != query.StringCondition_EQ || !InStringSlice(validMdmTypes, cnd.Value) {
		return InvalidArgumentErr("Invalid filter `%s` values can be `%v`", field, validMdmTypes)
	}
	return nil
}
