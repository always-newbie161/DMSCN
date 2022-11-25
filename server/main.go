package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"encoding/json"
	"flag"

	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/kaleido-io/kaleido-fabric-go/fabric"
	"github.com/kaleido-io/kaleido-fabric-go/kaleido"
)

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	ID             string `json:"ID"`
	Type		   string `json:"asset_type"`
	Name		   string `json:"name"`
	Bloodgroup	   string `json:"bloodgroup"`
	Stock          int    `json:"stock"`
	Owner          string `json:"owner"`
	AppraisedValue int    `json:"appraisedValue"`
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: myprog [inputfile]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	rand.Seed(time.Now().UTC().UnixNano())

	username := os.Getenv("USER_ID")
	if username == "" {
		username = "bts"
	}

	ccname := os.Getenv("CCNAME")
	if ccname == "" {
		ccname = "asset_transfer"
	}

	var err error
	var channel *kaleido.Channel
	var fabconnectClient *kaleido.FabconnectClient
	useFabconnect := os.Getenv("USE_FABCONNECT")

	if useFabconnect == "true" {
		fabconnectUrl := os.Getenv("FABCONNECT_URL")
		fabconnectClient = kaleido.NewFabconnectClient(fabconnectUrl, username)
		fmt.Println("Using Fabconnect for transaction submission")

		err := fabconnectClient.EnsureIdentity()
		if err != nil {
			fmt.Printf("Failed to ensure Fabconnect identity. %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Using Fabconnect identity: %s\n", username)
	} else {
		fmt.Println("Using the Fabric SDK for transaction submission")
		network := kaleido.NewNetwork()
		network.Initialize()
		config, err := fabric.BuildConfig(network)
		if err != nil {
			fmt.Printf("Failed to generate network configuration for the SDK: %s\n", err)
			os.Exit(1)
		}

		sdk1 := newSDK(config)
		defer sdk1.Close()

		wallet := kaleido.NewWallet(username, *network, sdk1)
		err = wallet.InitIdentity()
		if err != nil || wallet.Signer == nil {
			fmt.Printf("Failed to initiate wallet: %v\n", err)
			os.Exit(1)
		}

		fabric.AddTlsConfig(config, wallet.Signer)

		sdk2 := newSDK(config)
		defer sdk2.Close()

		channel = kaleido.NewChannel("default-channel", sdk2)
		err = channel.Connect(wallet.Signer.Identifier())
		if err != nil {
			fmt.Printf("Failed to connect to channel: %s\n", err)
			os.Exit(1)
		}
	}

	initChaincode := os.Getenv("INIT_CC")
	if initChaincode == "" {
		initChaincode = "false"
	} else {
		initChaincode = strings.ToLower(initChaincode)
	}

	if initChaincode == "true" {
		if useFabconnect == "true" {
			err = fabconnectClient.InitChaincode(ccname)
		} else {
			var Args [][]byte
			err = channel.InitChaincode(ccname, "InitLedger", Args)
		}
		if err != nil {
			fmt.Printf("Failed to initialize chaincode: %s\n", err)
		} else if useFabconnect == "true" {
			var batchWg sync.WaitGroup
			monitorFabconnectBatchReceipts(&batchWg, fabconnectClient, 0)
			batchWg.Wait()
			fabconnectClient.PrintFinalReport(1)
		}
	} else {
		var count int
		countStr := os.Getenv("TX_COUNT")
		if countStr != "" {
			count, err = strconv.Atoi(countStr)
			if err != nil {
				fmt.Printf("Failed to convert %s to integer", countStr)
				os.Exit(1)
			}
		} else {
			count = 1
		}
		if count > 50 {
			fmt.Println("Error: TX_COUNT cannot exceed 50")
			os.Exit(1)
		}
		var batches int
		batchStr := os.Getenv("BATCHES")
		if batchStr != "" {
			batches, err = strconv.Atoi(batchStr)
			if err != nil {
				fmt.Printf("Failed to convert %s to integer", batchStr)
				os.Exit(1)
			}
		} else {
			batches = 1
		}

		for i := 0; true; i++ {
			var wg sync.WaitGroup
			for j := 0; j < count; j++ {
				wg.Add(1)
				go func(idx int) {

					defer wg.Done()
					var fun string

					var id,newOwner string
					id = args[0]
					newOwner = args[1]
					fun = "TransferAsset"
					Args := [][]byte{[]byte(id), []byte(newOwner)}

					fmt.Printf("=> Batch %d: Send transaction %d of %d (Option: %d)\n", i+1, idx+1, count)
					err = channel.ExecChaincode(ccname, fun, Args)
					
					if err != nil {
						fmt.Printf("=> Batch %d: Failed to send transaction %d (Option: %d). %s\n", i+1, idx+1, err)
					}

				}(j)
			}
			wg.Wait()

			fmt.Printf("\nCompleted batch %d of %d\n\n", i+1, batches)

			if i < (batches-1) && useFabconnect != "true" {
				fmt.Println("Sleeping for 30 seconds before the next batch")
				time.Sleep(30 * time.Second)
			}
			
			break
		}

		if useFabconnect == "true" {
			var batchWg sync.WaitGroup
			for i := 0; i < batches; i++ {
				monitorFabconnectBatchReceipts(&batchWg, fabconnectClient, i)
			}
			batchWg.Wait()
			fabconnectClient.PrintFinalReport(batches)
		}

		fmt.Printf("\nAll Done!\n")
	}
}

