package pipes

import (
	"regexp"
)

var pipeRegex = regexp.MustCompile("\\/\\/\\s?@pipe.*")
const defaultFunctionName = "Invoke"

const (
	StageTypeFunction = iota
	StageTypeFunctor
)

type StageType int

type stage struct {
	name      string
	stateType StageType
	functor *Functor
	paramIndices []int
}
//
//// Pipe generator.
//type pipeGenerator struct {
//	file *jen.File
//	info *types.Info
//	functors map[*types.Named]*types.Signature
//}
//
//func NewPipeGenerator(file *jen.File, info *types.Info) *pipeGenerator {
//	return &pipeGenerator{
//		file:     file,
//		info:     info,
//		functors: map[*types.Named]*types.Signature{},
//	}
//}
//
//// Visit walks the AST.
//func (pg *pipeGenerator) Visit(node ast.Node) ast.Visitor {
//	switch node := node.(type) {
//	case *ast.GenDecl:
//		if !pg.isPipe(node) {
//			return pg
//		}
//
//		for _, spec := range node.Specs {
//			switch nodeSpec := spec.(type) {
//			case *ast.TypeSpec:
//				switch structType := nodeSpec.Type.(type) {
//				case *ast.StructType:
//					pg.genPipe(nodeSpec, structType)
//				}
//			}
//		}
//	}
//	return pg
//}
//
//func (pg *pipeGenerator) isPipe(node *ast.GenDecl) bool {
//	if node.Doc == nil {
//		return false
//	}
//
//	for _, c := range node.Doc.List {
//		if pipeRegex.MatchString(c.Text) {
//			return true
//		}
//	}
//
//	return false
//}
//
//func (pg *pipeGenerator) genPipe(nodeSpec *ast.TypeSpec, structType *ast.StructType) {
//	funcStmt := pg.file.Func().Params(jen.Id("p").Op("*").Id(nodeSpec.Name.Name)).Id(defaultFunctionName)
//
//	stages := make([]*stage, 0)
//	lastParamIndex := 0
//
//	for _, field := range structType.Fields.List {
//		fieldType, ok := pg.info.Types[field.Type]
//		if !ok {
//			continue
//		}
//
//		var signature *types.Signature = nil
//		var stageType StageType
//
//		switch fieldType := fieldType.Type.(type) {
//		case *types.Signature:
//			signature = fieldType
//			stageType = StageTypeFunction
//		case *types.Named:
//			signature = pg.getFunctor(fieldType)
//			stageType = StageTypeFunctor
//		}
//
//		if signature != nil {
//			for _, id := range field.Names {
//				var paramIndicies []int
//				params := signature.Params()
//				if params != nil {
//					paramIndicies = make([]int, params.Len())
//					for i := 0; i < params.Len(); i++ {
//						paramIndicies[i] = lastParamIndex + i
//					}
//					lastParamIndex += params.Len()
//				}
//
//				stages = append(stages, &stage{
//					name:      id.Name,
//					stateType: stageType,
//					signature: signature,
//					paramIndices: paramIndicies,
//				})
//			}
//		}
//
//
//		//fmt.Println(stages)
//		var _ = fieldType
//	}
//
//	params := stages[0].signature.Params()
//	//results := stages[len(stages) - 1].signature.Results()
//
//	prevInputs := make([]jen.Code, params.Len())
//
//	paramStmts := make([]jen.Code, 0)
//	for i := 0; i < params.Len(); i++ {
//		p := params.At(i)
//		t := p.Type()
//		pname := jen.Id(fmt.Sprintf("v%d", stages[0].paramIndices[i]))
//		paramStmts = append(paramStmts, pg.getQualifiedName(pname, t))
//		prevInputs[i] = jen.Id(fmt.Sprintf("v%d", stages[0].paramIndices[i]))
//	}
//
//	results := stages[len(stages) - 1].signature.Results()
//	resultTypes := make([]jen.Code, 0)
//	for i := 0; i < results.Len(); i++ {
//		p := results.At(i)
//		resultTypes = append(resultTypes, pg.getQualifiedName(jen.Empty(), p.Type()))
//	}
//
//	varCount := params.Len()
//	funcStmt.Params(paramStmts...).List(resultTypes...).BlockFunc(func(group *jen.Group) {
//		var output []jen.Code
//
//		for i := 0; i < len(stages); i++ {
//			results := stages[i].signature.Results()
//			output = make([]jen.Code, results.Len())
//			for j := 0; j < results.Len(); j++ {
//				output[j] = jen.Id(fmt.Sprintf("v%d", varCount))
//				varCount += 1
//			}
//			group.List(output...).Op(":=").Id(stages[i].name).Call(prevInputs...)
//
//			prevInputs = output
//		}
//
//		group.Return(output...)
//	})
//}
//
//func (pg *pipeGenerator) getQualifiedName(stmt *jen.Statement, t types.Type) *jen.Statement{
//	switch t := t.(type) {
//	case *types.Named:
//		return stmt.Qual(t.Obj().Pkg().Name(), t.Obj().Name())
//	case *types.Basic:
//		return stmt.Id(t.Name())
//	}
//
//	panic("cannot reach here")
//}
//
//func (pg *pipeGenerator) getFunctor(name *types.Named) *types.Signature {
//	signature, ok := pg.functors[name]
//	if !ok {
//		switch underlyingType := name.Underlying().(type) {
//		case *types.Interface:
//			if underlyingType.NumMethods() != 1 {
//				return nil
//			}
//
//			method := underlyingType.Method(0)
//			signature, ok = method.Type().(*types.Signature)
//			if !ok {
//				return nil
//			}
//
//			pg.functors[name] = signature
//			return signature
//		}
//	}
//
//	return signature
//}