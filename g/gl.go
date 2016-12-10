package g

import (
	"fmt"

	"github.com/goxjs/gl"
)

func glerror() error {
	if errcode := gl.GetError(); errcode != 0 {
		return fmt.Errorf("bind error: %d", errcode)
	}
	return nil
}
