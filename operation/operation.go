package operation

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

func ExeCmd(name string, args []string) (err error) {
	cmd := exec.Command(name, args...)
	_, err = cmd.Output()
	if err != nil {
		if eErr, ok := err.(*exec.ExitError); ok {
			err = fmt.Errorf("%s\n%s",
				eErr.ProcessState.String(),
				string(eErr.Stderr),
			)
		}
	}
	return
}

func DownLoad(url, filename string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Body == nil {
		return err
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	io.Copy(file, res.Body)
	return nil
}

func CopyFile(in, out string) error {
	b, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(out, b, 0777)
	if err != nil {
		return err
	}
	return nil
}

func Exist(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}
