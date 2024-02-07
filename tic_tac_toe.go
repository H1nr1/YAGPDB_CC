{{/* Regex: `\A-T(ic)? ?T(ac)? ?T(oe)?` */}}

{{with .ExecData}}
	{{if eq (dbGet .ID "ticTacToe").Value.runC .runC}}
		{{editMessage nil .msg "This game has expired."}}
		{{dbDel .ID "ticTacToe"}}
	{{end}}
	{{return}}
{{end}}

{{$args:=joinStr " " .CmdArgs}}

{{if inFold $args "start"}}
	{{$matchID:=randInt 1 100}}
	{{$content:=printf "%s wants to play Tic-Tac-Toe!\nAnyone can join with `-TicTacToe Join %d`" .User.Username $matchID}}
	{{$e:=sdict "title" "Tic-Tac-Toe" "description" "```   1   2   3\na    |   |   \n  -----------\nb    |   |   \n  -----------\nc    |   |   ```" "footer" (sdict "text" (print "ID: " $matchID)) "color" (randInt 0xFFFFFF)}}
	{{$u:=0}}{{with reFind `\d{17,19}` $args}}{{$u =userArg .}}{{end}}
	{{if $u}}
		{{$e.footer.Set "text" (printf "ID: %d; X: %s; O: %s" $matchID .User.Username $u.Username)}}
		{{$content =printf "%s, %s wants a match of tic-tac-toe with you!\nJoin with `-TicTacToe Join %d`" $u.Mention .User.Username $matchID}}
	{{end}}
	{{$msg:=sendMessageRetID nil (complexMessage "content" $content "embed" (cembed $e))}}
	{{dbSetExpire $matchID "ticTacToe" ($x:=sdict "msg" $msg "ID" $matchID "u1" (sdict "nil" .User "char" "X" "pos" cslice) "u2" (sdict "nil" $u "char" "O" "pos" cslice) "runC" 1) 180}}
	{{execCC .CCID nil 175 $x}}

{{else if inFold $args "join"}}
	{{$matchID:=reFind `\d{1,2}` $args|toInt}}
	{{with (dbGet $matchID "ticTacToe").Value}}
		{{if eq $.User.ID .u1.nil.ID}}Cannot join your own game{{return}}{{end}}
		{{printf "%s has joined %s's Tic-Tac-Toe game!\n%s may now add a position using `-TicTacToe <coordinate>`" $.User.Username .u1.nil.Username .u1.nil.Username}}
		{{.u2.Set "nil" $.User}}{{.Set "runC" (add .runC 1)}}{{dbSetExpire $matchID "ticTacToe" . 180}}
		{{execCC $.CCID nil 175 .}}
	{{else}}There is not a current game of Tic-Tac-Toe with ID `{{$matchID}}`{{return}}{{end}}
	

{{else}}
	{{if ($pos:=reFind `(?i)[abc][123]|[123][abc]|\d` $args)}}
		{{$db:=0}}{{$u:=0}}
		{{range dbTopEntries "ticTacToe" 100 0}}
			{{if eq .Value.u1.nil.ID $.User.ID}}{{$u ="u1"}}
			{{else if eq .Value.u2.nil.ID $.User.ID}}{{$u ="u2"}}
			{{else}}{{return}}{{end}}
			{{$db =.Value}}
	        {{end}}
	        {{if not $db}}No current Tic-Tac-Toe game found for {{.User.Username}}{{return}}{{end}}
		{{$notu:=0}}{{$ustr:=$u}}{{if eq $u "u1"}}{{$notu ="u2"}}{{else}}{{$notu ="u1"}}{{end}}
		{{$u =$db.Get $u}}{{$notu =$db.Get $notu}}
		{{if gt (len $u.pos) (len $notu.pos)}}It is not your turn.{{return}}{{end}}
	        {{$int:=toInt (reFind `\d` $pos)}}
	        {{if inFold $pos "b"}}{{$pos =add $int 3}}
		{{else if inFold $pos "c"}}{{$pos =add $int 6}}
		{{else}}{{$pos =$int}}{{end}}
	        {{if in ($db.u1.pos.AppendSlice $db.u2.pos) $pos}}This position is already taken! Please try again{{return}}{{end}}
		{{$descPos:=(dict 1 19 2 23 3 27 4 47 5 51 6 55 7 75 8 79 9 83).Get $pos}}
		{{$content:=printf "%s has taken position %d" .User.Username $pos}}
		{{$e:=structToSdict (index (getMessage nil $db.msg).Embeds 0)}}
		{{$matchID:=toInt (reFind `\d{1,2}` $e.Footer.Text)}}
		{{$e.Set "description" (joinStr $u.char (slice $e.Description 0 $descPos) (slice $e.Description (add $descPos 1)))}}
		{{$u.Set "pos" ($u.pos.Append $pos)}}
		{{$end:=false}}
		{{with $u.pos}}
			{{if or (and (in . 1) (in . 2) (in . 3)) (and (in . 4) (in . 5) (in . 6)) (and (in . 7) (in . 8) (in . 9)) (and (in . 1) (in . 4) (in . 7)) (and (in . 2) (in . 5) (in . 8)) (and (in . 3) (in . 6) (in . 9)) (and (in . 1) (in . 5) (in . 9)) (and (in . 3) (in . 5) (in . 7))}}
				{{$content =print $.User.Username " Wins!"}}
				{{dbDel $matchID "ticTacToe"}}{{$end =true}}
			{{end}}
		{{end}}
		{{if and (not $end) (eq (add (len $db.u1.pos) (len $db.u2.pos)) 9)}}
			{{$content ="No Winners Today."}}
			{{dbDel $matchID "ticTacToe"}}{{$end =true}}
		{{end}}
		{{deleteMessage nil $db.msg 0}}
		{{$msg:=sendMessageRetID nil (complexMessage "content" $content "embed" (cembed $e))}}
		{{if not $end}}
			{{$db.Set "msg" $msg}}{{$db.Set $ustr $u}}{{$db.Set "runC" (add $db.runC 1)}}
			{{dbSetExpire $matchID "ticTacToe" $db 180}}
			{{execCC .CCID nil 175 $db}}
		{{end}}

	{{else}}{{$x:=sendMessageRetID nil (cembed "title" "Invalid Syntax" "description" "Correct Syntax:\n`-TicTacToe Start`\n`-TicTacToe Join <User: @/ID>`\n`-TicTacToe <position>`" "color" 16711680)}}{{deleteMessage nil $x 15}}{{end}}
{{end}}
