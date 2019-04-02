package rockgo

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"sort"
)

type RockFFmpeg struct {
	ffmpegFile string
}

func (ffmpeg *RockFFmpeg) SetFFmpeg(ffmepgFile string) {
	ffmpeg.ffmpegFile = ffmepgFile
}

func (ffmpeg *RockFFmpeg) DoExec(args ...string) error {
	ffmpegCmd := exec.Command(ffmpeg.ffmpegFile, args...)
	outBytes, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return NewError("DoExec", args, err.Error(), string(outBytes))
	}
	return nil
}

// 视频转码
//
// videocodec 视频编码器，h264兼容最佳
//
// bitrate 比特率，建议2000k
//
// fps 帧率，建议 30
func (ffmpeg *RockFFmpeg) TranscodingVideo(sourceFile string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
	log.Println("开始视频转码",sourceFile,outFile)
	args := []string{"-i", sourceFile, "-vcodec", videocodec, "-b:v", bitrate, "-r", fps}
	args = append(args, parames...)
	args = append(args, outFile, "-y")
	return ffmpeg.DoExec(args...)
}

// 合并ts格式的video
func (ffmpeg *RockFFmpeg) ConcatVideos(videos []string, outFile string) error {
	log.Println("开始合并視頻", outFile)

	videosArgs := `concat:`
	videosLen := len(videos) - 1
	for index, item := range videos {
		item = filepath.FromSlash(item)
		videosArgs = videosArgs + item
		if index < videosLen {
			videosArgs = videosArgs + "|"
		}
	}
	videosArgs = videosArgs + ``
	args := []string{"-i", videosArgs, "-vcodec", "copy", outFile, "-y"}
	return ffmpeg.DoExec(args...)
}

// 合并未知格式的视频，先转码为ts格式
func (ffmpeg *RockFFmpeg) ConcatVideosTS(videos []string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
	outDir := filepath.Dir(outFile)
	fail := false
	for index, videoItem := range videos {
		log.Println("开始转码为 ts 格式", index, videoItem)
		itemVideoName := filepath.Base(videoItem)
		itemOut := filepath.Join(outDir, itemVideoName+".ts")
		err := ffmpeg.TranscodingVideo(videoItem, itemOut, videocodec, bitrate, fps, parames...)
		if err != nil {
			log.Println("转码失败", err)
			fail = true
			break
		}
		videos[index] = itemOut
	}
	if fail {
		return NewError("转码失败")
	}
	return ffmpeg.ConcatVideos(videos, outFile)
}

// 合并文件夹内的所有视频
func (ffmpeg *RockFFmpeg) ConcatDirVideos(dirpath string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
	dirpath = filepath.FromSlash(dirpath)
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	if len(files) <= 0 {
		return NewError("合并文件夹内视频失败_无文件", dirpath)
	}
	videos := make([]string, 0, len(files))
	for _, itemFile := range files {
		if itemFile.IsDir() {
			continue
		}
		videoName := itemFile.Name()
		videos = append(videos, filepath.Join(dirpath, videoName)) //filepath.Join(dirpath, videoName)
	}
	sort.Strings(videos)

	log.Println("文件夹内的视频文件合并顺序", dirpath)
	for index, itemVideo := range videos {
		log.Println(index, filepath.Base(itemVideo))
	}
	err = ffmpeg.ConcatVideos(videos, outFile)
	if err != nil {
		return err
	}

	targetOutFile := outFile + "_" + videocodec + bitrate + fps + filepath.Ext(outFile)
	return ffmpeg.TranscodingVideo(outFile, targetOutFile, videocodec, bitrate, fps, parames...)
}
