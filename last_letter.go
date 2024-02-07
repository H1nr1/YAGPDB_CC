{{/* Regex \A(-l(ast)?l(etter)?|\S+\z) */}}
{{/* Restrict to Last Letter channel */}}

{{/* configurable variables */}}
{{/* required */}}
{{ $staffRIDs := cslice }} {{/* staff Role IDs */}}
{{/* optional */}}
{{ $twice := false }} {{/* allow users to take a turn multiple times in a row; true/false */}}
{{ $repeats := false }} {{/* allow members to reuse words */}}
{{ $subWrong := true }} {{/* subtract length of word from score if wrong */}}
{{ $correctEmoji := "✅" }} {{/* emoji to react with if number is correct; if custom, use format name:id */}}
{{ $incorrectEmoji := "❌" }} {{/* emoji to react with if number is incorrect; if custom, use format name:id */}}
{{ $LBLength := 10 }} {{/* how many members to show on leaderboard when event ends */}}
{{/* end of configurable variables */}}

{{$db:=(dbGet 0 "lastLetter").Value}}

{{with .ExecData}}
	{{$msg:=getMessage nil .ID}}{{$out:=""}}
	{{if not $msg}} {{/* check if message was deleted */}}
		{{$out =printf "deleted their word of `%s` which was correct!" .content}}
	{{else if $msg.EditedTimestamp}} {{/* check if message was edited */}}
		{{$out =printf "edited their word of `%s` to `%s`..." .content $msg.Content}}
	{{end}}
	{{if $out}}
		{{sendMessage nil (cembed 
			"description" (printf "%s %s\nPlease send a word starting with `%s`" 
				$.User.Mention $out (len $db.word|add -1|slice $db.word)
			) 
			"color" (randInt 0xFFFFFF)
		)}}
	{{end}}
	{{return}}
{{end}}

{{$isStaff:=false}}{{range $staffRIDs}}{{- if hasRoleID .}}{{- $isStaff =true}}{{- end -}}{{end}}

{{if reFind `(?i)-l(ast)?l(etter)?` .Cmd}}
	{{$subCmd:=joinStr " " .CmdArgs|reFind `(?i)l(eader)?b(oard)?|start|end`|lower}}
	{{if not (eq $subCmd "lb" "leaderboard" "start" "end")}} {{/* Invalid subcommand */}}
		{{template "SE" "Syntax is `-LastLetter <Leaderboard|Start|End>`"}}
		{{return}}
	{{end}}
	{{if or (eq $subCmd "lb" "leaderboard") (and $isStaff (eq $subCmd "end"))}} {{/* LB or end subcommand */}}
 		{{$desc:=""}}{{$pos:=1}}
		{{range dbTopEntries "lastLetter" $LBLength 0}} {{/* build LB */}}
			{{- $desc =printf "%s\n#%-3d %4d - %-4v" $desc $pos (toInt .Value) (or (userArg .UserID) .UserID)}}
			{{- $pos =add $pos 1 -}}
		{{end}}
		{{if eq $subCmd "end"}} {{/* end subcommand */}}
			{{if not $db}}{{template "SE" "No event found. Use `-LastLetter Start` to begin an event."}}{{return}}{{end}}
			{{sendMessage nil (cembed 
				"author" (sdict 
					"icon_url" (.Guild.IconURL "512") 
					"name" "🍿 Last Letter - Event Concluded 🍿"
				) 
				"description" (printf "__Leaderboard__\n```Pos Score   User\n%s```\nEvent lasted %s" 
					$desc (currentTime.Sub (dbGet 0 "lastLetter").CreatedAt|humanizeDurationMinutes)
				) 
				"footer" (sdict "text" "Nice work everyone!") 
				"color" (randInt 0xFFFFFF)
			)}}
			{{$foo:=dbDelMultiple (sdict "pattern" "lastLetter") 100 0}}
		{{else}} {{/* LB subcommand */}}
			{{sendMessage nil (cembed 
				"author" (sdict 
					"icon_url" (.Guild.IconURL "512") 
					"name" "🍿 Last Letter - Leaderboard 🍿"
				) 
				"description" (printf "```Pos Score   User\n%s```" $desc) 
				"footer" (sdict "text" "Nice work everyone!") 
				"color" (randInt 0xFFFFFF)
			)}}
		{{end}}

	{{else if and $isStaff (eq $subCmd "start")}} {{/* start subcommand */}}
		{{if $db}}{{template "SE" "Event is already in progress. Use `-LastLetter End` to conclude the current event."}}{{return}}{{end}}
		{{dbSet 0 "lastLetter" (sdict 
			"user" 204255221017214977 
			"word" ($word:=lower noun) 
			"words" (cslice $word)
		)}}
		{{$eRepeats:=""}}{{$eTwice:=""}}
		{{if not $repeats}}{{$eRepeats ="\nCannot repeat previous words"}}{{end}}
		{{if not $twice}}{{$eTwice ="\nCannot take two turns in a row"}}{{end}}
		{{sendMessage nil (cembed 
			"author" (sdict 
				"icon_url" (.Guild.IconURL "512") 
				"name" "🍿 Last Letter - Event Started 🍿"
			) 
			"description" (printf "Last letter event has been started!\n\n**__Instructions__**:\nSend a word that starts with the last letter of the previous word\n**__Rules__**\nWord must be available in english dictionary\nYou will be scored according to the length of your word%s%s" $eRepeats $eTwice) 
			"footer" (sdict "text" (printf "The game will start with '%s'" $word)) 
			"color" (randInt 0xFFFFFF)
		)}}
	{{end}}
	{{deleteTrigger 3}}
	{{return}}
{{end}}

{{if not $db}}{{template "SE" (printf "An event has not been started!\nAsk <@&%d> to start an event using `-LastLetter Start`" (index $staffRIDs 0))}}{{return}}{{end}}
{{if and (not $twice) (eq $db.user .User.ID)}}{{template "SE" "Please wait your turn."}}{{return}}{{end}}
{{if eq (exec "dictionary" .Cmd|str) "Could not find a definition for that word."}}{{template "SE" "I couldn't find that word in my dictionary."}}{{return}}{{end}}
{{if and (not $repeats) (inFold $db.words .Cmd)}}{{template "SE" "Word has already been used. Please pick another."}}{{return}}{{end}}

{{if eq (slice .Cmd 0 1|lower) (len $db.word|add -1|slice $db.word)}} {{/* correct last-to-first letter */}}
	{{try}}{{addReactions $correctEmoji}}{{catch}}{{end}}
	{{$db.Set "user" .User.ID}}{{$db.Set "word" (lower .Cmd)}}
	{{$db.Set "words" ($db.words.Append (lower .Cmd))}}
	{{$s:=dbIncr .User.ID "lastLetter" (len .Cmd)}}
	{{execCC .CCID nil 15 .Message}} {{/* call to check for message deletion/edit */}}
{{else}} {{/* incorrect word */}}
	{{try}}{{addReactions $incorrectEmoji}}{{catch}}{{end}}
	{{if $subWrong}}{{$s:=dbIncr .User.ID "lastLetter" (len .Cmd|mult -1)}}{{end}}
	{{template "SE" (printf "`%s` doesn't start with the letter `%s`" 
		.Cmd (len $db.word|add -1|slice $db.word))
	}}
{{end}}
{{dbSet 0 "lastLetter" $db}}

{{define "SE"}}
	{{$ID:=sendMessageRetID nil (cembed 
		"description" . 
		"color" 16711680
	)}}
	{{deleteTrigger 3}}
	{{deleteMessage nil $ID 15}}
{{end}}
