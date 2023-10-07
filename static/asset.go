package static

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
)

var (
	staticBytesCache = map[string][]byte{}
)

func Asset(name string) []byte {
	bs, ok := staticBytesCache[name]
	if ok {
		unzipBS, err := util.UnGzipBytes(bs)
		if err != nil {
			util.Logger.Error("Asset["+name+"]异常", zap.Error(err))
			return nil
		}
		return unzipBS
	}
	return nil
}

func GetStaticNames() (names []string) {
	for name := range staticBytesCache {
		names = append(names, name)
	}
	return
}

func FindStatic(name string) (bs []byte, find bool) {
	bs, find = staticBytesCache[name]
	return
}
func SetAsset(distDir string, saveDir string, savePrefix string) (err error) {
	// 删除历史静态资源
	savedFileMap, err := util.LoadDirFiles(saveDir)
	if err != nil {
		return
	}

	for filename := range savedFileMap {
		if strings.HasPrefix(filename, savePrefix) {
			err = os.Remove(saveDir + "/" + filename)
			if err != nil {
				return
			}
		}
	}
	staticFileMap, err := util.LoadDirFiles(distDir)
	if err != nil {
		return
	}

	staticGroupMap := map[string]map[string][]byte{}
	for filename, bs := range staticFileMap {
		groupName := savePrefix

		if strings.Contains(filename, "/") {
			dirName := filename[0:strings.LastIndex(filename, "/")]
			dirName = strings.Replace(dirName, "/", "_", 3)
			if strings.Contains(dirName, "/") {
				dirName = dirName[0:strings.Index(dirName, "/")]
			}
			groupName += "_" + dirName
		}

		fileMap := staticGroupMap[groupName]
		if fileMap == nil {
			fileMap = map[string][]byte{}
			staticGroupMap[groupName] = fileMap
		}
		fileMap[filename] = bs
	}

	for groupName, fileMap := range staticGroupMap {

		fmt.Println("文件组[" + groupName + "]文件数量[" + fmt.Sprint(len(fileMap)) + "]")
		var f *os.File
		f, err = os.Create(saveDir + "/" + groupName + ".go")
		if err != nil {
			return
		}

		_, _ = f.WriteString("package static" + "\n")
		_, _ = f.WriteString("\n")
		_, _ = f.WriteString("\n")
		_, _ = f.WriteString("func init() {" + "\n")
		var zipBS []byte
		for filename, bs := range fileMap {
			zipBS, err = util.GzipBytes(bs)
			if err != nil {
				util.Logger.Error("SetAsset["+filename+"]异常", zap.Error(err))
				return err
			}
			fmt.Println("文件[" + filename + "]大小[" + fmt.Sprint(len(bs)) + "]压缩后大小[" + fmt.Sprint(len(zipBS)) + "]")

			_, _ = f.WriteString(`	staticBytesCache["` + filename + `"] = ` + "[]byte{")
			size := len(zipBS)
			for i, b := range zipBS {
				if i == size-1 {
					_, _ = f.WriteString(fmt.Sprintf("%d", b))
				} else {
					_, _ = f.WriteString(fmt.Sprintf("%d,", b))
				}
			}
			_, _ = f.WriteString("}")
			_, _ = f.WriteString("\n")
			_, _ = f.WriteString("\n")
			_, _ = f.WriteString("\n")
		}
		_, _ = f.WriteString("}" + "\n")
		_ = f.Close()
	}

	return
}
