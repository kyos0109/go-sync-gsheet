package main

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/sheets/v4"
)

type sheetsService struct {
	lock          sync.Mutex
	srv           *sheets.Service
	waitGroup     sync.WaitGroup
	spreadsheetId string
}

func main() {
	yamlConfig := GetConfig("config.yaml")

	log.Println("Every " + yamlConfig.Setting.Crontab.String() + " Start Syncing...")

	for {
		select {
		default:
			Go(yamlConfig)
			time.Sleep(yamlConfig.Setting.Crontab)
		}
	}
}

func Go(yamlConfig *yamlConfig) {
	var (
		ss  sheetsService
		err error
	)

	ctx := context.Background()
	client := AuthGoogle("credentials.json")
	ss.spreadsheetId = yamlConfig.Setting.SpreadsheetId

	ss.srv, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	for _, awsAuth := range yamlConfig.AWS.Auth {
		values := [][]interface{}{}
		sheetName := awsAuth.Account + "(" + awsAuth.Project + ")"
		rangeData := sheetName + "!A:XX"

		values = append(values, []interface{}{
			sheetName, "latest update", time.Now().Format("2006-01-02 15:04:05"), "", "", "", "",
		})

		values = append(values, []interface{}{
			"InstanceID", "Type", "TagName", "PrivateIpAddress", "PublicIpAddress", "State", "KeyName",
		})

		req := sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: sheetName,
				},
			},
		}

		rbb := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{&req},
		}

		ss.SendRequestService(rbb, ctx)
		time.Sleep(1 * time.Second)

		for _, ec2Des := range GetEC2List(&awsAuth) {
			for _, ins := range ec2Des.Instances {
				tagsName := FindEC2TagName(ins.Tags)
				data := []interface{}{
					ins.InstanceId, ins.InstanceType, tagsName, ins.PrivateIpAddress, ins.PublicIpAddress, ins.State.Name, ins.KeyName,
				}
				values = append(values, data)
			}
		}

		rb := &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "USER_ENTERED",
		}
		rb.Data = append(rb.Data, &sheets.ValueRange{
			Range:  rangeData,
			Values: values,
		})
		ss.SendRequestService(rb, ctx)
	}
	ss.waitGroup.Wait()
}

func (ss *sheetsService) SendRequestService(requestBatch interface{}, ctx context.Context) {
	switch rb := requestBatch.(type) {
	case *sheets.BatchUpdateSpreadsheetRequest:
		ss.waitGroup.Add(1)
		go func() {
			defer ss.waitGroup.Done()
			newSheetReq, err := ss.srv.Spreadsheets.BatchUpdate(ss.spreadsheetId, rb).Context(ctx).Do()
			if err != nil {
				log.Println(err)
			}
			log.Println(newSheetReq)
		}()
		break
	case *sheets.BatchUpdateValuesRequest:
		ss.waitGroup.Add(1)
		go func() {
			defer ss.waitGroup.Done()
			resp, err := ss.srv.Spreadsheets.Values.BatchUpdate(ss.spreadsheetId, rb).Context(ctx).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve data from sheet: %v", err)
			}
			log.Println(resp)
		}()
		break
	}
}
