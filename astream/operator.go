package astream

const (
	STREAM    NodeType = "Stream"
	FOREACH   NodeType = "FOREACH"
	MAP       NodeType = "Map"
	FILTER    NodeType = "FILTER"
	SORT      NodeType = "SORT"
	DISTINCT  NodeType = "DISTINCT"
	LIMIT     NodeType = "LIMIT"
	SKIP      NodeType = "SKIP"
	ALLMATCH  NodeType = "ALLMATCH"
	ANYMATCH  NodeType = "ANYMATCH"
	NONEMATCH NodeType = "NONEMATCH"
	COLLECT   NodeType = "COLLECT"
)

var handleFuncMap = map[NodeType]func(node *flowNode, resultChan FlowResultChan){
	STREAM:    handleStream,
	FOREACH:   handleForEach,
	MAP:       handleMap,
	FILTER:    handleFilter,
	SORT:      handleSort,
	DISTINCT:  handleDistinct,
	LIMIT:     handleLimit,
	SKIP:      handleSkip,
	ALLMATCH:  handleAllMatch,
	ANYMATCH:  handleAnyMatch,
	NONEMATCH: handleNoneMatch,
	COLLECT:   handleCollect,
}

type ForEachFunc func(a interface{})
type MapFunc func(a interface{}) interface{}
type FilterFunc func(a interface{}) bool
type SortFunc func(a, b interface{}) bool
type MatchFunc func(a interface{}) bool
