// Copyright 2016-2023, Pulumi Corporation.
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

package tests

import (
	"testing"

	"github.com/blang/semver"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/integration"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	statefulString "github.com/pulumi/pulumi-statefulstring/provider"
)

func TestStatefulStringCreate(t *testing.T) {
	prov := provider()

	response, err := prov.Create(p.CreateRequest{
		Urn: urn("StatefulString"),
		Properties: resource.PropertyMap{
			"string": resource.NewStringProperty("hello, world"),
			"triggers": resource.NewObjectProperty(resource.PropertyMap{
				"foo": resource.NewStringProperty("bar"),
			}),
		},
		Preview: false,
	})

	require.NoError(t, err)
	result := response.Properties["string"].StringValue()
	assert.Equal(t, "hello, world", result)
}

type ExpectedCreateResult struct {
}

func TestCreate(t *testing.T) {
	prov := provider()

	testCases := []struct {
		name           string
		properties     p.CreateRequest
		expectedResult p.CreateResponse
	}{
		{
			name: "Test case 1",
			properties: p.CreateRequest{
				Urn: urn("StatefulString"),
				Properties: resource.PropertyMap{
					"string": resource.NewStringProperty("hello, world"),
					"triggers": resource.NewObjectProperty(resource.PropertyMap{
						"foo": resource.NewStringProperty("bar"),
					}),
				},
				Preview: false,
			},
			expectedResult: p.CreateResponse{
				ID: "StatefulString",
				Properties: map[resource.PropertyKey]resource.PropertyValue{
					"string": resource.NewStringProperty("hello, world"),
					"triggers": resource.NewObjectProperty(resource.PropertyMap{
						"foo": resource.NewStringProperty("bar"),
					}),
				},
			},
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the resource
			response, err := prov.Create(tc.properties)
			require.NoError(t, err)

			// Check the result
			assert.Equal(t, tc.expectedResult.Properties, response.Properties)
			// assert.Equal(t, tc.expectedResult.ID, response.ID)
		})
	}
}

type ExpectedUpdateResult struct {
	String   string
	Triggers map[string]string
}

func TestStatefulStringUpdate(t *testing.T) {
	prov := provider()

	testCases := []struct {
		name           string
		initialProps   resource.PropertyMap
		updatedProps   resource.PropertyMap
		expectedResult ExpectedUpdateResult
	}{
		{
			name: "No string change with no trigger change",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("hello, world"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("hello, world"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "hello, world",
				Triggers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "String change with no trigger change",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("hello, world"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "hello, world",
				Triggers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "String change with trigger value change",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("hello, world"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "2",
				Triggers: map[string]string{
					"foo": "bar2",
				},
			},
		},
		{
			name: "No string change with trigger value change",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "1",
				Triggers: map[string]string{
					"foo": "bar2",
				},
			},
		},
		{
			name: "No string change with trigger key add",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "1",
				Triggers: map[string]string{
					"foo":  "bar",
					"foo2": "bar2",
				},
			},
		},
		{
			name: "String change with trigger key add",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "2",
				Triggers: map[string]string{
					"foo":  "bar",
					"foo2": "bar2",
				},
			},
		},
		{
			name: "String change with trigger key remove",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "2",
				Triggers: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "String change with trigger key remove and trigger value change",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedUpdateResult{
				String: "2",
				Triggers: map[string]string{
					"foo": "bar2",
				},
			},
		},
		// Add more test cases here...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the resource
			_, err := prov.Create(p.CreateRequest{
				Urn:        urn("StatefulString"),
				Properties: tc.initialProps,
				Preview:    false,
			})
			require.NoError(t, err)

			// Update the resource
			updateResponse, err := prov.Update(p.UpdateRequest{
				Urn:     urn("StatefulString"),
				Olds:    tc.initialProps,
				News:    tc.updatedProps,
				Preview: false,
			})
			require.NoError(t, err)

			// Check the result
			stringResult := updateResponse.Properties["string"].StringValue()

			// Get the triggers property
			triggersProp := updateResponse.Properties["triggers"]

			// Convert triggersProp to map[string]string
			triggersMap := make(map[string]string)
			for k, v := range triggersProp.ObjectValue() {
				triggersMap[string(k)] = v.StringValue()
			}

			assert.Equal(t, tc.expectedResult.String, stringResult)
			assert.Equal(t, tc.expectedResult.Triggers, triggersMap)
		})
	}
}

type ExpectedDiffResult struct {
	HasChanges   bool
	DetailedDiff map[string]p.PropertyDiff
}

