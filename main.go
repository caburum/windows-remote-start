package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/masterzen/winrm"
)

func main1() {
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "Hi")
	})

	http.HandleFunc("/winrm", func(w http.ResponseWriter, r *http.Request){
		endpoint := winrm.NewEndpoint("localhost", 5985, false, false, nil, nil, nil, 0)
		client, err := winrm.NewClient(endpoint, "Administrator", "")
		if err != nil {
			panic(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client.RunWithContext(ctx, "ipconfig /all", os.Stdout, os.Stderr)
		fmt.Fprintf(w, "WinRM")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main2() {
	godotenv.Load()
	user := os.Getenv("WINRM_USER")
	password := os.Getenv("WINRM_PASSWORD")

	winrm.DefaultParameters.TransportDecorator = func() winrm.Transporter {
		// winrm https module
		return &winrm.ClientAuthRequest{}
	}

	endpoint := winrm.NewEndpoint("localhost", 5986, true, true, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, user, password)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err2 := client.RunWithContext(ctx, "ipconfig /all", os.Stdout, os.Stderr)
	if err2 != nil {
		panic(err2)
	}
}

func main3() {
	godotenv.Load()
	host := os.Getenv("WINRM_HOST")
	user := os.Getenv("WINRM_USER")
	password := os.Getenv("WINRM_PASSWORD")

	caCert, err := os.ReadFile("./ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %q", err)
	}

	endpoint := winrm.NewEndpoint(host, 5986, true, false, caCert, nil, nil, 0)

	params := winrm.DefaultParameters
	enc, _ := winrm.NewEncryption("ntlm")
	params.TransportDecorator = func() winrm.Transporter { return enc }

	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	if err != nil {
		fmt.Println(err)
	}

	stdOut, stdErr, exitCode, err := client.RunCmdWithContext(context.Background(), "ipconfig /all")
	fmt.Printf("%d\n%v\n%s\n%s\n", exitCode, err, stdOut, stdErr)
	if err != nil || (len(stdOut) == 0 && len(stdErr) > 0) {
		_ = exitCode
		fmt.Println(err)
	} else {
		fmt.Println("Command Test Ok")
	}

	wmiQuery := `select * from Win32_ComputerSystem`
	psCommand := fmt.Sprintf(`$FormatEnumerationLimit=-1; Get-WmiObject -Query "%s" | Out-String -Width 4096`, wmiQuery)
	stdOut, stdErr, exitCode, err = client.RunPSWithContext(context.Background(), psCommand)
	fmt.Printf("%d\n%v\n%s\n%s\n", exitCode, err, stdOut, stdErr)
	if err != nil || (len(stdOut) == 0 && len(stdErr) > 0) {
		_ = exitCode
		fmt.Println(err)
	} else {
		fmt.Println("PowerShell Test Ok")
	}
}

func main4() {
	godotenv.Load()
	host := os.Getenv("WINRM_HOST")
	user := os.Getenv("WINRM_USER")
	password := os.Getenv("WINRM_PASSWORD")

	caCert, err := os.ReadFile("./ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %q", err)
	}

	endpoint := winrm.NewEndpoint(host, 5986, true, false, caCert, nil, nil, 0)

	params := winrm.DefaultParameters
	enc, _ := winrm.NewEncryption("ntlm")
	params.TransportDecorator = func() winrm.Transporter { return enc }

	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	if err != nil {
		fmt.Println(err)
	}

	stdOut, stdErr, exitCode, err := client.RunCmdWithContext(context.Background(), "ipconfig /all")
	fmt.Printf("%d\n%v\n%s\n%s\n", exitCode, err, stdOut, stdErr)
	if err != nil || (len(stdOut) == 0 && len(stdErr) > 0) {
		_ = exitCode
		fmt.Println(err)
	} else {
		fmt.Println("Command Test Ok")
	}
}

func main() {
	godotenv.Load()
	host := os.Getenv("WINRM_HOST")
	user := os.Getenv("WINRM_USER")
	password := os.Getenv("WINRM_PASSWORD")

	caCert, err := os.ReadFile("./ca.crt")
	if err != nil {
		log.Fatalf("failed to read ca cert: %q", err)
	}

	clientCert, err := os.ReadFile("./user.pem")
	if err != nil {
	log.Fatalf("failed to read client certificate: %q", err)
	}

	clientKey, err := os.ReadFile("./key.pem")
	if err != nil {
		log.Fatalf("failed to read client key: %q", err)
	}

	params := winrm.DefaultParameters
	// params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientAuthRequest{} } // doesn't work for cert auth?
	enc, _ := winrm.NewEncryption("ntlm")
	params.TransportDecorator = func() winrm.Transporter { return enc }

	endpoint := winrm.NewEndpoint(host, 5986, true, false, caCert, clientCert, clientKey, 0)

	client, err := winrm.NewClientWithParameters(endpoint, user, password, params)
	// client, err := winrm.NewClient(endpoint, user, password)
	if err != nil {
		fmt.Println(err)
	}

	stdOut, stdErr, exitCode, err := client.RunCmdWithContext(context.Background(), "whoami")
	fmt.Printf("%d\n%v\n%s\n%s\n", exitCode, err, stdOut, stdErr)
	if err != nil || (len(stdOut) == 0 && len(stdErr) > 0) {
		_ = exitCode
		fmt.Println(err)
	} else {
		fmt.Println("Command Test Ok")
	}
}