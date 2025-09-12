package services

import (
	"encoding/xml"

	"httpfromtcp/rootmod/internal/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type BussinesValidatable interface {
	Validate() error
}

type Parserable[S any] interface {
	ParseXmlSync(*gin.Context) (S, error)
}

type Parser[S any] struct {
}

func (p *Parser[S]) ParseXmlSync(c *gin.Context) error {

	var result S
	dat := xml.NewDecoder(c.Request.Body)
	if err := dat.Decode(&result); err != nil {
		errors.WriteAPIError(c, err)
		return err
	}

	if err := validate.Struct(result); err != nil {
		errors.WriteAPIError(c, err)
		return err
	}

	if vt, ok := any(&result).(BussinesValidatable); ok {
		if err := vt.Validate(); err != nil {
			errors.WriteAPIError(c, err)
			return err
		}
	}

	c.BindXML(result)

	return nil
}
