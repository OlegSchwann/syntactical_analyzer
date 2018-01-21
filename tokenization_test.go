package main

import (
	"reflect"
	"testing"
)

func TestTokenization1(t *testing.T) {
	input := "int as3();"
	expected := []string{"int", "as3", "(", ")", ";"}
	found, err := Tokenization(input)
	if err != nil {
		t.Errorf("ожидалось %t\n найдено %t\n, ошибка %s", expected, found, err.Error())
	} else {
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("ожидалось %t найдено %t", expected, found)
		}
	}
}

func TestTokenization2(t *testing.T) {
	input := "int *as3();"
	expected := []string{"int", "*", "as3", "(", ")", ";"}
	found, err := Tokenization(input)
	if err != nil {
		t.Errorf("ожидалось %t\n найдено %t\n, ошибка %s", expected, found, err.Error())
	} else {
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("ожидалось %t найдено %t", expected, found)
		}
	}
}

func TestTokenization3(t *testing.T) {
	input := "int as3(int a_34);"
	expected := []string{"int", "as3", "(", "int", "a_34", ")", ";"}
	found, err := Tokenization(input)
	if err != nil {
		t.Errorf("ожидалось %t\n найдено %t\n, ошибка %s", expected, found, err.Error())
	} else {
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("ожидалось %t найдено %t", expected, found)
		}
	}
}

func TestTokenization4(t *testing.T) {
	input := "int as3(int a_34[12]);"
	expected := []string{"int", "as3", "(", "int", "a_34", "[", "12", "]", ")", ";"}
	found, err := Tokenization(input)
	if err != nil {
		t.Errorf("ожидалось %t\n найдено %t\n, ошибка %s", expected, found, err.Error())
	} else {
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("ожидалось %t найдено %t", expected, found)
		}
	}
}

func TestTokenization5(t *testing.T) {
	input := "int as3(int a_34, float bjt, unsigned char car);"
	expected := []string{"int", "as3", "(", "int", "a_34", ",", "float", "bjt", ",", "unsigned", "char", "car", ")", ";"}
	found, err := Tokenization(input)
	if err != nil {
		t.Errorf("ожидалось %t\n найдено %t\n, ошибка %s", expected, found, err.Error())
	} else {
		if !reflect.DeepEqual(found, expected) {
			t.Errorf("ожидалось %t найдено %t", expected, found)
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
		new_offset, err := VariabelName(testCase.InArray, 0)
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
		{[]string{"unsigned", "int", "[", "]"}, 0}, //ошибка
		{[]string{"123_signed"}, 0},
		// согласно правилам misra с больше бвух разыменовываний быть не должно
		{[]string{"int", "*", "*", "*", "unsigned_123"}, 5},
		{[]string{"bool", "is_correct", "[", "16", "]", "[", "16", "]"}, 8},
		{[]string{"bool", "&", "by_link"}, 3}}
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
