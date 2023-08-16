{{/* Reaction: Added + Removed */}}

{{$errCID:=1036894710343671918}}
{{$Emoji:="🎉"}}

{{if and (eq .Reaction.Emoji.Name $Emoji) (in (dbGet 0 "GiveawayMsgIDs").Value .Message.ID)}}
	{{if ($G:=(dbGet ($GID:=toInt (reFind `\d+` (index .Message.Embeds 0).Footer.Text)) "Giveaways").Value)}}
		{{if .ReactionAdded}}
			{{$G.Set "Users" ($G.Users.Append .User.ID)}}{{dbSet $GID "Giveaways" $G}}
			{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (print "Your entry into the giveaway for **" $G.Prize "** has been confirmed!") "color" 65280)}}
		{{else}}
			{{$Users:=cslice}}{{range $G.Users}}{{- if ne . $.User.ID}}{{- $Users =$Users.Append .}}{{- end -}}{{end}}
			{{$G.Set "Users" $Users}}{{dbSet $GID "Giveaways" $G}}
			{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (print "Your entry into the giveaway for **" $G.Prize "** has been removed") "color" 16711680)}}
		{{end}}
	{{else}}
		{{sendDM (cembed "title" (print "Giveaway Hosted in " .Guild.Name) "description" (print "Your entry into the giveaway for **" $G.Prize "** could not be confirmed\nPlease re-react to [the giveaway](" .Message.Link ")\n*Staff have been notified of the failed entry*") "color" 16711680)}}
		{{sendMessage $errCID (printf "%s (%d)'s attempt to join the giveaway failed" .User.Username .User.ID)}}
	{{end}}
{{end}}
