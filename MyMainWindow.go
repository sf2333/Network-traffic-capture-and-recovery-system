package main

import (
	"fmt"
	"github.com/lxn/walk"
	"log"
	"os/exec"
	"strings"
)


type MyMainWindow struct {
	*walk.MainWindow
	Edit *walk.LineEdit
	ComboBox *walk.ComboBox
	TableModel *FileInfoModel
	StartPushButton *walk.PushButton
	EndPushButton *walk.PushButton
}

func (mw *MyMainWindow) SelectFile() {

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "可执行文件 (*.pcap)|*.pcap|所有文件 (*.*)|*.*"

	mw.Edit.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.Edit.SetText("")
		return
	} else if !ok {
		mw.Edit.SetText("")
		return
	}
	s := fmt.Sprintf("%s", dlg.FilePath)
	mw.Edit.SetText(s)
}

func (mw *MyMainWindow) StartAnalyze() {
	card := myTrim(mw.ComboBox.Text()," ")
	file := myTrim(mw.Edit.Text()," ")

	if card != ""{
		device = card
		pcapFileName = ""

		mw.StartPushButton.SetEnabled(false)
		mw.EndPushButton.SetEnabled(true)

		go capturePackets()
	}else if file != ""{
		pcapFileName = file
		device = ""

		mw.StartPushButton.SetEnabled(false)
		mw.EndPushButton.SetEnabled(true)

		capturePackets()

		mw.StartPushButton.SetEnabled(true)
		mw.EndPushButton.SetEnabled(false)
	}else{
		log.Println("请选择网卡或pcap文件！")
		return
	}



}

func (mw *MyMainWindow) EndAnalyze() {
	EndChan <- true


	mw.StartPushButton.SetEnabled(true)
	mw.EndPushButton.SetEnabled(false)
}

type FileInfo struct {
	Name string
}

type FileInfoModel struct {
	walk.SortedReflectTableModelBase
	dirPath string
	Item   []*FileInfo
}

var _ walk.ReflectTableModel = new(FileInfoModel)

func NewFileInfoModel() *FileInfoModel {
	fileInfoModel := new(FileInfoModel)
	fileInfoModel.Item = make([]*FileInfo,0)
	return new(FileInfoModel)
}

func (m *FileInfoModel) Items() interface{} {
	return m.Item
}

func (m *FileInfoModel) AddFilePath(filePath string) error {
	item := &FileInfo{
		Name: filePath,
	}

	m.Item = append(m.Item, item)

	m.PublishRowsReset()

	return nil
}

func myTrim(s string, prefix string) string{

	for strings.HasPrefix(s, prefix){
		s = s[len(prefix):]
	}

	for strings.HasSuffix(s,prefix){
		s = s[:len(s)-len(prefix)]
	}

	return s
}

func (mw *MyMainWindow) linkPython(){
	//首先生成cmd结构体,该结构体包含了很多信息，如执行命令的参数，命令的标准输入输出等
	cmd := exec.Command("cmd.exe", "/c", "start " + "H:\\exe\\main\\main.exe")
	//cmd := exec.Command("pycharm64","H:\\Virus analysis system\\main.py")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)

	}
}