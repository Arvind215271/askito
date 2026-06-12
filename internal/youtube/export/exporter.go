package export

type Exporter interface {
    Export(
        data ExportData,
    ) ([]byte, error)
}