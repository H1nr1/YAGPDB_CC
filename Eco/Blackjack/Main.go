{{/* Regex: `\A\$b(lack)?j` */}}

{{with .ExecData}}
  {{dbDel . "EcoBJ"}}
  {{editMessage nil . "*Blackjack has Expired*"}}
  {{deleteAllMessageReactions nil .}}
  {{return}}
{{end}}

{{$Args:=joinStr " " .CmdArgs}}

{{if ($Amt:=lower (reFind `(?i)all|half|\d+` $Args))}}
  {{$Amt =toInt (execTemplate "Amt" (sdict "Amt" $Amt "V" ($W:=toInt (dbGet .User.ID "Wallet").Value)))}}
  {{if gt $Amt $W}}{{template "Err" (sdict "T" "Fund" "A" "bet")}}
  {{else if lt $Amt 50}}{{template "Err" (sdict "T" "Min" "A" 50)}}
  {{else}}
    {{$:=dbIncr .User.ID "Wallet" (mult $Amt -1)}}
    {{$Sets:=cslice "❤️" "♦️" "♠️" "♣️"}}
    {{$Cards:=cslice "2" "3" "4" "5" "6" "7" "8" "9" "10" "J" "Q" "K" "A"}}
    {{$Values:=sdict "J" 10 "Q" 10 "K" 10 "A" 11}}
    {{$DC:=printf "`%s %s`" (index $Sets (randInt 4)) ($DV:=index $Cards (randInt 13))}}
    {{$UC1:=printf "`%s %s`" (index $Sets (randInt 4)) ($UV1:=index $Cards (randInt 13))}}
    {{$UC2:=printf "`%s %s`" (index $Sets (randInt 4)) ($UV2:=index $Cards (randInt 13))}}
    {{$UV:=add (or ($Values.Get $UV1) (toInt $UV1)) (or ($Values.Get $UV2) (toInt $UV2))}}
    {{if and (eq "A" $UV1 $UV2) (gt $UV 21)}}{{$UV =sub $UV 10}}{{end}}
    {{$e:=sdict 
      "title" "Blackjack" 
      "thumbnail" (sdict "url" "https://media.tenor.com/HjMiuoQ6KmoAAAAC/kumarhane-poker.gif") 
      "fields" (cslice 
        (sdict "name" "Dealer's Hand" "value" (printf "Cards: %s\nValue: `%d`" $DC ($DV =or ($Values.Get $DV) (toInt $DV))) "inline" true) 
        (sdict "name" (print .User.Globalname "'s Hand") "value" (printf "Cards: %s %s\nValue: `%d`" $UC1 $UC2 $UV) "inline" true)
      ) 
      "footer" (sdict "text" "This game will expire after 5 minutes") 
      "color" (randInt 0xFFFFFF)
    }}
    {{$end:=false}}
    {{if eq $UV 21}}
      {{$:=dbIncr .User.ID "Wallet" (mult $Amt 4)}}
      {{$e.Set "description" (printf "**%s has Blackjack and won %d coins!**" .User.Globalname (mult $Amt 4))}}{{$e.Del "footer"}}
      {{$end =true}}
    {{end}}
    {{$ID:=sendMessageRetID nil (cembed $e)}}
    {{if not $end}}
      {{addMessageReactions nil $ID "👊" "🧍‍♂️" "⏬"}}
      {{dbSet $ID "EcoBJ" (sdict "Opp" .User.ID "Amt" $Amt "D" (sdict "C" (cslice $DC) "V" $DV) "U" (sdict "C" (cslice $UC1 $UC2) "V" $UV))}}
      {{execCC .CCID nil 300 $ID}}
    {{end}}
  {{end}}
{{else}}{{template "Err" (sdict "T" "$Blackjack <Amount>")}}{{end}}

{{define "Amt"}}{{if reFind `\d+` .Amt}}{{return .Amt}}{{else if eq (lower .Amt) "all"}}{{return .V}}{{else if eq (lower .Amt) "half"}}{{return div .V 2}}{{end}}{{end}}

{{define "Err"}}{{with (sdict "$" (sdict "t" "Invalid Syntax" "d" (print "Syntax is `" .T "`")) "U" (sdict "t" "Invalid User" "d" "User is not a member of the server") "Fund" (sdict "t" "Insufficient Funds" "d" (print "Cannot " .A " more than you have!")) "Min" (sdict "t" "Insufficient Bet" "d" (print "Must bet " .A " or more coins!"))).Get (reFind `\$|U|Fund|Min` .T)}}{{sendMessage nil (cembed "title" .t "description" .d "color" 16711680)}}{{end}}{{end}}
