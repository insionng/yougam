// Code generated by protoc-gen-go.
// source: Filter.proto
// DO NOT EDIT!

package proto

import proto1 "github.com/insionng/yougam/libraries/golang/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = math.Inf

type FilterList_Operator int32

const (
	FilterList_MUST_PASS_ALL FilterList_Operator = 1
	FilterList_MUST_PASS_ONE FilterList_Operator = 2
)

var FilterList_Operator_name = map[int32]string{
	1: "MUST_PASS_ALL",
	2: "MUST_PASS_ONE",
}
var FilterList_Operator_value = map[string]int32{
	"MUST_PASS_ALL": 1,
	"MUST_PASS_ONE": 2,
}

func (x FilterList_Operator) Enum() *FilterList_Operator {
	p := new(FilterList_Operator)
	*p = x
	return p
}
func (x FilterList_Operator) String() string {
	return proto1.EnumName(FilterList_Operator_name, int32(x))
}
func (x *FilterList_Operator) UnmarshalJSON(data []byte) error {
	value, err := proto1.UnmarshalJSONEnum(FilterList_Operator_value, data, "FilterList_Operator")
	if err != nil {
		return err
	}
	*x = FilterList_Operator(value)
	return nil
}

type Filter struct {
	Name             *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	SerializedFilter []byte  `protobuf:"bytes,2,opt,name=serialized_filter" json:"serialized_filter,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto1.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}

func (m *Filter) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *Filter) GetSerializedFilter() []byte {
	if m != nil {
		return m.SerializedFilter
	}
	return nil
}

type ColumnCountGetFilter struct {
	Limit            *int32 `protobuf:"varint,1,req,name=limit" json:"limit,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ColumnCountGetFilter) Reset()         { *m = ColumnCountGetFilter{} }
func (m *ColumnCountGetFilter) String() string { return proto1.CompactTextString(m) }
func (*ColumnCountGetFilter) ProtoMessage()    {}

func (m *ColumnCountGetFilter) GetLimit() int32 {
	if m != nil && m.Limit != nil {
		return *m.Limit
	}
	return 0
}

type ColumnPaginationFilter struct {
	Limit            *int32 `protobuf:"varint,1,req,name=limit" json:"limit,omitempty"`
	Offset           *int32 `protobuf:"varint,2,opt,name=offset" json:"offset,omitempty"`
	ColumnOffset     []byte `protobuf:"bytes,3,opt,name=column_offset" json:"column_offset,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ColumnPaginationFilter) Reset()         { *m = ColumnPaginationFilter{} }
func (m *ColumnPaginationFilter) String() string { return proto1.CompactTextString(m) }
func (*ColumnPaginationFilter) ProtoMessage()    {}

func (m *ColumnPaginationFilter) GetLimit() int32 {
	if m != nil && m.Limit != nil {
		return *m.Limit
	}
	return 0
}

func (m *ColumnPaginationFilter) GetOffset() int32 {
	if m != nil && m.Offset != nil {
		return *m.Offset
	}
	return 0
}

func (m *ColumnPaginationFilter) GetColumnOffset() []byte {
	if m != nil {
		return m.ColumnOffset
	}
	return nil
}

type ColumnPrefixFilter struct {
	Prefix           []byte `protobuf:"bytes,1,req,name=prefix" json:"prefix,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *ColumnPrefixFilter) Reset()         { *m = ColumnPrefixFilter{} }
func (m *ColumnPrefixFilter) String() string { return proto1.CompactTextString(m) }
func (*ColumnPrefixFilter) ProtoMessage()    {}

func (m *ColumnPrefixFilter) GetPrefix() []byte {
	if m != nil {
		return m.Prefix
	}
	return nil
}

type ColumnRangeFilter struct {
	MinColumn          []byte `protobuf:"bytes,1,opt,name=min_column" json:"min_column,omitempty"`
	MinColumnInclusive *bool  `protobuf:"varint,2,opt,name=min_column_inclusive" json:"min_column_inclusive,omitempty"`
	MaxColumn          []byte `protobuf:"bytes,3,opt,name=max_column" json:"max_column,omitempty"`
	MaxColumnInclusive *bool  `protobuf:"varint,4,opt,name=max_column_inclusive" json:"max_column_inclusive,omitempty"`
	XXX_unrecognized   []byte `json:"-"`
}

func (m *ColumnRangeFilter) Reset()         { *m = ColumnRangeFilter{} }
func (m *ColumnRangeFilter) String() string { return proto1.CompactTextString(m) }
func (*ColumnRangeFilter) ProtoMessage()    {}

func (m *ColumnRangeFilter) GetMinColumn() []byte {
	if m != nil {
		return m.MinColumn
	}
	return nil
}

func (m *ColumnRangeFilter) GetMinColumnInclusive() bool {
	if m != nil && m.MinColumnInclusive != nil {
		return *m.MinColumnInclusive
	}
	return false
}

func (m *ColumnRangeFilter) GetMaxColumn() []byte {
	if m != nil {
		return m.MaxColumn
	}
	return nil
}

func (m *ColumnRangeFilter) GetMaxColumnInclusive() bool {
	if m != nil && m.MaxColumnInclusive != nil {
		return *m.MaxColumnInclusive
	}
	return false
}

type CompareFilter struct {
	CompareOp        *CompareType `protobuf:"varint,1,req,name=compare_op,enum=proto.CompareType" json:"compare_op,omitempty"`
	Comparator       *Comparator  `protobuf:"bytes,2,opt,name=comparator" json:"comparator,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *CompareFilter) Reset()         { *m = CompareFilter{} }
func (m *CompareFilter) String() string { return proto1.CompactTextString(m) }
func (*CompareFilter) ProtoMessage()    {}

func (m *CompareFilter) GetCompareOp() CompareType {
	if m != nil && m.CompareOp != nil {
		return *m.CompareOp
	}
	return CompareType_LESS
}

func (m *CompareFilter) GetComparator() *Comparator {
	if m != nil {
		return m.Comparator
	}
	return nil
}

type DependentColumnFilter struct {
	CompareFilter       *CompareFilter `protobuf:"bytes,1,req,name=compare_filter" json:"compare_filter,omitempty"`
	ColumnFamily        []byte         `protobuf:"bytes,2,opt,name=column_family" json:"column_family,omitempty"`
	ColumnQualifier     []byte         `protobuf:"bytes,3,opt,name=column_qualifier" json:"column_qualifier,omitempty"`
	DropDependentColumn *bool          `protobuf:"varint,4,opt,name=drop_dependent_column" json:"drop_dependent_column,omitempty"`
	XXX_unrecognized    []byte         `json:"-"`
}

func (m *DependentColumnFilter) Reset()         { *m = DependentColumnFilter{} }
func (m *DependentColumnFilter) String() string { return proto1.CompactTextString(m) }
func (*DependentColumnFilter) ProtoMessage()    {}

func (m *DependentColumnFilter) GetCompareFilter() *CompareFilter {
	if m != nil {
		return m.CompareFilter
	}
	return nil
}

func (m *DependentColumnFilter) GetColumnFamily() []byte {
	if m != nil {
		return m.ColumnFamily
	}
	return nil
}

func (m *DependentColumnFilter) GetColumnQualifier() []byte {
	if m != nil {
		return m.ColumnQualifier
	}
	return nil
}

func (m *DependentColumnFilter) GetDropDependentColumn() bool {
	if m != nil && m.DropDependentColumn != nil {
		return *m.DropDependentColumn
	}
	return false
}

type FamilyFilter struct {
	CompareFilter    *CompareFilter `protobuf:"bytes,1,req,name=compare_filter" json:"compare_filter,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *FamilyFilter) Reset()         { *m = FamilyFilter{} }
func (m *FamilyFilter) String() string { return proto1.CompactTextString(m) }
func (*FamilyFilter) ProtoMessage()    {}

func (m *FamilyFilter) GetCompareFilter() *CompareFilter {
	if m != nil {
		return m.CompareFilter
	}
	return nil
}

type FilterList struct {
	Operator         *FilterList_Operator `protobuf:"varint,1,req,name=operator,enum=proto.FilterList_Operator" json:"operator,omitempty"`
	Filters          []*Filter            `protobuf:"bytes,2,rep,name=filters" json:"filters,omitempty"`
	XXX_unrecognized []byte               `json:"-"`
}

func (m *FilterList) Reset()         { *m = FilterList{} }
func (m *FilterList) String() string { return proto1.CompactTextString(m) }
func (*FilterList) ProtoMessage()    {}

func (m *FilterList) GetOperator() FilterList_Operator {
	if m != nil && m.Operator != nil {
		return *m.Operator
	}
	return FilterList_MUST_PASS_ALL
}

func (m *FilterList) GetFilters() []*Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type FilterWrapper struct {
	Filter           *Filter `protobuf:"bytes,1,req,name=filter" json:"filter,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *FilterWrapper) Reset()         { *m = FilterWrapper{} }
func (m *FilterWrapper) String() string { return proto1.CompactTextString(m) }
func (*FilterWrapper) ProtoMessage()    {}

func (m *FilterWrapper) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type FirstKeyOnlyFilter struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *FirstKeyOnlyFilter) Reset()         { *m = FirstKeyOnlyFilter{} }
func (m *FirstKeyOnlyFilter) String() string { return proto1.CompactTextString(m) }
func (*FirstKeyOnlyFilter) ProtoMessage()    {}

type FirstKeyValueMatchingQualifiersFilter struct {
	Qualifiers       [][]byte `protobuf:"bytes,1,rep,name=qualifiers" json:"qualifiers,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *FirstKeyValueMatchingQualifiersFilter) Reset()         { *m = FirstKeyValueMatchingQualifiersFilter{} }
func (m *FirstKeyValueMatchingQualifiersFilter) String() string { return proto1.CompactTextString(m) }
func (*FirstKeyValueMatchingQualifiersFilter) ProtoMessage()    {}

func (m *FirstKeyValueMatchingQualifiersFilter) GetQualifiers() [][]byte {
	if m != nil {
		return m.Qualifiers
	}
	return nil
}

type FuzzyRowFilter struct {
	FuzzyKeysData    []*BytesBytesPair `protobuf:"bytes,1,rep,name=fuzzy_keys_data" json:"fuzzy_keys_data,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *FuzzyRowFilter) Reset()         { *m = FuzzyRowFilter{} }
func (m *FuzzyRowFilter) String() string { return proto1.CompactTextString(m) }
func (*FuzzyRowFilter) ProtoMessage()    {}

func (m *FuzzyRowFilter) GetFuzzyKeysData() []*BytesBytesPair {
	if m != nil {
		return m.FuzzyKeysData
	}
	return nil
}

type InclusiveStopFilter struct {
	StopRowKey       []byte `protobuf:"bytes,1,opt,name=stop_row_key" json:"stop_row_key,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *InclusiveStopFilter) Reset()         { *m = InclusiveStopFilter{} }
func (m *InclusiveStopFilter) String() string { return proto1.CompactTextString(m) }
func (*InclusiveStopFilter) ProtoMessage()    {}

func (m *InclusiveStopFilter) GetStopRowKey() []byte {
	if m != nil {
		return m.StopRowKey
	}
	return nil
}

type KeyOnlyFilter struct {
	LenAsVal         *bool  `protobuf:"varint,1,req,name=len_as_val" json:"len_as_val,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *KeyOnlyFilter) Reset()         { *m = KeyOnlyFilter{} }
func (m *KeyOnlyFilter) String() string { return proto1.CompactTextString(m) }
func (*KeyOnlyFilter) ProtoMessage()    {}

func (m *KeyOnlyFilter) GetLenAsVal() bool {
	if m != nil && m.LenAsVal != nil {
		return *m.LenAsVal
	}
	return false
}

type MultipleColumnPrefixFilter struct {
	SortedPrefixes   [][]byte `protobuf:"bytes,1,rep,name=sorted_prefixes" json:"sorted_prefixes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MultipleColumnPrefixFilter) Reset()         { *m = MultipleColumnPrefixFilter{} }
func (m *MultipleColumnPrefixFilter) String() string { return proto1.CompactTextString(m) }
func (*MultipleColumnPrefixFilter) ProtoMessage()    {}

func (m *MultipleColumnPrefixFilter) GetSortedPrefixes() [][]byte {
	if m != nil {
		return m.SortedPrefixes
	}
	return nil
}

type PageFilter struct {
	PageSize         *int64 `protobuf:"varint,1,req,name=page_size" json:"page_size,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PageFilter) Reset()         { *m = PageFilter{} }
func (m *PageFilter) String() string { return proto1.CompactTextString(m) }
func (*PageFilter) ProtoMessage()    {}

func (m *PageFilter) GetPageSize() int64 {
	if m != nil && m.PageSize != nil {
		return *m.PageSize
	}
	return 0
}

type PrefixFilter struct {
	Prefix           []byte `protobuf:"bytes,1,opt,name=prefix" json:"prefix,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *PrefixFilter) Reset()         { *m = PrefixFilter{} }
func (m *PrefixFilter) String() string { return proto1.CompactTextString(m) }
func (*PrefixFilter) ProtoMessage()    {}

func (m *PrefixFilter) GetPrefix() []byte {
	if m != nil {
		return m.Prefix
	}
	return nil
}

type QualifierFilter struct {
	CompareFilter    *CompareFilter `protobuf:"bytes,1,req,name=compare_filter" json:"compare_filter,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *QualifierFilter) Reset()         { *m = QualifierFilter{} }
func (m *QualifierFilter) String() string { return proto1.CompactTextString(m) }
func (*QualifierFilter) ProtoMessage()    {}

func (m *QualifierFilter) GetCompareFilter() *CompareFilter {
	if m != nil {
		return m.CompareFilter
	}
	return nil
}

type RandomRowFilter struct {
	Chance           *float32 `protobuf:"fixed32,1,req,name=chance" json:"chance,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *RandomRowFilter) Reset()         { *m = RandomRowFilter{} }
func (m *RandomRowFilter) String() string { return proto1.CompactTextString(m) }
func (*RandomRowFilter) ProtoMessage()    {}

func (m *RandomRowFilter) GetChance() float32 {
	if m != nil && m.Chance != nil {
		return *m.Chance
	}
	return 0
}

type RowFilter struct {
	CompareFilter    *CompareFilter `protobuf:"bytes,1,req,name=compare_filter" json:"compare_filter,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *RowFilter) Reset()         { *m = RowFilter{} }
func (m *RowFilter) String() string { return proto1.CompactTextString(m) }
func (*RowFilter) ProtoMessage()    {}

func (m *RowFilter) GetCompareFilter() *CompareFilter {
	if m != nil {
		return m.CompareFilter
	}
	return nil
}

type SingleColumnValueExcludeFilter struct {
	SingleColumnValueFilter *SingleColumnValueFilter `protobuf:"bytes,1,req,name=single_column_value_filter" json:"single_column_value_filter,omitempty"`
	XXX_unrecognized        []byte                   `json:"-"`
}

func (m *SingleColumnValueExcludeFilter) Reset()         { *m = SingleColumnValueExcludeFilter{} }
func (m *SingleColumnValueExcludeFilter) String() string { return proto1.CompactTextString(m) }
func (*SingleColumnValueExcludeFilter) ProtoMessage()    {}

func (m *SingleColumnValueExcludeFilter) GetSingleColumnValueFilter() *SingleColumnValueFilter {
	if m != nil {
		return m.SingleColumnValueFilter
	}
	return nil
}

type SingleColumnValueFilter struct {
	ColumnFamily      []byte       `protobuf:"bytes,1,opt,name=column_family" json:"column_family,omitempty"`
	ColumnQualifier   []byte       `protobuf:"bytes,2,opt,name=column_qualifier" json:"column_qualifier,omitempty"`
	CompareOp         *CompareType `protobuf:"varint,3,req,name=compare_op,enum=proto.CompareType" json:"compare_op,omitempty"`
	Comparator        *Comparator  `protobuf:"bytes,4,req,name=comparator" json:"comparator,omitempty"`
	FilterIfMissing   *bool        `protobuf:"varint,5,opt,name=filter_if_missing" json:"filter_if_missing,omitempty"`
	LatestVersionOnly *bool        `protobuf:"varint,6,opt,name=latest_version_only" json:"latest_version_only,omitempty"`
	XXX_unrecognized  []byte       `json:"-"`
}

func (m *SingleColumnValueFilter) Reset()         { *m = SingleColumnValueFilter{} }
func (m *SingleColumnValueFilter) String() string { return proto1.CompactTextString(m) }
func (*SingleColumnValueFilter) ProtoMessage()    {}

func (m *SingleColumnValueFilter) GetColumnFamily() []byte {
	if m != nil {
		return m.ColumnFamily
	}
	return nil
}

func (m *SingleColumnValueFilter) GetColumnQualifier() []byte {
	if m != nil {
		return m.ColumnQualifier
	}
	return nil
}

func (m *SingleColumnValueFilter) GetCompareOp() CompareType {
	if m != nil && m.CompareOp != nil {
		return *m.CompareOp
	}
	return CompareType_LESS
}

func (m *SingleColumnValueFilter) GetComparator() *Comparator {
	if m != nil {
		return m.Comparator
	}
	return nil
}

func (m *SingleColumnValueFilter) GetFilterIfMissing() bool {
	if m != nil && m.FilterIfMissing != nil {
		return *m.FilterIfMissing
	}
	return false
}

func (m *SingleColumnValueFilter) GetLatestVersionOnly() bool {
	if m != nil && m.LatestVersionOnly != nil {
		return *m.LatestVersionOnly
	}
	return false
}

type SkipFilter struct {
	Filter           *Filter `protobuf:"bytes,1,req,name=filter" json:"filter,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SkipFilter) Reset()         { *m = SkipFilter{} }
func (m *SkipFilter) String() string { return proto1.CompactTextString(m) }
func (*SkipFilter) ProtoMessage()    {}

func (m *SkipFilter) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type TimestampsFilter struct {
	Timestamps       []int64 `protobuf:"varint,1,rep,packed,name=timestamps" json:"timestamps,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *TimestampsFilter) Reset()         { *m = TimestampsFilter{} }
func (m *TimestampsFilter) String() string { return proto1.CompactTextString(m) }
func (*TimestampsFilter) ProtoMessage()    {}

func (m *TimestampsFilter) GetTimestamps() []int64 {
	if m != nil {
		return m.Timestamps
	}
	return nil
}

type ValueFilter struct {
	CompareFilter    *CompareFilter `protobuf:"bytes,1,req,name=compare_filter" json:"compare_filter,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}

func (m *ValueFilter) Reset()         { *m = ValueFilter{} }
func (m *ValueFilter) String() string { return proto1.CompactTextString(m) }
func (*ValueFilter) ProtoMessage()    {}

func (m *ValueFilter) GetCompareFilter() *CompareFilter {
	if m != nil {
		return m.CompareFilter
	}
	return nil
}

type WhileMatchFilter struct {
	Filter           *Filter `protobuf:"bytes,1,req,name=filter" json:"filter,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *WhileMatchFilter) Reset()         { *m = WhileMatchFilter{} }
func (m *WhileMatchFilter) String() string { return proto1.CompactTextString(m) }
func (*WhileMatchFilter) ProtoMessage()    {}

func (m *WhileMatchFilter) GetFilter() *Filter {
	if m != nil {
		return m.Filter
	}
	return nil
}

type FilterAllFilter struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *FilterAllFilter) Reset()         { *m = FilterAllFilter{} }
func (m *FilterAllFilter) String() string { return proto1.CompactTextString(m) }
func (*FilterAllFilter) ProtoMessage()    {}

func init() {
	proto1.RegisterEnum("proto.FilterList_Operator", FilterList_Operator_name, FilterList_Operator_value)
}
