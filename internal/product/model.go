package product

import (
	"math/rand/v2"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]"`
}

func NewProduct() *Product {
	return &Product{}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range len(b) {
		b[i] = letterRunes[rand.IntN(len(letterRunes))]
	}
	return string(b)
}
