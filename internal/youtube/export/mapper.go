package export

import (
	"reflect"
)

// It convert any data passed to it to exported Data Type.
//
// It would deal with Playlist and Videos only for current purpose.
// we convert this data to JSON. So we can simply use this data to convert to other format we might need in future.
func structToExportData(v any, fields map[string]bool) (ExportData, error) {
	//  this creates a map of string key and value as any
	out := make(ExportData)

	// reflect value

	// this gives us the actual value stored in the binary. Because to create our struct, GO allocate Bytes in sizes. And we use that to store that. So this is where we can check what is stored in each field of the struct. (actual data) 
	val := reflect.ValueOf(v)
	// here, gives us the struct definitoin. Like which field exist, what is the type of that field, etc. So we can simply use that to extract the field that we are required or given in the map to this function and filter those out... based on this typ.
	typ := reflect.TypeOf(v)

	// handle pointer safety
	// sometimes we mmight be passed a pointer. So instead of using a pointer, we get the actual data out of it... using Elem() function 
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	full := len(fields) == 0
	// if no field is set... Then simply return everything.

	// for each field present in struct... do a for loop .
	for i := 0; i < typ.NumField(); i++ {
		// get the field from the struct definition
		structField := typ.Field(i)

		// use json tag as export key (this is same as we write when defining structs)
		key := structField.Tag.Get("json")
		if key == "" {
			// if no value et. Use the field name itself
			key = structField.Name
		}

		// check if present or all keys have to be returned
		if full || fields[key] {
			// firstly, val store the actual data. Then we need to know which field to get the data from. Interface convert it to any data type. That is all
			out[key] = val.Field(i).Interface()
		}
	}

	return out, nil
}








// // this convert the Video data type to Exportable Data type. 
// func videoToExport(v Video, fields map[string]bool) ExportData {
// 	// this creates a map of string key and value as any
// 	out := make(ExportData)

// 	// if no field is set... Then simply return everything.
// 	full := len(fields) == 0

// 	// reflect value

// 	// this gives us the actual value stored in the binary. Because to create our struct, GO allocate Bytes in sizes.
// 	// So this is where we can check what is stored in each field of the struct. (actual data)
// 	val := reflect.ValueOf(v)

// 	// here gives us the struct definition. Like which field exist, what is the type of that field, etc.
// 	// So we can use that to extract only the fields that are required and filter them out.
// 	typ := reflect.TypeOf(v)

// 	// handle pointer safety
// 	// sometimes we might be passed a pointer. So instead of using a pointer, we get the actual data out of it using Elem()
// 	if val.Kind() == reflect.Ptr {
// 		val = val.Elem()
// 		typ = typ.Elem()
// 	}

// 	// for each field present in struct... do a loop
// 	for i := 0; i < typ.NumField(); i++ {

// 		// get the field definition from the struct
// 		sf := typ.Field(i)

// 		// use json tag as export key (same as struct definition)
// 		key := sf.Tag.Get("json")

// 		// if no json tag is present, fallback to struct field name
// 		if key == "" {
// 			key = sf.Name
// 		}

// 		// check if field is allowed or if we are exporting everything
// 		if full || fields[key] {

// 			// val holds actual data, so we extract the field value by index
// 			// Interface() converts it into empty interface (any type)
// 			out[key] = val.Field(i).Interface()
// 		}
// 	}

// 	return out
// }