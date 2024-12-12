package cmd

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version Command", func() {
	// There's not a great way to unit test dev vs merge vs release versions
	// It's more of an e2e test
	// We can atleast check the output format
	Context("version outputs", func() {
		It("default output should print just the version", func() {
			output := captureOutput(func() {
				// Setup the command
				cmd := versionCmd()

				// Call the command
				err := runVersion(cmd, []string{})
				Expect(err).To(BeNil())
			})

			Expect(output).To(Equal(fmt.Sprintf("Version: %s\n", version)))
		})

		It("verbose output should print version, date, and commit on separate lines", func() {
			output := captureOutput(func() {
				// Setup the command
				cmd := versionCmd()
				cmd.Flags().Set("verbose", "true")

				// Call the command
				err := runVersion(cmd, []string{})
				Expect(err).To(BeNil())
			})

			Expect(output).To(Equal(fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\n", version, commit, date)))
		})
	})
})
