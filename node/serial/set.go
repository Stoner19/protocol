package serial

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/Oneledger/protocol/node/log"
)

// Alloc a new variable of any type
func Alloc(dataType string, size int) interface{} {
	if dataType == "" {
		return nil
	}

	if size == -1 {
		return nil
	}

	entry := GetTypeEntry(dataType, size)

	var value reflect.Value

	switch entry.Category {
	case UNKNOWN:
		DumpTypes()
		log.Fatal("Unknown datatype", "dataType", dataType)

	case INTERFACE:
		// Need to handle this in set
		return nil

	case PRIMITIVE:
		// Don't need to alloc, only containers
		return nil

	case STRUCT:
		value = reflect.New(entry.DataType)
		if !value.IsValid() {
			log.Warn("Retrying as slice?")
			value = reflect.MakeSlice(entry.DataType, size, size)
			if !value.IsValid() {
				log.Warn("Retrying as byte array")
				value = reflect.ValueOf(make([]byte, size))
			}
		}
		return value.Interface()

	case MAP:
		smap := reflect.MakeMapWithSize(entry.RootType, size)
		value = reflect.New(smap.Type())
		value.Elem().Set(smap)
		return value.Interface()

	case SLICE:
		slice := reflect.MakeSlice(entry.DataType, size, size)
		value = reflect.New(slice.Type())
		value.Elem().Set(slice)
		return value.Interface()

	case ARRAY:
		//array := reflect.ArrayOf(size, entry.ValueType.DataType)
		//value = reflect.New(array)
		value = reflect.New(entry.RootType)
		return value.Interface()
	}

	log.Warn("Unknown Category", "dataType", dataType, "entry", entry)
	return nil
}

func concat(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

// Strip off the depth
func GetFieldName(name string) string {
	array := strings.Split(name, "-")
	result := concat(array[0 : len(array)-1]...)
	return result
}

// Slice and Array names are currently strings (associative arrays)
func GetFieldIndex(name string) int {
	array := strings.Split(name, "-")
	interim := concat(array[0 : len(array)-1]...)
	index, err := strconv.Atoi(interim)
	if err != nil {
		log.Fatal("Invalid Slice Index", "name", name, "interim", interim)
	}
	return index
}

// Set a structure with a given value, convert as necessary
func Set(parent interface{}, fieldName string, child interface{}) (status bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("Ignoring Set Panic", "r", r)
			log.Dump("Parameters", parent, fieldName, child)
			debug.PrintStack()
		}
	}()

	kind := GetBaseValue(parent).Kind()

	switch kind {

	case reflect.Struct:
		return SetStruct(parent, GetFieldName(fieldName), child)

	case reflect.Map:
		return SetMap(parent, GetFieldName(fieldName), child)

	case reflect.Slice:
		return SetSlice(parent, GetFieldIndex(fieldName), child)

	case reflect.Array:
		return SetArray(parent, GetFieldIndex(fieldName), child)
	}
	return false
}

// SetStructure takes a pointer to a structure and sets it.
func SetStruct(parent interface{}, fieldName string, child interface{}) bool {
	if child == nil {
		return false
	}

	// Convert the interfaces to structures
	//element := reflect.ValueOf(parent).Elem()
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Struct {
		log.Fatal("Not a structure", "element", element, "kind", element.Kind())
	}

	if !CheckValue(element, fieldName) {
		return false
	}

	field := element.FieldByName(fieldName)

	if !CheckValue(field, fieldName) {
		return false
	}

	if field.Type().Kind() == reflect.Interface {
		// When setting to a generic interface{}
		value := reflect.ValueOf(child)
		field.Set(value)

	} else {
		newValue := ConvertValue(child, field.Type())
		if newValue.Kind() == reflect.Int {
			field.SetInt(newValue.Int())
		} else {
			field.Set(newValue)
		}
	}
	return true
}

