{{/* Regex \A\w+ */}}
{{/* Restrict to One Word Story channel */}}

{{$sentence:=or 
  	(dbGet 0 "oneWordStory").Value 
 	(sdict 
  		"nil" "" 
  		"last" (sdict 
			"user" .BotUser.ID 
  			"msg" 0
		)
	)
 }}
 
{{if eq .User.ID $sentence.last.user}}
	Please wait your turn
	{{deleteTrigger 1}}
	{{deleteResponse 5}}
	{{return}}
{{end}}

{{if reFind .LinkRegex .Message.Content}}
	Links are not permitted
	{{deleteTrigger 1}}
	{{deleteResponse 5}}
	{{return}}
{{end}}

{{if gt (len .Args) 1}}
	Please only send one word at a time
	{{deleteTrigger 1}}
	{{deleteResponse 5}}
	{{return}}
{{end}}

{{if eq (exec "dictionary" .Cmd|str) "Could not find a definition for that word."}}
	Invalid word
	{{deleteTrigger 1}}
	{{deleteResponse 5}}
	{{return}}
{{end}}

{{deleteMessage nil $sentence.last.msg 0}}
{{index .Args 0|joinStr " " $sentence.nil|$sentence.Set "nil"}}

{{if hasSuffix .Message.Content "."}}
	{{printf "The sentence has concluded\nFinal sentence is: '%s'" $sentence.nil}}
	{{dbDel 0 "oneWordStory"}}
{{else}}
	{{sendMessageRetID nil (printf "The sentence is now: '%s'" $sentence.nil)|$sentence.last.Set "msg"}}
	{{dbSet 0 "oneWordStory" (sdict 
		"nil" $sentence.nil 
		"last" (sdict 
			"user" .User.ID 
			"msg" $sentence.last.msg
		)
	)}}
{{end}}
