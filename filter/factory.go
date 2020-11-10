package filter

import "strings"

/**
chain:
    name: xx接口处理流程
    filters:
      - filter: header_remover
        parameters:
          - {source: target: config:}
*/
type FilterChainDef struct {
	Name    string
	Filters []FilterDef
}

type FilterDef struct {
	Filter     string
	Parameters []SourceTarget
}

func Factory(def FilterChainDef) FilterChain {
	chain := FilterChain{}
	for _, v := range def.Filters {
		chain = append(chain, buildFilter(v))
	}

	return chain
}

func buildFilter(def FilterDef) Filter {
	var filter Filter
	filterName := strings.ToLower(def.Filter)
	switch filterName {
	case "header_remover":
		filter = &HeaderRemover{Parameters: def.Parameters}
	case "header_adder":
		filter = &HeaderAdd{Parameters: def.Parameters}
	case "parameter_adder":
		filter = &ParameterAdder{Parameters: def.Parameters}
	case "parameter_remover":
		filter = &ParamterRemover{Parameters: def.Parameters}
	case "parameter_rename":
		filter = &ParamterRename{Parameters: def.Parameters}

	}
	return filter
}
