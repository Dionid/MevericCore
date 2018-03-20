package mcdevicerpcmanager

import "mevericcore/mccommon"

type ProtocolManagerInterface interface {
	SendJSON(string, mccommon.JSONData) error
}
