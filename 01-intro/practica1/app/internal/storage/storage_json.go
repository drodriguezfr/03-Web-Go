package storage

// StorageDefault is the default storage implementation
type StorageJSON struct {
	// filePath is the path to the json file
	filePath string
	// layoutDate is the layout of the expiration date
	//layoutDate string
}

// ProductAttributesJSON is a struct that contains the information of a product
type ProductAttributesJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func NewStorageJSON(filePath string) *StorageJSON {
	// default config
	//defaultLayoutDate := time.DateOnly
	//if layoutDate != "" {
	//	defaultLayoutDate = layoutDate
	//}
	return &StorageJSON{
		filePath: filePath,
		//layoutDate: defaultLayoutDate,
	}
}
