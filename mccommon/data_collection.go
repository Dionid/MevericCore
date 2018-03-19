package mccommon

import "mevericcore/mcmongo"

type DataCollectionManagerSt struct {
	mcmongo.CollectionManagerBaseSt
}

type DataCollectionManagerInt interface {
	mcmongo.CollectionManagerBaseInterface
}