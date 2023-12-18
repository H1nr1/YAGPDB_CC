{{/* Regex: `\A\-g(ive)?a(wa)?y?` */}}

{{$emoji:="🎉"}} {{/* emoji to use as reaction */}}

{{if .ExecData}}
	{{$g:=(dbGet .ExecData "giveaways").Value}}
	{{$winners:=cslice}}
	{{$e:=sdict 
	  	"title" "🌟 Giveaway Ended! 🌟" 
		"thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") 
		"footer" (sdict 
			"text" (printf "ID: %#v | Ended " .ExecData)
		) 
	  	"timestamp" currentTime 
		"color" 13938487
	}}
	{{if $g.Users}}
		{{$winCount:=len $g.users|min $g.winCount|toInt}}
		{{while gt $winCount 0}}
			{{- $winner:=len $g.users|add 1|randInt|index $g.users}}
			{{- if not (in $winners $winner)}}
				{{- $winners =$winners.Append $winner}}
				{{- $winCount =sub $winCount 1}}
			{{- end -}}
		{{end}}
		{{$mentWinners:=""}}
		{{range $winners}}
			{{- $mentWinners =printf "%s<@%d>" $mentWinners . -}}
		{{end}}
		{{$e.Set "description" (printf "**__Prize:__ %s\n__Entries:__ %d\n__Winners:__** %s" 
			$g.prize (len $g.users) $mentWinners
		)}}
		{{editMessage nil $g.MID (cembed $e)}}
		{{printf "**__Prize:__ %s\n__Winners:__** %s" $g.prize $mentWinners}}
	{{else}}
		{{$e.Set "description" (printf "**__Prize:__ %s\n__Winners:__ No Entries 🥲**" $g.prize)}}
		{{editMessage nil $g.MID (cembed $e)}}
	{{end}}
	{{deleteAllMessageReactions nil $g.MID}}
	{{dbSetExpire .ExecData "giveaways" (sdict 
		"prize" $g.prize 
		"users" $g.users 
		"winners" $winners
	) 86400}}
	{{return}}
{{end}}

{{deleteTrigger 3}}
{{$args:=reFind `(?i)(s(tart)? (\d+(s|mo?|h|d|w|y)?)+ \d+ .+)|(e(nd)? \d+)|(re?r(oll)?( \d+){2})` (joinStr " " .CmdArgs)}}

{{if not $args}}
	{{$id:=sendMessageRetID nil (cembed 
		"title" "Syntax Error" 
		"description" "```Giveaway Start <Duration> <int:maxWinners> <string:prize>\nGiveaway End <int:ID>\nGiveaway Reroll <int:id> <int:count>```" 
		"color" 16711680
	)|deleteMessage nil $id}}
	{{return}}
{{end}}

{{$subCmd:=reFind `(?i)s(tart)?|e(nd)?|re?r(oll)?` $args|lower}}

{{if in $subCmd "s"}}
	{{$args =parseArgs 4 "" 
		(carg "string" "") 
		(carg "duration" "") 
		(carg "int" "") 
		(carg "string" "")
	}}
	{{$dur:=$args.Get 1}}
	{{$winCount:=$args.Get 2}}
	{{$prize:=$args.Get 3}}
	{{$id:=sendMessageRetID nil (cembed 
		"title" "🌟 Giveaway 🌟" 
		"description" (printf "**__Prize:__ %s\n__Winners:__ %d\n__Duration:__ %s\n\nReact with %s to enter the Giveaway**" 
			$prize $winCount (humanizeDurationMinutes $dur) $emoji
		) 
		"thumbnail" (sdict "url" "https://media.discordapp.net/attachments/1036894710343671918/1070947166241173514/GL_Giveaway.png?width=419&height=419") 
		"footer" (sdict 
			"text" (printf "ID: %d | Ends " ($GID:=dbIncr 0 "giveawayID" 1|toInt))
		) 
		"timestamp" (currentTime.Add $dur) 
		"color" 13938487
	)}}
	{{addMessageReactions nil $id $emoji}}
	{{dbSet $GID "giveaways" (sdict 
		"MID" $id 
		"prize" $prize 
		"winCount" $winCount 
		"users" cslice
	)}}
	{{scheduleUniqueCC .CCID nil $dur.Seconds $GID $GID}}

{{else if eq $subCmd "e" "end"}}
	{{$args =parseArgs 2 "" 
		(carg "string" "") 
		(carg "int" "")
	}}
	{{$GID:=$args.Get 1}}
	{{if not (dbGet $GID "giveaways")}}
		Invalid Giveaway ID
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{cancelScheduledUniqueCC .CCID $GID}}
	{{execCC .CCID nil 1 $GID}}

{{else if eq $subCmd "rr" "reroll" "rroll"}}
	{{$args =parseArgs 3 "" 
		(carg "string" "") 
		(carg "int" "") 
		(carg "int" "")
	}}
	{{$GID:=$args.Get 1}}
	{{$amt:=$args.Get 2}}
	{{$g:=(dbGet $GID "giveaways").Value}}
	{{if not $g}}
		Invalid Giveaway ID
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{if gt $amt (len $g.users)}}
		Reroll amount cannot be greater than users entered
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{if le (len $g.users) (len $g.winners)}}
		Users entered must be greater than winners
		{{deleteResponse}}
		{{return}}
	{{end}}
	{{$users:=cslice}}
	{{range $g.users}}
		{{- if not (in $g.winners .)}}
			{{- $users =$users.Append .}}
		{{- end -}}
	{{end}}
	{{while gt $amt 0}}
		{{$winner:=len $users|add 1|randInt 0|index $users}}
		{{$g.Set "winners" ($g.winners.Append $winner)}}
		{{$amt =sub $amt 1}}
	{{end}}
	{{$mentWinners:=""}}
	{{range $g.winners}}
		{{- $mentWinners =printf "%s <@%d>" $mentWinners . -}}
	{{end}}
	{{printf "*Giveaway was rerolled %d times*\n**__Prize:__ %s\n__New Winners:__** %s" 
		($args.Get 2) $g.prize $mentWinners
	}}
{{end}}
