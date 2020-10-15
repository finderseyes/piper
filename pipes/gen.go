package pipes

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"regexp"

	"github.com/finderseyes/piper/pipes/io"

	"github.com/dave/jennifer/jen"
	"github.com/pkg/errors"
)

var pipeRegex = regexp.MustCompile(`//\s?@pipe.*`)

const (
	defaultInvocationFunctionName = "Run"
	piperOutputFileName           = "piper_gen.go"
	piperErrorsPackage            = "github.com/finderseyes/piper/errors"
)

const (
	// StageTypeFunction ...
	StageTypeFunction = iota
	// StageTypeFunctor ...
	StageTypeFunctor
)

// StageType ...
type StageType int

type stage struct {
	name         string
	stateType    StageType
	functor      *Functor
	returnsError bool
	canSkip      bool
}

// pipeInfo contains information about a pipe.
type pipeInfo struct {
	returnsError           bool
	stages                 []*stage
	invocationFunctionName string
}

// Functor ...
type Functor struct {
	name      string
	signature *types.Signature
}

// Generator ...
type Generator struct {
	path          string
	file          *jen.File
	info          *types.Info
	functors      map[*types.Named]*Functor
	writerFactory io.WriterFactory
}

// NewGenerator ...
func NewGenerator(path string, opts ...Option) *Generator {
	generator := &Generator{
		path:     path,
		functors: map[*types.Named]*Functor{},
	}

	for _, opt := range opts {
		opt(generator)
	}

	if generator.writerFactory == nil {
		generator.writerFactory = io.NewStringWriterFactory()
	}

	return generator
}

