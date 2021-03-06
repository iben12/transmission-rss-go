package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

const (
	ValidXml string = `
		<?xml version="1.0" encoding="UTF-8"?>
			<rss version="2.0" xmlns:tv="https://showrss.info">
				<channel>
					<item>
						<title>Million Dollar Listing: Los Angeles 13x13 The Great British Cook-Off 720p</title>
						<link>magnet:?xt=urn:btih:ADF9E7857DD23232E829747A8A08E417E31622A9&amp;dn=Million+Dollar+Listing+Los+Angeles+S13E13+The+Great+British+Cook+Off+720p+WEBRip+x264+KOMPOST&amp;tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&amp;tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&amp;tr=udp%3A%2F%2Fexplodie.org%3A6969&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2960&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980</link>
						<guid isPermaLink="false">0e295e7b5eef198e75d4ab6d9f5c7c10116936de</guid>
						<pubDate>Fri, 10 Dec 2021 23:40:21 +0000</pubDate>
						<description>New episode: Million Dollar Listing Los Angeles S13E13 The Great British Cook Off 720p WEBRip x264 KOMPOST. Link: &lt;a href=&quot;magnet:?xt=urn:btih:ADF9E7857DD23232E829747A8A08E417E31622A9&amp;dn=Million+Dollar+Listing+Los+Angeles+S13E13+The+Great+British+Cook+Off+720p+WEBRip+x264+KOMPOST&amp;tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&amp;tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&amp;tr=udp%3A%2F%2Fexplodie.org%3A6969&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2960&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980&quot;&gt;magnet:?xt=urn:btih:ADF9E7857DD23232E829747A8A08E417E31622A9&amp;dn=Million+Dollar+Listing+Los+Angeles+S13E13+The+Great+British+Cook+Off+720p+WEBRip+x264+KOMPOST&amp;tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&amp;tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&amp;tr=udp%3A%2F%2Fexplodie.org%3A6969&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2960&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980&lt;/a&gt;</description>
						<tv:show_id>1499</tv:show_id>
						<tv:external_id>3333</tv:external_id>
						<tv:show_name>Million Dollar Listing: Los Angeles</tv:show_name>
						<tv:episode_id>153955</tv:episode_id>
						<tv:raw_title>Million Dollar Listing Los Angeles S13E13 The Great British Cook Off 720p WEBRip x264 KOMPOST</tv:raw_title>
						<tv:info_hash>ADF9E7857DD23232E829747A8A08E417E31622A9</tv:info_hash>
						<enclosure url="magnet:?xt=urn:btih:ADF9E7857DD23232E829747A8A08E417E31622A9&amp;dn=Million+Dollar+Listing+Los+Angeles+S13E13+The+Great+British+Cook+Off+720p+WEBRip+x264+KOMPOST&amp;tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&amp;tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&amp;tr=udp%3A%2F%2Fexplodie.org%3A6969&amp;tr=udp%3A%2F%2F9.rarbg.me%3A2960&amp;tr=udp%3A%2F%2F9.rarbg.to%3A2980" length="0" type="application/x-bittorrent" />
					</item>
				</channel>
			</rss>
		`

	InvalidXml string = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><rss>"
)

var (
	EpisodeExample trss.Episode = trss.Episode{
		ShowId:    "1",
		EpisodeId: "2",
		Title:     "From test",
		ShowTitle: "Testdriver",
		Link:      "magnet://hjfgsfgjfbjfbajsfb",
	}

	FeedItemExample trss.FeedItem = trss.FeedItem{
		ShowId:    "1",
		EpisodeId: "2",
		Title:     "From test",
		ShowTitle: "Testdriver",
		Link:      "magnet://hjfgsfgjfbjfbajsfb",
	}
)
