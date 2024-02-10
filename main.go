package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	drive "github.com/StollD/proton-drive"
	"github.com/StollD/webdav"
	"github.com/adrg/xdg"
	"gitlab.com/david_mbuvi/go_asterisks"
)

const (
	TokenFile  = "proton-webdav-bridge/tokens.json"
	AppVersion = "macos-drive@1.0.0-alpha.1+proton-webdav-bridge"
)

var (
	OptLogin  = false
	OptListen = "127.0.0.1:7984"
)

func doLogin() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter the username of your Proton Drive account.")
	fmt.Print("> ")

	user, err := reader.ReadString('\n')
	fmt.Println()

	if err != nil {
		return err
	}

	fmt.Println("Enter the password of your Proton Drive account.")
	fmt.Print("> ")

	pass, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	fmt.Println()

	if err != nil {
		return err
	}

	fmt.Println("Enter the mailbox password of your Proton Drive account.")
	fmt.Println("If you don't have a mailbox password, press enter.")
	fmt.Print("> ")

	mailbox, err := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	fmt.Println()

	if err != nil {
		return err
	}

	fmt.Println("Enter a valid 2FA token for your Proton Drive account.")
	fmt.Println("If you don't have 2FA setup, press enter.")
	fmt.Print("> ")

	twoFA, err := reader.ReadString('\n')
	fmt.Println()

	if err != nil {
		return err
	}

	credentials := drive.Credentials{
		Username:        user,
		Password:        string(pass),
		MailboxPassword: string(mailbox),
		TwoFA:           twoFA,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := drive.NewApplication(AppVersion)
	err = app.LoginWithCredentials(ctx, credentials)
	if err != nil {
		return err
	}

	err = storeTokens(*app.Tokens())
	if err != nil {
		return err
	}

	fmt.Println("Login successful.")
	return nil
}

func doListen() error {
	tokens, err := loadTokens()
	if err != nil {
		fmt.Println("Failed to load tokens!")
		fmt.Println("Run with --login to fix this!")
		fmt.Println()

		return err
	}

	fmt.Println("Waiting for network ...")
	WaitNetwork()

	fmt.Println("Connecting to Proton Drive ...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := drive.NewApplication(AppVersion)
	app.LoginWithTokens(&tokens)

	app.OnTokensUpdated(func(tokens *drive.Tokens) {
		err := storeTokens(*tokens)
		if err == nil {
			return
		}

		fmt.Println("Error storing tokens:", err)
	})

	app.OnTokensExpired(func() {
		fmt.Println("The stored tokens are no longer valid!")
		fmt.Println("Run with --login to fix this!")
		fmt.Println()

		os.Exit(1)
	})

	session := drive.NewSession(app)

	err = session.Init(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Connected!")

	return http.ListenAndServe(OptListen, &webdav.Handler{
		FileSystem: &ProtonFS{session: session},
		LockSystem: webdav.NewMemLS(),
	})
}

func loadTokens() (drive.Tokens, error) {
	var tokens drive.Tokens

	file, err := xdg.DataFile(TokenFile)
	if err != nil {
		return tokens, err
	}

	enc, err := os.ReadFile(file)
	if err != nil {
		return tokens, err
	}

	err = json.Unmarshal(enc, &tokens)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}

func storeTokens(tokens drive.Tokens) error {
	file, err := xdg.DataFile(TokenFile)
	if err != nil {
		return err
	}

	enc, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	return os.WriteFile(file, enc, 0600)
}

func main() {
	var err error = nil

	flag.BoolVar(&OptLogin, "login", OptLogin, "Run Proton Drive login")
	flag.StringVar(&OptListen, "listen", OptListen, "Which address the WebDAV server will listen to")
	flag.Parse()

	if OptLogin {
		err = doLogin()
	} else {
		err = doListen()
	}

	if err != nil {
		panic(err)
	}
}
