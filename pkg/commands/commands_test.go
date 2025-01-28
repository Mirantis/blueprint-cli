package commands

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"os"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mirantiscontainers/blueprint-cli/pkg/constants"
)

var _ = Describe("Commands", func() {
	Context("with version", func() {
		It("should be latest", func() {
			version := "latest"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal("https://github.com/mirantiscontainers/blueprint/releases/latest/download/blueprint-operator.yaml"))
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
		It("should be original remote uri", func() {
			version := "http://github.com"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal(version))
		})
		It("should be original file uri", func() {
			version := "file://~/bob/ross.yaml"
			uri, err := determineOperatorUri(version)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(Equal(version))
		})
		It("should error for an unknown value", func() {
			version := "13241"
			uri, err := determineOperatorUri(version)
			Expect(err).To(HaveOccurred())
			Expect(uri).To(Equal(""))
		})
	})

	Context("with image registry", Ordered, func() {
		var uris map[string]string
		remoteKey := "remote"
		localKey := "local"

		BeforeAll(func() {
			remoteURI := "https://github.com/mirantiscontainers/blueprint/releases/latest/download/blueprint-operator.yaml"

			testBopFile, err := os.CreateTemp("", "test-bop-*.yaml")
			Expect(err).ToNot(HaveOccurred())

			manifestBytes, err := downloadRemoteManifest(remoteURI)
			Expect(err).ToNot(HaveOccurred())

			n, err := testBopFile.Write(manifestBytes)
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(Equal(len(manifestBytes)))
			Expect(testBopFile.Close()).To(Succeed())

			localURI := fmt.Sprintf("file://%s", testBopFile.Name())

			uris = map[string]string{
				remoteKey: remoteURI,
				localKey:  localURI,
			}
		})

		AfterAll(func() {
			Expect(os.Remove(strings.TrimPrefix(uris[localKey], "file://"))).To(Succeed())
		})

		It("fails with an empty manifest", func() {
			_, needCleanup, err := setImageRegistry("", "registry.mirantis.com")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("empty BOP manifest URI"))
			Expect(needCleanup).To(BeFalse())
		})

		It("fails with a bad link for remote manifest", func() {
			_, needCleanup, err := setImageRegistry(uris[remoteKey]+"oops", "registry.mirantis.com")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to obtain BOP manifest"))
			Expect(needCleanup).To(BeFalse())
		})

		DescribeTable("should return original URI",
			func(testURIKey, registry string) {
				testURI := uris[testURIKey]
				uri, needCleanup, err := setImageRegistry(testURI, registry)
				Expect(err).ToNot(HaveOccurred())
				Expect(needCleanup).To(BeFalse())
				Expect(uri).To(Equal(testURI))
			},
			Entry("with empty registry and remote URI", remoteKey, ""),
			Entry("with empty registry and local URI", localKey, ""),
			Entry("with default registry and remote URI", remoteKey, constants.MirantisImageRegistry),
			Entry("with default registry and local URI", localKey, constants.MirantisImageRegistry),
		)

		DescribeTable("should update registry and return updated URI",
			func(testURIKey string) {
				testURI := uris[testURIKey]
				uri, needCleanup, err := setImageRegistry(testURI, "registry.mirantis.com")
				Expect(err).ToNot(HaveOccurred())
				defer os.Remove(strings.TrimPrefix(uri, "file://"))

				Expect(needCleanup).To(BeTrue())
				Expect(uri).ToNot(Equal(testURI))
				Expect(uri).To(HavePrefix("file://"))

				manifestBytes, err := readLocalManifest(strings.TrimPrefix(uri, "file://"))
				Expect(err).ToNot(HaveOccurred())
				Expect(string(manifestBytes)).NotTo(ContainSubstring(constants.MirantisImageRegistry))
				Expect(string(manifestBytes)).To(ContainSubstring("registry.mirantis.com"))
			},
			Entry("with remote URI", remoteKey),
			Entry("with local URI", localKey),
		)
	})

	It("detect image registry", func() {
		detected, err := detectDeployedRegistry([]corev1.Container{
			{
				Image: "ghcr.io/mirantiscontainers/kube-rbac-proxy:v1.0.0",
			},
			{
				Image: "ghcr.io/mirantiscontainers/blueprint-operator:v1.0.0",
			},
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(detected).To(Equal("ghcr.io/mirantiscontainers"))
	})
})
