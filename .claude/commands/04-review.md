---
description: Consolidate sprint status and generate validation guide
allowed-tools: Read, Task
argument-hint: ""
---

# Command: Sprint Review

## Description
This command consolidates the complete sprint status. It reads the original plan and all execution reports, then invokes the review agent to generate a document showing which tasks were completed and provides a validation guide for the user.

## Command Responsibilities
1. **Read** the `sprint/current/planning/readme.md` file (original plan)
2. **Read** all report files in `sprint/current/execution/*.md`
3. **Invoke** the `review` agent with all information
4. **Generate** a consolidated document with status and validation guide

## Execution Instructions

Please execute the following steps:

### Step 1: Validate input files
Verify that the `sprint/current/planning/readme.md` file exists. If it doesn't exist:
```
❌ Error: Sprint plan not found

Please execute first: /02-planning
```

Verify that files exist in `sprint/current/execution/`. If there are none:
```
ℹ️ Warning: No execution reports found

The sprint has no tasks executed yet.
Do you want to generate a status report anyway? (useful to see what's missing)
```

### Step 2: Read work plan
Read the complete `sprint/current/planning/readme.md` file.

### Step 3: Read all execution reports
List and read all files in `sprint/current/execution/*.md` (except rules.md if it exists).

Organize the reports chronologically to give them to the agent in order.

### Step 4: Invoke review agent
Use the Task tool with `subagent_type: "general-purpose"` to invoke the review agent.

Pass to the agent:
- **Complete prompt**: The agent's instructions (read `.claude/agents/review.md`)
- **Original plan**: Content of `sprint/current/planning/readme.md`
- **Execution reports**: All files read in step 3, in chronological order
- **Special instruction**: The agent must generate a final "User Validation Guide" section

### Step 5: Confirmation message
Once the agent completes its work, inform the user:
```
✅ Review completed successfully

📁 File generated:
- sprint/current/review/readme.md

📊 Report content:
- Original plan with tasks marked as completed ✅
- Summary of pending tasks
- Validation guide to test the sprint

📌 Next step:
- Read sprint/current/review/readme.md to see complete status
- Use the "Validation Guide" at the end of the document to test the application
- If everything is complete, execute /archive to archive this sprint
- If tasks are missing, execute /03-execution [phase] to continue
```

### Step 6: Show quick summary (optional)
Optionally, you can show a quick summary in the console:
```
📈 Sprint Summary:
├─ Total phases: X
├─ Completed phases: Y
├─ Total tasks: A
├─ Completed tasks: B
└─ Progress: ZZ%
```

## Important Notes
- This command gives **complete visibility** of sprint status
- The **validation guide** is crucial - it must be simple and practical for the user
- Allows making decisions about what to do next (continue, archive, or fix)
- Useful for presentations/demos showing work progress
