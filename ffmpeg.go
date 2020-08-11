package rockgo

import (
	"fmt"

	"os/exec"
)

type RockMedia struct {
	ffmpegpath string
}

func NewRockMedia(ffmpegPath string) *RockMedia {
	return &RockMedia{ffmpegpath: ffmpegPath}
}

/**
水平翻转视频
ffmpeg -i input.mkv -vf hflip hflip.mkv
*/
func (rock *RockMedia) HorizontalFlip(intputVideo string, outvideo string) {

}
func (rock *RockMedia) DoFFMPEG(arg ...string) error {
	cmd := exec.Command(rock.ffmpegpath, arg...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("combined out:\n%s\n", string(out))
		return err
	}
	return nil
}

/**
转为指定码率的webm格式视频
ffmpeg -i xx -c:v libvpx-vp9  -b:v 2000k -b:a 0 -c:a libopus xx.webm
*/
func (rock *RockMedia) TranscodingVP9(inputVideo string, bitrate string, outVideo string) error {
	return rock.DoFFMPEG("-i", inputVideo, "-b:v", bitrate, outVideo, "-y")
}

/**
视频截图
 ffmpeg -ss 0.1 -i E:\allin_video\ZhiNengWaiHu.mp4 -f image2 -vcodec mjpeg -vframes 1  -q:v 2 -y 1.jpg
-q:v 控制图片质量
*/
func (rockmedia *RockMedia) Screenshot(intputvideo, timePosition, outImg string) error {
	return rockmedia.DoFFMPEG("-ss", timePosition, "-i", intputvideo, "-f", "image2", "-vcodec", "libwebp", "-vframes", "1", "-y", outImg)
}
