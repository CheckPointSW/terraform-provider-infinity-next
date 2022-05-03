package utils

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SlicesDiff excepts to slices which represents old and new values
// THe function returns:
// 1. a slice of all values that exist in the new and not in the old slice (a.k.a "added")
// 2. a slice of all values that exist in the old and not in the new slice (a.k.a "removed")
func SlicesDiff[T any](oldValues, newValues []T) ([]T, []T) {
	newValuesMap := SliceToStringIndicatorsMap(newValues)
	oldValuesMap := SliceToStringIndicatorsMap(oldValues)

	var removed []T
	for k, v := range oldValuesMap {
		if _, ok := newValuesMap[k]; !ok {
			removed = append(removed, v)
		}
	}

	var added []T
	for k, v := range newValuesMap {
		if _, ok := oldValuesMap[k]; !ok {
			added = append(added, v)
		}
	}

	return added, removed
}

func Removed(oldValues, newValues []string) []string {
	var removed []string

	newValuesMap := make(map[string]struct{}, len(newValues))
	for _, v := range newValues {
		newValuesMap[v] = struct{}{}
	}

	oldValuesMap := make(map[string]struct{}, len(oldValues))
	for _, v := range oldValues {
		oldValuesMap[v] = struct{}{}
	}

	for k := range oldValuesMap {
		if _, ok := newValuesMap[k]; !ok {
			removed = append(removed, k)
		}
	}

	return removed
}

func Added(oldValues, newValues []string) []string {
	var added []string

	newValuesMap := make(map[string]struct{}, len(newValues))
	for _, v := range newValues {
		newValuesMap[v] = struct{}{}
	}

	oldValuesMap := make(map[string]struct{}, len(oldValues))
	for _, v := range oldValues {
		oldValuesMap[v] = struct{}{}
	}

	for k := range newValuesMap {
		if _, ok := oldValuesMap[k]; !ok {
			added = append(added, k)
		}
	}

	return added
}

// DiagError collects errors and appends diagnostics with diagnostic of err
func DiagError(summary string, err error, diags diag.Diagnostics) diag.Diagnostics {
	err = fmt.Errorf("%s: %w", summary, err)
	diags = append(diags, diag.FromErr(err)...)
	return diags
}

// SupressDiffIfExists shows diff only if resource is new
func SupressDiffIfExists(k, old, new string, d *schema.ResourceData) bool {
	return d.Id() != ""
}

// Filter takes 2 arguments:
// 1. A slice of generic type T
// 2. A function that excepts an argument of generic type T and returns a boolean
// The Filter function than applies the filterFunc to each element of the slice
// and returns a sub-slice which contains all element x where filterFunc(x) == true
func Filter[T any](slice []T, filterFunc func(T) bool) []T {
	ret := make([]T, 0, len(slice))
	for _, val := range slice {
		if filterFunc(val) {
			ret = append(ret, val)
		}
	}

	return ret
}

// Map takes 2 arguments:
// 1. A slice of generic type T
// 2. A function that excepts an argument of generic type T and returns a value of generic type U
// The Map function than applies the mapFunc to each element of the slice
// and returns a slice of the retruned values from the mapFunc
func Map[T, U any](slice []T, mapFunc func(T) U) []U {
	ret := make([]U, len(slice))
	for i, val := range slice {
		ret[i] = mapFunc(val)
	}

	return ret
}

// SliceToStringIndicatorsMap excepts a slice of generice type T
// The function returns a map from a string representation of each value in the slice
// to the value itself
// for example: [1, 2, 3] ==> {"1": 1, "2": 2, "3": 3}
func SliceToStringIndicatorsMap[T any](slice []T) map[string]T {
	ret := make(map[string]T)
	for _, val := range slice {
		key := fmt.Sprintf("%+v", val)
		ret[key] = val
	}

	return ret
}

// MustValueAs converts a the given argument to the desired type
// note that this conversion is not safe and may panic
func MustValueAs[T any](val any) T {
	return val.(T)
}

// MustSliceAs takes a slice as an argument and converts all it's members to the desired type T
func MustSliceAs[T any](slice []any) []T {
	return Map(slice, MustValueAs[T])
}

// MustUnmarshalAs marshals the argument and than unmarshals it to the desired type
// ignores error if marshal/unmarshal fails (prints the error and returns the empty value of the desired type)
func MustUnmarshalAs[T, U any](toConvert U) T {
	ret, err := UnmarshalAs[T](toConvert)
	if err != nil {
		fmt.Printf("%+v", err)
	}

	return ret
}

// UnmarshalAs marshals the argument and than unmarshals it to the desired type
// returns error upon failure
func UnmarshalAs[T, U any](toConvert U) (T, error) {
	var ret T
	toConvertBytes, err := json.Marshal(toConvert)
	if err != nil {
		return ret, fmt.Errorf("failed to marshal argument. Error: %+v. Input: %+v", err, toConvert)
	}

	if err := json.Unmarshal(toConvertBytes, &ret); err != nil {
		return ret, fmt.Errorf("failed to unmarshal argument to desired type. Error: %+v. Input: %+v", err, toConvert)
	}

	return ret, nil
}

// MustSchemaCollectionToSlice converts the argument to a slice of the desired type
// The argument must be of type schema.TypeSet or schema.TypeList
// otherwise, the function will panic
func MustSchemaCollectionToSlice[T any](listOrSet any) []T {
	var anySlice []any
	if anySet, isSet := listOrSet.(*schema.Set); isSet {
		anySlice = anySet.List()
	} else {
		anySlice = listOrSet.([]any)
	}

	return MustSliceAs[T](anySlice)
}

// MustResourceDataCollectionToSlice takes a ResourceData object and a fieldName
// the function tries to retrieve the field's value and expects it to be of type
// schema.TypeSet or schema.TypeList (otherwise it will panic!)
// this function then converts the set/list to a slice of the desired type
func MustResourceDataCollectionToSlice[T any](d *schema.ResourceData, fieldName string) []T {
	interfaceSetOrList, ok := d.GetOk(fieldName)
	if !ok {
		return nil
	}

	retSlice := MustSchemaCollectionToSlice[T](interfaceSetOrList)
	if len(retSlice) == 0 {
		return nil
	}

	return retSlice
}

// GetChangeWithParse checks if there has been any changes made to the value saved under the given key in the resourceData object
// if so, the functon applies the given parse function to the old and new values and returns the parsing results
// the function also returns a boolean which indicates weather a change has been made or not
func GetChangeWithParse[T any](d *schema.ResourceData, key string, parseFunc func(any) T) (T, T, bool) {
	if !d.HasChange(key) {
		var emptyRet T
		return emptyRet, emptyRet, false
	}

	oldVal, newVal := d.GetChange(key)
	return parseFunc(oldVal), parseFunc(newVal), true

}

// GetChangeWithParse checks if there has been any changes made to the value saved under the given key in the resourceData object
// if so, the functon converts the the old and new values to the desired type
// note that if the conversion fails the function will panic!
func MustGetChange[T any](d *schema.ResourceData, key string) (T, T, bool) {
	if !d.HasChange(key) {
		var emptyRet T
		return emptyRet, emptyRet, false
	}

	oldVal, newVal := d.GetChange(key)
	return oldVal.(T), newVal.(T), true
}
