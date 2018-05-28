package swagger

import (
	"fmt"

	"github.com/jancajthaml/rest-contract-test/model"
)

func NewSwagger(file string) (*model.Contract, error) {
	contract := new(model.Contract)

	fmt.Println("not implemented")

	/*
		contract.Source = file

		rootResource, err := ParseFile(file)
		if err != nil {
			return contract, err
		}

		contract.Name = rootResource.Title

		for path, v := range rootResource.Resources {
			walk(contract, path, &v)
		}

		contract.Type = rootResource.RAMLVersion
	*/

	return contract, nil
}
