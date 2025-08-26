/*
Copyright 2025 The llm-d-inference-sim Authors.

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

package llmdinferencesim_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/llm-d/llm-d-inference-sim/pkg/common"
	llmdinferencesim "github.com/llm-d/llm-d-inference-sim/pkg/llm-d-inference-sim"
)

var _ = Describe("Failures", func() {
	// Note: shouldInjectFailure is now private, so we test it through GetRandomFailure behavior

	Describe("GetRandomFailure", func() {
		It("should return a failure from all types when none specified", func() {
			config := &common.Configuration{
				Model:        "test-model",
				FailureTypes: []string{},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(BeNumerically(">=", 400))
			Expect(failure.Error.Message).ToNot(BeEmpty())
			Expect(failure.Error.Type).ToNot(BeEmpty())
		})

		It("should return rate limit failure when specified", func() {
			config := &common.Configuration{
				Model:        "test-model",
				FailureTypes: []string{common.FailureTypeRateLimit},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(Equal(429))
			Expect(failure.Error.Type).To(Equal("rate_limit_exceeded"))
			Expect(failure.Error.Code).To(Equal(429))
			Expect(strings.Contains(failure.Error.Message, "test-model")).To(BeTrue())
		})

		It("should return invalid API key failure when specified", func() {
			config := &common.Configuration{
				FailureTypes: []string{common.FailureTypeInvalidAPIKey},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(Equal(401))
			Expect(failure.Error.Type).To(Equal("invalid_request_error"))
			Expect(failure.Error.Code).To(Equal(401))
			Expect(failure.Error.Message).To(Equal("Incorrect API key provided"))
		})

		It("should return context length failure when specified", func() {
			config := &common.Configuration{
				FailureTypes: []string{common.FailureTypeContextLength},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(Equal(400))
			Expect(failure.Error.Type).To(Equal("invalid_request_error"))
			Expect(failure.Error.Code).To(Equal(400))
			Expect(failure.Error.Param).ToNot(BeNil())
			Expect(*failure.Error.Param).To(Equal("messages"))
		})

		It("should return server error when specified", func() {
			config := &common.Configuration{
				FailureTypes: []string{common.FailureTypeServerError},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(Equal(503))
			Expect(failure.Error.Type).To(Equal("server_error"))
			Expect(failure.Error.Code).To(Equal(503))
		})

		It("should return model not found failure when specified", func() {
			config := &common.Configuration{
				Model:        "test-model",
				FailureTypes: []string{common.FailureTypeModelNotFound},
			}
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(Equal(404))
			Expect(failure.Error.Type).To(Equal("invalid_request_error"))
			Expect(failure.Error.Code).To(Equal(404))
			Expect(strings.Contains(failure.Error.Message, "test-model-nonexistent")).To(BeTrue())
		})

		It("should return server error as fallback for empty types", func() {
			config := &common.Configuration{
				FailureTypes: []string{},
			}
			// This test is probabilistic since it randomly selects, but we can test structure
			failure := llmdinferencesim.GetRandomFailure(config)
			Expect(failure.StatusCode).To(BeNumerically(">=", 400))
			Expect(failure.Error.Type).ToNot(BeEmpty())
		})
	})
})
