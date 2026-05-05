# Agent Definition Template

Use this template to create new agent definition files (typically stored in `.gemini/agents/`).

```markdown
---
name: [agent-name]
description: [One-sentence description of the agent's role and when to use it.]
kind: local
tools:
  - [tool_name_1]
  - [tool_name_2]
model: inherit
temperature: 0.1
max_turns: 15
---

# [Agent Name]

## Role
[Define the persona. E.g., "A senior software architect specializing in distributed systems."]

## Mandate
[Clear, concise goals for the agent. Use bullet points.]
- Goal 1
- Goal 2

## Guidelines
[Specific constraints, style requirements, or operational rules.]
- Adhere to [standard]
- Prefer [pattern] over [anti-pattern]

## Tool Usage
[Instructions for specific tools or when to delegate.]
- Use `grep_search` to [action]
- Delegate [task] to [sub-agent]
```
