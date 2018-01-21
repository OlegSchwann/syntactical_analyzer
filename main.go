package main

import (
	"errors"
	"fmt"
	"regexp"
)

//Вариант 16. C++. Разработать грамматику и распознаватель прототипов функции.
//Считать, что параметры только стандартных (скалярных) типов. Например:
//int as3(int a_34, float bjt, unsigned char car);

//напишем функцию разбора, что бы не учитывать пробелы в нормальных функциях
func Tokenization(parsedString string) (result []string, err error) {
	whitespace, _ := regexp.Compile(`\s`) // пробельный символ
	serviceSymbol, _ := regexp.Compile(`[&\*\(\)\[\];,]`)
	alphabeticSymbol, _ := regexp.Compile(`[A-Za-z_]`)
	numericalSymbol, _ := regexp.Compile(`[0-9]`)
	current_token := ""
	for i, char := range parsedString {
		letter := string(char)
		switch {
		//если мы нашли пробельный символ
		case whitespace.MatchString(letter):
			//и новый токен ещё не начинал собираться
			if current_token == "" {
				//то пропускаем и берём следующий символ
			} else {
				//собранный токен добавляем к результату
				result = append(result, current_token)
				current_token = ""
			}
			//если мы нашли самостоятельный управляющий символ
		case serviceSymbol.MatchString(letter):
			//и новый токен ещё не начинал собираться
			if current_token == "" {

			} else {
				//старый собранный токен добавляем к результату
				result = append(result, current_token)
				current_token = ""
			}
			//загоняем служебный символ в результат
			result = append(result, letter)
			//если мы нашли букву a-zA-Z_
		case alphabeticSymbol.MatchString(letter):
			//добавляем букву
			current_token += letter
			//отошлём её, когда встретим элемент синтаксиса или пробел
			//если мы нашли цифру
		case numericalSymbol.MatchString(letter):
			//добавляем цифру к собираемому токену
			current_token += letter
		default:
			//кидаем ошибку о недопустимом символе
			err = errors.New(fmt.Sprintf("unknown character '%s' on position %d", letter, i))
		}
		//записываем в результат остаток, если он сохранился
	}
	if current_token != "" {
		result = append(result, current_token)
	}
	return result, err
} // и прошло тесты почти сразу, разве не прелесть?

// функция проверяет, являются ли токены встроенным скалярным типом,
// начиная включительно с переданного смещения.
// возвращает смещение, с которого начинается следующий элемент разбора и ошибку, если не является скалярным типом.
//─┬char──────────────────┬─ https://ru.wikipedia.org/wiki/Типы_данных_в_C
// ├bool──────────────────┤
// ├short───┬─────────────┤
// │        └int──────────┤
// ├int───────────────────┤
// ├float─────────────────┤
// ├double────────────────┤
// ├long────┬─────────────┤
// │        ├int──────────┤
// │        ├double───────┤
// │        └long┬────────┤
// │             └int─────┤
// ├signed──┬─────────────┤
// │        ├char─────────┤
// │        ├short┬───────┤
// │        │     └int────┤
// │        ├int──────────┤
// │        └long┬────────┤
// │             ├long┬───┤
// │             │    └int┤
// │             └int─────┤
// └unsigned┬─────────────┤
//          ├char─────────┤
//          ├short┬───────┤
//          │     └int────┤
//          ├int──────────┤
//          └long┬────────┤
//               ├long┬───┤
//               │    └int┤
//               └int─────┘
func builtInType(tockens []string, tockenOffset int) (resultTockenOffset int, err error) {
	switch tockens[tockenOffset] {
	case "char":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
	case "bool":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
	case "short":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
		switch tockens[tockenOffset] {
		case "int":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
		}
	case "int":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
	case "signed":
		fallthrough
	case "unsigned":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
		switch tockens[tockenOffset] {
		case "char":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
		case "short":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
			switch tockens[tockenOffset] {
			case "int":
				tockenOffset++
				if tockenOffset == len(tockens) {
					break
				}
			}
		case "int":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
		case "long":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
			switch tockens[tockenOffset] {
			case "int":
				tockenOffset++
				if tockenOffset == len(tockens) {
					break
				}
			case "long":
				tockenOffset++
				if tockenOffset == len(tockens) {
					break
				}
				switch tockens[tockenOffset] {
				case "int":
					tockenOffset++
					if tockenOffset == len(tockens) {
						break
					}
				}
			}
		}
	case "long":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
		switch tockens[tockenOffset] {
		case "int":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
		case "long":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
			switch tockens[tockenOffset] {
			case "int":
				tockenOffset++
				if tockenOffset == len(tockens) {
					break
				}
			}
		case "double":
			tockenOffset++
			if tockenOffset == len(tockens) {
				break
			}
		}
	case "float":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
	case "double":
		tockenOffset++
		if tockenOffset == len(tockens) {
			break
		}
	default:
		err = errors.New("не является типом")
	}
	resultTockenOffset = tockenOffset
	return resultTockenOffset, err
}

