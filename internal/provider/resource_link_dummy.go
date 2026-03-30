package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func dummySchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "Dummy interface (no configuration needed).",
		Attributes:  map[string]schema.Attribute{},
	}
}

func ifbSchemaBlock() schema.SingleNestedBlock {
	return schema.SingleNestedBlock{
		Description: "IFB (Intermediate Functional Block) interface (no configuration needed).",
		Attributes:  map[string]schema.Attribute{},
	}
}
