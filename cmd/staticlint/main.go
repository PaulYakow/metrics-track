/*
# Пакет staticlint представляет собой набор следующих статических анализаторов:
  - стандартные статические анализаторы пакета golang.org/x/tools/go/analysis/passes
  - анализаторы класса SA пакета staticcheck:
    SA1* - злоупотребления стандартной библиотекой,
    SA2* - проблемы с конкурентностью,
    SA3* - проблемы с тестированием,
    SA4* - бесполезный код (который ничего не делает по факту)
    SA5* - проблемы с корректностью
    SA6* - Проблемы с производительностью
    SA9* - Сомнительные кодовые конструкции (с высокой вероятностью могут быть ошибочными)
  - анализаторы класса S пакета staticcheck:
    S1001 - замена цикла for на вызов copy() для слайсов
    S1002 - опустить сравнения с true/false
    S1025 - не использовать fmt.Sprintf("%s", x) без необходимости
    S1028 - упрощать построение ошибок с помощью fmt.Errororf
  - анализаторы класса ST пакета staticcheck:
    ST1000 - неправильный или отсутствующий комментарий к пакету
    ST1005 - неправильно отформатированная строка ошибки
    ST1006 - неудачно выбранное имя получателя в методе
    ST1017 - не использовать условия вида `if 42 == x`
    ST1020 - комментарий экспортируемой функции должен начинаться с имени функции
    ST1021 - комментарий экспортируемого типа должен начинаться с имени типа
    ST1022 - комментарий экспортируемой переменной должен начинаться с имени переменной
    ST1023 - избыточный тип в объявлении переменной
  - анализаторы класса QF пакета staticcheck:
    QF1003 - преобразование if/else-if в switch
    QF1006 - поднять if+break в условие цикла
    QF1007 - объединение условного присвоения с объявлением переменной
  - анализатор `github.com/Antonboom/errname`:
    проверяет, что переменные ошибок имеют префикс Err, а типы ошибок - суффикс Error
  - анализатор `github.com/leonklingele/grouper`:
    анализ групп выражений (импорты, типы, переменные, константы)
  - анализатор `github.com/sashamelentyev/usestdlibvars`:
    обнаруживает возможность использования переменных/констант из стандартной библиотеки
  - анализатор ExitAnalyzer, который проверяет использование прямого вызова os.Exit в функции main пакета main

# Использование

	staticlint [package]

# Пример

	staticlint ./...
*/
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	errname "github.com/Antonboom/errname/pkg/analyzer"
	grouper "github.com/leonklingele/grouper/pkg/analyzer"
	usestdlibvars "github.com/sashamelentyev/usestdlibvars/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

const config = "config.json"

// ConfigData структура с конфигурацией анализаторов пакета staticcheck
//
// SA добавлены полностью. Для добавления дополнительных прописать необходимый в файле `config.json`
type ConfigData struct {
	Staticcheck []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), config))
	if err != nil {
		log.Fatal(err)
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	lintchecks := []*analysis.Analyzer{
		errname.New(),
		usestdlibvars.New(),
		grouper.New(),
		ExitAnalyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	}

	staticchecks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		staticchecks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		lintchecks = append(lintchecks, v.Analyzer)
	}

	for _, v := range quickfix.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			lintchecks = append(lintchecks, v.Analyzer)
		}
	}

	for _, v := range simple.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			lintchecks = append(lintchecks, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if staticchecks[v.Analyzer.Name] {
			lintchecks = append(lintchecks, v.Analyzer)
		}
	}

	multichecker.Main(lintchecks...)
}
