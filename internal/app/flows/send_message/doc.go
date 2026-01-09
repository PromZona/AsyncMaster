// ## Send Message Flow (FSM Diagram)
// This flow guides the user through selecting a recipient, writing a message body, and optionally adding a title.
// (State is stored in s.UserState)
//
//                        FlowStart (0)
//                            |
//                (Callback: "send" / handleInitialSend)
//                            v
//                        AwaitResipient (1)
//                            |
//                (Callback: "player_names" / handlePlayerName)
//                            v
//                        AwaitMessage (2)
//                            |
//                (Text Input: Message Body / handleMessageText)
//                            v
//                      AwaitTitleDecision (3)
//                       /           \
//                      /             \
//       (Callback: "yes_title")  (Callback: "no_title" / handleNoTitle)
//             /                        \
//            v                          v
//       AwaitTitle (4) ----------------> FINALIZE (finilize function)
//            |
//  (Text Input: Title / handleMessageTitle)
//            v
//         FINALIZE (finilize function)
//             |
//             v
//      [End: Session Done]

package sendmessage
