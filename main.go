package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"log"
	"github.com/cts3njitedu/healthful-heart/oauth_token"
	"google.golang.org/api/drive/v3"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello World</h1>")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>About</h1>")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/about", about)
	r.HandleFunc("/", index)
	http.Handle("/", r);
	fmt.Println("Server Starting...")
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

	http.ListenAndServe(GetPort(), nil)
}


func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "3000"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}

	fmt.Println("The port is: ",port)

	return ":" + port

}