package job

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	internaljob "d7y.io/dragonfly/v2/internal/job"
	"d7y.io/dragonfly/v2/manager/types"
)

func TestPreheat_GetImageLayers(t *testing.T) {
	tests := []struct {
		name   string
		args   types.PreheatArgs
		expect func(t *testing.T, layers []internaljob.PreheatRequest)
	}{
		{
			name: "get image layers with manifest url",
			args: types.PreheatArgs{
				URL:  "https://registry-1.docker.io/v2/dragonflyoss/busybox/manifests/1.35.0",
				Type: "image",
			},
			expect: func(t *testing.T, layers []internaljob.PreheatRequest) {
				assert := assert.New(t)
				assert.Equal(2, len(layers[0].URLs))
			},
		},
		{
			name: "get image layers with multi arch image layers",
			args: types.PreheatArgs{
				URL:      "https://registry-1.docker.io/v2/dragonflyoss/scheduler/manifests/v2.1.0",
				Platform: "linux/amd64",
			},
			expect: func(t *testing.T, layers []internaljob.PreheatRequest) {
				assert := assert.New(t)
				assert.Equal(5, len(layers[0].URLs))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			layers, err := GetImageLayers(context.Background(), tc.args, 30*time.Second, nil, true)
			if err != nil {
				t.Fatal(err)
			}

			tc.expect(t, layers)
		})
	}
}
