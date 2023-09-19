package pkg

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func RequestRegistry(url string, method string) (rpHeader http.Header, rpBody []byte, err error) {
	// Create a new HTTP request.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}

	// Set the custom header.
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	// Make the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()

	// Read the response body.
	rpBody, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	rpHeader = resp.Request.Response.Header.Clone()

	// Print the response body.
	fmt.Println(string(rpBody))
	return
}

func ShellCall(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()
	cmd.Start()
	<-done
	err := cmd.Wait()
	if err != nil {
		fmt.Println((err))
	}
}

func ShellCallResult(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()
	cmd.Start()
	<-done
	err := cmd.Wait()
	if err != nil {
		fmt.Println((err))
	}
	return ""
}

func ShellPipeStdin() {
	cmd := exec.Command("sh", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
		return
	}
	cmdReader, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	done := make(chan bool)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		done <- true
	}()

	go func() {
		defer stdin.Close()
		s1 := `
		echo 'hello world'
		pwd && hostname
		`
		io.WriteString(stdin, s1)
		io.WriteString(stdin, "echo 'done!'")
	}()
	err = cmd.Start()
	<-done
	if err != nil {
		log.Fatal(err)
	} else {
		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}
}
