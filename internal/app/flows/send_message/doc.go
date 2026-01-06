// State flow graph for "send message" scenario
// UserStateDefault
//     |
//     | (callback: send)
//     v
// UserStateAwaitResipient
//     |
//     | (callback: player_names)
//     v
// UserStateAwaitMessage
//     |
//     | (text)
//     v
// UserStateAwaitTitleDecision
//     |              |
//     | (yes_title)  | (no_title)
//     v              v
// UserStateAwaitTitle   finilize()
//     |
//     | (text)
//     v
// finilize()
//     |
//     v
// UserStateDefault

package sendmessage
