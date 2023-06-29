package main

import (
	"fmt"
	"os"

	charsets "github.com/emersion/go-message/charset"
	pst "github.com/mooijtech/go-pst/v6/pkg"
	"github.com/mooijtech/go-pst/v6/pkg/properties"
	"github.com/rotisserie/eris"
	"golang.org/x/text/encoding"
)

func main() {
	pst.ExtendCharsets(func(name string, enc encoding.Encoding) {
		charsets.RegisterEncoding(name, enc)
	})

	fmt.Println("Inititalizing...")

	reader, err := os.Open("./data/mails.pst")

	if err != nil {
		panic(fmt.Sprintf("Failed to open PST file: %+v\n", err))
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to open PST file: %+v\n", err))
	}

	pstFile, err := pst.New(reader)

	if err != nil {
		panic(fmt.Sprintf("Failed to open PST file: %+v\n", err))
	}
	defer func() {
		pstFile.Cleanup()

		if errClosing := reader.Close(); errClosing != nil {
			panic(fmt.Sprintf("Failed to close PST file: %+v\n", err))
		}
	}()

	if err := pstFile.WalkFolders(func(folder *pst.Folder) error {
		fmt.Printf("Walking folder: %s\n", folder.Name)

		messageIterator, err := folder.GetMessageIterator()

		if eris.Is(err, pst.ErrMessagesNotFound) {
			// Folder has no messages.
			return nil
		} else if err != nil {
			return err
		}

		// Iterate through messages.
		for messageIterator.Next() {
			message := messageIterator.Value()

			switch messageProperties := message.Properties.(type) {
			case *properties.Message:
				fmt.Printf("Subject: %s\n", messageProperties.GetSubject())
			default:
				fmt.Printf("Unknown message type\n")
			}
		}

		return messageIterator.Err()
	}); err != nil {
		panic(fmt.Sprintf("Failed to walk folders: %+v\n", err))
	}
}
