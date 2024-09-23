package commands

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commands", func() {
	Context("with version", func() {
		It("should be latest", func() {
			version := "latest"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal("https://github.com/MirantisContainers/blueprint/releases/download/latest/blueprint-operator.yaml"))
		})
		It("should be semver with a leading v", func() {
			version := "v1.2.3"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal("https://github.com/MirantisContainers/blueprint/releases/download/v1.2.3/blueprint-operator.yaml"))
		})
		It("should be semver without a leading v", func() {
			version := "1.2.3"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal("https://github.com/MirantisContainers/blueprint/releases/download/v1.2.3/blueprint-operator.yaml"))
		})
		It("should be original uri", func() {
			version := "http://github.com"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal("http://github.com"))
		})
		It("should error for an unknown value", func() {
			version := "13241"
			uri, err := determineOperatorUri(version)
			Expect(err).To(HaveOccurred())
			Expect(uri).To(Equal(""))
		})
	})

})
