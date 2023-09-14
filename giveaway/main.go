{{/* Regex: `\A\-g(ive)?a(wa)?y?` */}}

{{$Emoji:="🎉"}} {{/* Emoji to use as reaction */}}

{{if .ExecData}}
	{{$G:=(dbGet .ExecData "Giveaways").Value}}
	{{$Winners:=cslice}}
	{{$e:=sdict 
	  	"title" "🌟 Giveaway Ended! 🌟" 
		"thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") 
		"footer" (sdict "text" (print "ID: " .ExecData " | Ended ")) 
	  	"timestamp" currentTime 
		"color" 13938487
	}}
	{{if $G.Users}}
		{{$WinCount:=toInt (min $G.WinCount (len $G.Users))}}
		{{while gt $WinCount 0}}
			{{- $Winner:=index $G.Users (randInt (len $G.Users|add 1))}}
			{{- if not (in $Winners $Winner)}}
				{{- $Winners =$Winners.Append $Winner}}
				{{- $WinCount =sub $WinCount 1}}
			{{- end -}}
		{{end}}
		{{$MentWinners:=""}}
		{{range $Winners}}
			{{- $MentWinners =printf "%s<@%d>" $MentWinners . -}}
		{{end}}
		{{$e.Set "description" (printf "**__Prize:__ %s\n__Entries:__ %d\n__Winners:__** %s" $G.Prize (len $G.Users) $MentWinners)}}
		{{editMessage nil $G.MsgID (cembed $e)}}
		{{printf "**__Prize:__ %s\n__Winners:__** %s" $G.Prize $MentWinners}}
	{{else}}
		{{$e.Set "description" (printf "**__Prize:__ %s\n__Winners:__ No Entries 🥲**" $G.Prize)}}
		{{editMessage nil $G.MsgID (cembed $e)}}
	{{end}}
	{{deleteAllMessageReactions nil $G.MsgID}}
	{{dbSetExpire .ExecData "Giveaways" (sdict "Prize" $G.Prize "Users" $G.Users "Winners" $Winners) 86400}}
	{{return}}
{{end}}

{{deleteTrigger 3}}
{{$Args:=reFind `(?i)(s(tart)? (\d+(s|mo?|h|d|w|y)?)+ \d+ .+)|(e(nd)? \d+)|(re?r(oll)?( \d+){2})` (joinStr " " .CmdArgs)}}

{{if not $Args}}
	{{$ID:=sendMessage nil (cembed "title" "Syntax Error" "description" "```Giveaway Start <Duration> <int:MaxWinners> <string:Prize>\nGiveaway End <int:ID>\nGiveaway Reroll <int:ID> <int:Count>```" "color" 16711680)}}
	{{deleteMessage nil $ID 25}}
	{{return}}
{{end}}

{{$SubCmd:=reFind `(?i)s(tart)?|e(nd)?|re?r(oll)?` $Args|lower}}

{{if in $SubCmd "s"}}
	{{$Args =parseArgs 4 "" 
		(carg "string" "") 
		(carg "duration" "") 
		(carg "int" "") 
		(carg "string" "")
	}}
	{{$Dur:=$Args.Get 1}}
	{{$WinCount:=$Args.Get 2}}
	{{$Prize:=$Args.Get 3}}
	{{$ID:=sendMessageRetID nil (cembed 
		"title" "🌟 Giveaway 🌟" 
		"description" (printf "**__Prize:__ %s\n__Winners:__ %d\n__Duration:__ %s\n\nReact with %s to enter the Giveaway**" $Prize $WinCount (humanizeDurationMinutes $Dur) $Emoji) 
		"thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") 
		"footer" (sdict "text" (printf "ID: %d | Ends " ($GID:=dbIncr 0 "GiveawayID" 1|toInt))) 
		"timestamp" (currentTime.Add $Dur) 
		"color" 13938487
	)}}
	{{addMessageReactions nil $ID $Emoji}}
	{{dbSet $GID "Giveaways" (sdict "MsgID" $ID "Prize" $Prize "WinCount" $WinCount "Users" cslice)}}
	{{scheduleUniqueCC .CCID nil $Dur.Seconds $GID $GID}}

{{else if eq $SubCmd "e" "end"}}
	{{$Args =parseArgs 2 "" 
		(carg "string" "") 
		(carg "int" "")
	}}
	{{$GID:=$Args.Get 1}}
	{{if not (dbGet $GID "Giveaways")}}
		Invalid Giveaway ID.
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{cancelScheduledUniqueCC .CCID $GID}}
	{{execCC .CCID nil 0 $GID}}

{{else if eq $SubCmd "rr" "reroll" "rroll"}}
	{{$Args =parseArgs 3 "" 
		(carg "string" "") 
		(carg "int" "") 
		(carg "int" "")
	}}
	{{$GID:=$Args.Get 1}}
	{{$Amt:=$Args.Get 2}}
	{{$G:=(dbGet $GID "Giveaways").Value}}
	{{if not $G}}
		Invalid Giveaway ID.
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{if gt $Amt (len $G.Users)}}
		Reroll amount cannot be greater than users entered.
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{if le (len $G.Users) (len $G.Winners)}}
		Users entered must be greater than winners.
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{$Users:=cslice}}
	{{range $G.Users}}
		{{- if not (in $G.Winners .)}}
			{{- $Users =$Users.Append .}}
		{{- end -}}
	{{end}}
	{{while gt $Amt 0}}
		{{$Winner:=index $Users (randInt 0 (len $Users|add 1))}}
		{{$G.Set "Winners" ($G.Winners.Append $Winner)}}
		{{$Amt =sub $Amt 1}}
	{{end}}
	{{$MentWinners:=""}}
	{{range $G.Winners}}
		{{- $MentWinners =printf "%s <@%d>" $MentWinners . -}}
	{{end}}
	{{printf "*Giveaway was rerolled %d times*\n**__Prize:__ %s\n__New Winners:__** %s" ($Args.Get 2) $G.Prize $MentWinners}}
{{end}}
