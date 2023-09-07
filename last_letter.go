{{/* Regex: `\A(-l(ast)?l(etter)?|\S+\z)` */}}

{{/* Required */}}
{{ $StaffRIDs := cslice }} {{/* Staff Role IDs */}}
{{/* Optional */}}
{{ $Twice := false }} {{/* Allow users to take a turn multiple times in a row; true/false */}}
{{ $Repeats := false }} {{/* Allow members to reuse words */}}
{{ $SubWrong := true }} {{/* Subtract length of word from score if wrong */}}
{{ $CorrectEmoji := "✅" }} {{/* Emoji to react with if number is correct; if custom, use format name:id */}}
{{ $IncorrectEmoji := "❌" }} {{/* Emoji to react with if number is incorrect; if custom, use format name:id */}}
{{ $LBLength := 10 }} {{/* How many members to show on leaderboard when event ends */}}
{{/* End of configurable variables */}}

{{$db:=(dbGet 0 "LastLetter").Value}}

{{with .ExecData}}
	{{$foo:=""}}
	{{if not (getMessage nil .ID)}} {{/*Check if message was deleted */}}
		{{$foo =printf "deleted their word of `%s` which was correct!" .Content}}
	{{else if ne (getMessage nil .ID).Content .Content}} {{/* Check if message was edited */}}
		{{$foo =printf "edited their word of `%s` to `%s`..." .Content (getMessage nil .ID).Content}}
	{{end}}
	{{if $foo}}{{sendMessage nil (cembed "description" (printf "%s %s\nPlease send a word starting with `%s`" $.User.Mention $foo (slice $db.Word (len $db.Word|add -1))) "color" (randInt 0xFFFFFF))}}{{end}}
	{{return}}
{{end}}

{{$isStaff:=false}}{{range $StaffRIDs}}{{- if hasRoleID .}}{{- $isStaff =true}}{{- end -}}{{end}}

{{if reFind `(?i)-l(ast)?l(etter)?` .Cmd}}
	{{$SubCmd:=lower (reFind `(?i)l(eader)?b(oard)?|start|end` (joinStr " " .CmdArgs))}}
	{{if not (eq $SubCmd "lb" "leaderboard" "start" "end")}} {{/* Invalid subcommand */}}
		{{template "SE" (sdict "Desc" "Syntax is `-LastLetter <Leaderboard|Start|End>`")}}
		{{return}}
	{{end}}
	{{if or (eq $SubCmd "lb" "leaderboard") (and $isStaff (eq $SubCmd "end"))}} {{/* LB or end subcommand */}}
 		{{$Desc:=""}}{{$Place:=1}}
		{{range dbTopEntries "LastLetter" $LBLength 0}} {{/* Build LB */}}
			{{- $Desc =printf "%s\n#%-3d %4d - %-4v" $Desc $Place (toInt .Value) (or (userArg .UserID) .UserID)}}
			{{- $Place =add $Place 1 -}}
		{{end}}
		{{if eq $SubCmd "end"}} {{/* End subcommand */}}
			{{if not $db}}{{template "SE" (sdict "Desc" "No event found. Use `-LastLetter Start` to begin an event.")}}{{return}}{{end}}
			{{sendMessage nil (cembed "author" (sdict "icon_url" (.Guild.IconURL "512") "name" "🍿 Last Letter - Event Concluded 🍿") "description" (printf "__Leaderboard__\n```Pos Score   User\n%s```\nEvent lasted %s" $Desc (humanizeDurationMinutes (currentTime.Sub (dbGet 0 "LastLetter").CreatedAt))) "footer" (sdict "text" "Nice work everyone!") "timestamp" currentTime "color" (randInt 0xFFFFFF))}}
			{{$foo:=dbDelMultiple (sdict "pattern" "LastLetter") 100 0}}
		{{else}} {{/* LB subcommand */}}
			{{sendMessage nil (cembed "author" (sdict "icon_url" (.Guild.IconURL "512") "name" "🍿 Last Letter - Leaderboard 🍿") "description" (printf "```Pos Score   User\n%s```" $Desc) "footer" (sdict "text" "Nice work everyone!") "timestamp" currentTime "color" (randInt 0xFFFFFF))}}
		{{end}}
	{{else if and $isStaff (eq $SubCmd "start")}} {{/* Start subcommand */}}
		{{if $db}}{{template "SE" (sdict "Desc" "Event is already in progress. Use `-LastLetter End` to conclude the current event.")}}{{return}}{{end}}
		{{dbSet 0 "LastLetter" (sdict "User" 204255221017214977 "Word" ($Word:=lower noun) "Words" (cslice $Word))}}
		{{$ERepeats:=""}}{{$ETwice:=""}}
		{{if not $Repeats}}{{$ERepeats ="\nCannot repeat previous words"}}{{end}}
		{{if not $Twice}}{{$ETwice ="\nCannot take two turns in a row"}}{{end}}
		{{sendMessage nil (cembed "author" (sdict "icon_url" (.Guild.IconURL "512") "name" "🍿 Last Letter - Event Started 🍿") "description" (print "Last letter event has been started!\n\n**__Instructions__**:\nSend a word that starts with the last letter of the previous word\n**__Rules__**\nWord must be available in english dictionary\nYou will be scored according to the length of your word" $ERepeats $ETwice) "footer" (sdict "text" (printf "The game will start with '%s'" $Word)) "timestamp" currentTime "color" (randInt 0xFFFFFF))}}
	{{end}}
	{{deleteTrigger 3}}
	{{return}}
{{end}}

{{if not $db}}{{template "SE" (sdict "Desc" (print "An event has not been started!\nAsk <@&" (index $StaffRIDs 0) "> to start an event using `-LastLetter Start`"))}}{{return}}{{end}}
{{if and (not $Twice) (eq $db.User .User.ID)}}{{template "SE" (sdict "Desc" "Please wait your turn.")}}{{return}}{{end}}
{{if eq (exec "dictionary" .Cmd|str) "Could not find a definition for that word."}}{{template "SE" (sdict "Desc" "I couldn't find that word in my dictionary.")}}{{return}}{{end}}
{{if and (not $Repeats) (inFold $db.Words .Cmd)}}{{template "SE" (sdict "Desc" "Word has already been used. Please pick another.")}}{{return}}{{end}}

{{if eq (lower (slice .Cmd 0 1)) (slice $db.Word (sub (len $db.Word) 1))}} {{/* Correct last-to-first letter */}}
	{{try}}{{addReactions $CorrectEmoji}}{{catch}}{{end}}
	{{$db.Set "User" .User.ID}}{{$db.Set "Word" (lower .Cmd)}}{{$db.Set "Words" ($db.Words.Append (lower .Cmd))}}
	{{$s:=dbIncr .User.ID "LastLetter" (len .Cmd)}}
	{{execCC .CCID nil 15 .Message}} {{/* Call to check for message deletion/edit */}}
{{else}} {{/* Incorrect word */}}
	{{try}}{{addReactions $IncorrectEmoji}}{{catch}}{{end}}
	{{if $SubWrong}}{{$s:=dbIncr .User.ID "LastLetter" (mult (len .Cmd) -1)}}{{end}}
	{{template "SE" (sdict "Desc" (printf "`%s` doesn't start with the letter `%s`" .Cmd (slice $db.Word (len $db.Word|add -1))))}}
{{end}}
{{dbSet 0 "LastLetter" $db}}

{{define "SE"}}{{$ID:=sendMessageRetID nil (cembed "description" .Desc "color" 16711680)}}{{deleteTrigger 3}}{{deleteMessage nil $ID 15}}{{end}}