// Execute ...
func (g *Generator) Execute() error {
	if err := g.ensureDir(); err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, g.path, func(info os.FileInfo) bool {
		return info.Name() != "piper_gen.go"
	}, parser.ParseComments)
	if err != nil {
		return err
	}

	for name, pkg := range packages {
		piperGenFile := path.Join(g.path, piperOutputFileName)
		writer, err := g.writerFactory.CreateWriter(piperGenFile)
		if err != nil {
			return err
		}

		err = g.generatePackage(writer, fileSet, name, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) ensureDir() error {
	stat, err := os.Stat(g.path)
	if err != nil {
		return errors.Wrapf(err, "failed to retrieve stats for given path: %s", g.path)
	}

	if !stat.IsDir() {
		return errors.Errorf("not a directory: %s", g.path)
	}

	return nil
}

//
func (g *Generator) generatePackage(
	w io.ClosableWriter,
	fileSet *token.FileSet,
	pkgName string,
	pkg *ast.Package,
) error {
	// REF: https://stackoverflow.com/questions/55377694/how-to-find-full-package-import-from-callexpr
	files := make([]*ast.File, 0)
	for _, file := range pkg.Files {
		files = append(files, file)
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
		InitOrder:  make([]*types.Initializer, 0),
	}

	conf := types.Config{Importer: importer.ForCompiler(fileSet, "source", nil)}
	// TODO: probably need to handle error.
	// I decided to skip error checking since it most likely contains legitimate error, such as
	// when testing include a Pipe run, its Run function is not generated in the 1st place and hence checking
	// returns error.
	_, _ = conf.Check(pkgName, fileSet, files, info)

	// f := jen.NewFile(pkgName)
	f := jen.NewFilePath(pkgName)
	f.HeaderComment("Code generated by Piper. DO NOT EDIT.")
	// f.PackageComment("+build !piper")

	g.info = info
	g.file = f

	for _, file := range pkg.Files {
		// ast.Walk(NewPipeGenerator(f, info), file)
		ast.Walk(g, file)
	}

	_, _ = fmt.Fprintf(w, "%#v", f)

	return w.Close()
}

// Visit visits each node in an ASP tree.
func (g *Generator) Visit(node ast.Node) ast.Visitor {
	if decl, ok := node.(*ast.GenDecl); ok {
		if !g.isPipe(decl) {
			return g
		}

		for _, spec := range decl.Specs {
			if spec, ok := spec.(*ast.TypeSpec); ok {
				if structType, ok := spec.Type.(*ast.StructType); ok {
					g.genPipe(spec, structType)
				}
			}
		}
	}

	return g
}

func (g *Generator) isPipe(node *ast.GenDecl) bool {
	if node.Doc == nil {
		return false
	}

	for _, c := range node.Doc.List {
		if pipeRegex.MatchString(c.Text) {
			return true
		}
	}

	return false
}

// nolint:funlen,gocyclo // disable funlen and gocyclo linting since this is a complicated function.
func (g *Generator) genPipe(pipeSpec *ast.TypeSpec, pipeStruct *ast.StructType) {
	pipeInfo := g.getPipeInfo(pipeStruct)
	stages := pipeInfo.stages
	varCount := 0

	funcStmt := g.file.Func().Params(jen.Id("p").Op("*").Id(pipeSpec.Name.Name)).Id(pipeInfo.invocationFunctionName)

	// Generate input params.
	var vars []string
	{
		firstStage := stages[0]
		params := firstStage.functor.signature.Params()
		vars = make([]string, params.Len())

		for i := 0; i < params.Len(); i++ {
			vars[i] = fmt.Sprintf("v%d", varCount)
			varCount++
		}

		funcStmt.ParamsFunc(func(group *jen.Group) {
			for i := 0; i < params.Len(); i++ {
				g.getQualifiedName(group.Id(vars[i]), params.At(i).Type())
			}
		})
	}

	// Generate return types.
	{
		lastStage := stages[len(stages)-1]
		results := lastStage.functor.signature.Results()
		returnTypes := make([]jen.Code, 0, results.Len())

		for i := 0; i < results.Len(); i++ {
			t := results.At(i).Type()
			returnTypes = append(returnTypes, g.getQualifiedName(jen.Empty(), t))
		}

		if pipeInfo.returnsError && !lastStage.returnsError {
			returnTypes = append(returnTypes, jen.Error())
		}

		if len(returnTypes) == 1 {
			funcStmt.List(returnTypes...)
		} else if len(returnTypes) > 1 {
			funcStmt.Parens(jen.List(returnTypes...))
		}
	}

	// Determine if there's any stage can skip.
	{
		lastStage := stages[len(stages)-1]
		for i := 0; i < len(stages); i++ {
			stage := stages[i]
			stage.canSkip = g.canSkip(stage, lastStage)
		}
	}

	funcStmt.BlockFunc(func(group *jen.Group) {
		var resultVars []string
		lastStage := stages[len(stages)-1]

		for i := 0; i < len(stages); i++ {
			stage := stages[i]
			{
				results := stage.functor.signature.Results()
				resultsLen := results.Len()

				if stage.returnsError {
					resultsLen--
				}
				resultVars = make([]string, resultsLen)

				var lhs *jen.Statement

				// The function basically returns void
				if resultsLen <= 0 && !stage.returnsError {
					lhs = group.Empty().Id("p")
				} else {
					lhs = group.ListFunc(func(group *jen.Group) {
						for j := 0; j < resultsLen; j++ {
							resultVars[j] = fmt.Sprintf("v%d", varCount)
							group.Id(resultVars[j])
							varCount++
						}

						if stage.returnsError {
							group.Err()
						}
					}).Op(":=").Id("p")
				}

				var assignment *jen.Statement
				switch stage.stateType {
				case StageTypeFunction:
					assignment = lhs.Dot(stages[i].name)
				case StageTypeFunctor:
					assignment = lhs.Dot(stages[i].name).Dot(stages[i].functor.name)
				default:
					panic("should not reach here.")
				}

				assignment.CallFunc(func(group *jen.Group) {
					for j := 0; j < len(vars); j++ {
						group.Id(vars[j])
					}
				})

				vars = resultVars

				if stage.returnsError {
					group.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
						group.If(jen.Qual(piperErrorsPackage, "IsSkipped").Call(jen.Err())).BlockFunc(func(group *jen.Group) {
							if !stage.canSkip {
								// Last stage.
								if stage == lastStage {
									g.genStageReturns(stage, group, func(group *jen.Group, i int, isError bool) {
										if isError {
											group.Nil()
										} else {
											group.Id(resultVars[i])
										}
									})
								} else {
									g.genPipeReturnsDefault(pipeInfo, group, func(group *jen.Group) {
										group.
											Qual(piperErrorsPackage, "CannotSkip").
											Call(jen.Lit(stage.name), jen.Err())
									})
								}
							} else {
								g.genStageReturns(lastStage, group, func(group *jen.Group, j int, isError bool) {
									if isError {
										group.Err()
									} else {
										group.Id(resultVars[j])
									}
								})
							}
						}).Line()

						g.genPipeReturnsDefault(pipeInfo, group, func(group *jen.Group) {
							group.
								Qual(piperErrorsPackage, "NewError").
								Call(jen.Lit(stage.name), jen.Err())
						})
					}).Line()
				}
			}
		}

		if len(resultVars) > 0 || pipeInfo.returnsError {
			group.ReturnFunc(func(group *jen.Group) {
				for i := 0; i < len(resultVars); i++ {
					group.Id(resultVars[i])
				}

				if pipeInfo.returnsError {
					group.Nil()
				}
			})
		}
	}).Line()
}

