package cmd

import (
	"download/internal/download/http"
	"download/internal/file"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	url         string
	fileName    string
	concurrency int
	timeout     time.Duration
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from a given URL",
	Long: `This command allows you to download a file from a given URL.
You can also specify additional options such as tags, concurrency, and timeout.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 1. 初始化 FileManager
		fileManager := file.NewFileManager()

		// 2. 構建 Downloader
		downloader, err := http.NewHttpService(
			fileManager,
			http.WithURL(url),
			http.WithFileName(fileName),
			http.WithConcurrency(concurrency),
			http.WithTimeout(timeout),
		)
		if err != nil {
			fmt.Printf("Failed to initialize downloader: %v\n", err)
			return
		}

		// 3. 執行下載
		err = downloader.Downloader.Download()
		if err != nil {
			fmt.Printf("Download failed: %v\n", err)
			return
		}

		fmt.Println("Download complete!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := downloadCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	// 定義下載命令的參數
	downloadCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the file to download (required)")
	downloadCmd.Flags().StringVarP(&fileName, "fileName", "f", "download", "FileName of the file downloads (default: download)")
	downloadCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent downloads (default: 10)")
	downloadCmd.Flags().DurationVarP(&timeout, "timeout", "o", 30*time.Second, "Download timeout duration (default: 30s)")

	// 將 URL 設置為必填參數
	err := downloadCmd.MarkFlagRequired("url")
	if err != nil {
		fmt.Printf("Error marking flag as required: %v\n", err)
		os.Exit(1)
	}
}
