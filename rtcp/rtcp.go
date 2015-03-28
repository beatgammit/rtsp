package rtcp

import (
	"io"
	"io/ioutil"
)

func Handle(r io.Reader) error {
	io.Copy(ioutil.Discard, r)
	return nil
}
