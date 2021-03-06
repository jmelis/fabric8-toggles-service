package design

import (
	d "github.com/goadesign/goa/design"
	a "github.com/goadesign/goa/design/apidsl"
)

var featureSingle = JSONSingle(
	"Feature", "Holds a single feature",
	feature,
	nil)

var featureList = JSONList(
	"Feature", "Holds the list of features",
	feature,
	nil,
	nil)

var feature = a.Type("Feature", func() {
	a.Description(`JSONAPI for the feature object. See also http://jsonapi.org/format/#document-resource-object`)
	a.Attribute("id", d.String, "Id of feature", func() {
		a.Example("Feature name")
	})
	a.Attribute("type", d.String, "the 'features' type", func() {
		a.Example("features")
	})

	a.Attribute("attributes", featureAttributes)
	a.Required("id", "type", "attributes")
})

var featureAttributes = a.Type("FeatureAttributes", func() {
	a.Description(`JSONAPI store for all the "attributes" of a Feature. See also see http://jsonapi.org/format/#document-resource-object-attributes`)
	a.Attribute("description", d.String, "The description of the feature", func() {
		a.Example("Description of the feature")
	})
	a.Attribute("enabled", d.Boolean, "marks if the feature is globally enabled (prior to applying strategies)", func() {
		a.Example(true)
	})
	a.Attribute("user-enabled", d.Boolean, "marks if the feature is enabled for the current user", func() {
		a.Example(true)
	})
	a.Attribute("enablement-level", d.String, "The mimimum level of enablement for this feature. Empty/missing means that the feature is not accessible to the user", func() {
		a.Example("beta")
	})
	a.Required("description", "enabled", "user-enabled")
})

var _ = a.Resource("features", func() {
	a.BasePath("/features")

	a.Action("show", func() {
		a.Routing(
			a.GET("/:featureName"),
		)
		a.Params(func() {
			a.Param("featureName", d.String, "featureName")
		})
		a.Description("Show feature details.")
		a.Response(d.OK, featureSingle)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})

	a.Action("list", func() {
		a.Routing(
			a.GET(""),
		)
		a.Params(func() {
			a.Param("names", a.ArrayOf(d.String), "names")
		})
		a.Description("Show a list of features by their names.")
		a.Response(d.OK, featureList)
		a.Response(d.BadRequest, JSONAPIErrors)
		a.Response(d.Unauthorized, JSONAPIErrors)
		a.Response(d.NotFound, JSONAPIErrors)
		a.Response(d.InternalServerError, JSONAPIErrors)
	})
})
