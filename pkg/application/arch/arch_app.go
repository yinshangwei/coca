package arch

import (
	"github.com/phodal/coca/pkg/application/arch/tequila"
	"github.com/phodal/coca/pkg/domain"
)

type ArchApp struct {
}

func NewArchApp() ArchApp {
	return ArchApp{}
}

func (a ArchApp) Analysis(deps []domain.JClassNode, identifiersMap map[string]domain.JIdentifier) *tequila.FullGraph {
	fullGraph := &tequila.FullGraph{
		NodeList:     make(map[string]string),
		RelationList: make(map[string]*tequila.Relation),
	}

	for _, clz := range deps {
		if clz.Class == "Main" {
			continue
		}

		src := clz.Package + "." + clz.Class
		fullGraph.NodeList[src] = src

		for _, impl := range clz.Implements {
			relation := &tequila.Relation{
				From:  src,
				To:    impl,
				Style: "\"solid\"",
			}

			fullGraph.RelationList[relation.From+"->"+relation.To] = relation
		}

		addCallInField(clz, src, *fullGraph)
		addExtend(clz, src, *fullGraph)
		addCallInMethod(clz, identifiersMap, src, *fullGraph)
	}

	return fullGraph
}

func addCallInField(clz domain.JClassNode, src string, fullGraph tequila.FullGraph) {
	for _, field := range clz.MethodCalls {
		dst := field.Package + "." + field.Class
		relation := &tequila.Relation{
			From:  src,
			To:    dst,
			Style: "\"solid\"",
		}

		fullGraph.RelationList[relation.From+"->"+relation.To] = relation
	}
}

func addCallInMethod(clz domain.JClassNode, identifiersMap map[string]domain.JIdentifier, src string, fullGraph tequila.FullGraph) {
	for _, method := range clz.Methods {
		if method.Name == "main" {
			continue
		}

		// TODO: add implements, extends support
		for _, call := range method.MethodCalls {
			dst := call.Package + "." + call.Class
			if src == dst {
				continue
			}

			if _, ok := identifiersMap[dst]; ok {
				relation := &tequila.Relation{
					From:  src,
					To:    dst,
					Style: "\"solid\"",
				}

				fullGraph.RelationList[relation.From+"->"+relation.To] = relation
			}
		}
	}
}

func addExtend(clz domain.JClassNode, src string, fullGraph tequila.FullGraph) {
	if clz.Extend != "" {
		relation := &tequila.Relation{
			From:  src,
			To:    clz.Extend,
			Style: "\"solid\"",
		}

		fullGraph.RelationList[relation.From+"->"+relation.To] = relation
	}
}
