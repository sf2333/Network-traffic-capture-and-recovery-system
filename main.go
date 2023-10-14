package main

import (
	"RestoreFile/ExtractFile"
	"RestoreFile/capture"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

const (
	promiscuous bool = true //是否开启混杂模式
)

var (
	queryNetworkCard bool   //查询网卡信息
	device           string //网卡名称
	pcapFileName     string //pacp文件名
	mw               *MyMainWindow

	EndChan  chan bool
	DataChan chan []byte
	FilePath string
)

func parseParameters() {
	//flag.BoolVar(&queryNetworkCard, "i", false, "查询本机网卡的信息")
	//
	//device = flag.String("device", "", "要嗅探的网卡名称")
	//pcapFileName = flag.String("pcapFileName", "", "要解析的文件路径名称")
	//
	//flag.Parse()
	//
	//if *device == "" && *pcapFileName == "" {
	//	//*device = "\\Device\\NPF_{2CCCFA0A-FEE2-4688-BC5A-43A805A8DC67}"
	//	*pcapFileName = "saveHaveWord.pcap"
	//}
}

type NetworkCard struct {
	Id   int
	Name string
}

func main() {
	//解析命令行参数
	//parseParameters()
	FilePath, _ = os.Getwd()
	FilePath+="\\"


	DataChan = make(chan []byte, 10)
	EndChan = make(chan bool)

	//if queryNetworkCard {
	//	PrintNetworkCard()
	//	return
	//}

	//var waitGroup sync.WaitGroup

	//waitGroup.Add(1)

	go func() {
		disposeResult()
		//waitGroup.Done()
	}()

	//go capturePackets()

	StartGUI()

	//waitGroup.Wait()

}

func StartGUI() {
	var tableView *walk.TableView

	mw = &MyMainWindow{TableModel: NewFileInfoModel()}
	//mw.TableModel.SetDirPath()

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "网络流量抓取恢复与分析系统",
		MinSize:  Size{300, 200},
		Size:     Size{700, 500},
		Layout:   VBox{},

		Children: []Widget{
			GroupBox{

				Alignment: AlignHCenterVNear,
				//Layout: Grid{Columns: 4,Alignment:AlignHCenterVCenter},
				Layout: HBox{},
				Children: []Widget{
					GroupBox{
						Title: "NetCard",
						Layout:  HBox{},
						MinSize: Size{400, 50},
						Children: []Widget{
							ComboBox{
								AssignTo:      &mw.ComboBox,
								Value:         Bind("SpeciesId", SelRequired{}),
								BindingMember: "Id",
								DisplayMember: "Name",
								Model:         PrintNetworkCard(),
								MinSize:       Size{200, 10},
							},
						},
					},


					GroupBox{
						Title:   "Open PCAP File",
						Layout:  HBox{},
						MinSize: Size{400, 50},
						Children: []Widget{
							Label{
								Alignment: AlignHFarVCenter,
								Text:      "选择:",
								MaxSize:   Size{50, 200},
							},
							LineEdit{
								Enabled:  true,
								AssignTo: &mw.Edit,
								MaxSize:  Size{400, 200},
							},
							PushButton{
								Text:      "打开",
								OnClicked: mw.SelectFile, //点击事件响应函数
								MaxSize:   Size{50, 200},
							},
						},
					},
				},
			},


			GroupBox{
				Layout: HBox{MarginsZero: true},

				Children: []Widget{
					PushButton{
						AssignTo:  &mw.StartPushButton,
						Text:      "PE文件分析",
						MinSize: Size{Height: 20, Width: 200},
						OnClicked: mw.linkPython,
					},
				},
			},

			GroupBox{
				Layout: HBox{MarginsZero: true},

				Children: []Widget{
					PushButton{
						AssignTo:  &mw.StartPushButton,
						Text:      "开始",
						OnClicked: mw.StartAnalyze,
					},
					PushButton{
						AssignTo:  &mw.EndPushButton,
						Text:      "结束",
						OnClicked: mw.EndAnalyze,
						Enabled:   false,
					},
				},
			},
			GroupBox{
				Layout: VBox{},

				Children: []Widget{
					TableView{
						AssignTo:      &tableView,
						StretchFactor: 2,

						Columns: []TableViewColumn{
							{
								DataMember: "Name",
								Width:      600,
							},
						},
						Model: mw.TableModel,


						OnItemActivated: func() {

							if index := tableView.CurrentIndex(); index > -1 {
								name := mw.TableModel.Item[index].Name
								dirpath, _ := os.Getwd()
								exec.Command(`cmd`, `/c`, `explorer`, dirpath).Start()
								fmt.Println(name)
							}

						},
					},
				},
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

func capturePackets() {
	cap, err := capture.NewCapture(DataChan, EndChan)
	if err != nil {
		log.Fatal(err)
	}

	if device != "" {
		err = cap.SetCaptureSource(device, 1, promiscuous)
		if err != nil {
			log.Fatal(err)
		}
	} else if pcapFileName != "" {
		err = cap.SetCaptureSource(pcapFileName, 0, promiscuous)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("开始监听：")

	cap.StartCapture()
}

func PrintNetworkCard() []*NetworkCard {
	cap, err := capture.NewCapture(DataChan, nil)
	if err != nil {
		log.Fatal(err)
	}
	cards := cap.PrintDevices()

	networkCards := make([]*NetworkCard, 0)
	networkCards = append(networkCards, &NetworkCard{Id: 0, Name: ""})
	for i, v := range cards {
		networkCards = append(networkCards,
			&NetworkCard{Id: i + 1, Name: v,})
	}

	return networkCards

}

func disposeResult() {

	if !CheckFileIsExist(FilePath+"jpg\\") {
		_ = os.Mkdir(FilePath+"jpg\\", os.ModePerm)
	}

	if !CheckFileIsExist(FilePath+"png\\") {
		_ = os.Mkdir(FilePath+"png\\", os.ModePerm)
	}

	if !CheckFileIsExist(FilePath+"docx\\") {
		_ = os.Mkdir(FilePath+"docx\\", os.ModePerm)
	}

	i := 0
	for data := range DataChan {
		if len(data) <= 0 {
				continue
		}

		//filename := FilePath+"hex\\"+"test" + strconv.Itoa(i) + ".hex"
		//err := ioutil.WriteFile(filename, data, 0666)
		//if err != nil {
		//	log.Println("写入文件失败")
		//}

		data, length, outType := ExtractFile.ExtractFile(data, uint32(len(data)))
		if length != 0 {
			if outType == 0 {
				continue
			} else if outType == 1 {
				filename := FilePath+"jpg\\"+"test" + strconv.Itoa(i) + ".jpg"
				err := ioutil.WriteFile(filename, data, 0666)
				if err != nil {
					log.Println("写入文件失败")
				}
				log.Println(filename)
				go mw.TableModel.AddFilePath(filename)
				i++
			} else if outType == 2 {
				filename := FilePath+"docx\\"+"test" + strconv.Itoa(i) + ".docx"

				err := ioutil.WriteFile(filename, data, 0666)
				if err != nil {
					log.Println("写入文件失败")
				}
				log.Println(filename)

				go mw.TableModel.AddFilePath(filename)
				i++
			}else if outType == 3 {
				filename := FilePath+"png\\"+"test" + strconv.Itoa(i) + ".png"
				err := ioutil.WriteFile(filename, data, 0666)
				if err != nil {
					log.Println("写入文件失败")
				}
				log.Println(filename)

				go mw.TableModel.AddFilePath(filename)
				i++
			}

		}
		//log.Println(data)
	}
}


func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}

	return exist
}