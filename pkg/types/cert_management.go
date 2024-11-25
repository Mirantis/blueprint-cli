package types

import (
	"fmt"

	"github.com/mirantiscontainers/blueprint-operator/api/v1alpha1"
)

// CertManagement defines the desired state of cert-manager resources
type CertManagement struct {
	v1alpha1.CertManagement `json:",inline"`
}

// Validate checks the CertManagement structure and its children
func (c *CertManagement) Validate() error {
	for _, issuer := range c.Issuers {
		if issuer.Name == "" {
			return fmt.Errorf("issuer name cannot be empty")
		}
		if issuer.Namespace == "" {
			return fmt.Errorf("issuer namespace cannot be empty")
		}
	}

	for _, clusterIssuer := range c.ClusterIssuers {
		if clusterIssuer.Name == "" {
			return fmt.Errorf("cluster issuer name cannot be empty")
		}
	}

	for _, certificate := range c.Certificates {
		if certificate.Name == "" {
			return fmt.Errorf("certificate name cannot be empty")
		}
		if certificate.Namespace == "" {
			return fmt.Errorf("certificate namespace cannot be empty")
		}
		if certificate.Spec.IssuerRef.Name == "" {
			return fmt.Errorf("certificate issuer name cannot be empty")
		}
		if certificate.Spec.IssuerRef.Kind == "" {
			return fmt.Errorf("certificate issuer kind cannot be empty")
		}
	}

	return nil
}
