package main

import (
	"reflect"
	"testing"
)

func TestTokenization(t *testing.T) {
	valid_tockens := []struct {
		InString      string
		ResultArray   []string
	}{
		{"int as3(int a_34, float bjt, unsigned char car);",
			[]string{"int", "as3", "(", "int", "a_34", ",", "float", "bjt", ",", "unsigned", "char", "car", ")", ";"}},
		{"int as3(int a_34[12]);",
			[]string{"int", "as3", "(", "int", "a_34", "[", "12", "]", ")", ";"}},
		{"int as3(int a_34);",
			[]string{"int", "as3", "(", "int", "a_34", ")", ";"}},
		{"int *as3();",
			[]string{"int", "*", "as3", "(", ")", ";"}},
		{"int as3();", []string{"int", "as3", "(", ")", ";"}},
		{"int sort(int in_array, int swap(int *a, int *b))",
			[]string{"int", "sort", "(", "int", "in_array", ",", "int", "swap", "(", "int", "*", "a", ",", "int", "*", "b", ")", ")"}}}
	for _, testCase := range valid_tockens {
		result, err := Tokenization(testCase.InString)
		if err != nil {
			t.Errorf("обнаружена ошибка: " + err.Error())
		} else {
			if !reflect.DeepEqual(result, testCase.ResultArray) {
				t.Error("ожидалось ", testCase.ResultArray, " найдено ", result)
			}
		}
	}
}


func TestType(t *testing.T) {
	valid_tockens := []struct {
		InArray      []string
		ResultOffset int
	}{
		{[]string{"char"}, 1},
		{[]string{"short"}, 1},
		{[]string{"int"}, 1},
		{[]string{"signed"}, 1},
		{[]string{"unsigned"}, 1},
		{[]string{"long"}, 1},
		{[]string{"float"}, 1},
		{[]string{"double"}, 1},
		{[]string{"signed", "char"}, 2},
		{[]string{"unsigned", "char"}, 2},
		{[]string{"short", "int"}, 2},
		{[]string{"signed", "short"}, 2},
		{[]string{"unsigned", "short"}, 2},
		{[]string{"signed", "int"}, 2},
		{[]string{"unsigned", "int"}, 2},
		{[]string{"long", "int"}, 2},
		{[]string{"signed", "long"}, 2},
		{[]string{"unsigned", "long"}, 2},
		{[]string{"long", "long"}, 2},
		{[]string{"long", "double"}, 2},
		{[]string{"signed", "short", "int"}, 3},
		{[]string{"unsigned", "short", "int"}, 3},
		{[]string{"signed", "long", "int"}, 3},
		{[]string{"unsigned", "long", "int"}, 3},
		{[]string{"long", "long", "int"}, 3},
		{[]string{"signed", "long", "long"}, 3},
		{[]string{"unsigned", "long", "long"}, 3},
		{[]string{"signed", "long", "long", "int"}, 4},
		{[]string{"unsigned", "long", "long", "int"}, 4}}
	for _, testCase := range valid_tockens {
		new_offset, err := builtInType(testCase.InArray, 0)
		if err != nil {
			t.Error(err.Error(), " для ", testCase.InArray)
		}
		if new_offset != testCase.ResultOffset {
			t.Error("не совпало смещение для ", testCase.InArray)
		}
	}
}

func TestVariabelName(t *testing.T) {
	names := []struct {
		InArray      []string
		ResultOffset int
	}{
		{[]string{"char"}, 0},
		{[]string{"short"}, 0},
		{[]string{"int"}, 0},
		{[]string{"123_signed"}, 0},
		{[]string{"unsigned_123"}, 1},
		{[]string{"bjt"}, 1},
		{[]string{"as3"}, 1}}
	for _, testCase := range names {
		new_offset, err := variabelName(testCase.InArray, 0)
		if new_offset != testCase.ResultOffset {
			t.Error("не совпало смещение для ", testCase.InArray)
			if err != nil {
				t.Error("сообщение: ", err.Error())
			}
		}
	}
}