// SetMap takes a pointer to a structure and sets it.
func SetMap(parent interface{}, fieldName string, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Map {
		log.Fatal("Not a map", "element", element)
	}

	if !CheckValue(element, fieldName) {
		return false
	}

	//key := reflect.ValueOf(fieldName)
	keyType := GetBaseType(parent).Key()
	newKey := ConvertValue(fieldName, keyType)

	fieldType := GetBaseType(parent).Elem()
	newValue := ConvertValue(child, fieldType)

	if element.Type().Kind() == reflect.Interface {
		element.SetMapIndex(newKey, newValue)
	} else {
		element.SetMapIndex(newKey, newValue)
	}

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetSlice(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Slice {
		log.Fatal("Not a slice", "element", element)
	}

	fieldName := fmt.Sprintf("%s@%d", reflect.TypeOf(parent).Name, index)
	if !CheckValue(element, fieldName) {
		return false
	}

	if element.Index(index).Type().Kind() == reflect.Interface {
		// When setting to a generic interface{}

		value := reflect.ValueOf(child)
		element.Index(index).Set(value)

	} else {
		newValue := ConvertValue(child, element.Index(index).Type())
		element.Index(index).Set(newValue)
	}

	return true
}

// SetSlice takes a pointer to a structure and sets it.
func SetArray(parent interface{}, index int, child interface{}) bool {

	// Convert the interfaces to structures
	element := GetBaseValue(parent)

	if element.Kind() != reflect.Array {
		log.Fatal("Not a structure", "element", element)
	}

	fieldName := fmt.Sprintf("%s@%d", reflect.TypeOf(parent).Name, index)
	if !CheckValue(element, fieldName) {
		return false
	}

	cell := element.Index(index)
	newValue := ConvertValue(child, cell.Type())

	if element.Index(index).Type().Kind() == reflect.Interface {
		element.Index(index).Set(newValue)
	} else {
		element.Index(index).Set(newValue)
	}

	return true
}

func CheckValue(element reflect.Value, fieldName string) bool {
	if !element.IsValid() {
		log.Warn("Map is invalid (not writable)", "fieldName", fieldName, "element", element)
		debug.PrintStack()
		return false
	}

	if !element.CanSet() {
		log.Warn("Element not Settable", "element", element)
		return false
	}
	return true
}

// Convert any value to an arbitrary type
func ConvertValue(value interface{}, fieldType reflect.Type) reflect.Value {
	if value == nil {
		return reflect.ValueOf(nil)
	}

	typeOf := reflect.TypeOf(value)
	valueOf := reflect.ValueOf(value)

	// Remove any pointers, if they exist
	if typeOf.Kind() == reflect.Ptr {
		valueOf = valueOf.Elem()
		typeOf = reflect.TypeOf(valueOf)
	}

	switch typeOf.Kind() {

	case reflect.Float64:
		// JSON returns float64s for everything :-(
		result := ConvertNumber(fieldType, valueOf)
		return result

	case reflect.String:
		// TODO: Is this still necessary? Should be underlying type is an int?
		if fieldType.String() == "data.ChainType" {
			entry := GetTypeEntry(fieldType.String(), 1)
			interim, err := strconv.ParseInt(GetString(valueOf), 10, 0)
			if err != nil {
				log.Fatal("Failed to convert int")
			}
			valueof := reflect.New(entry.RootType)
			element := valueof.Elem()
			element.SetInt(interim)
			return element
		}
		if fieldType.Kind() == reflect.Int {
			var result int
			interim, err := strconv.ParseInt(GetString(valueOf), 10, 0)
			if err != nil {
				log.Fatal("Failed to convert int")
			}
			result = int(interim)
			return reflect.ValueOf(result)
		}
	}
	return valueOf
}

func GetString(value reflect.Value) string {
	if value.Kind() == reflect.String {
		return value.String()
	}
	return value.Elem().String()
}

// ConvertNumber handles JSON numbers as they are float64 from the parser
func ConvertNumber(fieldType reflect.Type, value reflect.Value) reflect.Value {

	// TODO: find a better way of handling types that are not structures
	/*
		if fieldType.String() == "data.ChainType" {
			entry := GetTypeEntry(fieldType.String(), 1)
			valueof := reflect.New(entry.DataType)
			var result int64
			result = int64(value.Float())
			element := valueof.Elem()
			element.SetInt(result)
			return element
		}
	*/
	/*
		if fieldType.String() == "action.Type" {
			entry := GetTypeEntry(fieldType.String(), 1)
			valueof := reflect.New(entry.DataType)

			var result int64
			result = int64(value.Float())

			element := valueof.Elem()
			element.SetInt(result)
			return element
		}
	*/

	// TODO: shouldn't be manaually creating big ints
	/*
		if fieldType.String() == "*big.Int" {
			log.Debug("Converting pointer to big.Int", "value", value, "fieldType", fieldType)
			converted := big.NewInt(int64(value.Float()))
			return reflect.ValueOf(converted)
		}
	*/

	switch fieldType.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(value.Float()))

	case reflect.Int8:
		return reflect.ValueOf(int8(value.Float()))

	case reflect.Int16:
		return reflect.ValueOf(int16(value.Float()))

	case reflect.Int32:
		return reflect.ValueOf(int32(value.Float()))

	case reflect.Int64:
		return reflect.ValueOf(int64(value.Float()))

	case reflect.Uint:
		return reflect.ValueOf(uint(value.Float()))

	case reflect.Uint8:
		return reflect.ValueOf(uint8(value.Float()))

	case reflect.Uint16:
		return reflect.ValueOf(uint16(value.Float()))

	case reflect.Uint32:
		return reflect.ValueOf(uint32(value.Float()))

	case reflect.Uint64:
		return reflect.ValueOf(uint64(value.Float()))

	case reflect.Float32:
		return reflect.ValueOf(float32(value.Float()))
	}
	return value
}