func (g *Generator) genPipeReturnsDefault(pipeInfo *pipeInfo, group *jen.Group, withError func(group *jen.Group)) {
	stages := pipeInfo.stages
	lastStage := stages[len(stages)-1]
	results := lastStage.functor.signature.Results()

	g.genStageReturns(lastStage, group, func(group *jen.Group, i int, isError bool) {
		if isError {
			withError(group)
		} else {
			g.genLit(group, results.At(i).Type())
		}
	})
}

func (g *Generator) genStageReturns(stage *stage, group *jen.Group,
	returnParam func(group *jen.Group, i int, isError bool),
) {
	results := stage.functor.signature.Results()
	resultsLen := results.Len()
	if stage.returnsError {
		resultsLen--
	}

	group.ReturnFunc(func(group *jen.Group) {
		for k := 0; k < resultsLen; k++ {
			returnParam(group, k, false)
		}

		returnParam(group, resultsLen, true)
	})
}

func (g *Generator) genLit(group *jen.Group, t types.Type) {
	switch t := t.(type) {
	case *types.Named:
		group.Id(t.Obj().Name()).Parens(jen.LitFunc(g.getZeroLit(t.Underlying())))
	case *types.Pointer:
		group.Nil()
	case *types.Basic:
		group.LitFunc(g.getZeroLit(t))
	default:
		panic("should not reach here")
	}
}

func (g *Generator) getZeroLit(t types.Type) func() interface{} {
	return func() interface{} {
		switch t := t.(type) {
		case *types.Named:
			return g.getZeroLit(t.Underlying())()
		case *types.Basic:
			// nolint:gomnd,exhaustive // disable since this is correct
			switch t.Kind() {
			case types.String:
				return ""
			case types.Float32, types.Float64:
				return 0.0
			default:
				return 0
			}
		default:
			panic("should not reach here")
		}
	}
}

func (g *Generator) getQualifiedName(stmt *jen.Statement, t types.Type) *jen.Statement {
	switch t := t.(type) {
	case *types.Named:
		pkg := t.Obj().Pkg()
		if pkg != nil {
			return stmt.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
		}

		return stmt.Id(t.Obj().Name())
	case *types.Basic:
		return stmt.Id(t.Name())
	case *types.Pointer:
		return g.getQualifiedName(stmt.Op("*"), t.Elem())
	}

	panic("cannot reach here")
}

