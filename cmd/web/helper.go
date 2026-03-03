package main

import (
	"errors"
	"fmt"

	"github.com/go-playground/form/v4"
	"github.com/labstack/echo/v4"
)

// decodePostForm decodes the POST form into the given dst struct.
func (b *backend) decodePostForm(c echo.Context, dst any) error {
	if err := c.Request().ParseForm(); err != nil {
		return err
	}
	err := b.decoder.Decode(dst, c.Request().PostForm)
	if err != nil {
		decErr, ok := errors.AsType[*form.InvalidDecoderError](err)
		if ok {
			panic(decErr.Error())
		}
		return fmt.Errorf("couldn't decode form: %w", err)
	}
	return nil
}
