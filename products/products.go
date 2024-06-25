package products

// KobbleProducts is the struct that holds the configuration for the Product
type KobbleProducts struct{}

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
