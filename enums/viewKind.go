package enums

import "errors"

type ViewKind string

const (
	ViewForest   ViewKind = "FOREST"
	ViewSea      ViewKind = "SEA"
	ViewLake     ViewKind = "LAKE"
	ViewMountain ViewKind = "MOUNTAIN"
	ViewOther    ViewKind = "OTHER"
)

func (viewkind *ViewKind) String() error {
	value := *viewkind

	switch *viewkind {
	case ViewForest, ViewSea, ViewLake, ViewMountain, ViewOther:
		*viewkind = value
		return nil
	}

	return errors.New("invalid product type value")
}
