package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPackageInfo(t *testing.T) {
	p, err := GetPackageInfo("react")
	if !assert.Nil(t, err) {
		t.Fatal(err.Error())
	}

	assert.Equal(t, "react", p.Name)
}
