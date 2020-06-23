/*
Copyright 2020 Red Hat, Inc. All rights reserved.

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

package e2e

import (
	"net/http"
	"testing"
)

func TestExporterRequestsTotal(t *testing.T) {
	expectedTotals := []int{5, 10}

	for _, expectedTotal := range expectedTotals {
		err := framework.ResetExporterMetrics()
		if err != nil {
			t.Fatalf("Could not reset exporter metrics: %v\n", err)
		}

		for i := 0; i < expectedTotal; i++ {
			resp, err := http.Get(framework.Exporter.EventServerURL + "/metrics")
			if err != nil {
				t.Fatalf("Could not get event metrics: %v\n", err)
			}
			err = resp.Body.Close()
			if err != nil {
				t.Logf("Could not close response body: %v\n", err)
			}
		}

		families, err := framework.GetExporterMetricFamilies()
		if err != nil {
			t.Fatalf("Could not get exporter metrics: %v\n", err)
		}

		requestsTotal, found := families["kube_events_exporter_requests_total"]
		if !found {
			t.Fatalf("kube_events_exporter_requests_total not found in metrics.\n")
		}

		requestsTotalValue := int(requestsTotal.Metric[0].Counter.GetValue())
		if requestsTotalValue != expectedTotal {
			t.Fatalf("kube_events_exporter_requests_total value is %d instead of %d.\n", requestsTotalValue, expectedTotal)
		}
	}
}