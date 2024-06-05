package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/epkgs/mini-blink/cmd/wke2go/parser"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "API",
	Short: "Home Moments API",
	Run: func(cmd *cobra.Command, args []string) {
		// do stuff here
	},
}

var file string
var out string

func init() {
	// Define a CLI flag for the config file
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "wke.h", "wke header file path")

	rootCmd.PersistentFlags().StringVarP(&out, "output", "o", "wke.gen.go", "output go file path")

	rootCmd.Execute()
}

var wkeWebDragDataStr = `
enum _wkeStorageType {
	StorageTypeString,
	StorageTypeFilename,
	StorageTypeBinaryData,
	StorageTypeFileSystemFile,
};
struct _wkeWebDragDataItem {
	enum _wkeStorageType storageType;
	wkeMemBuf* stringType;
	wkeMemBuf* stringData;
	wkeMemBuf* filenameData;
	wkeMemBuf* displayNameData;
	wkeMemBuf* binaryData;
	wkeMemBuf* title;
	wkeMemBuf* fileSystemURL;
	int64 fileSystemFileSize;
	wkeMemBuf* baseURL;
};
struct _wkeWebDragData {
	struct _wkeWebDragDataItem* m_itemList;
	int m_itemListLength;
	int m_modifierKeyState;
	wkeMemBuf* m_filesystemId;
};
typedef struct _wkeWebDragData wkeWebDragData;
`

func main() {
	byts, err := parser.Simplify(file)
	if err != nil {
		panic(err)
	}

	wd, _ := os.Getwd()

	// 处理额外情况
	re := regexp.MustCompile(`typedef\s+struct\s+_wkeWebDragData\s*\{[\s\S]*\}\s*wkeWebDragData;`)
	byts = re.ReplaceAll(byts, []byte(wkeWebDragDataStr))

	os.WriteFile(filepath.Join(wd, "wke.simplify.h"), []byte(byts), 0644)

	parsed := parser.Parse("blink", string(byts))

	// fmt.Println(string(byts))

	os.WriteFile(filepath.Join(wd, "wke.parsed.go"), []byte(parsed), 0644)
}
