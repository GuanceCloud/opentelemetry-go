// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package feed // import "go.opentelemetry.io/otel/exporters/guance/internal/feed"

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultRetryTimes = 6
	defaultMaxLength  = 32 * 1024 * 1024
	// defaultMaxLength = 128 // TODO 测试用的，回头删除
	defaultTimeout = 3 // TODO http超时多少秒合适啊？
)

// chan 接收数据结构
type FeedInfo struct {
	LineProto string
	URL       string
}

// 中间整理数据缓存结构
type tempInfo struct {
	lineProtos string
	urlStr     string
	overflow   bool
}

// 发送数据队列结构
type retryInfo struct {
	lineProtos string
	urlStr     string
	retryTimes int64
}

var (
	FeedCh     chan []FeedInfo // TODO 用 []*feedInfo 更合理吗？
	retryInfos []*retryInfo    // 发送数据队列
	wg         sync.WaitGroup
)

// RetryFeed 重试 Feed 数据
// - 定时执行
// - 每个发送包有一个默认最大字节数，例如32M。
// - 不同 url 开不同协程，不同url的数据在不同的切片里。
// - 同一个 url 如果数据超级大，开多个协程。
// - 全体协程都回来后，统一处理发送结果，删除发送成功的。
func feed() {
	tick := time.NewTicker(time.Second * 2) // 测试，回头改成1
	defer tick.Stop()
	dataCache := make([]FeedInfo, 0)
	for {
		select {
		case f := <-FeedCh:
			// 接收数据
			dataCache = append(dataCache, f...)
		case <-tick.C:
			if len(dataCache) > 0 {
				// 整理数据
				appendFeedInfos(dataCache)
				dataCache = make([]FeedInfo, 0)

				// 发送数据，
				for i := 0; i < len(retryInfos); i++ {
					wg.Add(1)
					go doSend(retryInfos[i])
				}
				wg.Wait()

				// 删除发送过的数据
				fmt.Println("执行eraseSendedInfo : ", len(retryInfos))
				eraseSendedInfo()
			}
		}
	}
}

// appendFeedInfos channel 数据追加到发送队列
func appendFeedInfos(data []FeedInfo) {

	// 查找定位函数
	findURLIdx := func(urlStr string, temps []tempInfo) int {
		for findIdx := 0; findIdx < len(temps); findIdx++ {
			if urlStr == temps[findIdx].urlStr && !temps[findIdx].overflow {
				return findIdx
			}
		}
		return -1
	}

	// 数据梳理、归并到中间变量temp
	tempInfos := make([]tempInfo, 0)
	for _, feedInfo := range data {
		urlIdx := findURLIdx(feedInfo.URL, tempInfos)
		if urlIdx == -1 {
			tempInfos = append(tempInfos, tempInfo{feedInfo.LineProto, feedInfo.URL, false})
		} else {
			if len(tempInfos[urlIdx].lineProtos+feedInfo.LineProto) >= defaultMaxLength {
				// 字节数超标
				tempInfos[urlIdx].overflow = true
				tempInfos = append(tempInfos, tempInfo{feedInfo.LineProto, feedInfo.URL, false})
			} else {
				tempInfos[urlIdx].lineProtos = tempInfos[urlIdx].lineProtos + feedInfo.LineProto
			}
		}
	}

	// 数据搬移到发送队列
	for _, temp := range tempInfos {
		retryInfos = append(retryInfos, &retryInfo{temp.lineProtos, temp.urlStr, defaultRetryTimes})
	}
}

// doSend 执行发送任务。
func doSend(info *retryInfo) {
	fmt.Println("进入 doSend 发送:")
	fmt.Println(info.lineProtos)
	defer wg.Done()
	req, err := http.NewRequest(http.MethodPost, info.urlStr, strings.NewReader(info.lineProtos))
	if err != nil {
		fmt.Println(" error: ", err)
		atomic.AddInt64(&info.retryTimes, -1)
		return
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*defaultTimeout)
	defer ctxCancel()
	req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(" error: ", err)
		atomic.AddInt64(&info.retryTimes, -1)
		return
	}

	if resp.StatusCode/100 != 2 {
		atomic.AddInt64(&info.retryTimes, -1)
		return
	}

	// 发送成功，设置retryTimes=0 准备删除。
	atomic.StoreInt64(&info.retryTimes, 0)
}

// eraseSendedInfo 删除发送成功、重试超标的数据。
func eraseSendedInfo() {
	// 创建一个新的骨架
	newRetryInfos := make([]*retryInfo, 0)
	for _, v := range retryInfos {
		if v.retryTimes > 0 {
			newRetryInfos = append(newRetryInfos, v)
		}
	}
	// 需要 retry 的数据，搬移过来
	retryInfos = newRetryInfos
	fmt.Println("eraseSendedInfo剩余", len(retryInfos))
}

func init() {
	FeedCh = make(chan []FeedInfo, 10) // TODO 这个不建议是0，否则可能会迟滞exporter的进程。多少合适咧？
	retryInfos = make([]*retryInfo, 0)
	go feed() // TODO 只会创建一个协程，这个没说错吧？
}
