/*
リファクタリング中
処理をslacklog packageに移動していく。
一旦、必要な処理はすべてslacklog packageから一時的にエクスポートするか、このファ
イル内で定義している。
*/

package subcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/slack-go/slack"
	"github.com/urfave/cli/v2"
	"github.com/vim-jp/slacklog-generator/internal/slacklog"
)

// DownloadEmoji downloads and save emoji image files.
func DownloadEmoji(c *cli.Context) error {
	slackToken := os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		return fmt.Errorf("$SLACK_TOKEN required")
	}

	var emojisDir string
	if c.Args().Present() {
		emojisDir = filepath.Clean(c.Args().Get(0))
	} else {
		emojisDir = filepath.Clean(filepath.Join("_logdata", "emoji"))
	}
	emojiJSONPath := filepath.Join(emojisDir, "emoji.json")
	if 1 < c.Args().Len() {
		emojiJSONPath = filepath.Clean(c.Args().Get(1))
	}

	api := slack.New(slackToken)

	emojis, err := api.GetEmoji()
	if err != nil {
		return err
	}

	d := slacklog.NewDownloader(slackToken)

	go generateEmojiFileTargets(d, emojis, emojisDir)

	err = outputSummary(emojis, emojiJSONPath)
	if err != nil {
		return err
	}

	err = d.Wait()
	if err != nil {
		return err
	}
	return nil
}

func generateEmojiFileTargets(d *slacklog.Downloader, emojis map[string]string, outputDir string) {
	defer d.CloseQueue()
	err := os.MkdirAll(outputDir, 0777)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %s: %s", outputDir, err)
		return
	}

	for name, url := range emojis {
		if strings.HasPrefix(url, "alias:") {
			continue
		}
		ext := filepath.Ext(url)
		path := filepath.Join(outputDir, name+ext)
		d.QueueDownloadRequest(
			url,
			path,
			false,
		)
	}
}

func outputSummary(emojis map[string]string, path string) error {
	exts := make(map[string]string, len(emojis))
	for name, url := range emojis {
		if strings.HasPrefix(url, "alias:") {
			exts[name] = url
			continue
		}
		ext := filepath.Ext(url)
		exts[name] = ext
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(exts)
	if err != nil {
		return err
	}
	return nil
}