// функция проверяет, является ли следующий токен названием переменной
// он должен состоять из букв и цифр, не начинаться с цифры и(sic!) не быть ключевым словом.
func VariabelName(tockens []string, tockenOffset int) (resultTockenOffset int, err error) {
	reservedNames := map[string]bool{
		"__abstract": true, "__alignof": true, "__asm": true, "__assume": true,
		"__based": true, "__box": true, "__cdecl": true, "__declspec": true,
		"__delegate": true, "__event": true, "__except": true, "__fastcall": true,
		"__finally": true, "__forceinline": true, "__gc": true, "__hook": true,
		"__identifier": true, "__if_exists": true, "__if_not_exists": true, "__inline": true,
		"__int16": true, "__int32": true, "__int64": true, "__int8": true,
		"__interface": true, "__leave": true, "__m128": true, "__m128d": true,
		"__m128i": true, "__m64": true, "__multiple_inheritance": true, "__nogc": true,
		"__noop": true, "__pin": true, "__property": true, "__raise": true,
		"__sealed": true, "__single_inheritance": true, "__stdcall": true, "__super": true,
		"__thiscall": true, "__try_cast": true, "__unaligned": true, "__unhook": true,
		"__uuidof": true, "__value": true, "__virtual_inheritance": true, "__w64": true,
		"__wchar_t": true, "wchar_t": true, "abstract": true, "array": true,
		"auto": true, "bool": true, "break": true, "case": true,
		"catch": true, "char": true, "const": true, "const_cast": true,
		"continue": true, "decltype": true, "default": true, "deprecated": true,
		"dllexport": true, "dllimport": true, "do": true, "double": true,
		"dynamic_cast": true, "else": true, "enum": true, "explicit": true,
		"extern": true, "false": true, "finally": true, "float": true,
		"for": true, "each": true, "in": true, "friend": true,
		"friend_as": true, "gcnew": true, "generic": true, "goto": true,
		"if": true, "initonly": true, "inline": true, "int": true,
		"interior_ptr": true, "literal": true, "long": true, "mutable": true,
		"naked": true, "namespace": true, "new": true, "noinline": true,
		"noreturn": true, "nothrow": true, "novtable": true, "nullptr": true,
		"private": true, "property": true, "protected": true, "public": true,
		"ref": true, "class": true, "register": true, "reinterpret_cast": true,
		"return": true, "safecast": true, "sealed": true, "selectany": true,
		"short": true, "signed": true, "sizeof": true, "static": true,
		"static_assert": true, "static_cast": true, "struct": true, "switch": true,
		"this": true, "thread": true, "throw": true, "true": true,
		"try": true, "typedef": true, "typeid": true, "typename": true,
		"union": true, "using": true, "uuid": true, "virtual": true,
		"void": true, "volatile": true, "while": true}
	_, isContained := reservedNames[tockens[tockenOffset]]
	if !isContained {
		nameForm, _ := regexp.Compile(`^[A-Z_a-z][0-9A-Z_a-z]*$`)
		if nameForm.MatchString(tockens[tockenOffset]) {
			tockenOffset++
		} else {
			err = errors.New("not a variabel name")
		}
	} else {
		err = errors.New("matches the keyword, can not be a variabel name")
	}
	resultTockenOffset = tockenOffset
	return resultTockenOffset, err
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// описанием переменной. Описание переменной состоит из имени и модификаторов
// ─тип_данных─┬↔┬─имя переменной─┬────↔────┬┬┬──────────────────────↔──────────────────────────┬─
//             ↓ ↑                ↓         ↑│↓                                                 ↑
//             ├&┤                └[┬число─]┘│└(─┬───────────────────↔───────────────────────┬─)┘
//             └*┘                  └]───────┘   └описание переменной┬──────────↔───────────┬┘
//                                                                   ↓                      ↑
//                                                                   └,─описание переменной─┘
func variableDescription(tockens []string, tockenOffset int) (resultTockenOffset int, err error) {
	resultTockenOffset = tockenOffset
	tockenOffset, err = builtInType(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid description, should start with type: " + err.Error())
		return resultTockenOffset, err
	} //началось с типа и ошибок нет
	// содержит модификаторы типа указатель или ссылка для C++, не обязательно
	if tockenOffset == len(tockens) {
		return resultTockenOffset, err
	}
	for tockens[tockenOffset] == "*" || tockens[tockenOffset] == "&" {
		tockenOffset++
	} //тип, модификатор
	tockenOffset, err = VariabelName(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid name of variable: " + err.Error())
		return resultTockenOffset, err
	} //тип, модификатор, имя
	//стандарт позволяет объявлять либо одномерный массив с неизвестной длинной, либо сколько угодно мерный со статической длинной
	if tockenOffset <= len(tockens)-2 {
		if tockens[tockenOffset] == "[" && tockens[tockenOffset+1] == "]" {
			tockenOffset = tockenOffset + 2
		} else {
			digit, _ := regexp.Compile(`\d*`)
			for tockenOffset <= len(tockens)-3 && tockens[tockenOffset] == "[" &&
				digit.MatchString(tockens[tockenOffset+1]) && tockens[tockenOffset+2] == "]" {
				tockenOffset = tockenOffset + 3
			}
		} //тип, модификатор, имя, массив
	}

	// TODO: функции чтоб можно было передавать как параметр
	resultTockenOffset = tockenOffset
	return resultTockenOffset, err
}

func main() {
	resultTockenOffset, _ := variableDescription([]string{"bool", "is_correct", "[", "16", "]", "[", "16", "]"}, 0)
	print(resultTockenOffset)
}