func TestVariableDescription(t *testing.T) {
	names := []struct {
		InArray      []string
		ResultOffset int
	}{
		{[]string{"char", "my_byte"}, 2},
		{[]string{"short", "i"}, 2},
		{[]string{"unsigned", "int", "[", "]"}, 2}, // на первый ошибочный токен
		{[]string{"123_signed"}, 0},
		// согласно правилам misra с больше двух разыменовываний быть не должно, однако язык это никак не ограничивает.
		{[]string{"int", "*", "*", "*", "unsigned_123"}, 5},
		{[]string{"bool", "is_correct", "[", "16", "]", "[", "16", "]"}, 8},
		{[]string{"bool", "&", "by_link"}, 3},
		{[]string{"int", "main", "(", ")"}, 4}, // первого раза тесты прошло!
		{[]string{"int", "main", "(", "("}, 3},
		{[]string{"int", "sub", "(", "int", "a", ",", "int", "b", ")"}, 9},
		{[]string{"int", "main", "(", "int", "argc", ",", "char", "*", "argv", "[", "]", ")"}, 12},
		{[]string{"int", "main", "(", "int", "argc", "char", "*", "argv", "[", "]", ")"}, 5},
		{[]string{"int", "sort", "(", "int", "in_array", ",", "int", "swap", "(", "int", "*", "a", ",", "no_int", "*", "b", ")", ")"}, 13}}
	for _, testCase := range names {
		new_offset, err := variableDescription(testCase.InArray, 0)
		if new_offset != testCase.ResultOffset {
			t.Error("не совпало смещение для ", testCase.InArray, " надо ", testCase.ResultOffset, " получено ", new_offset)
			if err != nil {
				t.Error("сообщение: ", err.Error())
			}
		}
	}
}

func TestParenthesesFunction(t *testing.T) {
	names := []struct {
		InArray      []string
		ResultOffset int
	}{
		{[]string{"(", "int", "a", ")", ""}, 4},
		{[]string{"(", "int", "a", ",", "int", "b", ")", ""}, 7},
		{[]string{"(", "int", "a", "(", ")", ")", ""}, 6},
		{[]string{"(", "int", "a", "[", "]", ",", "int", "b", "[", "]", ",", "int", "swap", "(", "int", "*", "a", ",", "int", "*", "b", ")", ""}, 22}}
	for _, testCase := range names {
		new_offset, err := parenthesesFunction(testCase.InArray, 0)
		if new_offset != testCase.ResultOffset {
			t.Error("не совпало смещение для ", testCase.InArray, " надо ", testCase.ResultOffset, " получено ", new_offset)
			if err != nil {
				t.Error("сообщение: ", err.Error())
			}
		}
	}
}

func TestFunctionDescription(t *testing.T) {
	names := []struct {
		InArray         []string
		ResultOffset    int
		ShouldHaveError bool
	}{
		{[]string{"int", "as3", "(", "int", "a_34", ",", "float", "bjt", ",", "unsigned", "char", "car", ")", ";"}, 14, false},
		{[]string{"int", "main", "(", ")", ";"}, 5, false},
		{[]string{"int", "main", "(", "int", "argc", ",", "char", "*", "argv", "[", "]", ")", ";"}, 13, false},
		{[]string{"int", "sort", "(", "int", "a", "[", "]", ",", "int", "lenght", ",", "int", "*", "sort", "(", "int", "*", "left", ",", "int", "*", "right", ")", ")", ";"}, 25, false},
		{[]string{"int", "main"}, 2, true},
		{[]string{"int", "main", ";"}, 2, true},
		{[]string{"int", "main", "[", "]"}, 4, true},
		{[]string{"int", "main", "(", ")"}, 4, true},
		{[]string{"int", "sort", "(", "int", "a", "[", "]", ",", "int", "lenght", ",", "int", "*", "sort", "(", "int", "*", "left", ",", "int", "*", "right", ")", ")", ")", ";"}, 24, true}}
	for _, testCase := range names {
		new_offset, err := FunctionDescription(testCase.InArray, 0)
		if new_offset != testCase.ResultOffset {
			t.Error("не совпало смещение для ", testCase.InArray, " надо ", testCase.ResultOffset, " получено ", new_offset)
			if err != nil {
				t.Error("сообщение: ", err.Error())
			}
		}
		if (err != nil) != testCase.ShouldHaveError {
			t.Error("не совпало наличие ошибки для ", testCase.InArray, "должна быть: ", testCase.ShouldHaveError, "а пришло ", err)
			if err != nil {
				t.Error("сообщение: ", err.Error())
			}
		}
	}
}
