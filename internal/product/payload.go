package product

type ProductCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ProductUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
