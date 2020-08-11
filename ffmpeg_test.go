package rockgo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFFMEPG(t *testing.T) {
	//encodePlanCoverImg()
	//checkPlanActionVideo()
	//log.Println(EncodeXORToString([]byte("all_plans"),"com.quantumfiture.MuscleMen"))
	xorJSON("all_plans", "/Users/lihong/Downloads/template.json", " /Users/lihong/ZHX/code/ChangeIntl/app/src/main/assets/OF5VKzwCJ18q")

	//rockMedia := NewRockMedia("ffmpeg")
	//videoCoverOutDir := "/Users/lihong/ZHX/code/ChangeIntl/app/src/dev/assets/imgs/videoCover"
	//
	//encodeAllVideos()
	//log.Println(EncodeXORToString([]byte("123汉子"),"1234"))
}
//fGZ6nP_I7-DE

func encodePlanCoverImg() {
	dirPath := "/Users/lihong/ZHX/code/ChangeIntl/app/src/dev/assets/imgs"
	outDir := "/Users/lihong/ZHX/code/ChangeIntl/app/src/main/assets/imgs"
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	index := 0
	xorKey := "com.quantumfiture.MuscleMen"

	for _, itemFile := range files {
		if itemFile.IsDir() || strings.HasPrefix(itemFile.Name(), ".") {
			continue
		}
		itemFilePath := filepath.Join(dirPath, itemFile.Name())
		name := strings.Split(itemFile.Name(), ".")[0]
		outName := ""
		if len(name) > 64 {
			hashSHA256 := sha256.New()
			hashSHA256.Write([]byte(name))
			outName = hex.EncodeToString(hashSHA256.Sum(nil))
		} else {
			outName = EncodeXORToString([]byte(name), xorKey)
		}
		dataBytes, err := ioutil.ReadFile(itemFilePath)
		if err != nil {
			log.Fatal(err, itemFilePath)
		}
		outFilepath := filepath.Join(outDir, outName)
		err = ioutil.WriteFile(outFilepath, dataBytes, 0777)
		if err != nil {
			log.Fatal(err, itemFilePath)
		}
		log.Println(index, name, outName)
	}
}

func screenshotVideo(videopath, outpath string, rockMedia *RockMedia) {
	err := rockMedia.Screenshot(videopath, "0.1", outpath)
	if err != nil {
		log.Fatalln("screenshot error", err, videopath)
	}
}

func encodeAllVideos() {
	videosDir := "/Users/lihong/Downloads/WEB_20200811"
	videoOutDir := videosDir + "_XOR"
	videoCoverOutDir := videosDir + "_Covers"
	files, err := ioutil.ReadDir(videosDir)
	if err != nil {
		log.Fatal(err, videosDir)
	}

	os.MkdirAll(videoOutDir, 0755)
	os.MkdirAll(videoCoverOutDir, 0755)
	rockMedia := NewRockMedia("ffmpeg")

	appIDKey := "com.quantumfiture.MuscleMen"
	for _, itemFile := range files {
		if itemFile.IsDir() || strings.HasPrefix(itemFile.Name(), ".") || strings.Contains(itemFile.Name(), "_") {
			log.Println(itemFile.Name())
			continue
		}
		code := strings.Split(itemFile.Name(), ".")[0]
		itemVideoFile := filepath.Join(videosDir, itemFile.Name())
		videoXOROutFile := filepath.Join(videoOutDir, itemFile.Name())
		coverName := EncodeXORToString([]byte(code), appIDKey)
		coverOutFile := filepath.Join(videoCoverOutDir, coverName)
		screenshotVideo(itemVideoFile, coverOutFile, rockMedia)
		xorVideo(code, itemVideoFile, videoXOROutFile)
	}
}

func xorVideo(videoCode, videoFile, outvideo string) {
	hashSHA256 := sha256.New()
	hashSHA256.Write([]byte(videoCode))
	xorKey := hex.EncodeToString(hashSHA256.Sum(nil))
	xorKey = strings.ToLower(xorKey)

	videoData, err := ioutil.ReadFile(videoFile)
	if err != nil {
		log.Fatalln("xorVideo error", videoFile, err)
	}
	log.Println("xorvideo", videoCode, xorKey)
	resBytes := EncodeXOR(videoData, xorKey)
	err = ioutil.WriteFile(outvideo, resBytes, os.ModePerm)
	if err != nil {
		log.Fatalln("xorVideo", err, videoFile)
	}
}
func xorJSON(key, jsonFile, outFile string) {
	outFile = jsonFile + ".xor"
	outName := EncodeXORToString([]byte(key), "com.quantumfiture.MuscleMen")
	hashSHA256 := sha256.New()
	hashSHA256.Write([]byte(outName))
	xorKey := hex.EncodeToString(hashSHA256.Sum(nil))
	jsonData, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatalln("xorjson error", jsonFile, err)
	}
	log.Println("xorjson", outName, jsonFile, xorKey)
	resBytes := EncodeXOR(jsonData, xorKey)
	err = ioutil.WriteFile(outFile, resBytes, 0777)
	if err != nil {
		log.Fatalln("xorjson", err, outFile)
	}

	//sourceBytes := EncodeXOR(resBytes, xorKey)
	//ioutil.WriteFile(outFile+".decode", sourceBytes, 0777)
}

