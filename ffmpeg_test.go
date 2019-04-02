package rockgo

import "testing"

var rockFFmpeg = &RockFFmpeg{}

func TestT_Start(t *testing.T) {
	rockFFmpeg.SetFFmpeg("ffmpeg")
	//testConcatVideos(t)

	dataMap := make(map[string]string)
	dataMap["F:\\AllInVideo\\pre第二期"] = "F:\\AllInVideo\\out\\pre2.mp4"

	for key, value := range dataMap {
		testConcatDirVideos(key, value, t)
	}
}

func testTranscodingVideo(t *testing.T) {
	sourceFile := "C:\\Users\\zhanghx\\Videos\\PRE\\out\\PRE_20190233.mp4"
	outFile := "C:\\Users\\zhanghx\\Videos\\PRE\\out\\RockFFmpeg_out.mp4"
	err := rockFFmpeg.TranscodingVideo(sourceFile, outFile, "h264", "2000k", "30", "-acodec", "aac")
	if err != nil {
		t.Fatal(err)
	}
}

func testConcatDirVideos(sourceDir string, outFile string, t *testing.T) {
	err := rockFFmpeg.ConcatDirVideos(sourceDir, outFile, "h264", "2000k", "30")
	if err != nil {
		t.Fatal(err)
	}
}

func testConcatVideos(t *testing.T) {
	videos := []string{"F:\\AllInVideo\\out\\2019YuanXiao\\00069.MTS.ts", "F:\\AllInVideo\\out\\2019YuanXiao\\00070.MTS.ts"}

	outFile := "F:\\AllInVideo\\out\\2019YuanXiao\\20190219YuanXiao.mp4"
	err := rockFFmpeg.ConcatVideos(videos, outFile)
	if err != nil {
		t.Fatal(err)
	}
}
