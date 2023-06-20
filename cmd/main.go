package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/filter"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type moji struct {
	text     string
	row      int
	col      int
	wishlist []string
}

func main() {
	// oneblock()
	// primeNumbers()

	// base := "こいびと"

	fmt.Println("start!")
	file, err := os.Create("application.log")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	log.SetOutput(file)

	findstr := "恋人"
	var fil []string
	base := chkString(findstr, fil)
	maxBoin := 11
	var hiragana = [][]string{
		[]string{"ア", "イ", "ウ", "エ", "オ"},
		[]string{"カ", "キ", "ク", "ケ", "コ"},
		[]string{"サ", "シ", "ス", "セ", "ソ"},
		[]string{"タ", "チ", "ツ", "テ", "ト"},
		[]string{"ナ", "ニ", "ヌ", "ネ", "ノ"},
		[]string{"ハ", "ヒ", "フ", "ヘ", "ホ"},
		[]string{"マ", "ミ", "ム", "メ", "モ"},
		[]string{"ヤ", "-", "ユ", "-", "ヨ"},
		[]string{"ラ", "リ", "ル", "レ", "ロ"},
		[]string{"ワ", "-", "-", "", "ヲ"},
		[]string{"ン", "-", "-", "-", "-"},
		[]string{"バ", "ビ", "ブ", "ベ", "ボ"},
	}
	fmt.Println(len(base))

	// 文字数を出す
	lens := utf8.RuneCountInString(base)
	fmt.Println(lens)

	// 配列を作る
	hako := make([]moji, lens)

	startIdx := 0
	endIdx := 3

	// 場所を取り出す
	for ii := 0; ii < lens; ii++ {
		if ii > 0 {
			startIdx += 3
			endIdx += 3
		}
		// 0->0,3
		// 1->3,6
		// 2->6,9
		fmt.Println(base[startIdx:endIdx])

		hako[ii].text = base[startIdx:endIdx]

		for key, valMother := range hiragana {

			for keyCol, valChild := range valMother {
				// fmt.Printf("%v,%v:%v \n", key, keyCol, valChild)

				if hako[ii].text == valChild {
					hako[ii].row = key
					hako[ii].col = keyCol
				}
			}
		}
	}
	boin := ""
	// 母音
	for _, v := range hako {

		// fmt.Printf("【母音】%v:%v \n", k, hiragana[0][v.col])

		boin += hiragana[0][v.col]
	}
	fmt.Println(boin)

	// 考えられる組み合わせをすべて展開する

	strItem := ""
	alllens := 1
	for ptnIdx := 0; ptnIdx < len(hako); ptnIdx++ {
		alllens *= maxBoin
	}

	// 検索ワードを１つずつ取り出す
	for k, v := range hako {
		lst := make([]string, maxBoin)
		// 候補
		for j := 0; j < maxBoin; j++ {

			// 未定義は除く
			if hiragana[j][v.col] != "-" {
				strItem = hiragana[j][v.col]
				// log.Printf("【%v.候補】:%v \n", k, strItem)
			}
			lst = append(lst, strItem)
		}

		setUnq := func(s []string) []string {
			m := make(map[string]bool)
			var result []string
			for _, v := range s {
				if _, ok := m[v]; !ok {
					m[v] = true
					result = append(result, v)
				}
			}
			return result
		}

		hako[k].wishlist = setUnq(lst)
		fmt.Println(hako[k].wishlist)
	}

	items := make([][]string, lens)
	for i := 0; i < lens; i++ {
		items[i] = hako[i].wishlist
	}
	patternALL(items, lens, base)
	// for k, v := range hako {
	// 	fmt.Printf("No%v:%v \n", k, v)
	// }
}

// patternALL is 文字数が一致したすべての組み合わせを表示する、
func patternALL(s [][]string, orgLen int, orgStr string) {

	var result []string
	for _, v := range s {
		if len(result) == 0 {
			result = v
		} else {
			var temp []string
			for _, r := range result {
				for _, w := range v {
					temp = append(temp, r+w)
				}
			}
			result = temp
		}
	}
	chkString(orgStr, filter.POS{"名詞"})
	for _, v := range result {
		chk := utf8.RuneCountInString(v)

		if orgLen == chk {
			log.Println(v)
			chkString(v, filter.POS{"名詞"})
		}

	}
}

// chkString is 辞書チェック
func chkString(item string, fill []string) string {

	ret := ""
	t, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		panic(err)
	}

	// tokens := t.Tokenize(item)
	// フィルタリングする

	if len(fill) > 0 {
		log.Printf("Filter-ON:%v\n", item)
		// 名詞フィルター
		var nounFilter = filter.NewPOSFilter(fill)
		tokens := t.Tokenize(item)
		nounFilter.Keep(&tokens)
		for _, token := range tokens {
			yomi, flg := token.Reading()
			if flg {
				log.Printf("%v \n", yomi)
				ret = yomi
			}

			features := strings.Join(token.Features(), ",")
			log.Printf("フィルタ有:%s|%v\n", token.Surface, features)
		}
	} else {

		log.Printf("Filter-OFF:%v\n", item)
		tokens := t.Tokenize(item)
		for _, token := range tokens {
			yomi, flg := token.Reading()
			if flg {
				log.Printf("%v \n", yomi)
				ret = yomi
			}

			features := strings.Join(token.Features(), ",")
			log.Printf("フィルタなし:%s|%v\n", token.Surface, features)
		}
	}

	return ret
}
