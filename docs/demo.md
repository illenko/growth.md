# growth.md Demo Script

This demo showcases all current functionality of growth.md CLI.

## Prerequisites

Build the binary first:
```bash
make build
```

## 1. Initialize a New Growth Repository

```bash
# Create a test directory
mkdir -p /tmp/growth-demo
cd /tmp/growth-demo

# Initialize growth repository
# This creates directory structure and config
./growth init .

# Initialize git repository for version control
git init
git config user.name "Demo User"
git config user.email "demo@example.com"
git add .
git commit -m "Initial growth.md setup"
```

## 2. Create Skills

```bash
# Create Python skill
./growth skill create "Python Programming" \
  --category "Programming Languages" \
  --level intermediate \
  --tags python,backend

# Create Go skill
./growth skill create "Go Programming" \
  --category "Programming Languages" \
  --level beginner \
  --tags go,backend,concurrency

# Create Machine Learning skill
./growth skill create "Machine Learning" \
  --category "AI/ML" \
  --level beginner \
  --tags ml,ai,data-science

# Create Docker skill
./growth skill create "Docker" \
  --category "DevOps" \
  --level intermediate \
  --tags docker,containers,devops

# List all skills
./growth skill list

# List skills by category
./growth skill list --category "Programming Languages"

# View a specific skill
./growth skill view skill-001
```

## 3. Create Learning Resources

```bash
# Add a Python resource
./growth resource create "Python Crash Course Book" \
  --skill-id skill-001 \
  --type book \
  --author "Eric Matthes" \
  --hours 40

# Add a Go course
./growth resource create "Learn Go with Tests" \
  --skill-id skill-002 \
  --type course \
  --url "https://quii.gitbook.io/learn-go-with-tests" \
  --hours 30

# Add an ML course
./growth resource create "Fast.ai Practical Deep Learning" \
  --skill-id skill-003 \
  --type course \
  --url "https://course.fast.ai" \
  --hours 50

# Start working on a resource
./growth resource start resource-001

# List all resources
./growth resource list

# List in-progress resources
./growth resource list --status in-progress

# View resource details
./growth resource view resource-001
```

## 4. Create Goals

```bash
# Create a goal to become a backend engineer
./growth goal create "Become Senior Backend Engineer" \
  --priority high \
  --target 2026-12-31

# Create an ML goal
./growth goal create "Build Production ML System" \
  --priority medium \
  --target 2026-06-30

# List all goals
./growth goal list

# View goal details
./growth goal view goal-001
```

## 5. Create Learning Paths

```bash
# Create a manual learning path
./growth path create "Backend Engineering Mastery" \
  --type manual

# Create another path
./growth path create "ML Engineering Path" \
  --type manual

# Link paths to goals
./growth goal add-path goal-001 path-001
./growth goal add-path goal-002 path-002

# List all paths
./growth path list

# View path details
./growth path view path-001
```

## 6. Create Milestones

```bash
# Create skill-level milestone
./growth milestone create "Master Python Basics" \
  --type skill-level \
  --ref-type skill \
  --ref-id skill-001

# Create goal-level milestone
./growth milestone create "Complete Backend Fundamentals" \
  --type goal-level \
  --ref-type goal \
  --ref-id goal-001

# Create path-level milestone
./growth milestone create "Finish Backend Path Phase 1" \
  --type path-level \
  --ref-type path \
  --ref-id path-001

# List all milestones
./growth milestone list

# Achieve a milestone
./growth milestone achieve milestone-001 \
  --proof "https://github.com/myuser/python-projects"

# View milestone
./growth milestone view milestone-001
```

## 7. Log Progress

```bash
# Log today's progress
./growth progress log \
  --hours 3 \
  --mood motivated \
  --skills skill-001,skill-002
# Enter summary when prompted, then press Ctrl+D or type '.' on a new line

# Log progress for a specific date
./growth progress log \
  --date 2025-12-20 \
  --hours 2 \
  --mood focused

# List all progress logs
./growth progress list

# View specific progress log
./growth progress view progress-001
```

## 8. Search Functionality

```bash
# Search across all entities
./growth search python

# Search for specific entity type
./growth search backend --type goal

# Search for ML-related items
./growth search "machine learning"
```

## 9. Overview & Statistics

```bash
# View dashboard overview
./growth overview

# View detailed statistics
./growth stats
```

## 10. Git Integration

```bash
# View git status (shows auto-committed changes)
./growth git status

# View commit log
./growth git log

# View more commits
./growth git log --count 20

# Check git log from terminal to see auto-commits
git log --oneline -10
```

## 11. Update Operations

```bash
# Update a skill's level
./growth skill edit skill-002 --level intermediate

# Mark a resource as complete
./growth resource complete resource-001

# Update a goal's priority
./growth goal edit goal-002 --priority high

# Check git log to see auto-commits for updates
./growth git log --count 5
```

## 12. Different Output Formats

```bash
# View skills as JSON
./growth skill list --format json

# View goals as YAML
./growth goal list --format yaml

# View resources as table (default)
./growth resource list --format table
```

## 13. Filtering and Querying

```bash
# Filter skills by level
./growth skill list --level intermediate

# Filter skills by category and level
./growth skill list --category "Programming Languages" --level beginner

# Filter resources by skill
./growth resource list --skill-id skill-001

# Filter goals by status
./growth goal list --status active

# Filter goals by priority
./growth goal list --priority high
```

## 14. Cleanup (Optional)

```bash
# Delete a milestone
./growth milestone delete milestone-003

# Delete a resource
./growth resource delete resource-003

# Check git log shows delete operations
./growth git log --count 3
```

## Summary

You've now seen all major features:
- ✅ Repository initialization with git
- ✅ CRUD operations for all entity types (Skills, Goals, Resources, Paths, Milestones, Progress)
- ✅ Entity relationships (linking goals with skills and paths)
- ✅ Resource lifecycle (start, complete)
- ✅ Milestone achievements with proof
- ✅ Progress logging with daily summaries
- ✅ Search across all entities
- ✅ Overview dashboard and statistics
- ✅ Git auto-commit integration
- ✅ Multiple output formats (table, JSON, YAML)
- ✅ Filtering and querying

All changes are automatically committed to git when auto-commit is enabled in config!
