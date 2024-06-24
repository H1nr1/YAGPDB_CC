{{/* Reaction: Added + Removed */}}

{{ $gEmoji := "🎉" }} {{/* emoji to use as reaction */}}
{{ $errCID := .Channel.ID }} {{/* ID of channel to send errors */}}

{{if and (eq .Reaction.Emoji.APIName $gEmoji) .Message.Embeds}}
	{{if in (index .Message.Embeds 0).Title "Giveaway"}}
		{{$e:=sdict 
			"title" (printf "Giveaway Hosted in %s" .Guild.Name) 
			"color" 16711680
		}}
		{{if $g:=(dbGet ($gID:=(index .Message.Embeds 0).Footer.Text|reFind `\d+`|toInt) "giveaways").Value}}
			{{if .ReactionAdded}}
				{{$g.Set "users" ($g.users.Append .User.ID)}}
				{{$e.Set "description" (printf "Your entry into the giveaway for **%s** has been confirmed!" $g.prize)}}
				{{$e.Set "color" 65280}}
			{{else}}
				{{if not (in $g.users .User.ID)}}{{return}}{{end}}
				{{$users:=cslice}}
				{{range $g.users}}
					{{- if ne . $.User.ID}}
						{{- $users =$users.Append .}}
					{{- end -}}
				{{end}}
				{{$g.Set "users" $users}}
				{{$e.Set "description" (printf "Your entry into the giveaway for **%s** has been removed" 
					$g.prize
				)}}
			{{end}}
			{{dbSet $gID "giveaways" $g}}
		{{else}}
			{{$e.Set "description" (printf "Your entry into the giveaway for **%s** could not be confirmed\nPlease re-react to [the giveaway](<%s>)\n*Staff have been notified of the failed entry*" 
				$g.prize .Message.Link
			)}}
			{{sendMessage $errCID (printf "%s (%d)'s attempt to join giveaway `%d` failed" 
				.User.Username .User.ID $gID
			)}}
		{{end}}
		{{sendDM (cembed $e)}}
	{{end}}
{{end}}
