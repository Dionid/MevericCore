package mccommon

type JSONData interface {
	MarshalJSON() ([]byte, error)
}
