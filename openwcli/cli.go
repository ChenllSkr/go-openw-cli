package openwcli

import (
	"fmt"
	"github.com/asdine/storm"
	"github.com/blocktree/OpenWallet/console"
	"github.com/blocktree/OpenWallet/log"
	"github.com/blocktree/OpenWallet/openwallet"
	"github.com/blocktree/OpenWallet/owtp"
	"github.com/blocktree/go-openw-sdk/openwsdk"
	"github.com/coreos/bbolt"
	"path/filepath"
	"time"
)

type CLI struct {
	config *Config             //工具配置
	db     *openwallet.StormDB //本地数据库
	api    *openwsdk.APINode   //api
}

// 初始化工具
func NewCLI(c *Config) (*CLI, error) {

	if len(c.appkey) == 0 {
		return nil, fmt.Errorf("appkey is empty. ")
	}

	if len(c.appid) == 0 {
		return nil, fmt.Errorf("appid is empty. ")
	}

	if len(c.remoteserver) == 0 {
		return nil, fmt.Errorf("remoteserver is empty. ")
	}

	dbfile := filepath.Join(c.dbdir, c.appid+".db")

	//加载数据
	db, err := openwallet.OpenStormDB(
		dbfile,
		storm.BoltOptions(0600, &bolt.Options{Timeout: 5 * time.Second}),
	)
	if err != nil {
		return nil, err
	}

	cli := &CLI{
		config: c,
		db:     db,
	}

	//配置日志
	SetupLog(c.logdir, "openwcli.log", false)

	keychain, err := cli.GetKeychain()
	if keychain != nil {
		cli.setupAPISDK(keychain)
	}

	return cli, nil
}

//setupAPI 配置APISDK
func (cli *CLI) setupAPISDK(keychain *Keychain) error {

	if keychain != nil {
		cert, _ := keychain.Certificate()
		sdkConfig := &openwsdk.APINodeConfig{
			AppID:           cli.config.appid,
			AppKey:          cli.config.appkey,
			HostNodeID:      "openw-server",
			ConnectType:     owtp.HTTP,
			Address:         cli.config.remoteserver,
			EnableSignature: true,
			Cert:            cert,
		}

		apiSDK := openwsdk.NewAPINode(sdkConfig)
		cli.api = apiSDK
	}

	return nil
}

//checkConfig 检查配置加载完
func (cli *CLI) checkConfig() error {

	if cli.config == nil {
		return fmt.Errorf("config is not loaded. ")
	}

	if cli.db == nil {
		return fmt.Errorf("database is not loaded. ")
	}
	return nil
}

//RegisterFlow 注册节点流程
func (cli *CLI) RegisterFlow() error {

	var (
		confirm bool
		keychain *Keychain
	)

	if check := cli.checkConfig(); check != nil {
		return check
	}

	keychain, err := cli.GetKeychain()
	if keychain != nil {
		//已经存在，提示是否需要覆盖
		confirm, _ = console.Stdin.PromptConfirm("The keychain already exist, do you want to regenerate current keychain?")
	} else {
		confirm = true
	}

	if confirm {
		//生成keychain
		keychain, err = cli.GenKeychain()
		if err != nil {
			return err
		}

		log.Info("Create keychain successfully.")

		//打印密钥对
		printKeychain(keychain)
	}

	//配置APISDK
	err = cli.setupAPISDK(keychain)
	if err != nil {
		return err
	}


	//登记节点
	err = cli.RegisterOnServer()
	if err != nil {
		return err
	}

	log.Info("Register node on opew-server successfully.")

	return nil
}

//GetNodeInfo 获取节点信息
func (cli *CLI) GetNodeInfoFlow() error {

	keychain, err := cli.GetKeychain()
	if err != nil {
		return err
	}

	printKeychain(keychain)

	return nil
}

//printKeychain 打印证书钥匙串
func printKeychain(keychain *Keychain) {
	//打印证书信息
	log.Notice("--------------- PRIVATE KEY ---------------")
	log.Notice(keychain.PrivateKey)
	log.Notice("--------------- PUBLIC KEY ---------------")
	log.Notice(keychain.PublicKey)
	log.Notice("--------------- NODE ID ---------------")
	log.Notice(keychain.NodeID)
}

//NewWalletFlow 创建钱包流程
func (cli *CLI) NewWalletFlow() error {

	var (
		password string
		name     string
		err      error
	)

	if cli.api == nil {
		return err
	}

	// 等待用户输入钱包名字
	name, err = console.InputText("Enter wallet's name: ", true)

	// 等待用户输入密码
	password, err = console.InputPassword(false, 3)

	err = cli.CreateWalletOnServer(name, password)
	if err != nil {
		return err
	}

	return nil
}

//ListWalletFlow
func (cli *CLI) ListWalletFlow() error {
	//TODO: WIP
	wallets, _ := cli.GetWalletsOnServer()
	cli.printWalletList(wallets)
	return nil
}

//NewAccountFlow
func (cli *CLI) NewAccountFlow() error {
	//TODO: WIP
	return nil
}

//ListAccountFlow
func (cli *CLI) ListAccountFlow() error {
	//TODO: WIP
	return nil
}

//NewAddressFlow
func (cli *CLI) NewAddressFlow() error {
	//TODO: WIP
	return nil
}

//SearchAddressFlow
func (cli *CLI) SearchAddressFlow() error {
	//TODO: WIP
	return nil
}

//TransferFlow
func (cli *CLI) TransferFlow() error {
	//TODO: WIP
	return nil
}

//SetSumFlow
func (cli *CLI) SetSumFlow() error {
	//TODO: WIP
	return nil
}

//StartSumFlow
func (cli *CLI) StartSumFlow() error {
	//TODO: WIP
	return nil
}
