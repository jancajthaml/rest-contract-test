package parser

//import "fmt"
import (
	"fmt"

	"github.com/jancajthaml/rest-contract-test/parser/raml/v08"
	"github.com/jancajthaml/rest-contract-test/parser/raml/v10"
)

type Response struct {
	Example string
	Schema  string
}

type Request struct {
	Example string
	Schema  string
	Query   string
}

type Endpoint struct {
	Path      string
	Method    string
	Responses []Response
	Headers   string
	Request   Request
}

type Contract struct {
	Source    string
	Type      string
	Name      string
	Endpoints []Endpoint
}

func (contract *Contract) FromFile(file string) error {

	contract.Source = file

	switch GetDocumentType(file) {

	// INFO does not work with includes but good MVP for now
	case "RAML 0.8":
		contract.Type = "RAML 0.8"

		apiDefinition, err := v08.RAMLv08(file)
		if err != nil {
			return err
		}

		contract.Name = apiDefinition.Title

		//fmt.Printf("+------------------------------------------------------------------------\n")
		//fmt.Printf("| RAML %s\n", file)
		//fmt.Printf("+------------------------------------------------------------------------\n")
		//fmt.Printf("| title: %s\n", apiDefinition.Title)
		//fmt.Printf("+------------------------------------------------------------------------\n")

		//endpoints =
		//var endpoints []Endpoint

		// Iterate and print all endpoints
		for k, v := range apiDefinition.Resources {
			if v.Get != nil {
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "GET",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
			if v.Head != nil {
				//fmt.Printf("| HEAD    | %s\n", k)
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "HEAD",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
			if v.Post != nil {
				//fmt.Printf("| POST    | %s\n", k)
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "POST",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
			if v.Put != nil {
				//fmt.Printf("| PUT     | %s\n", k)
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "PUT",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
			if v.Patch != nil {
				//fmt.Printf("| PATCH   | %s\n", k)
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "PATCH",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
			if v.Delete != nil {
				//fmt.Printf("| DELETE  | %s\n", k)
				contract.Endpoints = append(contract.Endpoints, Endpoint{
					Path:      k,
					Method:    "DELETE",
					Responses: nil,
					Headers:   "",
					Request:   Request{},
				})
			}
		}

		return nil

	// INFO Does not work
	case "RAML 1.0":
		contract.Type = "RAML 1.0"

		_, err := v10.RAMLv10(file)
		if err != nil {
			return err
		}

		//fmt.Println(apiDefinition)
		//return new(Contract)
		return nil

	default:
		contract.Type = "Invalid"

		return fmt.Errorf("unsupported document")
	}

}
