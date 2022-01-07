package astream

type NodeType string

const (
	opStream    NodeType = "Stream"
	opForeach   NodeType = "Foreach"
	opMap       NodeType = "Map"
	opFlatMap   NodeType = "FlatMap"
	opReduce    NodeType = "Reduce"
	opFilter    NodeType = "Filter"
	opSort      NodeType = "Sort"
	opDistinct  NodeType = "Distinct"
	opLimit     NodeType = "Limit"
	opSkip      NodeType = "Skip"
	opAllMatch  NodeType = "AllMatch"
	opAnyMatch  NodeType = "AnyMatch"
	opNoneMatch NodeType = "NoneMatch"
	opEmpty     NodeType = "Empty"
	opCollect   NodeType = "Collect"
	opCount     NodeType = "Count"
	opMax       NodeType = "Max"
	opMin       NodeType = "Min"
)

var handleFuncMap = map[NodeType]func(node *flowNode, resultChan flowResultChan){
	opStream:    handleStream,
	opForeach:   handleForEach,
	opMap:       handleMap,
	opFlatMap:   handleFlatMap,
	opReduce:    handleReduce,
	opFilter:    handleFilter,
	opSort:      handleSort,
	opDistinct:  handleDistinct,
	opLimit:     handleLimit,
	opSkip:      handleSkip,
	opAllMatch:  handleAllMatch,
	opAnyMatch:  handleAnyMatch,
	opNoneMatch: handleNoneMatch,
	opEmpty:     handleEmpty,
	opCollect:   handleCollect,
	opCount:     handleCount,
	opMax:       handleMax,
	opMin:       handleMin,
}

type ForEachFunc func(a interface{})
type MapFunc func(a interface{}) interface{}
type FlatMapFunc func(a interface{}) *Flow
type FilterFunc func(a interface{}) bool
type SortFunc func(a, b interface{}) bool
type MatchFunc func(a interface{}) bool
type ReduceFunc func(a, b interface{}) interface{}
type MaxFunc func(a, b interface{}) interface{}
type MinFunc func(a, b interface{}) interface{}
