package mcdevicerpcmanager

//import "mevericcore/mccommon"
//
//type DeviceCreator func() mccommon.DeviceBaseModelInterface
//
//type DevicesModelAndCollectionSt struct {
//	DeviceCreator DeviceCreator
//	CollectionManager mccommon.DevicesCollectionManagerInterface
//}
//
//type DeviceModelsAndCollectionsManagerSt struct {
//	DevicesByTypeName map[string]*DevicesModelAndCollectionSt
//}
//
//func CreateNewDeviceModelsAndCollectionsManager() *DeviceModelsAndCollectionsManagerSt {
//	return &DeviceModelsAndCollectionsManagerSt{
//		map[string]*DevicesModelAndCollectionSt{},
//	}
//}
//
//func (this *DeviceModelsAndCollectionsManagerSt) RegisterNewDeviceType(deviceType string, deviceModel DeviceCreator, collectionManager mccommon.DevicesCollectionManagerInterface) {
//	this.DevicesByTypeName[deviceType] = &DevicesModelAndCollectionSt{
//		deviceModel,
//		collectionManager,
//	}
//}
//
//func (this *DeviceModelsAndCollectionsManagerSt) GetDeviceModelByDeviceType(deviceType string) (deviceModel mccommon.DeviceBaseModelInterface, collectionManager mccommon.DevicesCollectionManagerInterface) {
//	if devAndCol, ok := this.DevicesByTypeName[deviceType]; ok {
//		return devAndCol.DeviceCreator(), devAndCol.CollectionManager
//	}
//	return nil, nil
//}