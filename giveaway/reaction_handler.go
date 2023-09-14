{{/* Reaction: Added + Removed */}}

{{$Emoji:="🎉"}} {{/* Emoji to use as reaction */}}
{{$errCID:=.Channel.ID}} {{/* ID of channel to send errors */}}

{{if and (eq .Reaction.Emoji.Name $Emoji) .Message.Embeds}}
	{{if ($G:=(dbGet ($GID:=toInt (reFind `\d+` (index .Message.Embeds 0).Footer.Text)) "Giveaways").Value)}}
		{{if .ReactionAdded}}
			{{$G.Set "Users" ($G.Users.Append .User.ID)}}{{dbSet $GID "Giveaways" $G}}
			{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (printf "Your entry into the giveaway for **%s** has been confirmed!" $G.Prize) "color" 65280)}}
		{{else}}
			{{$Users:=cslice}}{{range $G.Users}}{{- if ne . $.User.ID}}{{- $Users =$Users.Append .}}{{- end -}}{{end}}
			{{$G.Set "Users" $Users}}{{dbSet $GID "Giveaways" $G}}
			{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (printf "Your entry into the giveaway for **%s** has been removed" $G.Prize) "color" 16711680)}}
		{{end}}
	{{else}}
		{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (printf "Your entry into the giveaway for **%s** could not be confirmed\nPlease re-react to [the giveaway](<%s>)\n*Staff have been notified of the failed entry*" $G.Prize .Message.Link) "color" 16711680)}}
		{{sendMessage $errCID (printf "%s (%d)'s attempt to join giveaway `%d` failed" .User.Username .User.ID $GID)}}
	{{end}}
{{end}}
