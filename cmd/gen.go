package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate app code",
	
	Run: func(cmd *cobra.Command, args []string) {
		// download
		url := "https://github.com/gingerxman/gin/archive/master.zip"
		target := "./__gin_master.zip"
		fmt.Println(fmt.Sprintf("[download] Download %s...", url))
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		f, err := os.Create("./__gin_master.zip")
		if err != nil {
			panic(err)
		}
		size, err := io.Copy(f, res.Body)
		
		if err != nil {
			panic(err)
		} else {
			fmt.Println(fmt.Sprintf("[download] download %d bytes into ./__gin_master.zip - SUCCESS", size))
		}
		
		// unzip
		unzipDir := "./_gen_workspace"
		osCmd := exec.Command("rm", "-rf", unzipDir)
		err = osCmd.Run()
		if err != nil {
			panic(err)
		}
		
		osCmd = exec.Command("unzip", "-d", unzipDir, target)
		var out bytes.Buffer
		osCmd.Stdout = &out
		err = osCmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Printf("[unzip] unzip to %s - SUCCESS\n", unzipDir)
		
		// run gen.sh
		genShPath := fmt.Sprintf("%s/gin-master/pytool/code_generator/gen.sh", unzipDir)
		fmt.Printf("[gen] run %s to generate code...\n", genShPath)
		osCmd = exec.Command("bash", genShPath)
		out.Reset()
		osCmd.Stdout = &out
		err = osCmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", out.String())
		fmt.Println(">>>>>>>>>>>>>>>>>>>> Generate Code Success <<<<<<<<<<<<<<<<<<<<")
		fmt.Println("RESULT DIR: _gen_workspace/_generate")
	},
}


func init() {
	rootCmd.AddCommand(genCmd)
}