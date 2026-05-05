# Agent Design Principles

## 1. Focused Scope
Agents are most effective when they have a narrow, well-defined scope. Avoid "God agents" that do everything.

## 2. Descriptive Metadata
The `description` in the frontmatter is critical. It's how the orchestrator decides which agent to invoke.
- **Good**: "Specialized agent for analyzing and refactoring database schemas."
- **Bad**: "Agent for database tasks."

## 3. Configuration & Optimization
Fine-tune the agent's behavior using frontmatter fields:
- **`tools`**: Only provide the tools necessary for the agent's mandate. Minimizing available tools reduces the chance of hallucinations and improves performance.
- **`model`**: Use `inherit` by default. Specify a specific model (e.g., `gemini-2.0-flash-exp`) only if the task requires higher reasoning or a specific capability not met by the base model.
- **`temperature`**: Use lower values (e.g., `0.1`) for precise, deterministic tasks (coding, refactoring). Use higher values (e.g., `0.7`) for creative tasks (naming, brainstorming).
- **`max_turns`**: Limit the conversation length to prevent infinite loops or excessive context usage. A typical range is `10` to `20` turns.

## 4. Clear Mandates
Use the imperative mood for mandates. Tell the agent exactly what it is responsible for.
- "Maintain the integrity of the auth service."
- "Ensure all unit tests follow the table-driven pattern."

## 5. Operational Guidelines
Provide specific rules that reflect the project's standards. This reduces the need for repetitive prompting.

## 6. Tool Stewardship
Define which tools the agent should prefer and how it should use them. If an agent should never use a certain tool (e.g., `run_shell_command` for a read-only researcher), state it explicitly in the guidelines or exclude it from the `tools` list.