func TestDiff(t *testing.T) {
	prov := provider()

	tests := []struct {
		name           string
		initialProps   resource.PropertyMap
		updatedProps   resource.PropertyMap
		expectedResult ExpectedDiffResult
	}{
		{
			name: "String Same__Triggers Same",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges:   false,
				DetailedDiff: map[string]p.PropertyDiff{},
			},
		},
		{
			name: "String Change__Triggers Same",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges:   false,
				DetailedDiff: map[string]p.PropertyDiff{},
			},
		},
		{
			name: "String Change__Triggers Change Value",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar3"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"string": {
						Kind: p.DiffKind("update"),
					},
					"triggers.foo2": {
						Kind: p.DiffKind("update"),
					},
				},
			},
		},
		{
			name: "String Change__Triggers Delete Key",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"string": {
						Kind: p.DiffKind("update"),
					},
					"triggers.foo2": {
						Kind: p.DiffKind("delete"),
					},
				},
			},
		},
		{
			name: "String Change__Triggers Add Key",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo": resource.NewStringProperty("bar"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("2"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"string": {
						Kind: p.DiffKind("update"),
					},
					"triggers.foo2": {
						Kind: p.DiffKind("add"),
					},
				},
			},
		},
		{
			name: "String Same__Triggers Change Value",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar3"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"triggers.foo2": {
						Kind: p.DiffKind("update"),
					},
				},
			},
		},
		{
			name: "String Same__Triggers Add and Delete Key",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo3": resource.NewStringProperty("bar3"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"triggers.foo2": {
						Kind: p.DiffKind("delete"),
					},
					"triggers.foo3": {
						Kind: p.DiffKind("add"),
					},
				},
			},
		},
		{
			name: "String Same__Triggers Add and Delete Key and Update Value",
			initialProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("bar"),
					"foo2": resource.NewStringProperty("bar2"),
				}),
			},
			updatedProps: resource.PropertyMap{
				"string": resource.NewStringProperty("1"),
				"triggers": resource.NewObjectProperty(resource.PropertyMap{
					"foo":  resource.NewStringProperty("barX"),
					"foo3": resource.NewStringProperty("bar3"),
				}),
			},
			expectedResult: ExpectedDiffResult{
				HasChanges: true,
				DetailedDiff: map[string]p.PropertyDiff{
					"triggers.foo": {
						Kind: p.DiffKind("update"),
					},
					"triggers.foo2": {
						Kind: p.DiffKind("delete"),
					},
					"triggers.foo3": {
						Kind: p.DiffKind("add"),
					},
				},
			},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDiff, err := prov.Diff(p.DiffRequest{
				// ID:            "",
				Urn:           urn("StatefulString"),
				Olds:          tt.initialProps,
				News:          tt.updatedProps,
				IgnoreChanges: nil,
			})
			require.NoError(t, err)
			// Check the result
			assert.Equal(t, tt.expectedResult.HasChanges, gotDiff.HasChanges)
			assert.Equal(t, tt.expectedResult.DetailedDiff, gotDiff.DetailedDiff)
		})
	}
}

// func TestCheck(t *testing.T) {
// 	t.Log("Starting TestCheck")
// 	prov := provider()

// 	testCases := []struct {
// 		name           string
// 		request        p.CheckRequest
// 		expectedResult p.CheckResponse
// 	}{
// 		{
// 			name: "String with 2 Triggers",
// 			request: p.CheckRequest{
// 				Urn:  urn("StatefulString"),
// 				Olds: resource.PropertyMap{},
// 				News: resource.PropertyMap{
// 					"string": resource.NewStringProperty("1"),
// 					"triggers": resource.NewObjectProperty(resource.PropertyMap{
// 						"foo":  resource.NewStringProperty("barX"),
// 						"foo3": resource.NewStringProperty("bar3"),
// 					}),
// 				},
// 			},
// 			expectedResult: p.CheckResponse{
// 				Inputs: resource.PropertyMap{
// 					"string": resource.NewStringProperty("1"),
// 					"triggers": resource.NewObjectProperty(resource.PropertyMap{
// 						"foo":  resource.NewStringProperty("barX"),
// 						"foo3": resource.NewStringProperty("bar3"),
// 					}),
// 				},
// 				Failures: nil,
// 			},
// 		},
// 		{
// 			name: "No String with 2 Triggers",
// 			request: p.CheckRequest{
// 				Urn:  urn("StatefulString"),
// 				Olds: resource.PropertyMap{},
// 				News: resource.PropertyMap{
// 					"string":   resource.NewStringProperty("1"),
// 					"triggers": resource.NewObjectProperty(resource.PropertyMap{}),
// 				},
// 			},
// 			expectedResult: p.CheckResponse{
// 				Inputs: resource.PropertyMap{
// 					"string": resource.NewStringProperty(""),
// 				},
// 				Failures: nil,
// 			},
// 		},
// {
// 	name: "String with 0 Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"string": resource.NewStringProperty("1"),
// 			"triggers": resource.NewObjectProperty(resource.PropertyMap{
// 				"foo":  resource.NewStringProperty("barX"),
// 				"foo3": resource.NewStringProperty("bar3"),
// 			}),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "No String with 0 Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"string":   resource.NewStringProperty("1"),
// 			"triggers": resource.NewObjectProperty(resource.PropertyMap{}),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "String with missing Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"string": resource.NewStringProperty("1"),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "No String with missing Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"string": resource.NewStringProperty("1"),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "Missing String with Missing Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "Missing String with 0 Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"triggers": resource.NewObjectProperty(resource.PropertyMap{}),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// {
// 	name: "Missing String with 1 Triggers",
// 	request: p.CheckRequest{
// 		Urn:  urn("StatefulString"),
// 		Olds: resource.PropertyMap{},
// 		News: resource.PropertyMap{
// 			"triggers": resource.NewObjectProperty(resource.PropertyMap{
// 				"foo": resource.NewStringProperty("bar"),
// 			}),
// 		},
// 	},
// 	expectedResult: p.CheckResponse{
// 		Inputs: resource.PropertyMap{
// 			"property1": resource.NewStringProperty("value1"),
// 		},
// 		Failures: nil,
// 	},
// },
// }

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Call the Check function
// 			response, err := prov.Check(tc.request)
// 			require.NoError(t, err)

// 			// Check the result
// 			assert.Equal(t, tc.expectedResult, response)
// 		})
// 	}
// }

// urn is a helper function to build an urn for running integration tests.
func urn(typ string) resource.URN {
	return resource.NewURN("stack", "proj", "",
		tokens.Type("test:index:"+typ), "name")
}

// Create a test server.
func provider() integration.Server {
	return integration.NewServer(statefulString.Name, semver.MustParse("1.0.0"), statefulString.Provider())
}
