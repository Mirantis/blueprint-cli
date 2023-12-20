package types

import (
	"fmt"
	"runtime"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var (
	// Hack to get the path to this file to use as an "existing" config file
	_, thisFile, _, _ = runtime.Caller(0)
)

// TestKubernetesConfigPath tests the Validate method of the Kubernetes type
func TestKubernetesConfigPath(t *testing.T) {
	tests := map[string]struct {
		path string
		want types.GomegaMatcher
	}{
		"config doesn't exist": {path: "/some/file", want: Equal(fmt.Errorf("config file \"/some/file\" does not exist: stat /some/file: no such file or directory"))},
		"config exists":        {path: thisFile, want: BeNil()},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Set up the test environment
			g := NewWithT(t)

			// Run the method under test
			kubernetes := Kubernetes{
				ConfigPath: tc.path,
			}
			actual := kubernetes.Validate()

			// Check the results
			g.Expect(actual).Should(tc.want)

		})
	}
}
