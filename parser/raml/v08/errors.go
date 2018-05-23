package v08

import (
	"fmt"
	"strings"

	yaml "github.com/advance512/yaml"
)

type RamlError struct {
	Errors []string
}

func (e *RamlError) Error() string {
	return fmt.Sprintf("Error parsing RAML:\n  %s\n",
		strings.Join(e.Errors, "\n  "))
}

func populateRAMLError(ramlError *RamlError,
	yamlErrors *yaml.TypeError) {

	for _, currErr := range yamlErrors.Errors {
		ramlError.Errors =
			append(ramlError.Errors, convertYAMLError(currErr))
	}
}

func convertYAMLError(yamlError string) string {
	if strings.Contains(yamlError, "cannot unmarshal") {
		yamlErrorParts := strings.Split(yamlError, " ")

		if len(yamlErrorParts) >= 7 {

			fmt.Println(yamlError)

			var ok bool
			var source string
			var target string
			var targetName string
			line := yamlErrorParts[1]
			line = line[:len(line)-1]

			if source, ok = yamlTypeToName[yamlErrorParts[4]]; !ok {
				source = yamlErrorParts[4]
			}
			fmt.Println("source: ", source)

			if source == "string" {
				source = fmt.Sprintf("string (got %s)", yamlErrorParts[5])
				target = yamlErrorParts[7]
			} else {
				target = yamlErrorParts[6]

			}
			if targetName, ok = ramlTypeNames[target]; !ok {
				targetName = target
			}

			target, _ = ramlTypes[target]

			return fmt.Sprintf("line %s: %s cannot be of "+
				"type %s, must be %s", line, targetName, source, target)

		}
	}

	return fmt.Sprintf("YAML error, %s", yamlError)
}

var yamlTypeToName map[string]string = map[string]string{
	"!!seq":       "sequence",
	"!!map":       "mapping",
	"!!int":       "integer",
	"!!str":       "string",
	"!!null":      "null",
	"!!bool":      "boolean",
	"!!float":     "float",
	"!!timestamp": "timestamp",
	"!!binary":    "binary",
	"!!merge":     "merge",
}

var ramlTypeNames map[string]string = map[string]string{
	"string": "string value",
	"int":    "numeric value",
	"raml.NamedParameter":       "named parameter",
	"raml.HTTPCode":             "HTTP code",
	"raml.HTTPHeader":           "HTTP header",
	"raml.Header":               "header",
	"raml.Documentation":        "documentation",
	"raml.Body":                 "body",
	"raml.Response":             "response",
	"raml.DefinitionParameters": "definition parameters",
	"raml.DefinitionChoice":     "definition choice",
	"raml.Trait":                "trait",
	"raml.ResourceTypeMethod":   "resource type method",
	"raml.ResourceType":         "resource type",
	"raml.SecuritySchemeMethod": "security scheme method",
	"raml.SecurityScheme":       "security scheme",
	"raml.Method":               "method",
	"raml.Resource":             "resource",
	"raml.APIDefinition":        "API definition",
}

var ramlTypes map[string]string = map[string]string{
	"string": "string",
	"int":    "integer",
	"raml.NamedParameter":       "mapping",
	"raml.HTTPCode":             "integer",
	"raml.HTTPHeader":           "string",
	"raml.Header":               "mapping",
	"raml.Documentation":        "mapping",
	"raml.Body":                 "mapping",
	"raml.Response":             "mapping",
	"raml.DefinitionParameters": "mapping",
	"raml.DefinitionChoice":     "string or mapping",
	"raml.Trait":                "mapping",
	"raml.ResourceTypeMethod":   "mapping",
	"raml.ResourceType":         "mapping",
	"raml.SecuritySchemeMethod": "mapping",
	"raml.SecurityScheme":       "mapping",
	"raml.Method":               "mapping",
	"raml.Resource":             "mapping",
	"raml.APIDefinition":        "mapping",
}
