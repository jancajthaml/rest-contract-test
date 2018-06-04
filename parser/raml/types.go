// Copyright (c) 2016-2018, Jan Cajthaml <jan.cajthaml@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raml

// FIXME info query parameters and headers are redundant data-types
type QueryParameters struct {
	Data map[string]NamedParameter
}

func (ref *QueryParameters) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make(map[string]NamedParameter)

	if err = unmarshaler(ref.Data); err == nil {
		return
	}

	data := make(map[string]interface{}, 0)
	if err = unmarshaler(&data); err == nil {
		for k, v := range data {
			ref.Data[k] = NamedParameter{
				Type: v,
			}
		}

		return
	}

	return
}

type Headers struct {
	Data map[string]NamedParameter
}

func (ref *Headers) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make(map[string]NamedParameter)

	if err = unmarshaler(ref.Data); err == nil {
		return
	}

	data := make(map[string]interface{}, 0)
	if err = unmarshaler(&data); err == nil {
		for k, v := range data {
			ref.Data[k] = NamedParameter{
				Type: v,
			}
		}

		return
	}

	return
}

type NamedParameter struct {
	Name        string
	DisplayName string `yaml:"displayName"`
	Description string
	Type        interface{}
	Enum        []string `yaml:"enum,flow"`
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
	Headers        *Headers                  `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
}

type Bodies struct {
	Referenced            *string
	DefaultSchema         interface{}               `yaml:"schema"`
	DefaultDescription    string                    `yaml:"description"`
	DefaultExample        interface{}               `yaml:"example"`
	DefaultFormParameters map[string]NamedParameter `yaml:"formParameters"`
	Headers               *Headers                  `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
	ForMIMEType           map[string]Body           `yaml:",regexp:.*"`
}

type LiteralBodies struct {
	DefaultSchema         interface{}               `yaml:"schema"`
	DefaultDescription    string                    `yaml:"description"`
	DefaultExample        interface{}               `yaml:"example"`
	DefaultFormParameters map[string]NamedParameter `yaml:"formParameters"`
	Headers               *Headers                  `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
	ForMIMEType           map[string]Body           `yaml:",regexp:.*"`
}

type Response struct {
	HTTPCode    int
	Description string
	Headers     *Headers `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
	Bodies      Bodies   `yaml:"body"`
}

type Traits struct {
	Data map[string]*Trait
}

func (ref *Traits) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make(map[string]*Trait)

	if err = unmarshaler(ref.Data); err == nil {
		return
	}

	data := make([]map[string]*Trait, 0)
	if err = unmarshaler(&data); err == nil {
		for _, subset := range data {
			for k, v := range subset {
				ref.Data[k] = v
			}
		}

		return
	}

	return
}

func (ref *Bodies) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	literal := new(LiteralBodies)
	if err = unmarshaler(literal); err == nil {
		ref.DefaultSchema = literal.DefaultSchema
		ref.DefaultDescription = literal.DefaultDescription
		ref.DefaultExample = literal.DefaultExample
		ref.DefaultFormParameters = literal.DefaultFormParameters
		ref.Headers = literal.Headers
		ref.ForMIMEType = literal.ForMIMEType
		return
	}

	data := new(string)
	if err = unmarshaler(data); err == nil {
		ref.Referenced = data
		return
	}

	return
}

type ResourceTypes struct {
	Data map[string]interface{}
}

func (ref *ResourceTypes) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make(map[string]interface{})

	if err = unmarshaler(ref.Data); err == nil {
		return
	}

	data := make([]map[string]interface{}, 0)
	if err = unmarshaler(&data); err == nil {
		for _, subset := range data {
			for k, v := range subset {
				ref.Data[k] = v
			}
		}

		return
	}

	return
}

type Schemas struct {
	Data map[string]interface{}
}

func (ref *Schemas) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make(map[string]interface{})

	if err = unmarshaler(ref.Data); err == nil {
		return
	}

	data := make([]map[string]interface{}, 0)
	if err = unmarshaler(&data); err == nil {
		for _, subset := range data {
			for k, v := range subset {
				ref.Data[k] = v
			}
		}

		return
	}

	return
}

type BaseURI struct {
	Data string
}

func (ref *BaseURI) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	simple := new(string)
	if err = unmarshaler(simple); err == nil {
		ref.Data = *simple
		return
	}
	composite := make(map[string]interface{})
	if err = unmarshaler(composite); err == nil {
		if val, ok := composite["value"]; ok {
			ref.Data = val.(string)
		}
		return
	}

	return
}

type MediaType struct {
	Data []string
}

