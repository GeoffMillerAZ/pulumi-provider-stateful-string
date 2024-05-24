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

package provider

import (
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

// Version is initialized by the Go linker to contain the semver of this build.
var Version string

const Name string = "statefulString"

func Provider() p.Provider {
	// We tell the provider what resources it needs to support.
	// In this case, a single custom resource.
	return infer.Provider(infer.Options{
		Resources: []infer.InferredResource{
			infer.Resource[StatefulString, StatefulStringArgs, StatefulStringState](),
		},
		ModuleMap: map[tokens.ModuleName]tokens.ModuleName{
			"provider": "index",
		},
	})
}

// Each resource has a controlling struct.
// Resource behavior is determined by implementing methods on the controlling struct.
// The `Create` method is mandatory, but other methods are optional.
// - Check: Remap inputs before they are typed.
// - Diff: Change how instances of a resource are compared.
// - Update: Mutate a resource in place.
// - Read: Get the state of a resource from the backing provider.
// - Delete: Custom logic when the resource is deleted.
// - Annotate: Describe fields and set defaults for a resource.
// - WireDependencies: Control how outputs and secrets flows through values.
type StatefulString struct{}

// Each resource has an input struct, defining what arguments it accepts.
type StatefulStringArgs struct {
	// Fields projected into Pulumi must be public and hava a `pulumi:"..."` tag.
	// The pulumi tag doesn't need to match the field name, but it's generally a
	// good idea.
	String   string            `pulumi:"string"`
	Triggers map[string]string `pulumi:"triggers"`
}

// Each resource has a state, describing the fields that exist on the created resource.
type StatefulStringState struct {
	// It is generally a good idea to embed args in outputs, but it isn't strictly necessary.
	StatefulStringArgs
	// Here we define a required output called result.
}

// All resources must implement Create at a minimum.
func (ss StatefulString) Create(ctx p.Context, name string, input StatefulStringArgs, preview bool) (id string, output StatefulStringState, err error) {
	id = name
	output = StatefulStringState{
		StatefulStringArgs: input,
	}
	err = nil

	return id, output, nil
}

type checkTriggerDiffAndUpdateResult struct {
	triggerChanged     bool
	changeMap          map[string]p.PropertyDiff
	statefulStringArgs StatefulStringArgs
}

func checkTriggerDiffAndUpdate(olds StatefulStringState, news StatefulStringArgs) (result checkTriggerDiffAndUpdateResult, err error) {
	// Assume no triggers have changed initially
	r := checkTriggerDiffAndUpdateResult{
		triggerChanged: false,
		changeMap:      map[string]p.PropertyDiff{},
		statefulStringArgs: StatefulStringArgs{
			String:   olds.String,
			Triggers: news.Triggers,
		},
	}

	// 1. Check if any new triggers have values different from old triggers or are newly added
	for newKey, newValue := range news.Triggers {
		oldValue, exists := olds.Triggers[newKey]
		fullKey := "triggers." + newKey
		if !exists {
			r.changeMap[fullKey] = p.PropertyDiff{
				Kind:      p.DiffKind("add"),
				InputDiff: false,
			}
			// If a new trigger is added
			r.triggerChanged = true
		} else if newValue != oldValue {
			// If an existing trigger's value has changed
			r.triggerChanged = true
			r.changeMap[fullKey] = p.PropertyDiff{
				Kind:      p.DiffKind("update"),
				InputDiff: false,
			}
		}
	}

	// 2. Check if any old triggers have been removed
	for oldKey := range olds.Triggers {
		if _, exists := news.Triggers[oldKey]; !exists {
			fullKey := "triggers." + oldKey
			// If an old trigger is removed
			r.triggerChanged = true
			r.changeMap[fullKey] = p.PropertyDiff{
				Kind:      p.DiffKind("delete"),
				InputDiff: false,
			}
		}
	}

	// If a trigger has changed, update the string and triggers
	if r.triggerChanged {
		r.statefulStringArgs = news
		if news.String != olds.String {
			r.changeMap["string"] = p.PropertyDiff{
				Kind:      p.DiffKind("update"),
				InputDiff: false,
			}
		}
	}

	return r, nil
}

func (ss StatefulString) Update(ctx p.Context, name string, olds StatefulStringState, news StatefulStringArgs, preview bool) (output StatefulStringState, err error) {
	d, _ := checkTriggerDiffAndUpdate(olds, news)

	// If no triggers have changed, return the old string but with new triggers
	return StatefulStringState{
		StatefulStringArgs: d.statefulStringArgs,
	}, nil
}

func (ss StatefulString) Diff(ctx p.Context, name string, olds StatefulStringState, news StatefulStringArgs) (p.DiffResponse, error) {
	d, _ := checkTriggerDiffAndUpdate(olds, news)

	return p.DiffResponse{
		HasChanges:   d.triggerChanged,
		DetailedDiff: d.changeMap,
	}, nil
}

func (ss StatefulString) Check(ctx p.Context, name string, oldInputs resource.PropertyMap, newInputs resource.PropertyMap) (p.CheckResponse, error) {
	print("XXXXXXXXXXXX CHECK XXXXXXXXXXXXXXX")
	print("XXXXXXXXXXXX CHECK XXXXXXXXXXXXXXX")
	print("XXXXXXXXXXXX CHECK XXXXXXXXXXXXXXX")
	print("XXXXXXXXXXXX CHECK XXXXXXXXXXXXXXX")
	// Initialize an empty slice to hold any failures
	// failures := []p.CheckFailure{}

	// // Check if the Triggers map is empty
	// if len(news.Triggers) == 0 {
	// 	failures = append(failures, p.CheckFailure{
	// 		Property: "Triggers",
	// 		Reason:   "Triggers map cannot be empty",
	// 	})
	// }

	// // If there are any failures, return them
	// if len(failures) > 0 {
	// 	return p.CheckResponse{Failures: failures}, nil
	// }

	// If there are no failures, return a successful CheckResponse
	return p.CheckResponse{}, nil
}
