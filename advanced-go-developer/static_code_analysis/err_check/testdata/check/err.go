package error

import "fmt"

func errFunc() (int, error) {
	return 0, fmt.Errorf("error")
}

func errCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	errFunc()              // want "expression returns unchecked error"
	res, _ := errFunc()    // want "assignment with unchecked error"
	fmt.Println(errFunc()) // want "expression returns unchecked error"
	_ = res
}
