package v08

type Any interface{}

type HTTPCode int

type HTTPHeader string

type NamedParameter struct {
	Name string

	DisplayName string `yaml:"displayName"`

	Description string

	Type Any

	Enum []Any `yaml:",flow"`

	Pattern *string

	MinLength *int `yaml:"minLength"`

	MaxLength *int `yaml:"maxLength"`

	Minimum *float64

	Maximum *float64

	Example Any

	Repeat *bool

	Required bool

	Default Any

	format Any `ramlFormat:"Named parameters must be mappings. Example: userId: {displayName: 'User ID', description: 'Used to identify the user.', type: 'integer', minimum: 1, example: 5}"`
}

type Header NamedParameter

type Documentation struct {
	Title string `yaml:"title"`

	Content Any `yaml:"content"`
}

type Body struct {
	mediaType string `yaml:"mediaType"`

	Schema Any `yaml:"schema"`

	Description string `yaml:"description"`

	Example Any `yaml:"example"`

	FormParameters map[string]NamedParameter `yaml:"formParameters"`

	Headers map[HTTPHeader]Header `yaml:"headers"`
}

type Bodies struct {
	DefaultSchema Any `yaml:"schema"`

	DefaultDescription string `yaml:"description"`

	DefaultExample Any `yaml:"example"`

	DefaultFormParameters map[string]NamedParameter `yaml:"formParameters"`

	Headers map[HTTPHeader]Header `yaml:"headers"`

	ForMIMEType map[string]Body `yaml:",regexp:.*"`
}

type Response struct {
	HTTPCode HTTPCode

	Description string

	Headers map[HTTPHeader]Header `yaml:"headers"`

	Bodies Bodies `yaml:"body"`
}

type DefinitionParameters map[Any]Any

type DefinitionChoice struct {
	Name       Any
	Parameters DefinitionParameters
}

func (dc *DefinitionChoice) UnmarshalYAML(unmarshaler func(interface{}) error) error {

	simpleDefinition := new(Any)
	parameterizedDefinition := make(map[Any]DefinitionParameters)

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
	Name string

	Usage string

	Description string

	Bodies Bodies `yaml:"body"`

	Headers map[HTTPHeader]Header `yaml:"headers"`

	Responses map[HTTPCode]Response `yaml:"responses"`

	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	Protocols []string `yaml:"protocols"`

	OptionalBodies Bodies `yaml:"body?"`

	OptionalHeaders map[HTTPHeader]Header `yaml:"headers?"`

	OptionalResponses map[HTTPCode]Response `yaml:"responses?"`

	OptionalQueryParameters map[string]NamedParameter `yaml:"queryParameters?"`
}

type ResourceTypeMethod struct {
	Name string

	Description string

	Bodies Bodies `yaml:"body"`

	Headers map[HTTPHeader]Header `yaml:"headers"`

	Responses map[HTTPCode]Response `yaml:"responses"`

	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	Protocols []string `yaml:"protocols"`
}

type ResourceType struct {
	Name string

	Usage string

	Description string

	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	Get *ResourceTypeMethod `yaml:"get"`

	Head *ResourceTypeMethod `yaml:"head"`

	Post *ResourceTypeMethod `yaml:"post"`

	Put *ResourceTypeMethod `yaml:"put"`

	Delete *ResourceTypeMethod `yaml:"delete"`

	Patch *ResourceTypeMethod `yaml:"patch"`

	OptionalUriParameters map[string]NamedParameter `yaml:"uriParameters?"`

	OptionalBaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters?"`

	OptionalGet *ResourceTypeMethod `yaml:"get?"`

	OptionalHead *ResourceTypeMethod `yaml:"head?"`

	OptionalPost *ResourceTypeMethod `yaml:"post?"`

	OptionalPut *ResourceTypeMethod `yaml:"put?"`

	OptionalDelete *ResourceTypeMethod `yaml:"delete?"`

	OptionalPatch *ResourceTypeMethod `yaml:"patch?"`
}

type SecuritySchemeMethod struct {
	Bodies Bodies `yaml:"body"`

	Headers map[HTTPHeader]Header `yaml:"headers"`

	Responses map[HTTPCode]Response `yaml:"responses"`

	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`
}

type SecurityScheme struct {
	Name string

	Description string

	Type Any

	DescribedBy SecuritySchemeMethod

	Settings map[string]Any

	Other map[string]string
}

type Method struct {
	Name string

	Description string

	SecuredBy []DefinitionChoice `yaml:"securedBy"`

	Headers map[HTTPHeader]Header `yaml:"headers"`

	Protocols []string `yaml:"protocols"`

	QueryParameters map[string]NamedParameter `yaml:"queryParameters"`

	Bodies Bodies `yaml:"body"`

	Responses map[HTTPCode]Response `yaml:"responses"`

	Is []DefinitionChoice `yaml:"is"`
}

type Resource struct {
	URI string

	Parent *Resource

	DisplayName string

	Description string

	SecuredBy []DefinitionChoice `yaml:"securedBy"`

	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	Type *DefinitionChoice `yaml:"type"`

	Is []DefinitionChoice `yaml:"is"`

	Get *Method `yaml:"get"`

	Head *Method `yaml:"head"`

	Post *Method `yaml:"post"`

	Put *Method `yaml:"put"`

	Delete *Method `yaml:"delete"`

	Patch *Method `yaml:"patch"`

	Nested map[string]*Resource `yaml:",regexp:/.*"`

	//Parent *Resource // FIXME double link
}

type APIDefinition struct {
	RAMLVersion string `yaml:"raml_version"`

	Title string `yaml:"title"`

	Version string `yaml:"version"`

	BaseUri string

	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`

	UriParameters map[string]NamedParameter `yaml:"uriParameters"`

	Protocols []string `yaml:"protocols"`

	MediaType string `yaml:"mediaType"`

	Schemas []map[string]Any

	SecuritySchemes []map[string]SecurityScheme `yaml:"securitySchemes"`

	SecuredBy []DefinitionChoice `yaml:"securedBy"`

	Documentation []Documentation `yaml:"documentation"`

	Traits []map[string]Trait `yaml:"traits"`

	ResourceTypes []map[string]ResourceType `yaml:"resourceTypes"`

	Resources map[string]Resource `yaml:",regexp:/.*"`
}

func (r *APIDefinition) GetResource(path string) *Resource {
	return nil
}
