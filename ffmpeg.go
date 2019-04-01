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
	err := ffmpegCmd.Start()
	if err != nil {
		return err
	}
	err = ffmpegCmd.Wait()

	if err != nil {
		if failErr, ok := err.(*exec.ExitError); ok {
			return failErr
		}
		return err
	}
	return nil
}

// 视频转码
//
// videocodec 视频编码器，h264兼容最佳
//
// bitrate 比特率，建议500k
//
// fps 帧率，建议 30
func (ffmpeg *RockFFmpeg) TranscodingVideo(sourceFile string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
	args := []string{"-i", sourceFile, "-vcodec", videocodec, "-b:v", bitrate, "-r", fps}
	args = append(args, parames...)
	args = append(args, outFile, "-y")
	return ffmpeg.DoExec(args...)
}

// 合并ts格式的video
func (ffmpeg *RockFFmpeg) ConcatTSVideos(videos []string, outFile string) error {
	videosArgs := `concat:"`
	for _, item := range videos {
		videosArgs = videosArgs + item + "|"
	}
	videosArgs = videosArgs + `"`
	args := []string{"-i", videosArgs, "-c", "copy", outFile, "-y"}
	return ffmpeg.DoExec(args...)
}

// 合并未知格式的视频，先转码为ts格式
func (ffmpeg *RockFFmpeg) ConcatVideos(videos []string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
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
	return ffmpeg.ConcatTSVideos(videos, outFile)
}

// 合并文件夹内的所有视频
func (ffmpeg *RockFFmpeg) ConcatDirVideos(dirpath string, outFile string, videocodec string, bitrate string, fps string, parames ...string) error {
	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}

	if len(files) <= 0 {
		return NewError("合并文件夹内视频失败_无文件", dirpath)
	}
	videos := make([]string, len(files))
	for index, itemFile := range files {
		if itemFile.IsDir() {
			continue
		}
		videoName := itemFile.Name()

		videos[index] = filepath.Join(dirpath, videoName)
	}
	sort.Strings(videos)

	log.Println("文件夹内的视频文件合并顺序",dirpath)
	for index,itemVideo:=range videos{
		log.Println(index,filepath.Base(itemVideo))
	}
	return ffmpeg.ConcatVideos(videos, outFile, videocodec, bitrate, fps, parames...)
}
