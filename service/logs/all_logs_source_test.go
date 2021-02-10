// dbaas-controller
// Copyright (C) 2020 Percona LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

// Package logs contains implementation of API for getting logs out of
// Kubernetes cluster workloads.
package logs

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	controllerv1beta1 "github.com/percona-platform/dbaas-api/gen/controller"
)

func TestLimitLines(t *testing.T) {
	type testCase struct {
		limit    int
		input    []*controllerv1beta1.Logs
		expected []*controllerv1beta1.Logs
	}
	testCases := []testCase{
		testCase{
			limit: 10,
			input: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"a", "b", "c", "d"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{},
				},
			},
			expected: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"a", "b", "c", "d"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{},
				},
			},
		},
		testCase{
			limit: 10,
			input: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"a", "b", "c", "d", "e", "f", "g"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"h", "i", "j"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"l", "m", "o", "p", "q", "r", "s"},
				},
			},
			expected: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"d", "e", "f", "g"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"h", "i", "j"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"q", "r", "s"},
				},
			},
		},
		testCase{
			limit: 10,
			input: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"a", "b", "c", "d", "e", "f", "g", "l", "m", "o", "p", "q", "r", "s"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"h"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"i"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"j"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"k"},
				},
			},
			expected: []*controllerv1beta1.Logs{
				&controllerv1beta1.Logs{
					Logs: []string{"m", "o", "p", "q", "r", "s"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"h"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"i"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"j"},
				},
				&controllerv1beta1.Logs{
					Logs: []string{"k"},
				},
			},
		},
	}
	for _, tc := range testCases {
		limitLines(tc.input, tc.limit)
		assert.Truef(t, equal(tc.input, tc.expected), "expected %v\ngot %v", tc.expected, tc.input)
	}
}

func equal(a, b []*controllerv1beta1.Logs) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual((*a[i]).Logs, (*b[i]).Logs) {
			log.Println((*a[i]).Logs, "!=", (*b[i]).Logs)
			return false
		}
	}
	return true
}
