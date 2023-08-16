{{/* Regex: `\A\w+` */}}

{{$Sentence:=or (dbGet 0 "Sentence").Value (sdict "nil" "" "LastUser" 204255221017214977 "LastResp" 0)}}

{{if eq .User.ID $Sentence.LastUser}}
	Please wait your turn!
	{{deleteTrigger 5}}{{deleteResponse 5}}
	{{return}}
{{else if ne (len .Args) 1}}
	Please send one word at a time
	{{deleteTrigger 5}}{{deleteResponse 5}}
	{{return}}
{{end}}

{{deleteMessage nil $Sentence.LastResp 0}}
{{$Sentence.Set "nil" (joinStr " " $Sentence.nil (index .Args 0))}}
{{if reFind `\.$` .Message.Content}}
	{{sendMessage nil (cembed "title" "The sentence has concluded" "description" (printf "Final sentence was: `%s`" $Sentence.nil) "color" (randInt 0xFFFFFF))}}
	{{dbDel 0 "Sentence"}}
{{else}}
	{{$Sentence.Set "LastResp" (sendMessageRetID nil (printf "The sentence is now: `%s`" $Sentence.nil))}}
	{{dbSet 0 "Sentence" (sdict "nil" $Sentence.nil "LastUser" .User.ID "LastResp" $Sentence.LastResp)}}
{{end}}
