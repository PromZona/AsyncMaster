// ## Registration Flow (FSM Diagram)
// This flow is designed to only accept sequential Text input to collect registration data. It explicitly rejects all Callbacks.
// (State is stored in s.UserState)
//
//                         FlowStart (Implied / Initial State)
//                             |
//                  (Registration Middleware)
//                             v
//                       AwaitPassword (0)
//                             |
//                  (Text Input: Password / handlePassword)
//                             v
//                       AwaitCodename (1)
//                             |
//                  (Text Input: Codename / handlePlayerName)
//                             v
//                  [End: User Registered, Session Done]

package registration
