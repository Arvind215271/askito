package fields

func DefaultFields() []string {
    out := make([]string, 0)

    out = append(out, FieldID)
    out = append(out, MetadataFields...)
    out = append(out, DescriptionFields...)
    out = append(out, TranscriptFields...)
    out = append(out, SignalFields...)
    // optionally:

    return out
}