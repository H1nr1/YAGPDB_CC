{{/* Regex: `\A\-g(ive)?a(wa)?y?` */}}

{{$Emoji:="🎉"}}

{{$err:=false}}

{{if .ExecData}}
	{{$G:=(dbGet .ExecData "Giveaways").Value}}
	{{$MsgIDsNew:=cslice}}{{$Winners:=cslice}}
	{{$e:=sdict "title" "🌟 Giveaway Ended! 🌟" "thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") "footer" (sdict "text" (print "ID: " .ExecData " | Ended ")) "timestamp" currentTime "color" 13938487}}
	{{if $G.Users}}
		{{$WinCount:=toInt (min $G.WinCount (len $G.Users))}}
		{{while gt $WinCount 0}}
			{{- $Winner:=index $G.Users (randInt (len $G.Users|add 1))}}
			{{- if not (in $Winners $Winner)}}
				{{- $Winners =$Winners.Append $Winner}}
				{{- $WinCount =sub $WinCount 1}}
			{{- end -}}
		{{end}}
		{{$MentionWinners:=""}}{{range $Winners}}{{- $MentionWinners =print $MentionWinners (userArg .).Mention -}}{{end}}
		{{$e.Set "description" (printf "**__Prize:__ %s\n__Entries:__ %d\n__Winners:__** %s" $G.Prize (len $G.Users) $MentionWinners)}}
		{{sendMessage nil (print "**__Prize:__ " $G.Prize "\n__Winners:__** " $MentionWinners)}}
		{{range (dbGet 0 "GiveawayMsgIDs").Value}}{{- if ne $G.MsgID .}}{{- $MsgIDsNew =$MsgIDsNew.Append .}}{{- end -}}{{end}}
	{{else}}{{$e.Set "description" (printf "**__Prize:__ %s\n__Winners:__ No Entries 🥲**" $G.Prize)}}{{end}}
	{{editMessage nil $G.MsgID (cembed $e)}}
	{{deleteAllMessageReactions nil $G.MsgID}}
	{{dbSet 0 "GiveawayMsgIDs" $MsgIDsNew}}
	{{dbSetExpire .ExecData "Giveaways" (sdict "Prize" $G.Prize "Users" $G.Users "Winners" $Winners) 86400}}
	{{return}}
{{end}}

{{if .CmdArgs}}
	{{$Args =slice ($Args:=print .CmdArgs) 1 (sub (len $Args) 1)}}
	{{if reFind `(?i)s(tart)? (\d+(s|mo?|h|d|w|y)?)+ \d+ .+` $Args}}
		{{$Args =parseArgs 4 "" (carg "string" "SubCmd") (carg "duration" "Duration") (carg "int" "WinCount") (carg "string" "Prize")}}
		{{$Dur:=$Args.Get 1}}{{$WinCount:=$Args.Get 2}}{{$Prize:=$Args.Get 3}}
		{{$ID:=sendMessageRetID nil (cembed "title" "🌟 Giveaway 🌟" "description" (print "**__Prize:__ " $Prize "\n__Winners:__ " $WinCount "\n__Duration:__ " (humanizeDurationMinutes $Dur) "\n\nReact with " $Emoji " to enter the Giveaway**") "thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") "footer" (sdict "text" (print "ID: " ($GID := toInt (dbIncr 0 "GiveawayID" 1)) " | Ends ")) "timestamp" (currentTime.Add $Dur) "color" 13938487)}}
		{{addMessageReactions nil $ID $Emoji}}
		{{dbSet $GID "Giveaways" (sdict "MsgID" $ID "Prize" $Prize "WinCount" $WinCount "Users" cslice)}}
		{{dbSet 0 "GiveawayMsgIDs" ((dbGet 0 "GiveawayMsgIDs").Value.Append $ID)}}
		{{scheduleUniqueCC .CCID nil $Dur.Seconds $GID $GID}}
	{{else if reFind `(?i)e(nd)? \d+` $Args}}
		{{$Args =parseArgs 2 "" (carg "string" "SubCmd") (carg "int" "GID")}}{{$GID:=$Args.Get 1}}
		{{if dbGet $GID "Giveaways"}}{{cancelScheduledUniqueCC .CCID $GID}}{{execCC .CCID nil 0 $GID}}
		{{else}}Invalid Giveaway ID.{{end}}
	{{else if reFind `(?i)re?r(oll)?( \d+){2}` $Args}}
		{{$Args =parseArgs 3 "" (carg "string" "SubCmd") (carg "int" "GID") (carg "int" "Amount")}}
		{{$GID:=$Args.Get 1}}{{$Amount:=$Args.Get 2}}
		{{if ($G:=(dbGet $GID "Giveaways").Value)}}
			{{if and (le $Amount (len $G.Users)) (gt (len $G.Users) (len $G.Winners))}}
				{{$Users:=cslice}}{{range $G.Users}}{{- if not (in $G.Winners .)}}{{- $Users =$Users.Append .}}{{- end -}}{{end}}
				{{while gt $Amount 0}}
					{{$Winner:=index $Users (randInt 0 (len $Users|add 1))}}
					{{$G.Set "Winners" ($G.Winners.Append $Winner)}}
					{{$Amount =sub $Amount 1}}
				{{end}}
				{{$MentionWinners:=""}}{{range $G.Winners}}{{- $MentionWinners =joinStr " " $MentionWinners (userArg .).Mention -}}{{end}}
				{{sendMessage nil (print "*Giveaway was rerolled " ($Args.Get 2) " times*\n**__Prize:__ " $G.Prize "\n__New Winners:__** " $MentionWinners)}}
			{{else}}Reroll amount cannot be greater than users entered. Users entered must be greater than winners.{{end}}
		{{else}}Invalid Giveaway ID.{{end}}
	{{else}}{{$err =true}}{{end}}
{{else}}{{$err =true}}{{end}}

{{if $err}}{{sendMessage nil (cembed "title" "Syntax Error" "description" "```Giveaway Start <Duration> <int:MaxWinners> <string:Prize>\nGiveaway End <int:ID>\nGiveaway Reroll <int:ID> <int:Count>```" "color" 16711680)}}{{end}}

{{deleteTrigger 3}}