func (ref *MediaType) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	var simple string
	if err = unmarshaler(&simple); err == nil {
		ref.Data = make([]string, 1)
		ref.Data[0] = simple
		return
	}

	if err = unmarshaler(&(ref.Data)); err == nil {
		return
	}

	return
}

type DefinitionChoice struct {
	Name       string
	Parameters map[interface{}]interface{}
}

type Reference struct {
	Data []string
}

func (ref *Reference) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	ref.Data = make([]string, 0)

	var anything interface{}
	if err = unmarshaler(&anything); err == nil {

		switch hinted := anything.(type) {

		case []interface{}:
			for _, v := range hinted {
				switch typed := v.(type) {

				case map[interface{}]interface{}:
					for k := range typed {
						ref.Data = append(ref.Data, k.(string))
					}

				case interface{}:
					ref.Data = append(ref.Data, typed.(string))
				}
			}

		case interface{}:
			ref.Data = append(ref.Data, hinted.(string))

		}
	}

	return
}

func (ref *DefinitionChoice) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	simpleDefinition := new(string)
	parameterizedDefinition := make(map[string]map[interface{}]interface{})

	if err = unmarshaler(simpleDefinition); err == nil {
		ref.Name = *simpleDefinition
		ref.Parameters = nil
		return
	}

	if err = unmarshaler(parameterizedDefinition); err == nil {
		for choice, params := range parameterizedDefinition {
			ref.Name = choice
			ref.Parameters = params
		}
	}

	return err
}

type Trait struct {
	Bodies                  Bodies                 `yaml:"body"`
	Headers                 *Headers               `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
	Responses               map[int]Response       `yaml:"responses"`
	QueryParameters         *QueryParameters       `yaml:"queryParameters"`
	Protocols               []string               `yaml:"protocols"`
	OptionalBodies          Bodies                 `yaml:"body?"`
	OptionalHeaders         map[string]interface{} `yaml:"headers?"`
	OptionalResponses       map[int]Response       `yaml:"responses?"`
	OptionalQueryParameters map[string]interface{} `yaml:"queryParameters?"`
}

type ResourceTypeMethod struct {
	Name        string
	Description string
	Bodies      Bodies `yaml:"body"`

	Headers         *Headers         `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`
	Responses       map[int]Response `yaml:"responses"`
	QueryParameters *QueryParameters `yaml:"queryParameters"`
	Protocols       []string         `yaml:"protocols"`
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

// FIXME name differently
type SecuritySchemeMethod struct {
	Bodies Bodies `yaml:"body"`

	Headers *Headers `yaml:"headers"` //map[string]NamedParameter `yaml:"headers"`

	Responses       map[int]Response `yaml:"responses"`
	QueryParameters *QueryParameters `yaml:"queryParameters"`
}

type SecurityScheme struct {
	Name        string
	Description string
	Type        interface{}
	DescribedBy SecuritySchemeMethod `yaml:"describedBy"`
	Settings    map[string]interface{}
	Other       map[string]string
}

type Method struct {
	Name            string
	Description     string
	Protocols       []string         `yaml:"protocols"`
	SecuredBy       *Reference       `yaml:"securedBy"`
	Headers         *Headers         `yaml:"headers"`
	QueryParameters *QueryParameters `yaml:"queryParameters"`
	Bodies          Bodies           `yaml:"body"`
	Responses       map[int]Response `yaml:"responses"`
	Is              *Reference       `yaml:"is"`
}

type Resource struct {
	SecuredBy         *Reference                `yaml:"securedBy"`
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`
	UriParameters     map[string]NamedParameter `yaml:"uriParameters"`
	Type              *DefinitionChoice         `yaml:"type"`
	Is                *Reference                `yaml:"is"`
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
	BaseUri           *BaseURI                  `yaml:"baseUri"`
	BaseUriParameters map[string]NamedParameter `yaml:"baseUriParameters"`
	UriParameters     map[string]NamedParameter `yaml:"uriParameters"`
	Protocols         []string                  `yaml:"protocols"`
	MediaType         *MediaType                `yaml:"mediaType"`
	SecuritySchemes   map[string]SecurityScheme `yaml:"securitySchemes"`

	SecuredBy     *Reference      `yaml:"securedBy"`
	Documentation []Documentation `yaml:"documentation"`

	// FIXME these three are same data structures and thus now redundant
	Schemas       *Schemas       `yaml:"schemas"`
	Traits        *Traits        `yaml:"traits"`
	ResourceTypes *ResourceTypes `yaml:"resourceTypes"`

	Resources map[string]Resource `yaml:",regexp:/.*"`
}
