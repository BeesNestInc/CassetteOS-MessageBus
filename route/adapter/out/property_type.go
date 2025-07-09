package out

import (
	"github.com/BeesNestInc/CassetteOS-MessageBus/codegen"
	"github.com/BeesNestInc/CassetteOS-MessageBus/model"
)

func PropertyTypeAdapter(propertyType model.PropertyType) codegen.PropertyType {
	return codegen.PropertyType{
		Name: propertyType.Name,
	}
}
