package ExtractFile

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"testing"
)

func TestExtractFile(t *testing.T) {
	data, err := ioutil.ReadFile("jpg.hex")
	if err != nil {
		log.Println(err)
		return
	}

	result, _, outType := ExtractFile(data, uint32(len(data)))

	fmt.Println(outType)
	err = ioutil.WriteFile("test.jpg", result, 0666)
	if err != nil {
		fmt.Println(err)
	}

}

func TestOpenDir(t *testing.T) {

	// 以经典的 C:\Program Files 为例
	exec.Command(`cmd`, `/c`, `explorer`, `C:\Program Files`).Start()

}
