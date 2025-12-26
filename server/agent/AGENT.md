Agent: ConversationAgent

Responsability: 
- Maintain a conversation with a user using an LLM and perserve context with durable storage

Status: 
- idle
- thinking 
- responding

State:
- status
- history

Events:
- UserMessage {text}
- LLMResponse {text}
- Timeout

Actions:
- CallLLM {prompt}
- SendToClient {text} 
- PersistState

Rule:
On UserMessage in IDLE:
  → append message to history
  → CallLLM
  → status = THINKING

On LLMResponde in THINKING:
  → append response to history
  → SendToClient
  → PersistState
  → status = RESPONDING 

On TimeOut in RESPONDING:
 → status = IDLE 

On TimeOut in THINKING: 
 → status =  IDLE 