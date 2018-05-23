package parser

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

func fillResponses(method *v08.Method) []Response {

	//for k, v := range method.Responses {
	//fmt.Printf("%d -> example: %s , schema: %s\n", k, v.Bodies.DefaultExample, v.Bodies.DefaultSchema)
	//}

	// FIXME TBD

	return nil
}

func (contract *Contract) appendEndpoint(path, method string, endpoint *v08.Method) {
	res := Endpoint{
		Path:      path,
		Method:    method,
		Responses: fillResponses(endpoint),
		Headers:   "",
		Request:   Request{},
	}

	contract.Endpoints = append(contract.Endpoints, res)
}

func (contract *Contract) walk(path string, resource *v08.Resource) {
	var foundSomething = false

	if resource.Get != nil {
		foundSomething = true
		contract.appendEndpoint(path, "GET", resource.Get)
	}
	if resource.Head != nil {
		foundSomething = true
		contract.appendEndpoint(path, "HEAD", resource.Head)
	}
	if resource.Post != nil {
		foundSomething = true
		contract.appendEndpoint(path, "POST", resource.Post)
	}
	if resource.Put != nil {
		foundSomething = true
		contract.appendEndpoint(path, "PUT", resource.Put)
	}
	if resource.Patch != nil {
		foundSomething = true
		contract.appendEndpoint(path, "PATCH", resource.Patch)
	}
	if resource.Delete != nil {
		foundSomething = true
		contract.appendEndpoint(path, "DELETE", resource.Delete)
	}

	if !foundSomething {
		// FIXME no method meand GET and construct headers + response manually
	}

	for k, v := range resource.Nested {
		contract.walk(path+k, v)
	}
}

func (contract *Contract) FromFile(file string) error {

	contract.Source = file

	switch GetDocumentType(file) {

	// INFO does not work with includes but good MVP for now
	case "RAML 0.8":
		contract.Type = "RAML 0.8"

		rootResource, err := v08.RAMLv08(file)
		if err != nil {
			return err
		}

		contract.Name = rootResource.Title

		for k, v := range rootResource.Resources {
			contract.walk(k, &v)
		}

		return nil

	// INFO Does not work
	case "RAML 1.0":
		contract.Type = "RAML 1.0"

		_, err := v10.RAMLv10(file)
		if err != nil {
			return err
		}

		return nil

	default:
		contract.Type = "Invalid"

		return fmt.Errorf("unsupported document")
	}

}
