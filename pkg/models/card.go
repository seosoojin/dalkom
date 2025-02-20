package models

type Card struct {
	ID           string   `json:"id" bson:"_id"`
	Name         string   `json:"name" bson:"name" indexed:"true"`
	ShortName    string   `json:"short_name" bson:"short_name"`
	Type         CardType `json:"type" bson:"type" indexed:"true"`
	ImageUrl     string   `json:"image_url" bson:"image_url"`
	GroupID      string   `json:"group_id" bson:"group_id" indexed:"true"`
	CollectionID string   `json:"collection_id" bson:"collection_id" indexed:"true"`
	IdolIDs      []string `json:"idol_ids" bson:"idol_ids" indexed:"true"`
}

type EnrichedCard struct {
	ID         string     `json:"id" bson:"_id"`
	Name       string     `json:"name" bson:"name"`
	ShortName  string     `json:"short_name" bson:"short_name"`
	Type       CardType   `json:"type" bson:"type"`
	FmtType    string     `json:"fmt_type" bson:"fmt_type"`
	ImageUrl   string     `json:"image_url" bson:"image_url"`
	Group      Group      `json:"group" bson:"group"`
	Collection Collection `json:"collection" bson:"collection"`
	Idols      []Idol     `json:"idols" bson:"idols"`
}

type CardType string

const (
	CardTypeRegular   CardType = "regular"
	CardTypePOB       CardType = "pob"
	CardTypeEvent     CardType = "event"
	CardTypeSpecial   CardType = "special"
	CardTypeTrading   CardType = "trading"
	CardTypeLimited   CardType = "limited"
	CardTypeMerch     CardType = "merch"
	CardTypeLuckyDraw CardType = "lucky_draw"
)

var ShortTypesMap = map[CardType]string{
	CardTypeRegular:   "R",
	CardTypePOB:       "POB",
	CardTypeEvent:     "Event",
	CardTypeSpecial:   "S",
	CardTypeTrading:   "TC",
	CardTypeLimited:   "Lim",
	CardTypeMerch:     "Merch",
	CardTypeLuckyDraw: "LD",
}
