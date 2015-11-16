package assetrepo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAssetrepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Assetrepo Suite")
}
