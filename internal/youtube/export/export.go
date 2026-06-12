package export

import (	
	
)

const SchemaVersion = "1.0"




// // This is to filter any field that have asked to be contained.
// //
// // Other fields are removed 
// func filterFields(
// 	data ExportData,
// 	fields []string,
// ) (ExportData, error) {

	 
// 	filtered := make(
// 		ExportData,
// 		len(fields),
// 	)

// 	for _, field := range fields {

// 		value, ok := data[field]
// 		// check if the field exist in data. If not return an Error.
		
// 		if !ok {

// 			return nil,
// 				Err.Export.InvalidField().
// 					AddField(
// 						field,
// 						"field does not exist",
// 					)
// 		}

// 		filtered[field] = value
// 	}

// 	return filtered, nil
// }
