package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type ByName []os.FileInfo

/**
получить длину a
 */
func (a ByName) Len() int{
	return len(a)
}

/**
меняем местами 2 значения
 */
func (a ByName) Swap(i, j int){
	a[i], a[j] = a[j], a[i]
}

/**
сравниваем 2 значения
 */

func (a ByName) Less(i, j int) bool {
	return a[i].Name() < a[j].Name()
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
/**
вызов построения дерева категорий
@out - вывод
@path - путь к директории обхода файлов
@printFiles - вывод файлов
 */
func dirTree(out io.Writer, path string, printFiles bool) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	return buildTree(1, out, dir, printFiles, []bool{})
}

/**
 рекурсивное построение дерева категорий
@level - уровень вложенности
@out - вывод
@dir - текущая директория
@printFiles - печать файлов
 */
func buildTree(level int, out io.Writer, dir *os.File, printFiles bool, t []bool) error {
	files, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	if !printFiles {
		// удаляем файлы
		var dirs []os.FileInfo
		for _, f := range files {
			if f.IsDir() {
				dirs = append(dirs, f)
			}
		}
		files = dirs
	}

	//сортируем названия файлов
	sort.Sort(ByName(files))
	for i, f := range files {
		var newt []bool
		var str string
		if i != len(files)-1 {
			str = tree(t) + "├───" + f.Name()
			newt = append(t, true)
		} else {
			str = tree(t) + "└───" + f.Name()
			newt = append(t, false)
		}

		//если это файл, а не дирректория, то считаем размер файла
		// если файл пустой - выводим 0b
		if !f.IsDir() {
			sizefile := int(f.Size())
			if sizefile != 0 {
				str += " (" + strconv.Itoa(sizefile) + "b)"
			} else {
				str += " (0b)"
			}
		}
		fmt.Fprintf(out, "%s\n", str)

		//если это директория - то ищем в ней директории и файлы
		if f.IsDir() {
			parentDir, err := os.Open(dir.Name() + "/" + f.Name())
			if err != nil {
				return err
			}
			//проверяем, есть ли у директории потомки
			if err := buildTree(level+1, out, parentDir, printFiles, newt); err != nil {
				return err
			}
		}
	}
	return nil
}
/**
рисуем отступы
 */
func tree(t []bool) string {
	var str string
	for _, b := range t {
		if b {
			str += "│\t"
		} else {
			str += "\t"
		}
	}
	return str
}