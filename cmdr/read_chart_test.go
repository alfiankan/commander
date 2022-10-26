package cmdr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadChartJson(t *testing.T) {

	cmdr := NewCmdr()
	charts := cmdr.readCharts("/Users/alfiankan/development/repack/commander/charts")

	for _, v := range charts {
		fmt.Println(v)
	}

	assert.True(t, len(charts) > 0)
}
