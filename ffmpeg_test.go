package rockgo

import "testing"

var rockFFmpeg = &RockFFmpeg{}

func TestT_Start(t *testing.T) {
	rockFFmpeg.SetFFmpeg("ffmpeg")
	testConcatDirVideos(t)
}

func testTranscodingVideo(t *testing.T) {
	sourceFile := "C:\\Users\\zhanghx\\Videos\\PRE\\out\\PRE_20190233.mp4"
	outFile := "C:\\Users\\zhanghx\\Videos\\PRE\\out\\RockFFmpeg_out.mp4"
	err := rockFFmpeg.TranscodingVideo(sourceFile, outFile, "h264", "500k", "30", "-acodec", "aac")
	if err != nil {
		t.Fatal(err)
	}
}

func testConcatDirVideos(t *testing.T) {
	sourceDir := "C:\\Users\\zhanghx\\Videos\\ZhiNengWaiHu\\out"
	outFile := "C:\\Users\\zhanghx\\Videos\\ZhiNengWaiHu\\out\\test.mp4"
	err := rockFFmpeg.ConcatDirVideos(sourceDir, outFile, "h264", "200k", "30")
	if err != nil {
		t.Fatal(err)
	}

}
