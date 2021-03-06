package main

import (
	"errors"
	"fmt"
	"regexp"
	"bufio"
	"os"
	"strings"
)

//Вариант 16. C++. Разработать грамматику и распознаватель прототипов функции.
//Считать, что параметры только стандартных (скалярных) типов. Например:
//int as3(int a_34, float bjt, unsigned char car);

//Надо было сразу подумать о поведении в случае ошибки.
//Надо сделать так, что бы озвращал смещение на ошибку и её описание.
//тогда в результате можно будет подчеркнуть ошибку красным.
//На будущее: думать не только о интерфейсе, но и интерфейсе на случай ошибки

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
// принимает массив токенов, смещение на проверяемое слово.
// возвращает смещение + 1 и отсутствие ошибки или пепеданное смещение и
func variabelName(tockens []string, tockenOffset int) (resultTockenOffset int, err error) {
	reservedNames := map[string]bool{
		"__abstract":    true, "__alignof": true, "__asm": true, "__assume": true,
		"__based":       true, "__box": true, "__cdecl": true, "__declspec": true,
		"__delegate":    true, "__event": true, "__except": true, "__fastcall": true,
		"__finally":     true, "__forceinline": true, "__gc": true, "__hook": true,
		"__identifier":  true, "__if_exists": true, "__if_not_exists": true, "__inline": true,
		"__int16":       true, "__int32": true, "__int64": true, "__int8": true,
		"__interface":   true, "__leave": true, "__m128": true, "__m128d": true,
		"__m128i":       true, "__m64": true, "__multiple_inheritance": true, "__nogc": true,
		"__noop":        true, "__pin": true, "__property": true, "__raise": true,
		"__sealed":      true, "__single_inheritance": true, "__stdcall": true, "__super": true,
		"__thiscall":    true, "__try_cast": true, "__unaligned": true, "__unhook": true,
		"__uuidof":      true, "__value": true, "__virtual_inheritance": true, "__w64": true,
		"__wchar_t":     true, "wchar_t": true, "abstract": true, "array": true,
		"auto":          true, "bool": true, "break": true, "case": true,
		"catch":         true, "char": true, "const": true, "const_cast": true,
		"continue":      true, "decltype": true, "default": true, "deprecated": true,
		"dllexport":     true, "dllimport": true, "do": true, "double": true,
		"dynamic_cast":  true, "else": true, "enum": true, "explicit": true,
		"extern":        true, "false": true, "finally": true, "float": true,
		"for":           true, "each": true, "in": true, "friend": true,
		"friend_as":     true, "gcnew": true, "generic": true, "goto": true,
		"if":            true, "initonly": true, "inline": true, "int": true,
		"interior_ptr":  true, "literal": true, "long": true, "mutable": true,
		"naked":         true, "namespace": true, "new": true, "noinline": true,
		"noreturn":      true, "nothrow": true, "novtable": true, "nullptr": true,
		"private":       true, "property": true, "protected": true, "public": true,
		"ref":           true, "class": true, "register": true, "reinterpret_cast": true,
		"return":        true, "safecast": true, "sealed": true, "selectany": true,
		"short":         true, "signed": true, "sizeof": true, "static": true,
		"static_assert": true, "static_cast": true, "struct": true, "switch": true,
		"this":          true, "thread": true, "throw": true, "true": true,
		"try":           true, "typedef": true, "typeid": true, "typename": true,
		"union":         true, "using": true, "uuid": true, "virtual": true,
		"void":          true, "volatile": true, "while": true}
	_, isContained := reservedNames[tockens[tockenOffset]]
	if !isContained {
		nameForm, _ := regexp.Compile(`^[A-Z_a-z][0-9A-Z_a-z]*$`)
		if nameForm.MatchString(tockens[tockenOffset]) {
			tockenOffset++
		} else {
			err = errors.New("сontains illegal characters or starts with a number")
		}
	} else {
		err = errors.New("matches the keyword, can not be a variabel name")
	}
	resultTockenOffset = tockenOffset
	return resultTockenOffset, err
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// модификатором переменной
// ─┬↔┬─
//  ↓ ↑
//  └*┘
// принимает массив токенов и смещение, с которого ищет,
// возвращает смещение на следующий после описания токен и отсутствие ошибки или смещение на ошибку и её описание
func variableModifier(tockens []string, tockenOffset int) (int, error) {
	// содержит модификаторы типа указатель или ссылка для C++, не обязательно
	for tockenOffset != len(tockens) && (tockens[tockenOffset] == "*" || tockens[tockenOffset] == "&") {
		tockenOffset++
	} //тип, модификатор
	return tockenOffset, nil
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// модификатором переменной. стандарт позволяет объявлять либо одномерный массив с неизвестной длинной,
// либо сколько угодно мерный со статической длинной
// ─┬────↔────┬┬─
//  ↓         ↑│
//  └[┬число─]┘│
//    └]───────┘
// принимает массив токенов и смещение, с которого ищет,
// возвращает смещение на следующий после описания токен и отсутствие ошибки или смещение на ошибку и её описание
func arrayBrackets(tockens []string, tockenOffset int) (int, error) {
	if tockenOffset <= len(tockens)-2 {
		if tockens[tockenOffset] == "[" && tockens[tockenOffset+1] == "]" {
			tockenOffset = tockenOffset + 2
		} else {
			digit, _ := regexp.Compile(`\d*`)
			for tockenOffset <= len(tockens)-3 && tockens[tockenOffset] == "[" &&
				digit.MatchString(tockens[tockenOffset+1]) && tockens[tockenOffset+2] == "]" {
				tockenOffset = tockenOffset + 3
			}
		}
	}
	return tockenOffset, nil
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// модификатором переменной "функция".
// ─(─┬───────────────────↔───────────────────────┬─)─
//    └описание переменной┬──────────↔───────────┬┘      //косвенная рекурсия
//                        ↓                      ↑
//                        └,─описание переменной─┘
// принимает массив токенов и смещение, с которого ищет,
// возвращает смещение на следующий после описания токен и отсутствие ошибки или смещение на ошибку и её описание
func parenthesesFunction(tockens []string, tockenOffset int) (int, error) {
	// можно передавать функцию
	var err error
	if tockens[tockenOffset] == "(" {
		tockenOffset++
		//тип, модификатор, название, начало функции
		if tockenOffset == len(tockens) {
			err = errors.New("not a valid description, pairwise bracket, end detected")
			return tockenOffset, err
		}
		// не рекурсивный случай - просто закрывающая скобка после открывающей
		if tockens[tockenOffset] == ")" {
			tockenOffset++
		} else { // рекурсивный случай - после скобки что-то идёт
			tockenOffset, err = variableDescription(tockens, tockenOffset)
			if err != nil {
				err = errors.New("not valid function parameter declaration, not valid argument: " + err.Error())
				return tockenOffset, err
			}
			if tockenOffset == len(tockens) {
				err = errors.New("not a valid description, pairwise bracket, end detected")
				return tockenOffset, err
			} // прошло нормально первое. Теперь запятая и сколько угодно новых объявлений через запятую
			for tockens[tockenOffset] == "," {
				tockenOffset++
				if tockenOffset == len(tockens) {
					err = errors.New("not a valid description, pairwise bracket, end detected")
					return tockenOffset, err
				}
				tockenOffset, err = variableDescription(tockens, tockenOffset)
				if err != nil {
					err = errors.New("not valid function parameter declaration, not valid argument: " + err.Error())
					return tockenOffset, err
				}
			}
			if tockens[tockenOffset] == ")" {
				tockenOffset++
			} else {
				err = errors.New("unknown after the announcement, a closing parenthesis was expected, " +
					"perhaps you forgot the comma")
				return tockenOffset, err
			}
		}
	} else { // корректное описание, скобка закрыта.
		err = errors.New("invalid argument description, no opening parenthesis")
	}
	return tockenOffset, err
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// описанием переменной. Описание переменной состоит из имени и модификаторов
// ─тип_данных─┬↔┬─имя переменной─┬────↔────┬┬┬──────────────────────↔──────────────────────────┬─
//             ↓ ↑                ↓         ↑│↓                                                 ↑
//             ├&┤                └[┬число─]┘│└(─┬───────────────────↔───────────────────────┬─)┘
//             └*┘                  └]───────┘   └описание переменной┬──────────↔───────────┬┘
//                                                                   ↓                      ↑
//                                                                   └,─описание переменной─┘
// принимает массив токенов и смещение, с которого ищет,
// возвращает смещение на следующий после описания токен и отсутствие ошибки или смещение на ошибку и её описание
func variableDescription(tockens []string, tockenOffset int) (int, error) {
	var err error
	tockenOffset, err = builtInType(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid description, should start with type: " + err.Error())
		return tockenOffset, err
	} //началось с типа и ошибок нет
	tockenOffset, err = variableModifier(tockens, tockenOffset)
	if tockenOffset == len(tockens) {
		err = errors.New("not a valid description, missing name, end detected")
		return tockenOffset, err
	}
	tockenOffset, err = variabelName(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid description, missing name of variable: " + err.Error())
		return tockenOffset, err
	} //тип, возможно модификатор, имя
	tockenOffset, err = arrayBrackets(tockens, tockenOffset)
	if tockenOffset == len(tockens) {
		return tockenOffset, nil
	}
	// можно передавать функцию, но она не обязательно должна быть.
	newTockenOffset, err := parenthesesFunction(tockens, tockenOffset)
	if newTockenOffset != tockenOffset && err != nil {
		// значит не просто не нашли объявление, а нашли открывающую скобку и ошибку.
		err = errors.New("not a valid description, invalid argument description: " + err.Error())
		tockenOffset = newTockenOffset
		return tockenOffset, err
	} else {
		tockenOffset = newTockenOffset
	}
	return tockenOffset, nil
}

// функция проверяет, является ли массив токенов, начиная с переданного смещения включительно
// описанием функции. Описание функции состоит из имени, модификаторов, скобок с описанием параметров
// ─тип_данных─┬↔┬─имя переменной─┬────↔────┬┬─(─┬───────────────────↔───────────────────────┬─)─;─
//             ↓ ↑                ↓         ↑│   ↓                                           ↑
//             ├&┤                └[┬число─]┘│   └описание переменной┬──────────↔───────────┬┘
//             └*┘                  └]───────┘                       ↓                      ↑
//                                                                   └,─описание переменной─┘
// принимает массив токенов и смещение, с которого ищет,
// возвращает смещение на следующий после описания токен и отсутствие ошибки или смещение на ошибку и её описание
func FunctionDescription(tockens []string, tockenOffset int) (int, error) {
	var err error
	tockenOffset, err = builtInType(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid description, should start with type: " + err.Error())
		return tockenOffset, err
	} //началось с типа и ошибок нет
	tockenOffset, err = variableModifier(tockens, tockenOffset)
	if tockenOffset == len(tockens) {
		err = errors.New("not a valid description, missing name, end detected")
		return tockenOffset, err
	}
	tockenOffset, err = variabelName(tockens, tockenOffset)
	if err != nil {
		err = errors.New("not a valid description, missing name of variable: " + err.Error())
		return tockenOffset, err
	} //тип, возможно модификатор, имя
	tockenOffset, err = arrayBrackets(tockens, tockenOffset)
	if tockenOffset == len(tockens) {
		err = errors.New("expected function argument, found the end of string")
		return tockenOffset, err
	}
	// можно передавать функцию
	tockenOffset, err = parenthesesFunction(tockens, tockenOffset)
	if err != nil {
		err = errors.New("invalid function description, invalid argument description: " + err.Error())
		return tockenOffset, err
	}
	if tockenOffset == len(tockens) {
		err = errors.New("expected ';', found the end of string")
		return tockenOffset, err
	}
	if tockens[tockenOffset] == ";" {
		tockenOffset++
	} else {
		err = errors.New("unknown token, expected ';' after the description of the function parameters")
		return tockenOffset, err
	}
	return tockenOffset, nil
}

func main() {
	//приметивная консольная версия
	fmt.Println("Введите описание С++ функции:")
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	for { // break the loop if text == "q"
		scanner.Scan()
		text = scanner.Text()
		tokenArray, _ := Tokenization(text)
		resultTokenOffset, err := FunctionDescription(tokenArray, 0)
		if (resultTokenOffset == len(tokenArray)) && (err == nil) {
			fmt.Println("Это валидная функция на С++")
		} else {
			out_put := strings.Join(tokenArray[0:resultTokenOffset], " ") +
				" \x1B[31m" + tokenArray[resultTokenOffset] + "\x1B[0m " +
				strings.Join(tokenArray[resultTokenOffset+1:len(tokenArray)], " ")
			fmt.Println(out_put, "\n", err.Error())
		}
	}
}
