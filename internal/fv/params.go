package fv

import "io"

type Params struct {
	Body *io.Reader
	ID string
}
