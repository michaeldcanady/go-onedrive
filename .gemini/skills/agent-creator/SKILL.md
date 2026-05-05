---
name: agent-creator
description: Facilitates the creation and design of specialized sub-agents. Use when you need to define a new agent with proper frontmatter (tools, model, temperature, max_turns), role, mandate, and guidelines.
---

# Agent Creator

## Overview

This skill provides a standardized workflow for designing and implementing specialized agents within the Gemini CLI ecosystem. It ensures that agents have clear roles, precise mandates, optimized toolsets, and consistent configuration.

## Workflow

1. **Define the Scope**: Identify the specific domain or task the agent will handle.
2. **Draft the Frontmatter**: Write a clear `name` and a single-line `description` for triggering.
3. **Configure the Agent**: Select necessary `tools` and set `model`, `temperature`, and `max_turns` based on the task complexity.
4. **Design the Persona**: Use [references/design-principles.md](references/design-principles.md) to define a focused role.
5. **Use the Template**: Populate the [references/agent-definition-template.md](references/agent-definition-template.md) with mandates and guidelines.
6. **Verify Wiring**: Ensure the agent is correctly registered and reachable by the orchestrator.

## Quick Start

### 1. Identify the need
"I need an agent that specifically handles refactoring our gRPC proto files."

### 2. Use the template
Read [references/agent-definition-template.md](references/agent-definition-template.md) and copy its content to a new file in `.gemini/agents/grpc-expert.md`.

### 3. Apply design principles
Consult [references/design-principles.md](references/design-principles.md) to ensure the `description`, `tools`, and `mandate` are optimized for agentic orchestration.

## Advanced Usage

For complex multi-agent systems, ensure that mandates are mutually exclusive and collectively exhaustive (MECE) to avoid confusion during delegation.
