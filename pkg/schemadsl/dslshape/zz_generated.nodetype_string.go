// Code generated by "stringer -type=NodeType -output zz_generated.nodetype_string.go"; DO NOT EDIT.

package dslshape

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NodeTypeError-0]
	_ = x[NodeTypeFile-1]
	_ = x[NodeTypeComment-2]
	_ = x[NodeTypeUseFlag-3]
	_ = x[NodeTypeDefinition-4]
	_ = x[NodeTypeCaveatDefinition-5]
	_ = x[NodeTypeCaveatParameter-6]
	_ = x[NodeTypeCaveatExpression-7]
	_ = x[NodeTypeRelation-8]
	_ = x[NodeTypePermission-9]
	_ = x[NodeTypeTypeReference-10]
	_ = x[NodeTypeSpecificTypeReference-11]
	_ = x[NodeTypeCaveatReference-12]
	_ = x[NodeTypeTraitReference-13]
	_ = x[NodeTypeUnionExpression-14]
	_ = x[NodeTypeIntersectExpression-15]
	_ = x[NodeTypeExclusionExpression-16]
	_ = x[NodeTypeArrowExpression-17]
	_ = x[NodeTypeIdentifier-18]
	_ = x[NodeTypeNilExpression-19]
	_ = x[NodeTypeCaveatTypeReference-20]
}

const _NodeType_name = "NodeTypeErrorNodeTypeFileNodeTypeCommentNodeTypeUseFlagNodeTypeDefinitionNodeTypeCaveatDefinitionNodeTypeCaveatParameterNodeTypeCaveatExpressionNodeTypeRelationNodeTypePermissionNodeTypeTypeReferenceNodeTypeSpecificTypeReferenceNodeTypeCaveatReferenceNodeTypeTraitReferenceNodeTypeUnionExpressionNodeTypeIntersectExpressionNodeTypeExclusionExpressionNodeTypeArrowExpressionNodeTypeIdentifierNodeTypeNilExpressionNodeTypeCaveatTypeReference"

var _NodeType_index = [...]uint16{0, 13, 25, 40, 55, 73, 97, 120, 144, 160, 178, 199, 228, 251, 273, 296, 323, 350, 373, 391, 412, 439}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
