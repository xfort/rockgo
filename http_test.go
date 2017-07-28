package rockgo

import (
	"testing"
	"log"
	"net/http"
)

func TestWeiPHP(t *testing.T) {
	rockhttp := NewRockHttp()
	resByte, err, response := rockhttp.GetBytes("http://2.wei.xiushuang.com/index.php?s=/addon/WeiSite/WeiSite/lists/cate_id/4.html", nil)

	log.Println(err, response.Status, string(resByte))

	request, _ := http.NewRequest("get", "", nil)
}
