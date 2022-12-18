package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

// nil は don't care として扱う
func QuineMcCluskey(in [][]int, out [][]*int) ([][]SignificantGroup, error) {
	if len(in) != len(out) {
		return nil, errors.New("len(in) must be equal to len(out)")
	}
	n := len(in)
	if n == 0 {
		return nil, nil
	}
	if !all(in, func(t []int) bool { return len(t) == len(in[0]) }) || !all(out, func(t []*int) bool { return len(t) == len(out[0]) }) {
		return nil, errors.New("")
	}
	_, outM := len(in[0]), len(out[0])

	table := make([][]SignificantGroup, outM)
	// out の列ごとに Quine-McCluskey をやる
	for i := 0; i < outM; i++ {
		table[i] = quineMcCluskeyWith1out(in, getColumnFrom2d(out, i))
	}

	return table, nil
}

var dfsCount int

type SignificantGroup struct {
	significant []*int
	indexList   []int
}

func dfs(significantGroupList []SignificantGroup) []SignificantGroup {
	dfsCount++
	if dfsCount == 10 {
		panic("too many call of dfs")
	}

	newSignificantGroupListByID := map[string]SignificantGroup{} // 項
	for i := 0; i < len(significantGroupList)-1; i++ {
		for j := i + 1; j < len(significantGroupList); j++ {
			if calcHammingDistance(significantGroupList[i].significant, significantGroupList[j].significant) != 1 {
				continue
			}
			significantGroup := SignificantGroup{
				significant: combination(significantGroupList[i].significant, significantGroupList[j].significant),
				indexList:   append(significantGroupList[i].indexList, significantGroupList[j].indexList...),
			}
			slices.Sort(significantGroup.indexList)
			id := strings.Join(lo.Map(significantGroup.indexList, func(item, _ int) string { return strconv.Itoa(item) }), ",")
			newSignificantGroupListByID[id] = significantGroup
		}
	}
	newSignificantGroupList := lo.Values(newSignificantGroupListByID)

	// 単一項(index)が所属している組み合わせ項の index の list
	belongListByIndex := map[int][]int{}
	for i, significantGroup := range newSignificantGroupList {
		for _, index := range significantGroup.indexList {
			belongListByIndex[index] = append(belongListByIndex[index], i)
		}
	}

	mustSignificantGroupListByIndex := map[int]SignificantGroup{}
	for index, belongList := range belongListByIndex {
		_, exists := mustSignificantGroupListByIndex[index]
		if len(belongList) == 1 && !exists {
			mustSignificantGroupListByIndex[index] = newSignificantGroupList[belongList[0]]
		}
	}
	mustSignificantGroupList := lo.Values(mustSignificantGroupListByIndex)

	// 必須項のみだったらこれ以上再帰をする必要がないので return
	if len(mustSignificantGroupList) == len(newSignificantGroupList) {
		return mustSignificantGroupList
	}

	fmt.Printf("===must-significant-list(%d)===\n", dfsCount)
	printSignificantTable(mustSignificantGroupList)
	fmt.Println()
	fmt.Printf("===new-significant-list(%d)===\n", dfsCount)
	printSignificantTable(newSignificantGroupList)
	fmt.Println()

	mustSignificantList2 := dfs(newSignificantGroupList)

	return append(mustSignificantGroupList, mustSignificantList2...)
}

func quineMcCluskeyWith1out(in [][]int, out []*int) []SignificantGroup {
	// out が1となるような input の行を追加していく
	significantList := make([]SignificantGroup, 0)
	for i := 0; i < len(in); i++ {
		if out[i] == nil {
			significantList = append(significantList, SignificantGroup{
				significant: lo.Map(in[i], func(a int, _ int) *int {
					return &a
				}),
				indexList: []int{i},
			})
		} else if *out[i] == 1 {
			significantList = append(significantList, SignificantGroup{
				significant: lo.Map(in[i], func(a int, _ int) *int {
					return &a
				}),
				indexList: []int{i},
			})
		}
	}

	fmt.Println("===in(out=1)===")
	printSignificantTable(significantList)

	return dfs(significantList)
}

func calcHammingDistance(a, b []*int) int {
	count := 0
	for i := 0; i < len(a); i++ {
		if a[i] == nil && b[i] == nil {
			continue
		} else if (a[i] == nil && b[i] != nil) || (a[i] != nil && b[i] == nil) {
			count++
		} else if *a[i] != *b[i] {
			count++
		}
	}
	return count
}

// 0101, 0111 → 01-1
func combination(a, b []*int) []*int {
	combi := make([]*int, len(a))
	for i := 0; i < len(a); i++ {
		if (a[i] == nil && b[i] != nil) || (a[i] != nil && b[i] == nil) {
			panic("pattern -0, -1, 0-, 1- are not allowed")
		}
		if a[i] == nil && b[i] == nil {
			combi[i] = nil
		} else if *a[i] != *b[i] {
			combi[i] = nil
		} else {
			combi[i] = a[i]
		}
	}
	return combi
}

func printSignificant(label string, row SignificantGroup) {
	dispList := lo.Map(row.significant, func(item *int, _ int) string {
		if item == nil {
			return "-"
		} else {
			return strconv.Itoa(*item)
		}
	})
	indexList := lo.Map(row.indexList, func(index int, _ int) string {
		return strconv.Itoa(index)
	})

	fmt.Println(label, strings.Join(dispList, " | "), "("+strings.Join(indexList, ", ")+")")
}

func printSignificantTable(a []SignificantGroup) {
	for _, row := range a {
		printSignificant("", row)
	}
}
