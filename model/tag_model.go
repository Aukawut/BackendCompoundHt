package model

type Tags struct {
	TagNo       string `json:"tagNo"`
	PartNo      string `json:"partNo"`
	Qty         any    `json:"qty"`
	Batch       int    `json:"batch"`
	Basket      int    `json:"basket"`
	Lot         string `json:"lotNo"`
	OriginalTag string `json:"originalTag"`
}

type SaveTagsBody struct {
	Tags     []Tags `json:"tags"`
	ActionBy string `json:"actionBy"`
	Location string `json:"location"`
}

type IdLists struct {
	Id int `json:"id"`
}

type TagsCancelBody struct {
	Ids      []IdLists `json:"ids"`
	ActionBy string    `json:"actionBy"`
}

type TagsList struct {
	Id         int
	OLD_TAG    string
	TAG_NO     string
	PART_NO    string
	QTY        float64
	BATCH      int
	BASKET     int
	PO_NO      string
	LOT_NO     string
	LOCATION   string
	ACTIVE     string
	CREATED_AT string
	CREATED_BY string
	UPDATED_AT any
	UPDATED_BY any
}
