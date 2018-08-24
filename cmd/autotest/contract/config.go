package contract

import (
	"sync"
	"github.com/inconshreveable/log15"
	"github.com/BurntSushi/toml"
	"gitlab.33.cn/chain33/chain33/cmd/autotest/testcase"
)


var (

	fileLog = log15.New()
	stdLog = log15.New()
)


//contract type
/*
bty,
token,
trade,

 */

type TestCaseFile struct{

	Contract string `toml:"contract"`
	Filename string `toml:"filename"`
}



type TestCaseConfig struct {

	CliCommand string `toml:"cliCmd"`
	CheckSleepTime int `toml:"checkSleepTime"`
	CheckTimeout int `toml:"checkTimeout"`
	TestCaseFileArr []TestCaseFile `toml:"TestCaseFile"`
}




type TestRunner interface {

	RunTest(tomlFile string, wg *sync.WaitGroup)
}




func InitConfig(logfile string){

	fileLog.SetHandler(log15.Must.FileHandler(logfile, log15.LogfmtFormat()))

}



func DoTestOperation(configFile string){

	var wg sync.WaitGroup
	var configConf TestCaseConfig

	if _, err := toml.DecodeFile(configFile, &configConf); err != nil {

		stdLog.Error("ErrTomlDecode", "Error", err.Error())
		return
	}

	testcase.Init(configConf.CliCommand, configConf.CheckSleepTime, configConf.CheckTimeout)

	stdLog.Info("[=====================BeginTest====================]")
	fileLog.Info("[=====================BeginTest====================]")

	for _, caseFile := range configConf.TestCaseFileArr {

		filename := caseFile.Filename

		switch caseFile.Contract {

		case "init":

			//init需要优先进行处理, 阻塞等待完成
			new(TestInitConfig).RunTest(filename, &wg)
			break

		case "bty":

			wg.Add(1)
			go new(TestBtyConfig).RunTest(filename, &wg)

		case "token":

			wg.Add(1)
			go new(TestTokenConfig).RunTest(filename, &wg)

		case "trade":

			wg.Add(1)
			go new(TestTradeConfig).RunTest(filename, &wg)

		}

	}

	wg.Wait()

	stdLog.Info("[=====================EndTest======================]")
	fileLog.Info("[=====================EndTest======================]")
}