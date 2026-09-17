/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package commentstart

import (
	"slices"

	"golang.org/x/tools/go/analysis"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/initializer"
	"sigs.k8s.io/kube-api-linter/pkg/analysis/registry"
)

func init() {
	registry.DefaultRegistry().RegisterLinter(Initializer())
}

// Initializer returns the AnalyzerInitializer for this
// Analyzer so that it can be added to the registry.
func Initializer() initializer.AnalyzerInitializer {
	return initializer.NewConfigurableInitializer(
		name,
		initAnalyzer,
		true,
		validateConfig,
	)
}

func initAnalyzer(cfg *Config) (*analysis.Analyzer, error) {
	return newAnalyzer(cfg), nil
}

// validateConfig validates the commentstart linter configuration.
func validateConfig(cfg *Config, fldPath *field.Path) field.ErrorList {
	if cfg == nil {
		return field.ErrorList{}
	}

	fieldErrors := field.ErrorList{}

	seen := make(map[string]struct{}, len(cfg.ExcludePrefixes))

	for i, prefix := range cfg.ExcludePrefixes {
		if prefix == "" {
			fieldErrors = append(fieldErrors, field.Invalid(
				fldPath.Child("excludePrefixes").Index(i),
				prefix,
				"must not be empty",
			))

			continue
		}

		if _, ok := seen[prefix]; ok {
			fieldErrors = append(fieldErrors, field.Duplicate(
				fldPath.Child("excludePrefixes").Index(i),
				prefix,
			))

			continue
		}

		seen[prefix] = struct{}{}

		if slices.Contains(defaultExcludePrefixes, prefix) {
			fieldErrors = append(fieldErrors, field.Invalid(
				fldPath.Child("excludePrefixes").Index(i),
				prefix,
				"is already a built-in default prefix",
			))
		}
	}

	return fieldErrors
}