func (g *Generator) getFunctor(name *types.Named) *Functor {
	functor, ok := g.functors[name]
	if !ok {
		if underlyingType, ok := name.Underlying().(*types.Interface); ok {
			if underlyingType.NumMethods() != 1 {
				return nil
			}

			method := underlyingType.Method(0)
			signature, ok := method.Type().(*types.Signature)
			if !ok {
				return nil
			}

			functor = &Functor{
				name:      method.Name(),
				signature: signature,
			}

			g.functors[name] = functor
			return functor
		}
	}

	return functor
}

// nolint:funlen,gocyclo // disable since this function is complicated.
func (g *Generator) getPipeInfo(pipeStruct *ast.StructType) *pipeInfo {
	fields := pipeStruct.Fields.List
	stages := make([]*stage, 0, len(fields))
	returnsError := false
	invocationFunctionName := defaultInvocationFunctionName

	for _, field := range fields {
		fieldType, ok := g.info.Types[field.Type]
		if !ok {
			panic("should not reach here.")
		}

		var functor *Functor = nil
		var stageType StageType

		switch fieldType := fieldType.Type.(type) {
		case *types.Signature:
			functor = &Functor{
				name:      "",
				signature: fieldType,
			}
			stageType = StageTypeFunction
		case *types.Named:
			functor = g.getFunctor(fieldType)
			stageType = StageTypeFunctor
		default:
			panic("should not reach here.")
		}

		for _, id := range field.Names {
			results := functor.signature.Results()
			var isError = make([]bool, results.Len())
			errorCount := 0
			for i := 0; i < results.Len(); i++ {
				rt := results.At(i).Type()

				if rt, ok := rt.(*types.Named); ok {
					if rt.Obj().Pkg() == nil && rt.Obj().Name() == "error" {
						isError[i] = true
						errorCount++
					}
				}
			}

			stageReturnsError, err := g.doesReturnError(functor)
			if err != nil {
				panic("should not reach here.")
			}

			if stageReturnsError {
				returnsError = true
			}

			stages = append(stages, &stage{
				name:         id.Name,
				stateType:    stageType,
				functor:      functor,
				returnsError: stageReturnsError,
			})
		}
	}

	return &pipeInfo{
		returnsError:           returnsError,
		stages:                 stages,
		invocationFunctionName: invocationFunctionName,
	}
}

// doesReturnFunctor checks if a Functor returns error parameter.
// Piper only support functors with no error parameters or one error parameters as the last parameter.
func (g *Generator) doesReturnError(functor *Functor) (bool, error) {
	results := functor.signature.Results()
	returnsError := false

	for i := 0; i < results.Len(); i++ {
		rt := results.At(i).Type()

		if rt, ok := rt.(*types.Named); ok {
			if rt.Obj().Pkg() == nil && rt.Obj().Name() == "error" {
				if returnsError {
					return false, errors.New("error parameter must be the last parameter")
				}

				returnsError = true
			}
		}
	}

	return returnsError, nil
}

func (g *Generator) canSkip(stage, lastStage *stage) bool {
	if stage == lastStage || !stage.returnsError {
		return false
	}

	stageFunctor := stage.functor
	lastStageFunctor := lastStage.functor

	resultsLen := stageFunctor.signature.Results().Len()
	lastResultsLen := lastStageFunctor.signature.Results().Len()

	// Ignore error params.
	if stage.returnsError {
		resultsLen--
	}

	// Ignore error params.
	if lastStage.returnsError {
		lastResultsLen--
	}

	// NOTE: a stage can be skipped if its return params is a subset of last stage return params.
	if lastResultsLen > resultsLen {
		return false
	}

	for i := 0; i < lastResultsLen; i++ {
		a := stageFunctor.signature.Results().At(i)
		b := lastStageFunctor.signature.Results().At(i)

		if a.Type() != b.Type() {
			return false
		}
	}

	return true
}
