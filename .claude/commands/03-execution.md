---
description: Execute planned sprint tasks (all, specific phase, or task)
allowed-tools: Read, Task
argument-hint: "[phase-N|task-N.M]"
---

# Command: Sprint Execution

## Description
This command executes the planned sprint tasks. It can execute the entire plan or specific phases/tasks according to provided arguments. It reads the execution plan and optionally the project rules, then invokes the execution agent.

## Syntax
```bash
/03-execution              # Execute entire plan
/03-execution phase-1      # Execute only phase 1
/03-execution task-2.3     # Execute only task 3 of phase 2
```

## Command Responsibilities
1. **Read** the `sprint/current/planning/readme.md` file
2. **Filter** content according to arguments (if provided)
3. **Read** the `sprint/current/execution/rules.md` file (if it exists)
4. **Invoke** the `execution` agent with tasks and rules
5. **Allow limited access** for the agent to analysis/planning folders for additional context

## Execution Instructions

Please execute the following steps:

### Step 1: Validate input file
Verify that the `sprint/current/planning/readme.md` file exists. If it doesn't exist:
```
❌ Error: Sprint plan not found

Please execute first: /02-planning
```

### Step 2: Read work plan
Read the complete `sprint/current/planning/readme.md` file.

### Step 3: Process arguments (if any)
If the user provided arguments (e.g., `phase-1`, `task-2.3`):
- Extract from the plan only the section corresponding to that phase/task
- Verify the dependencies of that phase/task
- If there are uncompleted dependencies, warn the user but allow continuation

If NO arguments:
- Use the complete plan

### Step 4: Verify project rules
Verify if the `sprint/current/execution/rules.md` file exists:
```bash
If exists → Read it and pass to the agent
If NOT exists → Continue without rules (agent will use best practices)
```

### Step 5: Invoke execution agent
Use the Task tool with `subagent_type: "general-purpose"` to invoke the execution agent.

Pass to the agent:
- **Complete prompt**: The agent's instructions (read `.claude/agents/execution.md`)
- **Tasks to execute**: Complete or filtered plan according to step 3
- **Project rules**: Content of rules.md (if it exists)
- **Special permissions**:
  - Can read files from `sprint/current/analysis/` and `sprint/current/planning/` for additional context
  - Can write/modify files in the project root folder
  - CANNOT touch the `.claude/` folder
  - CANNOT touch the `sprint/` folder except to write reports in `sprint/current/execution/`

### Step 6: Confirmation message
Once the agent completes its work, inform the user:
```
✅ Execution completed successfully

📁 Report generated:
- sprint/current/execution/[phase-step]-[timestamp].md

✅ Validations performed:
- Code compiled correctly
- Tests executed (if applicable)

📌 Next step:
- Execute /04-review to see consolidated sprint status
- Or execute /03-execution [another-phase] to continue with other tasks
```

If there were compilation errors or failed tests:
```
⚠️ Execution completed with warnings

📁 Report generated:
- sprint/current/execution/[phase-step]-[timestamp].md

⚠️ Problems detected:
[List problems]

📌 Recommendation:
Review the execution report and fix problems before continuing
```

## Important Notes
- This command allows **modular execution** - you can execute specific phases/tasks
- The agent **validates that code compiles** before marking the task as complete
- **Project rules** are optional but recommended for consistency
- Each execution generates a **separate report** with timestamp for traceability
