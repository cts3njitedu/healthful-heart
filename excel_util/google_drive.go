package excel_util

import (
	"github.com/cts3njitedu/healthful-heart/oauth_token"
	"google.golang.org/api/drive/v3"
	"fmt"
	"log"
)

func dummyExcel() {
	client, err:=oauth_token.GetClient()
	srv, err := drive.New(client);

	
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	read, err := srv.Files.List().PageSize(10).
                Fields("nextPageToken, files(id, name, mimeType, parents)").Do()
	if err != nil {
			log.Fatalf("Unable to retrieve files: %v", err)
	}

	fmt.Println("Files:")
	if len(read.Files) == 0 {
			fmt.Println("No files found.")
	} else {
			for _, i := range read.Files {
				if i.MimeType == "application/vnd.google-apps.folder" {
					fmt.Printf("%s (%s) (%s) (%v)\n", i.Name, i.Id, i.MimeType,i.Parents)
				}
					
			}
	}
	// spreadsheetId:="10S85mXce3fCXTk0fpayhyz2gSCL_JGSzmHlVyyacz6I"

	// readRange:="Sheet1!A1:EB108"
	// resp,err:=srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()

	// if err != nil {
	// 	log.Fatalf("Unable to retrieve data from sheet: %v", err)
	// }

	// if len(resp.Values) == 0 {
	// 		fmt.Println("No data found.")
	// } else {
	// 	fmt.Println("Name, Major:")
	// 	for _, row := range resp.Values {
	// 			// Print columns A and E, which correspond to indices 0 and 4.
	// 			fmt.Println(row)
	// 	}
	// }			

	// if err != nil {
	// 	log.Fatalf("Unable to retrieve Drive client: %v", err)
	// }

	// read, err := srv.Files.List().PageSize(10).
	// 	Fields("nextPageToken, files(id, name)").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve files: %v", err)
	// }
	// fmt.Println("Files:")
	// if len(read.Files) == 0 {
	// 	fmt.Println("No files found.")
	// } else {
	// 	for _, i := range read.Files {
	// 			fmt.Printf("%s (%s)\n", i.Name, i.Id)
	// 	}
	// }
	fmt.Println(client,err)
}





