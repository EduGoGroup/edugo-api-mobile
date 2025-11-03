---
description: Analyze sprint requirements and generate architectural documentation
allowed-tools: Read, Task
argument-hint: ""
---

# Command: Sprint Analysis

## Description
This command initiates the sprint analysis process. It reads the `sprint/current/readme.md` file containing sprint requirements and passes it to the specialized analysis agent to generate detailed architectural documentation.

## Command Responsibilities
1. **Read** the `sprint/current/readme.md` file
2. **Validate** that the file exists and contains information
3. **Invoke** the `analysis` agent passing the content as context
4. **Prevent** the agent from reading files directly (strict isolation)

## Execution Instructions

Please execute the following steps:

### Step 1: Validate input file
Verify that the `sprint/current/readme.md` file exists. If it doesn't exist, inform the user that they must create the file with sprint requirements.

### Step 2: Read content
Read the complete `sprint/current/readme.md` file and keep it in context.

### Step 3: Invoke analysis agent
Use the Task tool with `subagent_type: "general-purpose"` to invoke the analysis agent.

Pass to the agent:
- **Complete prompt**: The content of readme.md preceded by the agent's instructions (read the `.claude/agents/analysis.md` file)
- **Readme context**: All content read in step 2
- **Explicit restriction**: The agent MUST NOT read files, only work with the provided context

### Step 4: Confirmation message
Once the agent completes its work, inform the user:
```
✅ Analysis completed successfully

📁 Files generated in sprint/current/analysis/:
- architecture.md (Architecture diagram with validated Mermaid)
- data-model.md (ER diagram and database structure)
- process-diagram.md (System process flows)
- readme.md (Executive summary for planning)

📌 Next step: Execute /02-planning to generate the task plan
```

## Important Notes
- This command acts as an **interaction layer** between the user and the agent
- The analysis agent is **completely isolated** - it only receives what this command passes to it
- This guarantees total control over what information the agent processes
