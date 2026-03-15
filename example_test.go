package shak_test

import (
	"fmt"

	shak "github.com/yashx/shak"
	"github.com/yashx/shak/rule"
	"github.com/yashx/shak/validation"
)

func ExampleValidate() {
	age := 150
	err := shak.Validate(age, rule.Min(0), rule.Max(120))
	fmt.Println(err)
	// Output:
	// value 150 is greater than maximum 120
}

func ExampleValidateAll() {
	score := 150
	errs := shak.ValidateAll(score, rule.Min(200), rule.Max(100))
	for _, err := range errs {
		fmt.Println(err)
	}
	// Output:
	// value 150 is less than minimum 200
	// value 150 is greater than maximum 100
}

func ExampleRunValidation() {
	a := Address{Street: "123", City: ""}
	err := shak.RunValidation(a.Validation())
	fmt.Println(err)
	// Output:
	// street: length is 3 bytes, expected at least 5
}

func ExampleRunFullValidation() {
	a := Address{Street: "123", City: ""}
	errs := shak.RunFullValidation(a.Validation())
	for _, err := range errs {
		fmt.Println(err)
	}
	// Output:
	// street: length is 3 bytes, expected at least 5
	// city: value must not be blank
}

func ExampleSatisfies() {
	age := 10
	fmt.Println(shak.Satisfies(age, rule.Min(0), rule.Max(120)))
	// Output:
	// true
}

func ExampleSatisfiesValidation() {
	a := Address{Street: "1 Main St", City: "Austin"}
	fmt.Println(shak.SatisfiesValidation(a.Validation()))
	// Output:
	// true
}

func ExampleRunValidation_nested() {
	o := BasicOrder{ID: 1, Address: Address{Street: "", City: "Austin"}}
	err := shak.RunValidation(o.Validation())
	fmt.Println(err)
	// Output:
	// address.street: value must not be blank
}

func ExampleRunValidation_nestedOuterRule() {
	o := ServicedOrder{ID: 1, Address: Address{Street: "1 Oak Ave", City: "Miami"}}
	err := shak.RunValidation(o.Validation())
	fmt.Println(err)
	// Output:
	// address: delivery not available in "Miami"
}

func ExampleRunValidation_errorFields() {
	a := Address{Street: "", City: "Austin"}
	err := shak.RunValidation(a.Validation())
	fmt.Println(err)
	fmt.Println(err.Path)
	fmt.Println(err.Message)
	// Output:
	// street: value must not be blank
	// street
	// value must not be blank
}

func ExampleRunValidation_errorPath() {
	w := Warehouse{
		Name: "Main",
		Sections: []map[string]int{
			{"bolts": 100, "nuts": 50},
			{"bolts": -3, "nuts": 200},
		},
	}
	err := shak.RunValidation(w.Validation())
	fmt.Println(err)
	// Output:
	// sections[1][bolts]: value -3 is less than minimum 0
}

func ExampleValidate_forEachMapValue() {
	scores := map[string]int{"alice": 95, "bob": -1}
	err := shak.Validate(scores, rule.ForEachValue(&scores, rule.Min(0)))
	fmt.Println(err)
	// Output:
	// [bob]: value -1 is less than minimum 0
}

func ExampleValidate_forEachSlice() {
	prices := []float64{9.99, -1.00, 49.99}
	err := shak.Validate(prices, rule.ForEach(&prices, rule.Positive[float64]()))
	fmt.Println(err)
	// Output:
	// [1]: value -1 is not positive
}

func ExampleValidate_unique() {
	tags := []string{"go", "go", "validation"}
	err := shak.Validate(tags, rule.Unique(&tags))
	fmt.Println(err)
	// Output:
	// duplicate value go found
}

func ExampleValidate_isValidSlice() {
	products := []Product{
		{Name: "bolt", Price: 1.99},
		{Name: "", Price: -5},
	}
	err := shak.Validate(products, rule.ForEach(&products, rule.IsValid[Product]()))
	fmt.Println(err)
	// Output:
	// [1].name: value must not be blank
}

func ExampleValidate_isValidMap() {
	catalog := map[string]Product{
		"A1": {Name: "bolt", Price: 1.99},
		"A2": {Name: "", Price: -5},
	}
	err := shak.Validate(catalog, rule.ForEachValue(&catalog, rule.IsValid[Product]()))
	fmt.Println(err)
	// Output:
	// [A2].name: value must not be blank
}

func ExampleValidate_customRule() {
	err := shak.Validate(10, Divisible(3))
	fmt.Println(err)
	// Output:
	// 10 is not divisible by 3
}

func ExampleValidate_and() {
	validName := rule.And(rule.NotBlank[string](), rule.MinLengthString[string](3))
	err := shak.Validate("  ", validName)
	fmt.Println(err)
	// Output:
	// value must not be blank
}

func ExampleValidate_or() {
	validInput := rule.Or(rule.LengthString[string](3), rule.MinLengthString[string](20))
	err := shak.Validate("hello", validInput)
	fmt.Println(err)
	// Output:
	// length is 5 bytes, expected at least 20
}

func ExampleValidate_withMessage() {
	err := shak.Validate("", rule.NotBlank[string]().WithMessage("username is required"))
	fmt.Println(err)
	// Output:
	// username is required
}

func ExampleRunValidation_withMessageOnRange() {
	type Person struct {
		Age int
	}
	personValidation := func(p Person) validation.Validation {
		return validation.NewValidations(
			validation.Value("age", p.Age,
				rule.Range(0, 120).WithMessage("age must be between 0 and 120"),
			),
		)
	}
	err := shak.RunValidation(personValidation(Person{Age: 150}))
	fmt.Println(err)
	// Output:
	// age: age must be between 0 and 120
}

func ExampleRunValidation_ifCondition() {
	o := GiftOrder{IsGift: true, GiftMessage: ""}
	err := shak.RunValidation(o.Validation())
	fmt.Println(err)
	// Output:
	// giftMessage: value must not be blank
}

func ExampleRunValidation_ref() {
	n := -1
	c := Config{MaxRetries: &n}
	err := shak.RunValidation(c.Validation())
	fmt.Println(err)
	// Output:
	// maxRetries: value -1 is less than minimum 0
}

func ExampleValidate_deref() {
	ptrRule := rule.Ref(rule.Positive[int]())
	err := shak.Validate(-1, rule.Deref(ptrRule))
	fmt.Println(err)
	// Output:
	// value -1 is not positive
}
