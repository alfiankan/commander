package cmdr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadChartJson(t *testing.T) {

	cmdr := NewCmdr("/Users/alfiankan/development/repack/commander/charts")
	charts := cmdr.readCharts()

	for _, v := range charts {
		fmt.Println(v)
	}

	assert.True(t, len(charts) > 0)
}

func TestListingCharts(t *testing.T) {

	cmdr := NewCmdr("/Users/alfiankan/development/repack/commander/charts")

	cmdr.listViewCharts()
}
