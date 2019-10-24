package offerings

import (
	"fmt"
	"hello/products"
	supplier "hello/suppliers"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/bxcodec/faker"
	"github.com/gocarina/gocsv"
)

// FakeOffering ...
type FakeOffering struct {
	ID             string `faker:"uuid_digit"`
	CreatedDate    int64  `faker:"unix_time"`
	ActiveDate     int64  `faker:"unix_time"`
	ExpirationDate int64  `faker:"unix_time"`
}

type Record struct {
	ID             string  `csv:"id"`
	CreatedDate    int64   `csv:"created_date"`
	ActiveDate     int64   `csv:"active_date"`
	ExpirationDate int64   `csv:"expiration_date"`
	UnitRetail     float64 `csv:"unit_retail"`
	UnitCost       float64 `csv:"unit_cost"`
	UnitOfMeasure  string  `csv:"uom"`
	ProductKey     string  `csv:"product_key"`
	SupplierKey    string  `csv:"supplier_key"`
}

// CustomGenerator ...
func CustomGenerator() {
	faker.AddProvider("productName", func(v reflect.Value) (interface{}, error) {
		return "Useful " + strings.Title(faker.Word()) + " Tool", nil
	})
}

func generateCost(retail float64) float64 {
	n := rand.Float64()
	return ((retail - (retail * n)) * 100) / 100
}

func generateRetail() float64 {
	max := rand.Intn(500) + 1
	n := rand.Float64()
	return ((float64(max) + n) * 100) / 100
}

func pickUnitOfMeasure() string {
	units := []string{
		"EA.",
		"Box",
		"Case",
		"Kit",
		"Pallet",
	}
	n := rand.Int() % len(units)
	return units[n]
}

func pickRandomProductID(products []*products.Record) string {
	n := rand.Int() % len(products)
	return products[n].ID
}

func generateOfferingCSV() {
	rand.Seed(time.Now().UnixNano())
	// supplier.CreateSupplierCSV()
	// products.CreateProductsCSV()
	suppliersFile, err := os.OpenFile("suppliers.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	productsFile, err := os.OpenFile("products.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer suppliersFile.Close()
	defer productsFile.Close()
	suppliers := []*supplier.Record{}
	products := []*products.Record{}

	if err := gocsv.UnmarshalFile(suppliersFile, &suppliers); err != nil { // Load clients from file
		panic(err)
	}

	if err := gocsv.UnmarshalFile(productsFile, &products); err != nil { // Load clients from file
		panic(err)
	}
	CustomGenerator()
	offerings := []Record{}
	rand.Seed(time.Now().Unix())
	for _, supplier := range suppliers {
		n := (rand.Int() % 30) + 1
		for i := 0; i < n; i++ {
			fakeOffering := FakeOffering{}
			faker.FakeData(&fakeOffering)
			retail := generateRetail()
			offering := Record{
				ID:             fakeOffering.ID,
				ActiveDate:     fakeOffering.ActiveDate,
				CreatedDate:    fakeOffering.CreatedDate,
				ExpirationDate: fakeOffering.ExpirationDate,
				UnitRetail:     retail,
				UnitCost:       generateCost(retail),
				UnitOfMeasure:  pickUnitOfMeasure(),
				ProductKey:     pickRandomProductID(products),
				SupplierKey:    supplier.ID,
			}
			fmt.Printf("%+v", offering)
			offerings = append(offerings, offering)

		}
	}

	file, err := os.Create("offerings.csv")
	defer file.Close()
	err = gocsv.MarshalFile(&offerings, file) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}
