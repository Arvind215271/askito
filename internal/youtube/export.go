// ./internal/youtube/export.go

package youtube

import (
	"encoding/json"
	"reflect"
)

const SchemaVersion = "1.0"



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

// this convert the Video data type to Exportable Data type. 
func videoToExport(v Video, fields map[string]bool) ExportData {
	// this creates a map of string key and value as any
	out := make(ExportData)

	// if no field is set... Then simply return everything.
	full := len(fields) == 0

	// reflect value

	// this gives us the actual value stored in the binary. Because to create our struct, GO allocate Bytes in sizes.
	// So this is where we can check what is stored in each field of the struct. (actual data)
	val := reflect.ValueOf(v)

	// here gives us the struct definition. Like which field exist, what is the type of that field, etc.
	// So we can use that to extract only the fields that are required and filter them out.
	typ := reflect.TypeOf(v)

	// handle pointer safety
	// sometimes we might be passed a pointer. So instead of using a pointer, we get the actual data out of it using Elem()
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// for each field present in struct... do a loop
	for i := 0; i < typ.NumField(); i++ {

		// get the field definition from the struct
		sf := typ.Field(i)

		// use json tag as export key (same as struct definition)
		key := sf.Tag.Get("json")

		// if no json tag is present, fallback to struct field name
		if key == "" {
			key = sf.Name
		}

		// check if field is allowed or if we are exporting everything
		if full || fields[key] {

			// val holds actual data, so we extract the field value by index
			// Interface() converts it into empty interface (any type)
			out[key] = val.Field(i).Interface()
		}
	}

	return out
}

// This is to filter any field that have asked to be contained.
//
// Other fields are removed 
func filterFields(
	data ExportData,
	fields []string,
) (ExportData, error) {

	 
	filtered := make(
		ExportData,
		len(fields),
	)

	for _, field := range fields {

		value, ok := data[field]
		// check if the field exist in data. If not return an Error.
		
		if !ok {

			return nil,
				Err.Export.InvalidField().
					AddField(
						field,
						"field does not exist",
					)
		}

		filtered[field] = value
	}

	return filtered, nil
}

// this is common function that we will be using to export playlist and convert it to a simple format that can be used by an export TYPE like JSON, CSV, etc.
//
// It is the filter layer actually. We are already getting the data in our Domain Model.
// The only thing left is to filter what is needed from Video ONLY.
func BuildPlaylistExport(
	playlist Playlist,
	videoFields []string,
) (ExportData, error) {

	fieldSet := make(map[string]bool, len(videoFields))
	for _, f := range videoFields {
		fieldSet[f] = true
	}

	videos := make([]any, 0, len(playlist.Videos))

	for _, v := range playlist.Videos {

		// ONLY Video struct is filtered
		videoData, err := structToExportData(v.Video, fieldSet)
		if err != nil {
			return nil, Err.Export.MarshalFailed().Wrap(err)
		}

		// PlaylistVideo metadata is ALWAYS preserved
		videoData["position"] = v.Position
		videoData["added_at"] = v.AddedAt

		videos = append(videos, videoData)
	}

	// Playlist itself is NOT filtered
	return ExportData{
		"id":             playlist.ID,
		"title":          playlist.Title,
		"description":    playlist.Description,
		"channel_id":     playlist.ChannelID,
		"channel_title":  playlist.ChannelTitle,
		"thumbnail_url":  playlist.ThumbnailURL,
		"item_count":     playlist.ItemCount,
		"privacy_status": playlist.PrivacyStatus,
		"published_at":   playlist.PublishedAt,

		"videos": videos,
	}, nil
}

// this is common function that we will be using to export video and convert it to a simple format that can be used by an export TYPE like JSON, CSV, etc.
//
// It is the filter layer actually. We are already getting the data in our Domain Model.
// The only thing left is to filter what is needed from Video ONLY.
func BuildVideoExport(
	video Video,
	fields []string,
) (ExportData, error) {

	fieldSet := make(map[string]bool, len(fields))
	for _, f := range fields {
		fieldSet[f] = true
	}

	// ONLY Video is filterable
	data, err := structToExportData(video, fieldSet)
	if err != nil {
		return nil, Err.Export.MarshalFailed().Wrap(err)
	}

	return data, nil
}

// to export the ExportData as JSON format to the user.
func ExportJSON(data ExportData) ([]byte, error) {

	resp := ExportResponse{
		SchemaVersion: SchemaVersion,
		Format:        FormatJSON,
		Data:          data,
	}

	b, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		return nil, Err.Export.MarshalFailed().Wrap(err)
	}

	return b, nil
}