//func TestTolow(t *testing.T) {
//	dirPath := "/Users/lihong/Downloads/WEB_20200803"
//	//outDir := "/Users/lihong/ZHX/code/ChangeIntl/app/src/dev/assets/temp_2000"
//	files, err := ioutil.ReadDir(dirPath)
//	if err != nil {
//		t.Fatal(err)
//	}
//	for _, itemFile := range files {
//		if itemFile.IsDir() || strings.HasPrefix(itemFile.Name(), ".") || !strings.Contains(itemFile.Name(), "webm") {
//			log.Println(itemFile.Name())
//			continue
//		}
//		itemFilePath := dirPath + "/" + itemFile.Name()
//		itemFileOut := dirPath + "/" + strings.ToLower(itemFile.Name())
//		err := os.Rename(itemFilePath, itemFileOut)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//}

func checkPlanActionVideo() {
	iosVideosData := `"bisf21", "bisf22", "bisl19", "bisl7", "bisl9", "bk2", "bk5", "bwab18", "bwab19", "bwab3", "co1", "co19", "co5", "ct15", "ct16", "ct20", "ct21", "fd11", "fd12", "fd8", "gl12", "gl14", "gl15", "gl16", "gl3", "gl6", "gl8", "hami11", "lg18", "lg20", "lg24", "lg4", "lg40", "ply36", "ply42", "ply45", "ply50", "ply51", "rd13", "sl1", "sl10", "sl12", "sl16", "sl17", "sl18", "sl2", "sl27", "sl28", "sl29", "sl30", "sl31", "sl7", "sl8", "sl9", "su11", "su12", "su16", "su17", "su18", "su20", "su21", "su26", "su28", "su29", "su3", "su4", "su5", "tris14", "tris45"`
	iosVideosData = strings.ReplaceAll(iosVideosData, `"`, "")
	iosVideosData = strings.ReplaceAll(iosVideosData, " ", "")

	iosCodes := strings.Split(iosVideosData, ",")
	log.Println("IOS codes", len(iosCodes), iosCodes)

	videosDir := "/Users/lihong/Downloads/WEB_20200807"
	jsonFile := "/Users/lihong/Downloads/template.json"
	jsonBytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err, jsonFile)
	}
	var respData map[string]interface{}
	err = json.Unmarshal(jsonBytes, &respData)
	if err != nil {
		log.Fatal("1", err)
	}
	respData = respData["data"].(map[string]interface{})
	dataList := respData["list"].([]interface{})
	androidlossVideos := make(map[string]int)
	for _, item := range dataList {
		contentList := item.(map[string]interface{})["contentList"].([]interface{})
		for _, itemContent := range contentList {
			sectionList := itemContent.(map[string]interface{})["sections"].([]interface{})
			for _, itemSection := range sectionList {
				actions := itemSection.(map[string]interface{})["actions"].([]interface{})
				for _, itemAction := range actions {
					actionCode := itemAction.(map[string]interface{})["code"].(string)
					codes := strings.Split(actionCode, ",")
					for _, itemCode := range codes {
						actionFile := filepath.Join(videosDir, itemCode+".webm")
						_, err = os.Stat(actionFile)
						if err != nil {
							if androidlossVideos[itemCode] == 0 {
								androidlossVideos[itemCode] = 1
								log.Println(itemCode, err, item.(map[string]interface{})["nameZH"], item.(map[string]interface{})["nameZH"], item.(map[string]interface{})["id"])
							}
						} else {
							androidlossVideos[itemCode] = 10
						}
					}
				}
			}
		}
	}
	androidIndex := 0
	for code, status := range androidlossVideos {
		if status == 1 {
			androidIndex++
			fmt.Println("android filter", androidIndex, code)
		}
	}

	for _, ioscode := range iosCodes {
		if androidlossVideos[ioscode] == 1 {
			androidlossVideos[ioscode] = 2
		} else if androidlossVideos[ioscode] == 0 {
			androidlossVideos[ioscode] = 1
		}
	}
	//androidCodes := make([]string, 0)
	index := 0
	for code, status := range androidlossVideos {
		if status == 1 {
			//Android 丢失
			index++
			fmt.Println("android", index, code)
		} else if status == 10 {
			//正常
			//fmt.Println("ok", code)
		} else if status == 2 {
			//fmt.Println("ios", code)
		}
		//androidCodes = append(androidCodes, key)
	}
	//sort.Strings(androidCodes)
	//fmt.Println("android codes", len(androidCodes))
	//
	//for _, value := range androidCodes {
	//
	//	fmt.Println(value)
	//}
}
