/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

func replaceRefContainerASTType(newNode, parent AST) {
	parent.(*RefContainer).ASTType = newNode.(AST)
}
func replaceRefContainerASTImplementationType(newNode, parent AST) {
	parent.(*RefContainer).ASTImplementationType = newNode.(*Leaf)
}

type replaceRefSliceContainerASTElements int

func (r *replaceRefSliceContainerASTElements) replace(newNode, parent AST) {
	parent.(*RefSliceContainer).ASTElements[int(*r)] = newNode.(AST)
}
func (r *replaceRefSliceContainerASTElements) inc() {
	*r++
}

type replaceRefSliceContainerASTImplementationElements int

func (r *replaceRefSliceContainerASTImplementationElements) replace(newNode, parent AST) {
	parent.(*RefSliceContainer).ASTImplementationElements[int(*r)] = newNode.(*Leaf)
}
func (r *replaceRefSliceContainerASTImplementationElements) inc() {
	*r++
}
func replaceValueContainerASTType(newNode, parent AST) {
	parent.(*ValueContainer).ASTType = newNode.(AST)
}
func replaceValueContainerASTImplementationType(newNode, parent AST) {
	parent.(*ValueContainer).ASTImplementationType = newNode.(*Leaf)
}

func replaceValValueSliceContainerASTElements(idx int) func(newNode, parent AST) {
	return func(newNode, parent AST) {
		parent.(ValueSliceContainer).ASTElements[idx] = newNode.(AST)
	}
}

type replaceValueSliceContainerValASTImplementationElements int

func (r *replaceValueSliceContainerValASTImplementationElements) replace(newNode, parent AST) {
	parent.(ValueSliceContainer).ASTImplementationElements[int(*r)] = newNode.(*Leaf)
}
func (r *replaceValueSliceContainerValASTImplementationElements) inc() {
	*r++
}

type replaceValueSliceContainerASTElements int

func (r *replaceValueSliceContainerASTElements) replace(newNode, parent AST) {
	parent.(*ValueSliceContainer).ASTElements[int(*r)] = newNode.(AST)
}
func (r *replaceValueSliceContainerASTElements) inc() {
	*r++
}

type replaceValueSliceContainerASTImplementationElements int

func (r *replaceValueSliceContainerASTImplementationElements) replace(newNode, parent AST) {
	parent.(*ValueSliceContainer).ASTImplementationElements[int(*r)] = newNode.(*Leaf)
}
func (r *replaceValueSliceContainerASTImplementationElements) inc() {
	*r++
}
func (a *application) apply(parent, node AST, replacer replacerFunc) {
	if node == nil || isNilValue(node) {
		return
	}
	saved := a.cursor
	a.cursor.replacer = replacer
	a.cursor.node = node
	a.cursor.parent = parent
	if a.pre != nil && !a.pre(&a.cursor) {
		a.cursor = saved
		return
	}
	switch n := node.(type) {
	case *Leaf:
	case *RefContainer:
		a.apply(node, n.ASTType, replaceRefContainerASTType)
		a.apply(node, n.ASTImplementationType, replaceRefContainerASTImplementationType)
	case *RefSliceContainer:
		replacerASTElements := replaceRefSliceContainerASTElements(0)
		replacerASTElementsB := &replacerASTElements
		for _, item := range n.ASTElements {
			a.apply(node, item, replacerASTElementsB.replace)
			replacerASTElementsB.inc()
		}
		replacerASTImplementationElements := replaceRefSliceContainerASTImplementationElements(0)
		replacerASTImplementationElementsB := &replacerASTImplementationElements
		for _, item := range n.ASTImplementationElements {
			a.apply(node, item, replacerASTImplementationElementsB.replace)
			replacerASTImplementationElementsB.inc()
		}
	case ValueContainer:
		a.apply(node, n.ASTType, replacePanic("ValueContainer ASTType"))
		a.apply(node, n.ASTImplementationType, replacePanic("ValueContainer ASTImplementationType"))
	case *ValueContainer:
		a.apply(node, n.ASTType, replaceValueContainerASTType)
		a.apply(node, n.ASTImplementationType, replaceValueContainerASTImplementationType)
	case ValueSliceContainer:
		for idx, item := range n.ASTElements {
			a.apply(node, item, replaceValValueSliceContainerASTElements(idx))
		}
		replacerASTImplementationElements := replaceValueSliceContainerValASTImplementationElements(0)
		replacerASTImplementationElementsB := &replacerASTImplementationElements
		for _, item := range n.ASTImplementationElements {
			a.apply(node, item, replacerASTImplementationElementsB.replace)
			replacerASTImplementationElementsB.inc()
		}
	case *ValueSliceContainer:
		replacerASTElements := replaceValueSliceContainerASTElements(0)
		replacerASTElementsB := &replacerASTElements
		for _, item := range n.ASTElements {
			a.apply(node, item, replacerASTElementsB.replace)
			replacerASTElementsB.inc()
		}
		replacerASTImplementationElements := replaceValueSliceContainerASTImplementationElements(0)
		replacerASTImplementationElementsB := &replacerASTImplementationElements
		for _, item := range n.ASTImplementationElements {
			a.apply(node, item, replacerASTImplementationElementsB.replace)
			replacerASTImplementationElementsB.inc()
		}
	}
	if a.post != nil && !a.post(&a.cursor) {
		panic(abort)
	}
	a.cursor = saved
}
