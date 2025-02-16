@extern(embed)
@if(step3)
package holos

config: _

// Normalized values with no overrides.
component: Values: flattenedValues["flattened-values/customers/\(config.customer)/clusters/\(config.cluster)/values.yaml"]

// Load the flattened value files into CUE.
flattenedValues: _ @embed(glob=flattened-values/customers/*/clusters/*/values.yaml)
