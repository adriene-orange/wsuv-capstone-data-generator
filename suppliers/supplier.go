package supplier

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bxcodec/faker"
	"github.com/gocarina/gocsv"
)

// FakeSupplier ...
type FakeSupplier struct {
	ID          string `faker:"uuid_digit"`
	CreatedDate int64  `faker:"unix_time"`
	ActiveDate  int64  `faker:"unix_time"`
}

type Record struct {
	ID          string `csv:"id"`
	Name        string `csv:"supplier_name"`
	CreatedDate int64  `csv:"created_date"`
	ActiveDate  int64  `csv:"active_date"`
}

func customSupplierGenerator(currentWords map[string]bool) string {
	tag := []string{
		"CO.",
		"Inc.",
		"Supplier",
		"& Sons",
		"Company",
		"Incorporate",
	}
	n := rand.Int() % len(tag)
	word := faker.Word()
	for currentWords[word] == true {
		word = faker.Word()
	}
	currentWords[word] = true
	return strings.Title(word) + strings.Title(faker.Word()) + " " + tag[n]
}

func CreateSupplierCSV() {
	rand.Seed(time.Now().UnixNano())
	file, err := os.Create("suppliers.csv")
	defer file.Close()
	suppliers := []Record{}
	currentWords := make(map[string]bool)
	for i := 0; i < 100; i++ {
		fakeSupplier := FakeSupplier{}
		faker.FakeData(&fakeSupplier)
		supplier := Record{
			ID:          fakeSupplier.ID,
			Name:        customSupplierGenerator(currentWords),
			ActiveDate:  fakeSupplier.ActiveDate,
			CreatedDate: fakeSupplier.CreatedDate,
		}
		fmt.Printf("%+v", supplier)
		suppliers = append(suppliers, supplier)
	}
	err = gocsv.MarshalFile(&suppliers, file) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}
