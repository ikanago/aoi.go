package main

import (
	"github.com/bwmarrin/discordgo"
)

// HelpMessageEmbeds is contents of `help` command.
var HelpMessageEmbeds = []*discordgo.MessageEmbedField{
	{
		Name: "カスタム絵文字作成",
		Value: "アオイチャンは条件にあったTweetを #subscription に転送してくれます\n" +
			"以下のコマンドでどんなTweetを取得するかというフィルタを設定できます\n" +
			"`ID`:        Tweetを取得したいアカウントのTwitter ID\n" +
			"`KEYWORDS`:  転送されるツイートが含むべきキーワード(OR検索)を空白区切りで入力してください\n" +
			"`CHANNEL:    ツイートを転送するチャンネル(デフォルトは subscription)`\n" +
			"キーワードを設定せずにフィルタを作った場合，そのアカウントのすべてのツイートが転送されます\n" +
			"`@Aoi tweet create ID KEYWORDS`:  新しいフィルタを作ります\n" +
			"`@Aoi tweet add    ID KEYWORDS`:  既に存在するIDに対してキーワードを追加します\n" +
			"`@Aoi tweet remove ID KEYWORDS`:  既に存在するIDに対してキーワードを削除します\n" +
			"`@Aoi tweet delete ID         `:  そのIDに関するフィルタをすべて削除します\n" +
			"`@Aoi tweet change ID        ` :  そのIDのツイートを今見ているチャンネルに送信するようにします\n" +
			"`@Aoi tweet show`              :  現在登録されているフィルタの一覧を表示します",
	},
	{
		Name: "アイデアの記録",
		Value: "あとから一覧にして見返せるように発言を各チャンネルごとに記録できます\n" +
			"`TEXT`:  記録したい発言\n" +
			"`@Aoi memo TEXT`:  発言を記録します 空白や改行を入れても問題ありません\n" +
			"`@Aoi memo show`:  そのチャンネルで記録された発言を一覧表示します",
	},
	{
		Name: "Ping",
		Value: "`@Aoi ping`\n" +
			"と打つと元気よく Pong! と返事をします",
	},
	{
		Name: "ヘルプ",
		Value: "使い方が分からない? そんなときは\n" +
			"`@Aoi help`\n" +
			"と打ってみましょう! きっとすぐに使えるようになりますよ!",
	},
}
