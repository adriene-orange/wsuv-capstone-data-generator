package main

import (
	"encoding/json"
	"fmt"
	"hello/offerings"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker"
	"github.com/gocarina/gocsv"
)

type Record struct {
	ID           string `csv:"id"`
	OfferingKeys string `csv:"offering_keys"`
	Type         string `csv:"type"`
	Tiers        string `csv:"tiers"`
	SupplierKey  string `csv:"supplier_key"`
	ProductKey   string `csv:"product_key"`
}

// Discount ...
type Discount struct {
	OfferingKeys []string
	Type         string
	SupplierKey  string
	ProductKey   string
	Tiers        []discountTierRecord
}

type discountTierRecord struct {
	DiscountPercentage float64
	UOM                string
	MinQty             int
	MaxQty             int
}

func pickDiscountType() string {
	discountTypes := []string{
		"SUPPLIER_DISCOUNT",
		"PRODUCT_DISCOUNT",
		"BULK_DISCOUNT",
	}
	n := rand.Int() % len(discountTypes)
	return discountTypes[n]
}

func getTiersByDiscountType(offering *offerings.Record) discountTierRecord {
	initialMax := 9999
	return discountTierRecord{
		DiscountPercentage: rand.Float64(),
		MaxQty:             rand.Intn(initialMax) + 1,
		MinQty:             1,
		UOM:                offering.UnitOfMeasure,
	}
}

func generateBulkDiscountTiers(offering *offerings.Record) []discountTierRecord {
	initialDiscount := rand.Float64()
	initialMin := 1
	initialMax := 50
	bulkTiers := []discountTierRecord{}
	for i := 0; i < 4; i++ {
		bulkTiers = append(bulkTiers, discountTierRecord{
			DiscountPercentage: (initialDiscount * float64(10*i)),
			MaxQty:             initialMax * (i + 1),
			MinQty:             initialMin * (i + 50),
			UOM:                offering.UnitOfMeasure,
		})
	}
	return bulkTiers
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// supplier.CreateSupplierCSV()
	// products.CreateProductsCSV()
	offeringsFile, err := os.OpenFile("offerings.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer offeringsFile.Close()
	offerings := []*offerings.Record{}

	if err := gocsv.UnmarshalFile(offeringsFile, &offerings); err != nil { // Load clients from file
		panic(err)
	}
	discounts := []Record{}
	currentDiscounts := make(map[string]*Discount)
	for _, offering := range offerings {
		// should handle supplier discounts
		_, supplierDiscount := currentDiscounts[offering.SupplierKey]
		newRecord := false
		giveDiscount := rand.Intn(2)
		if supplierDiscount {
			supplierDiscount := currentDiscounts[offering.SupplierKey]
			currentDiscounts[offering.SupplierKey].OfferingKeys = append(supplierDiscount.OfferingKeys, offering.ID)
			currentDiscounts[offering.SupplierKey].Tiers = append(supplierDiscount.Tiers, getTiersByDiscountType(offering))
		} else if giveDiscount > 0 {
			newRecord = true
		}
		if newRecord {
			discountType := pickDiscountType()
			if discountType == "SUPPLIER_DISCOUNT" {
				currentDiscounts[offering.SupplierKey] = &Discount{
					OfferingKeys: []string{offering.ID},
					Tiers:        []discountTierRecord{getTiersByDiscountType(offering)},
					Type:         discountType,
					ProductKey:   "",
					SupplierKey:  offering.SupplierKey,
				}
			} else if discountType == "BULK_DISCOUNT" {
				currentDiscounts[offering.ProductKey] = &Discount{
					OfferingKeys: []string{offering.ID},
					Tiers:        generateBulkDiscountTiers(offering),
					Type:         discountType,
					ProductKey:   offering.ProductKey,
					SupplierKey:  "",
				}
			} else {
				currentDiscounts[offering.ProductKey] = &Discount{
					OfferingKeys: []string{offering.ID},
					Tiers:        []discountTierRecord{getTiersByDiscountType(offering)},
					Type:         discountType,
					ProductKey:   offering.ProductKey,
					SupplierKey:  "",
				}
			}
		}
	}

	for _, currentDiscount := range currentDiscounts {
		serializedOfferingKeys, _ := json.Marshal(currentDiscount.OfferingKeys)
		serializedTiers, _ := json.Marshal(currentDiscount.Tiers)
		discountRecord := Record{
			ID:           faker.UUIDDigit(),
			OfferingKeys: string(serializedOfferingKeys),
			Tiers:        string(serializedTiers),
			Type:         currentDiscount.Type,
			ProductKey:   currentDiscount.ProductKey,
			SupplierKey:  currentDiscount.SupplierKey,
		}
		fmt.Printf("%+v", discountRecord)
		discounts = append(discounts, discountRecord)
	}

	file, err := os.Create("discounts.csv")
	defer file.Close()
	err = gocsv.MarshalFile(&discounts, file) // Use this to save the CSV back to the file
	if err != nil {
		panic(err)
	}
}
