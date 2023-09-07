{{/* Reaction: Added + Removed */}}

{{if and (eq .Reaction.Emoji.Name "👊" "🧍‍♂️" "⏬") .ReactionAdded (eq (or ($db:=(dbGet .Message.ID "EcoBJ").Value).Opp 0) .User.ID)}}
	{{$Sets:=cslice "❤️" "♦️" "♠️" "♣️"}}{{$Cards:=cslice "2" "3" "4" "5" "6" "7" "8" "9" "10" "J" "Q" "K" "A"}}
	{{$Values:=sdict "J" 10 "Q" 10 "K" 10 "A" 11}}
	{{$e:=structToSdict (index .Message.Embeds 0)}}{{$e.Set "description" ""}}
	{{$Amt:=0}}{{$end:=false}}
	{{if eq .Reaction.Emoji.Name "⏬"}}
		{{if gt (mult $db.Amt 2) ($W:=toInt (dbGet .User.ID "Wallet").Value)}}
			{{$ID:=sendMessageRetID nil (cembed "title" "Insufficient Funds" "description" (print "Cannot bet more than you have!\nPlease withdraw " (sub (mult $db.Amt 2) $W) " coins and re-react") "color" 16711680)}}
			{{deleteMessage nil $ID 10}}
			{{deleteMessageReaction nil .Message.ID .User.ID .Reaction.Emoji.Name}}
			{{return}}
		{{else}}
			{{$db.Set "Amt" (mult $db.Amt 2)}}
			{{$e.Set "description" (print .User.Globalname " doubled down taking their final card.")}}
		{{end}}
	{{end}}
	{{if eq .Reaction.Emoji.Name "👊" "⏬"}}
		{{$db.U.Set "C" ($db.U.C.Append (printf "`%s %s`" (index $Sets (randInt 4)) ($V:=index $Cards (randInt 13))))}}
		{{$db.U.Set "V" (add $db.U.V (or ($Values.Get $V) (toInt $V)))}}
		{{if gt $db.U.V 21}}
			{{$V:=0}}{{$Ace:=0}}
			{{range $db.U.C}}
				{{- $V =add $V (or ($Values.Get ($c:=reFind `\d+|J|Q|K|A` .)) (toInt $c))}}
				{{- if eq $c "A"}}{{- $Ace =add $Ace 1}}{{- end -}}
			{{end}}
			{{$x:=$Ace}}
			{{while $x}}
				{{- if gt $V 21}}{{- $V =sub $V 10}}{{- end}}
				{{- $x =sub $x 1 -}}
			{{end}}
			{{$db.U.Set "V" $V}}
			{{if gt $V 21}}
				{{$e.Set "description" (printf "%s\n%s busted, losing **%s coins**" $e.description .User.Globalname ($db.Amt|humanizeThousands))}}
				{{$end =true}}
			{{end}}
		{{else if eq $db.U.V 21}}
			{{$Amt =mult $db.Amt 4}}
			{{$e.Set "description" (printf "%s\n%s has Blackjack, gaining **%s coins!**" $e.description .User.Globalname (mult $db.Amt 4|humanizeThousands))}}
			{{$end =true}}
		{{end}}
	{{end}}
	{{if or (eq .Reaction.Emoji.Name "🧍‍♂️") (and (eq .Reaction.Emoji.Name "⏬") (lt $db.U.V 21))}}
		{{while lt $db.D.V 17}}
			{{- $db.D.Set "C" ($db.D.C.Append (printf "`%s %s`" (index $Sets (randInt 4)) ($V:=index $Cards (randInt 13))))}}
			{{- $db.D.Set "V" ($V =add $db.D.V (or ($Values.Get $V) (toInt $V))) -}}
		{{end}}
		{{if or (gt $db.U.V $db.D.V) (gt $db.D.V 21)}}
			{{$Amt =mult $db.Amt 2}}
			{{$e.Set "description" (printf "%s\n%s won, gaining **%s coins**!" $e.description .User.Globalname (mult $db.Amt 2|humanizeThousands))}}
		{{else if eq $db.U.V $db.D.V}}
			{{$Amt =$db.Amt}}
			{{$e.Set "description" (printf "%s\n%s drew with the dealer\nNo gains nor losses." $e.description .User.Globalname)}}
		{{else}}
			{{$e.Set "description" (printf "%s\n%s lost their bet of **%s coins**." $e.description .User.Globalname ($db.Amt|humanizeThousands))}}
		{{end}}
		{{$end =true}}
	{{end}}
	{{dbSet .Message.ID "EcoBJ" $db}}
	{{deleteMessageReaction nil .Message.ID .BotUser.ID "⏬"}}
	{{deleteMessageReaction nil .Message.ID .User.ID .Reaction.Emoji.Name}}
	{{$e.Set "fields" (cslice 
		(sdict "name" "Dealer's Hand" "value" (printf "Cards: %s\nValue: `%d`" (joinStr " " $db.D.C) $db.D.V) "inline" true) 
		(sdict "name" (print .User.Globalname "'s Hand") "value" (printf "Cards: %s\nValue: `%d`" (joinStr " " $db.U.C) $db.U.V) "inline" true)
	)}}
	{{editMessage nil .Message.ID (cembed $e)}}
	{{if $end}}
		{{$s:=dbIncr .User.ID "Wallet" $Amt}}
		{{deleteAllMessageReactions nil .Message.ID}}
		{{dbDel .Message.ID "EcoBJ"}}
	{{end}}
{{end}}
