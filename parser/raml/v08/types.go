package v08

type NamedParameter struct {
	Name        string
	DisplayName string `yaml:"displayName"`
	Description string
	Type        interface{}
	Enum        []interface{} `yaml:",flow"`
	Pattern     *string
	MinLength   *int `yaml:"minLength"`
	MaxLength   *int `yaml:"maxLength"`
	Minimum     *float64
	Maximum     *float64
	Example     interface{}
	Repeat      *bool
	Required    bool
	Default     interface{}
}

type Documentation struct {
	Title   string      `yaml:"title"`
	Content interface{} `yaml:"content"`
}

type Body struct {
	mediaType      string                    `yaml:"mediaType"`
	Schema         interface{}               `yaml:"schema"`
	Description    string                    `yaml:"description"`
	Example        interface{}               `yaml:"example"`
	FormParameters map[string]NamedParameter `yaml:"formParameters"`
	Headers        map[string]NamedParameter `yaml:"headers"`
}

type Bodies struct {
	DefaultSchema         interface{}               `yaml:"schema"`
	DefaultDescription    string                    `yaml:"description"`
	DefaultExample        interface{}               `yaml:"example"`
	DefaultFormParameters map[string]NamedParameter `yaml:"formParameters"`
	Headers               map[string]NamedParameter `yaml:"headers"`
	ForMIMEType           map[string]Body           `yaml:",regexp:.*"`
}

type Response struct {
	HTTPCode    int
	Description string
	Headers     map[string]NamedParameter `yaml:"headers"`
	Bodies      Bodies                    `yaml:"body"`
}

type DefinitionChoice struct {
	Name       interface{}
	Parameters map[interface{}]interface{}
}

func (dc *DefinitionChoice) UnmarshalYAML(unmarshaler func(interface{}) error) error {

	simpleDefinition := new(interface{})
	parameterizedDefinition := make(map[interface{}]map[interface{}]interface{})

	var err error
	if err = unmarshaler(simpleDefinition); err == nil {
		dc.Name = *simpleDefinition
		dc.Parameters = nil
	} else if err = unmarshaler(parameterizedDefinition); err == nil {
		for choice, params := range parameterizedDefinition {
			dc.Name = choice
			dc.Parameters = params
		}
	}

	return err
}

type Trait struct {
	Name                    string
	Usage                   string
	Description             string
	Bodies                  Bodies                    `yaml:"body"`
	Headers                 map[string]NamedParameter `yaml:"headers"`
	Responses               map[int]Response          `yaml:"responses"`
	QueryParameters         map[string]NamedParameter `yaml:"queryParameters"`
	Protocols               []string                  `yaml:"protocols"`
	OptionalBodies          Bodies                    `yaml:"body?"`
	OptionalHeaders         map[string]NamedParameter `yaml:"headers?"`
	OptionalResponses       map[int]Response          `yaml:"responses?"`
	OptionalQueryParameters map[string]NamedParameter `yaml:"queryParameters?"`
}

type ResourceTypeMethod struct {
	Name            string
	Description     string
	Bodies          Bodies                    `yaml:"body"`
	Headers         map[string]NamedParameter `yaml:"headers"`
	Responses       map[int]Response          `yaml:"responses"`
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`
	Protocols       []string                  `yaml:"protocols"`
}

type ResourceType struct {
	Name                      string
	Usage                     string
	Description               string
	UriParameters             map[string]NamedParameter `yaml:"uriParameters"`
	BaseUriParameters         map[string]NamedParameter `yaml:"baseUriParameters"`
	Get                       *ResourceTypeMethod       `yaml:"get"`
	Head                      *ResourceTypeMethod       `yaml:"head"`
	Post                      *ResourceTypeMethod       `yaml:"post"`
	Put                       *ResourceTypeMethod       `yaml:"put"`
	Delete                    *ResourceTypeMethod       `yaml:"delete"`
	Patch                     *ResourceTypeMethod       `yaml:"patch"`
	OptionalUriParameters     map[string]NamedParameter `yaml:"uriParameters?"`
	OptionalBaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters?"`
	OptionalGet               *ResourceTypeMethod       `yaml:"get?"`
	OptionalHead              *ResourceTypeMethod       `yaml:"head?"`
	OptionalPost              *ResourceTypeMethod       `yaml:"post?"`
	OptionalPut               *ResourceTypeMethod       `yaml:"put?"`
	OptionalDelete            *ResourceTypeMethod       `yaml:"delete?"`
	OptionalPatch             *ResourceTypeMethod       `yaml:"patch?"`
}

type SecuritySchemeMethod struct {
	Bodies          Bodies                    `yaml:"body"`
	Headers         map[string]NamedParameter `yaml:"headers"`
	Responses       map[int]Response          `yaml:"responses"`
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`
}

type SecurityScheme struct {
	Name        string
	Description string
	Type        interface{}
	DescribedBy SecuritySchemeMethod
	Settings    map[string]interface{}
	Other       map[string]string
}

type Method struct {
	Name            string
	Description     string
	SecuredBy       []DefinitionChoice        `yaml:"securedBy"`
	Headers         map[string]NamedParameter `yaml:"headers"`
	Protocols       []string                  `yaml:"protocols"`
	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`
	Bodies          Bodies                    `yaml:"body"`
	Responses       map[int]Response          `yaml:"responses"`
	Is              []DefinitionChoice        `yaml:"is"`
}

type Resource struct {
	SecuredBy         []DefinitionChoice        `yaml:"securedBy"`
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`
	UriParameters     map[string]NamedParameter `yaml:"uriParameters"`
	Type              *DefinitionChoice         `yaml:"type"`
	Is                []DefinitionChoice        `yaml:"is"`
	Get               *Method                   `yaml:"get"`
	Head              *Method                   `yaml:"head"`
	Post              *Method                   `yaml:"post"`
	Put               *Method                   `yaml:"put"`
	Delete            *Method                   `yaml:"delete"`
	Patch             *Method                   `yaml:"patch"`
	Nested            map[string]*Resource      `yaml:",regexp:/.*"`
}

type APIDefinition struct {
	RAMLVersion       string                    `yaml:"raml_version"`
	Title             string                    `yaml:"title"`
	Version           string                    `yaml:"version"`
	BaseUri           string                    `yaml:"baseUri"`
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`
	UriParameters     map[string]NamedParameter `yaml:"uriParameters"`
	Protocols         []string                  `yaml:"protocols"`
	MediaType         string                    `yaml:"mediaType"`
	Schemas           []map[string]interface{}
	SecuritySchemes   []map[string]SecurityScheme `yaml:"securitySchemes"`
	SecuredBy         []DefinitionChoice          `yaml:"securedBy"`
	Documentation     []Documentation             `yaml:"documentation"`
	Traits            []map[string]Trait          `yaml:"traits"`
	ResourceTypes     []map[string]ResourceType   `yaml:"resourceTypes"`
	Resources         map[string]Resource         `yaml:",regexp:/.*"`
}
