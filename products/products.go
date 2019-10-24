package products

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/bxcodec/faker"
	"github.com/gocarina/gocsv"
)

// FakeProduct ...
type FakeProduct struct {
	ID              string `faker:"uuid_digit"`
	LongDescription string `faker:"paragraph"`
	CreatedDate     int64  `faker:"unix_time"`
	ActiveDate      int64  `faker:"unix_time"`
}

type Record struct {
	ID              string `csv:"id"`
	Name            string `csv:"product_name"`
	LongDescription string `csv:"long_description"`
	CreatedDate     int64  `csv:"created_date"`
	ActiveDate      int64  `csv:"active_date"`
}

func customProductGenerator(currentWords map[string]bool) string {
	descriptor := []string{
		"Useful",
		"Heavy",
		"Large",
		"Multiple",
		"Steel",
		"Wooden",
		"Small",
	}
	object := []string{
		"Tool",
		"Bucket",
		"Case",
		"Kit",
		"Siding",
		"Fencing",
		"Box",
	}
	n := rand.Int() % len(descriptor)
	word := faker.Word()
	for currentWords[word] == true {
		word = faker.Word()
	}
	currentWords[word] = true
	return descriptor[n] + " " + strings.Title(word) + " " + object[n]
}

func CreateProductsCSV() {
	file, err := os.Create("products.csv")
	defer file.Close()
	products := []Record{}
	currentWords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		fakeProduct := FakeProduct{}
		faker.FakeData(&fakeProduct)
		product := Record{
			ID:              fakeProduct.ID,
			Name:            customProductGenerator(currentWords),
			LongDescription: fakeProduct.LongDescription,
			ActiveDate:      fakeProduct.ActiveDate,
			CreatedDate:     fakeProduct.CreatedDate,
		}
		fmt.Printf("%+v", product)
		products = append(products, product)
	}
	err = gocsv.MarshalFile(&products, file) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}
