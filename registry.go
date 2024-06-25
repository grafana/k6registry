// Package k6registry contains the data model of the k6 extensions registry.
package k6registry

//go:generate go run github.com/atombender/go-jsonschema@v0.16.0 -p k6registry --only-models -o registry_gen.go registry.schema.yaml
//go:generate sh -c "go run github.com/mikefarah/yq/v4@v4.44.2 -o=json -P registry.schema.yaml > registry.schema.json"
//go:generate go run github.com/szkiba/mdcode@v0.2.0 update
