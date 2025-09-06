package domain

type Stock struct {
	SkuID      int64  `json:"sku"`
	TotalCount uint32 `json:"total_count"`
	Reserved   uint32 `json:"reserved"`
}