func monitorFabconnectBatchReceipts(batchWg *sync.WaitGroup, fabconnectClient *kaleido.FabconnectClient, batch int) {
	batchWg.Add(1)
	go func(idx int) {
		defer batchWg.Done()
		fmt.Printf("=> Batch %d: Start Monitoring for transaction receipts\n", idx+1)
		fabconnectClient.MonitorBatch(idx)
	}(batch)
}

func newSDK(config map[string]interface{}) *fabsdk.FabricSDK {
	configProvider, err := fabric.NewConfigProvider(config)
	if err != nil {
		fmt.Printf("Failed to create config provider from config map: %s\n", err)
		os.Exit(1)
	}

	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		fmt.Printf("Failed to instantiate an SDK: %s\n", err)
		os.Exit(1)
	}
	return sdk
}

func ListAllFuncs() {
	
	var functions[7]string

	functions[0] = "1. Create an Asset"
	functions[1] = "2. Delete an Asset"
	functions[2] = "3. Transfer an Asset"
	functions[3] = "4. Read an Asset"
	functions[4] = "5. Update an Asset"
	functions[5] = "6. Get All Assets"

	for i:=0; i<len(functions); i++ {
		fmt.Println(functions[i])
	}
}

func printJSON(asset Asset) {
	fmt.Println("Asset{")
	fmt.Println("ID:", asset.ID)
	fmt.Println("Type:", asset.Type)
	fmt.Println("Name:", asset.Name)
	fmt.Println("Blood Group:", asset.Bloodgroup)
	fmt.Println("Stock:", asset.Stock)
	fmt.Println("Owner:", asset.Owner)
	fmt.Println("Appraisedvalue:", asset.AppraisedValue)
	fmt.Println("}")
}

func strtoArray(data []string) []*Asset{
	
	data[0] = data[0][2:]

	for i:=0 ; i<len(data)-1; i++ {
		data[i] = data[i] + "}"
	}
	data[len(data)-1] = data[len(data)-1][:len(data[len(data)-1])-3]

	var assets []*Asset
	for i:=0; i<len(data); i++ {

		var asset Asset
		d := data[i]

		d = strings.Replace(d, "\\", "", -1)
		err := json.Unmarshal([]byte(d), &asset)

		if err != nil {
			panic(err)
		}

		assets = append(assets, &asset)
	}

	return assets
}
