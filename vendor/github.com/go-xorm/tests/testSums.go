package tests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/go-xorm/xorm"
)

type SumStruct struct {
	Int   int
	Float float32
}

var (
	cases = []SumStruct{
		{1, 6.2},
		{2, 5.3},
		{92, -0.2},
	}
)

func sumCases() (i int, f float32) {
	for _, v := range cases {
		i += v.Int
		f += v.Float
	}
	return
}

func testSumSetUp(engine *xorm.Engine, t *testing.T) {
	err := engine.DropTables(new(SumStruct))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = engine.CreateTables(new(SumStruct))
	if err != nil {
		t.Error(err)
		panic(err)
	}

	_, err = engine.Insert(cases)
	if err != nil {
		t.Error(err)
		panic(err)
	}
}

func testSum(engine *xorm.Engine, t *testing.T) {
	testSumSetUp(engine, t)
	testSumOne(engine, t)
	testSums(engine, t)
}

func isFloatEq(i, j float64, precision int) bool {
	return fmt.Sprintf("%."+strconv.Itoa(precision)+"f", i) == fmt.Sprintf("%."+strconv.Itoa(precision)+"f", j)
}

func testSumOne(engine *xorm.Engine, t *testing.T) {
	colInt := engine.ColumnMapper.Obj2Table("Int")
	colFloat := engine.ColumnMapper.Obj2Table("Float")

	i, f := sumCases()

	sumInt, err := engine.Sum(new(SumStruct), "`"+colInt+"`")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if int(sumInt) != i {
		err = fmt.Errorf("sum result is %d, expect %d", int(sumInt), i)
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Sum %s %d!!!\n", colInt, int(sumInt))

	sumFloat, err := engine.Sum(new(SumStruct), "`"+colFloat+"`")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	if !isFloatEq(sumFloat, float64(f), 2) {
		err = fmt.Errorf("sum result is %f, expect %f", sumFloat, f)
		t.Error(err)
		panic(err)
	}
	fmt.Printf("Sum %s %f!!!\n", colFloat, sumFloat)
}

func testSums(engine *xorm.Engine, t *testing.T) {
	colInt := engine.ColumnMapper.Obj2Table("Int")
	colFloat := engine.ColumnMapper.Obj2Table("Float")

	sums, err := engine.Sums(new(SumStruct), "`"+colInt+"`", "`"+colFloat+"`")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	i, f := sumCases()

	if int(sums[0]) != i {
		err = fmt.Errorf("sum result is %d, expect %d", int(sums[0]), i)
		t.Error(err)
		panic(err)
	}

	if !isFloatEq(sums[1], float64(f), 2) {
		err = fmt.Errorf("sum result is %f, expect %f", sums[1], f)
		t.Error(err)
		panic(err)
	}

	fmt.Printf("Sum %s %d, %s, %f!!!\n", colInt, int(sums[0]), colFloat, sums[1])
}
