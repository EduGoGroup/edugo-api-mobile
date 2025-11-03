---
description: Archive current sprint and prepare project for new cycle
allowed-tools: Bash, Read, Write
argument-hint: ""
---

# Command: Archive Sprint

## Description
This command archives the current sprint and prepares the project for a new cycle. It renames the sprint/current folder with a timestamp, moves it to the sprint/archived folder, and creates a new clean sprint/current folder.

## Command Responsibilities
1. **Validate** that content exists to archive
2. **Generate** unique name based on date/time
3. **Move** the current sprint/current folder to sprint/archived
4. **Create** new clean sprint/current folder with initial structure
5. **Report** the process to the user

## File Naming Format
```
sprint-YYYY-MM-DD-HHMM
Example: sprint-2025-10-31-1430
```

## Execution Instructions

Please execute the following steps:

### Step 1: Validate sprint/current folder
Verify that the `sprint/current/` folder exists and contains files. If it's empty or doesn't exist:
```
ℹ️ Warning: The sprint/current folder is empty or doesn't exist

There's nothing to archive. Do you want to create a new sprint/current folder anyway?
```

### Step 2: Generate archive name
Generate the name using current date/time format:
```bash
# Format: sprint-YYYY-MM-DD-HHMM
# Example: sprint-2025-10-31-1430

# Use date command to generate timestamp
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")
ARCHIVE_NAME="sprint-${TIMESTAMP}"
```

### Step 3: Verify name doesn't exist
Verify that a folder with that name doesn't already exist in `sprint/archived/`:
```bash
# If it exists, add a numeric suffix
# sprint-2025-10-31-1430-1
# sprint-2025-10-31-1430-2
# etc.
```

### Step 4: Move sprint/current folder to archive
```bash
# Ensure sprint/archived/ exists
mkdir -p sprint/archived

# Move sprint/current to sprint/archived with new name
mv sprint/current "sprint/archived/${ARCHIVE_NAME}"
```

### Step 5: Create new clean sprint/current folder
```bash
# Create clean structure
mkdir -p sprint/current/analysis
mkdir -p sprint/current/planning
mkdir -p sprint/current/execution
mkdir -p sprint/current/review
```

### Step 6: Confirmation message
```
✅ Sprint archived successfully

📦 Archive created:
- sprint/archived/sprint-2025-10-31-1430/

📁 New sprint/current folder created with clean structure:
- sprint/current/
  ├─ analysis/
  ├─ planning/
  ├─ execution/
  └─ review/

📌 Next step:
1. Create a new sprint/current/readme.md file with new sprint requirements
2. Execute /01-analysis to start the new cycle

💡 The previous sprint is available at:
   sprint/archived/sprint-2025-10-31-1430/
```

### Step 7: Template creation suggestion (optional)
Optionally, ask the user if they want to create an empty readme.md template:
```
Do you want to create a sprint/current/readme.md template file for the new sprint?

If the user accepts, create:
```markdown
# Sprint: [Sprint Name]

## Description
[Brief description of sprint objective]

## Requirements
- [ ] Requirement 1
- [ ] Requirement 2
- [ ] Requirement 3

## Context
[Additional relevant information]

## Expected Deliverables
1. [Deliverable 1]
2. [Deliverable 2]

## Restrictions/Considerations
- [Restriction 1]
- [Restriction 2]
```
```

## Important Notes
- This command has **full read/write permissions**
- Uses **timestamp** to avoid overwriting previous files
- Preserves **complete history** of sprints in the sprint/archived folder
- Guarantees there's always a clean sprint/current folder ready for a new cycle
- DO NOT archive if the user is in the middle of an important sprint (ask for confirmation if sprint/current has recent content)
