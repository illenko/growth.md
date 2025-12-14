# growth.md: Technical Whitepaper

**A Git-Native Career Development & Learning Path Manager for Software Engineers**

---

## Abstract

**growth.md** is an open-source, local-first career development platform designed specifically for software engineers. Built on Git and Markdown, it provides a structured yet flexible framework for tracking skills, defining career goals, generating AI-powered learning paths, and measuring professional growth over time. Unlike centralized learning platforms, growth.md prioritizes data ownership, privacy, and developer-friendly workflows through CLI-first tooling and seamless integration with AI assistants via the Model Context Protocol (MCP).

This whitepaper presents the technical architecture, entity model, and design principles that enable engineers to version-control their career progression with the same rigor applied to software development.

---

## Table of Contents

1. [Introduction](#1-introduction)
2. [Core Concepts & Entity Model](#2-core-concepts--entity-model)
3. [Entity Relationships & Hierarchy](#3-entity-relationships--hierarchy)
4. [Technical Architecture](#4-technical-architecture)
5. [File Structure & Format Specifications](#5-file-structure--format-specifications)
6. [CLI Interface Design](#6-cli-interface-design)
7. [AI Integration Strategy](#7-ai-integration-strategy)
8. [Development Roadmap](#8-development-roadmap)
9. [Competitive Analysis](#9-competitive-analysis)
10. [Open Source Strategy](#10-open-source-strategy)
11. [Conclusion](#11-conclusion)

---

## 1. Introduction

### 1.1 Problem Statement

Software engineers face several challenges in managing their career development:

- **Fragmented Learning**: Resources scattered across platforms (Udemy, Coursera, YouTube, blogs)
- **Lack of Personalization**: Generic learning paths that ignore individual background and goals
- **No Progress Visibility**: Difficulty tracking what has been learned and what remains
- **Platform Lock-in**: Career data trapped in proprietary systems
- **No Version Control**: Career progression treated as mutable state rather than versioned history

### 1.2 Solution Overview

growth.md addresses these challenges through:

1. **Git-Native Design**: Every skill, goal, and learning path is a versioned Markdown file
2. **Local-First Architecture**: Complete data ownership with offline capability
3. **AI-Powered Personalization**: Context-aware learning path generation via LLMs
4. **Developer Workflow Integration**: CLI-first tooling with MCP support for AI assistants
5. **Community-Driven**: Shareable learning paths via standard Git workflows

### 1.3 Design Principles

- **Plain Text Everything**: Human-readable Markdown files, no proprietary formats
- **Single Source of Truth**: Clear ownership hierarchy with no circular dependencies
- **Composability**: Reusable entities (Skills, Resources) across multiple contexts
- **Privacy by Default**: Local storage, optional sharing
- **YAGNI**: Start simple, extend when needed

---

## 2. Core Concepts & Entity Model

The growth.md system is built on six primary entities and two cross-cutting entities. Each entity serves a distinct purpose in modeling career development.

### 2.1 Primary Entities

#### 2.1.1 Goal (Top-Level Entity)

**Definition**: A high-level career objective representing what the engineer wants to achieve.

**Purpose**:
- Provides direction and motivation for learning activities
- Acts as the root of the dependency tree
- Links to Learning Paths that define strategies to achieve the goal

**Examples**:
- "Become Machine Learning Engineer"
- "Master Distributed Systems Architecture"
- "Transition to Engineering Management"

**Key Attributes**:
- Unique identifier (e.g., `goal-001`)
- Status: `active`, `completed`, `archived`
- Priority: `high`, `medium`, `low`
- Target completion date
- Success criteria (linked Milestones)
- References to Learning Paths (unidirectional)

**Ownership**: Goals own their associated Learning Paths.

---

#### 2.1.2 Learning Path (Strategy Entity)

**Definition**: A structured, sequenced plan for achieving a specific Goal, typically AI-generated but can be manually created.

**Purpose**:
- Defines the "how" for achieving a Goal
- Breaks down complex objectives into manageable Phases
- Provides personalized sequencing based on current skill level
- Can be shared and reused across multiple Goals

**Examples**:
- "ML Engineer Track (Fast.ai Approach)"
- "Backend Systems Mastery (Go-focused)"
- "Frontend Development (React Ecosystem)"

**Key Attributes**:
- Unique identifier (e.g., `path-001`)
- Type: `ai-generated` or `manual`
- Status: `active`, `archived`, `completed`
- Ordered list of Phases
- Generation metadata (date, AI model used, context)

**Important Design Decision**:
- Paths do NOT store which Goals reference them (unidirectional relationship)
- This enables true reusability: the same "Python Fundamentals" path can be referenced by "Become ML Engineer" and "Become Backend Developer" goals
- Reverse lookup ("which goals use this path?") can be computed when needed

**Ownership**: Learning Paths own their Phases.

---

#### 2.1.3 Phase (Sequential Step Entity)

**Definition**: A distinct stage within a Learning Path, representing a cohesive block of learning activity.

**Purpose**:
- Groups related Skills into logical progression steps
- Provides time-boxing and progress tracking boundaries
- Associates Milestones with specific learning stages
- Enables parallel skill development within a phase

**Examples**:
- "Phase 1: Foundations (2 months)"
- "Phase 2: Core Machine Learning (3 months)"
- "Phase 3: Production ML Systems (4 months)"

**Key Attributes**:
- Unique identifier (e.g., `phase-001`)
- Parent Learning Path ID
- Order/sequence number (1, 2, 3...)
- Estimated duration (optional)
- Required skills with target proficiency levels
- Associated Milestones (completion criteria)

**Relationship to Skills**:
- Many-to-Many: A Phase requires multiple Skills, and a Skill can appear in multiple Phases/Paths

**Ownership**: Phases are owned by Learning Paths.

---

#### 2.1.4 Skill (Atomic Competency Entity)

**Definition**: A single, measurable technical or professional competency that exists globally across the system.

**Purpose**:
- Represents the atomic unit of learning
- Maintains single source of truth for proficiency level
- Reusable across multiple Goals, Paths, and Phases
- Links to specific learning Resources

**Examples**:
- "Python Programming"
- "Docker & Containerization"
- "System Design"
- "GraphQL APIs"

**Key Attributes**:
- Unique identifier (e.g., `skill-002`)
- Category: `programming`, `math`, `ml`, `system-design`, `soft-skills`, etc.
- Proficiency level: `beginner`, `intermediate`, `advanced`, `expert`
- Status: `not-started`, `learning`, `mastered`
- Tags for searchability
- List of associated Resources

**Design Decision: Flat Structure (MVP)**:
- Skills are stored as flat, independent entities
- No parent/child hierarchy in MVP
- Future enhancement: Introduce skill trees (e.g., "Machine Learning" → "Supervised Learning" → "Linear Regression")
- Rationale: YAGNI principle, reduces complexity in MVP

**Global Nature**:
- Skills exist independently of Goals and Paths
- A skill's proficiency level is global (e.g., "Python: Intermediate" applies everywhere)
- Phases reference skills but don't own them

**Ownership**: Skills own their Resources.

---

#### 2.1.5 Resource (Learning Material Entity)

**Definition**: A specific learning material (book, course, article, video, project) that teaches a Skill.

**Purpose**:
- Provides actionable learning content
- Tracks completion status
- Enables resource recommendations and sharing
- Links theoretical skills to practical materials

**Examples**:
- "Book: Designing Data-Intensive Applications (Martin Kleppmann)"
- "Course: Fast.ai Practical Deep Learning for Coders"
- "Article: Stripe's API Design Philosophy"
- "Project: Build a Distributed Key-Value Store"

**Key Attributes**:
- Unique identifier (e.g., `resource-001`)
- Type: `book`, `course`, `video`, `article`, `project`, `documentation`
- Primary Skill ID (which skill this resource teaches)
- Status: `not-started`, `in-progress`, `completed`
- URL or reference
- Tags for discoverability
- Estimated time investment
- Notes/reviews

**Relationship to Skills**:
- Belongs to one primary Skill
- Can be tagged for discoverability across multiple skill areas
- Example: "System Design Interview Course" belongs to "System Design" skill but tags "Distributed Systems", "Databases", "Caching"

**Ownership**: Resources are owned by Skills.

---

### 2.2 Cross-Cutting Entities

#### 2.2.1 Milestone (Achievement Marker)

**Definition**: A significant, measurable achievement that marks progress toward a Goal, completion of a Path, or mastery of a Skill.

**Purpose**:
- Provides concrete success criteria
- Enables celebration of progress
- Acts as checkpoint for AI path re-evaluation
- Can be used for portfolio building

**Examples**:
- "First ML model deployed to production"
- "Completed AWS Solutions Architect certification"
- "Contributed to Kubernetes open-source project"
- "Led system design review for payment microservice"

**Key Attributes**:
- Unique identifier (e.g., `milestone-001`)
- Type: `goal-level`, `path-level`, `skill-level`
- Reference ID: points to a Goal, Path, or Skill
- Status: `pending`, `achieved`
- Achievement date (when marked complete)
- Evidence/proof (optional link to project, certification, etc.)

**Flexibility**:
- Milestones can reference any entity in the system
- A Goal has success criteria Milestones
- A Phase can have completion Milestones
- A Skill can have mastery Milestones

**Ownership**: Milestones reference entities but are stored independently for cross-cutting visibility.

---

#### 2.2.2 Progress Log (Time-Based Journal)

**Definition**: A chronological record of learning activities, accomplishments, and reflections over a specific time period.

**Purpose**:
- Track daily/weekly learning activities
- Provide retrospective visibility
- Enable velocity and consistency metrics
- Record context that may be lost over time

**Examples**:
- Weekly log: "Completed Fast.ai Lesson 3, deployed digit classifier"
- Daily log: "Studied backpropagation math, still confused on chain rule"
- Milestone log: "Passed TensorFlow certification exam!"

**Key Attributes**:
- Unique identifier (e.g., `progress-2024-w50`)
- Time period: week-of date or specific date
- Skills worked on (list of Skill IDs)
- Resources used (list of Resource IDs)
- Milestones achieved (list of Milestone IDs)
- Free-form notes (Markdown)
- Time invested (optional)

**Flexibility**:
- Can reference any number of Skills, Resources, Milestones
- Supports both structured data (IDs) and unstructured reflection (notes)
- Enables both quantitative analytics and qualitative insights

**Ownership**: Progress logs are independent, time-indexed entities.

---

## 3. Entity Relationships & Hierarchy

### 3.1 Ownership Hierarchy (Top-Down)

The system follows a strict parent-child ownership model to prevent circular dependencies:

```
Goal (root)
  ├── Learning Path (1:N - one goal, multiple paths)
  │     ├── Phase (1:N - one path, multiple ordered phases)
  │     │     └── Skill Reference (M:N - phase requires multiple skills)
  │     │           └── Resource (1:N - one skill, multiple resources)
```

**Key Principles**:
1. Each child entity belongs to exactly ONE parent
2. Parents can have MULTIPLE children
3. No circular references in the ownership tree

---

### 3.2 Reference Relationships

While ownership is unidirectional, some entities use references for discoverability:

| Relationship | Type | Direction | Implementation |
|-------------|------|-----------|----------------|
| Goal → Learning Path | Reference | Unidirectional | Goal stores list of Path IDs |
| Phase → Skill | Reference | Many-to-Many | Phase stores list of Skill IDs |
| Skill → Resource | Ownership | One-to-Many | Resource stores Skill ID |
| Milestone → Entity | Reference | Polymorphic | Milestone stores entity type + ID |
| Progress → Entity | Reference | Many-to-Many | Progress stores lists of entity IDs |

---

### 3.3 Example Relationship Flow

Consider the goal: **"Become Machine Learning Engineer"**

```
┌─────────────────────────────────────────────────────────────┐
│ Goal: Become ML Engineer (goal-001)                        │
│ - Status: active                                            │
│ - Target: 2025-12-31                                        │
│ - Learning Paths: [path-001, path-002]                     │
│ - Milestones: [milestone-001, milestone-002]                │
└────────────────────┬────────────────────────────────────────┘
                     │
      ┌──────────────┴──────────────┐
      │                             │
      ▼                             ▼
┌──────────────────┐       ┌──────────────────┐
│ Path 001:        │       │ Path 002:        │
│ Fast.ai Track    │       │ Academic Track   │
└────┬─────────────┘       └────┬─────────────┘
     │                          │
     │ Phases: [p1, p2, p3]     │ Phases: [p4, p5]
     ▼                          ▼
┌──────────────────┐       ┌──────────────────┐
│ Phase 1:         │       │ Phase 4:         │
│ Foundations      │       │ Theory Deep-Dive │
└────┬─────────────┘       └────┬─────────────┘
     │                          │
     │ Skills: [s1, s2, s3]     │ Skills: [s2, s6, s7]
     ▼                          ▼
┌──────────────────┐       ┌──────────────────┐
│ Skill 002:       │       │ Skill 007:       │
│ Python           │       │ Linear Algebra   │
│ Level: Inter.    │       │ Level: Beginner  │
└────┬─────────────┘       └────┬─────────────┘
     │                          │
     │ Resources: [r1, r2]      │ Resources: [r5]
     ▼                          ▼
┌──────────────────┐       ┌──────────────────┐
│ Resource 001:    │       │ Resource 005:    │
│ Book: Fluent     │       │ 3Blue1Brown      │
│ Python           │       │ Linear Algebra   │
└──────────────────┘       └──────────────────┘
```

**Cross-Cutting Entities**:

```
┌──────────────────────────────────────────────────────────┐
│ Milestone 001: Deploy First ML Model                    │
│ - Type: goal-level                                       │
│ - Reference: goal-001                                    │
│ - Status: pending                                        │
└──────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────┐
│ Progress Log: 2024-W50                                   │
│ - Skills worked: [skill-002, skill-007]                  │
│ - Resources used: [resource-001, resource-005]           │
│ - Milestones: []                                         │
│ - Notes: "Math is hard! Taking longer than expected."   │
└──────────────────────────────────────────────────────────┘
```

---

### 3.4 Data Integrity Rules

1. **Cascade Behavior**:
   - Deleting a Goal → prompts to archive or delete associated Paths
   - Deleting a Path → prompts to archive or delete associated Phases
   - Deleting a Skill → warns if referenced by active Phases

2. **Orphan Prevention**:
   - Phases cannot exist without a parent Path
   - Resources should warn if Skill is deleted

3. **Reference Validation**:
   - System validates that referenced IDs exist
   - Broken references flagged in CLI with repair suggestions

---

## 4. Technical Architecture

### 4.1 Technology Stack

| Component | Technology | Rationale |
|-----------|-----------|-----------|
| **Language** | Go | Excellent CLI tooling, static binaries, performance, cross-platform compilation |
| **CLI Framework** | Cobra | Industry standard, used by kubectl, GitHub CLI, Hugo |
| **Storage** | Git + Markdown | Version control, human-readable, no vendor lock-in |
| **TUI** | Bubble Tea | Modern, declarative terminal UI framework |
| **AI Integration** | OpenAI/Anthropic API + MCP | Flexible LLM support, AI assistant integration |
| **Config** | YAML | Human-friendly configuration format |
| **Testing** | Go standard library + testify | Comprehensive testing without external dependencies |

### 4.2 Why Go?

**Advantages**:
1. **Single Binary Distribution**: No runtime dependencies, easy installation
2. **Cross-Platform**: Compile for Linux, macOS, Windows from single codebase
3. **Performance**: Fast startup time, low memory footprint
4. **Excellent Stdlib**: File I/O, Git operations, HTTP clients built-in
5. **Strong CLI Ecosystem**: Cobra, Viper, Bubble Tea widely adopted
6. **Static Typing**: Catch errors early, better refactoring support
7. **Concurrency**: Goroutines for parallel AI requests, file processing

**Examples of Successful Go CLI Tools**:
- kubectl (Kubernetes CLI)
- gh (GitHub CLI)
- terraform
- docker
- hugo

### 4.3 System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    growth CLI (Cobra)                       │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────┐  │
│  │  skill   │  │  goal    │  │  path    │  │ progress  │  │
│  │ commands │  │ commands │  │ commands │  │ commands  │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └─────┬─────┘  │
└───────┼─────────────┼─────────────┼──────────────┼─────────┘
        │             │             │              │
        ▼             ▼             ▼              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Core Domain Layer                         │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Entity Models (Skill, Goal, Path, Phase, etc.)     │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Business Logic (validation, relationships)          │  │
│  └──────────────────────────────────────────────────────┘  │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
┌──────────────┐ ┌─────────────┐ ┌────────────┐
│   Storage    │ │  AI Engine  │ │  Git Ops   │
│   Layer      │ │             │ │            │
│              │ │             │ │            │
│ - File I/O   │ │ - OpenAI    │ │ - Commit   │
│ - Markdown   │ │ - Anthropic │ │ - Branch   │
│ - YAML       │ │ - MCP       │ │ - Status   │
│ - Indexing   │ │             │ │            │
└──────┬───────┘ └──────┬──────┘ └─────┬──────┘
       │                │              │
       ▼                ▼              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Filesystem (growth/)                     │
│  ├── skills/                                                │
│  ├── goals/                                                 │
│  ├── paths/                                                 │
│  ├── resources/                                             │
│  ├── milestones/                                            │
│  ├── progress/                                              │
│  ├── config.yml                                             │
│  └── .git/                                                  │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                Optional: MCP Server                         │
│  - Exposes growth.md to Claude/AI assistants               │
│  - Natural language interaction                             │
│  - Context-aware suggestions                                │
└─────────────────────────────────────────────────────────────┘
```

### 4.4 Project Structure

```
growth.md/
├── cmd/
│   └── growth/
│       └── main.go              # CLI entry point
├── internal/
│   ├── cli/                     # Cobra commands
│   │   ├── skill.go
│   │   ├── goal.go
│   │   ├── path.go
│   │   ├── resource.go
│   │   ├── milestone.go
│   │   ├── progress.go
│   │   └── board.go
│   ├── core/                    # Domain models
│   │   ├── skill.go
│   │   ├── goal.go
│   │   ├── path.go
│   │   ├── phase.go
│   │   ├── resource.go
│   │   ├── milestone.go
│   │   └── progress.go
│   ├── storage/                 # File operations
│   │   ├── markdown.go          # MD parsing/writing
│   │   ├── filesystem.go        # File I/O
│   │   └── index.go             # Entity indexing
│   ├── ai/                      # AI integration
│   │   ├── openai.go
│   │   ├── anthropic.go
│   │   ├── pathgen.go           # Path generation
│   │   └── mcp/                 # MCP server
│   ├── git/                     # Git operations
│   │   └── operations.go
│   └── tui/                     # Terminal UI
│       ├── board.go             # Overview dashboard
│       └── components/          # Reusable UI components
├── pkg/                         # Public libraries (if needed)
├── examples/                    # Example learning paths
│   ├── ml-engineer/
│   ├── backend-specialist/
│   └── frontend-developer/
├── docs/
│   ├── growth-md-concept.md     # This document
│   ├── getting-started.md
│   ├── cli-reference.md
│   └── entity-model.md
├── tests/
│   ├── integration/
│   └── fixtures/
├── go.mod
├── go.sum
├── Makefile
├── LICENSE                      # MIT
└── README.md
```

---

## 5. File Structure & Format Specifications

### 5.1 User Data Directory Structure

When a user runs `growth init`, the following structure is created:

```
~/growth/  (or custom path)
├── .git/                        # Git repository
├── config.yml                   # User configuration
├── skills/
│   ├── skill-001-python.md
│   ├── skill-002-machine-learning.md
│   └── skill-003-distributed-systems.md
├── goals/
│   ├── goal-001-become-ml-engineer.md
│   └── goal-002-master-system-design.md
├── paths/
│   ├── path-001-ml-engineer-fastai.md
│   └── path-002-ml-engineer-academic.md
├── resources/
│   ├── resource-001-book-hands-on-ml.md
│   ├── resource-002-course-fastai.md
│   └── resource-003-project-mnist.md
├── milestones/
│   ├── milestone-001-first-ml-model.md
│   └── milestone-002-aws-cert.md
└── progress/
    ├── 2024-12-week-50.md
    └── 2024-12-week-51.md
```

### 5.2 File Naming Convention

- Format: `{type}-{id}-{slug}.md`
- ID: Zero-padded 3-digit number
- Slug: Kebab-case, URL-safe, human-readable
- Examples:
  - `skill-001-python.md`
  - `goal-002-master-system-design.md`
  - `path-005-ml-engineer-track.md`

### 5.3 Markdown File Format

All entity files use YAML frontmatter + Markdown body.

#### 5.3.1 Goal File

```markdown
---
id: goal-001
title: Become Machine Learning Engineer
status: active
priority: high
created: 2024-12-13T10:00:00Z
updated: 2024-12-14T15:30:00Z
targetDate: 2025-12-31
learningPaths:
  - path-001
  - path-002
milestones:
  - milestone-001
  - milestone-002
  - milestone-003
tags:
  - ml
  - career-change
  - ai
---

# Goal: Become Machine Learning Engineer

## Motivation

Transition from backend/fintech engineering to AI/ML field. I see AI as the defining technology of the next decade, and I want to be at the forefront of building intelligent systems.

## Success Criteria

- [ ] Build and deploy 3 production ML models
- [ ] Contribute to a major open-source ML project (TensorFlow, PyTorch, etc.)
- [ ] Pass ML system design interviews at FAANG companies
- [ ] Land ML Engineer role at AI-first company

## Current Progress

**Overall: 15%**

- Completed: Python fundamentals refresher
- In Progress: Fast.ai Practical Deep Learning (Lesson 3/8)
- Blocked: Need stronger linear algebra foundation

## Timeline

- **Q1 2025**: Foundations (Math, Python for ML, Basic ML concepts)
- **Q2 2025**: Deep Learning & Computer Vision
- **Q3 2025**: Production ML Systems, MLOps
- **Q4 2025**: Interview prep & job search

## Notes

Started journey after attending an AI conference. Backend experience with distributed systems should translate well to ML infrastructure.
```

---

#### 5.3.2 Learning Path File

```markdown
---
id: path-001
title: ML Engineer Track (Fast.ai Approach)
type: ai-generated
status: active
created: 2024-12-13T10:00:00Z
updated: 2024-12-14T15:30:00Z
generatedBy: claude-opus-4-5
generationContext: |
  Background: 5 years backend engineering (Java, Go, microservices)
  Current skills: Strong in distributed systems, basic Python
  Learning style: Project-first, hands-on
phases:
  - phase-001
  - phase-002
  - phase-003
tags:
  - ml
  - fast-ai
  - practical
---

# Learning Path: ML Engineer (Fast.ai Approach)

*AI-generated learning path personalized for your backend engineering background*

## Path Philosophy

This path uses the **top-down** approach popularized by Fast.ai: build real models first, understand theory later. Ideal for engineers who learn by doing.

## Your Starting Point

**Strengths** (leverage these):
- Distributed systems architecture
- API design & microservices
- Git, Docker, CI/CD
- Strong debugging & problem-solving skills

**Current Level**:
- Python: Intermediate
- Math: Rusty (need refresher)
- ML: Beginner

**Gaps to Fill**:
- Linear Algebra & Calculus fundamentals
- Statistics & Probability
- Neural networks theory
- PyTorch/TensorFlow

---

## Phase 1: Foundations (2 months)

**Goal**: Build your first working models while learning Python ML stack

### Required Skills
- `skill-002`: Python (target: Intermediate)
- `skill-007`: Math for ML (target: Beginner)
- `skill-009`: NumPy & Pandas (target: Intermediate)

### Projects
1. Exploratory Data Analysis on Kaggle dataset
2. Digit recognition (MNIST)
3. Tabular data prediction (fraud detection - use your fintech background!)

### Milestone
- `milestone-005`: Deploy first ML model to production

### Estimated Time
8-10 hours/week for 8 weeks

---

## Phase 2: Deep Learning (3 months)

**Goal**: Master neural networks and computer vision

### Required Skills
- `skill-003`: Machine Learning (target: Intermediate)
- `skill-012`: Neural Networks (target: Intermediate)
- `skill-015`: PyTorch (target: Intermediate)

### Projects
1. Image classifier (cats vs dogs)
2. Object detection
3. Transfer learning application

### Milestone
- `milestone-006`: Complete Fast.ai Part 1

---

## Phase 3: Production ML (4 months)

**Goal**: Learn MLOps and deploy real systems

### Required Skills
- `skill-020`: MLOps (target: Intermediate)
- `skill-021`: Model Serving (target: Intermediate)
- `skill-008`: Distributed Systems (target: Advanced) - *leverage existing knowledge*

### Projects
1. End-to-end ML pipeline
2. Model monitoring & A/B testing
3. Scalable inference service

### Milestone
- `milestone-001`: First production ML system

---

## Alternative Paths

If this path feels too fast-paced, consider:
- `path-002`: ML Engineer (Academic Track) - theory-first approach
- `path-003`: ML Engineer (Google ML Crash Course) - structured, slower pace

## AI Revision Notes

This path was generated on 2024-12-13 based on your profile. As you make progress, run `growth path update path-001` to regenerate with latest context.
```

---

#### 5.3.3 Skill File

```markdown
---
id: skill-002
title: Python
category: programming
level: intermediate
status: learning
created: 2024-12-10T08:00:00Z
updated: 2024-12-14T12:00:00Z
resources:
  - resource-001
  - resource-002
  - resource-008
tags:
  - programming
  - ml
  - backend
---

# Skill: Python

## Current Level: Intermediate

**Self-Assessment**:
- ✅ Strong: Functions, classes, modules, standard library
- ✅ Comfortable: List comprehensions, decorators, context managers
- 🟡 Learning: Type hints, async/await, metaclasses
- ❌ Weak: NumPy, Pandas (ML-specific libraries)

## Learning Goals

**Short-term** (1 month):
- Master NumPy for array operations
- Learn Pandas for data manipulation
- Understand Jupyter notebook workflow

**Long-term** (3 months):
- Deep dive into PyTorch internals
- Contribute to Python ML libraries
- Write idiomatic, performant Python code

## Resources

### In Progress
- [ ] Book: "Fluent Python" (Chapter 5/18) → `resource-001`
- [x] Course: Real Python - Intermediate Python → `resource-002`

### Planned
- [ ] Course: Python for Data Analysis → `resource-008`

## Projects Using This Skill

1. **Fraud Detection Model** (In Progress)
   - Uses: Pandas for data preprocessing, scikit-learn for modeling
   - Status: 60% complete
   - Link: `projects/fraud-detection/`

2. **API for ML Model Serving** (Planned)
   - Uses: FastAPI, Pydantic, async Python
   - Goal: Serve PyTorch models at scale

## Notes

Coming from Go background, Python feels slower to execute but faster to prototype. Need to embrace dynamic typing while leveraging type hints for clarity.

## Referenced By

- Phase 1 in `path-001` (ML Engineer Track)
- Phase 2 in `path-003` (Backend to ML Transition)
```

---

#### 5.3.4 Resource File

```markdown
---
id: resource-001
title: "Book: Fluent Python (2nd Edition)"
type: book
skillId: skill-002
status: in-progress
created: 2024-12-10T08:00:00Z
updated: 2024-12-14T12:00:00Z
url: https://www.oreilly.com/library/view/fluent-python-2nd/9781492056348/
author: Luciano Ramalho
estimatedHours: 40
tags:
  - python
  - advanced
  - programming
---

# Resource: Fluent Python (2nd Edition)

## Overview

Comprehensive guide to writing effective, Pythonic code. Covers advanced topics like metaprogramming, concurrency, and type hints.

## Why This Resource

- **Depth**: Goes beyond basics to advanced patterns
- **Relevance**: 2nd edition covers Python 3.10+ features
- **Credibility**: Industry-standard reference

## Progress

**Current**: Chapter 5 (Data Class Builders)
**Started**: 2024-12-10
**Target Completion**: 2025-01-31

### Chapters Completed
- [x] Chapter 1: The Python Data Model
- [x] Chapter 2: An Array of Sequences
- [x] Chapter 3: Dictionaries and Sets
- [x] Chapter 4: Unicode Text vs Bytes
- [ ] Chapter 5: Data Class Builders (50% done)
- [ ] ... (13 more chapters)

## Key Takeaways

1. **Dunder methods**: `__getitem__`, `__len__` make custom classes feel native
2. **Data classes**: Use `@dataclass` instead of manual `__init__`
3. **Unicode handling**: Always decode bytes early, encode late

## Application

Applied data classes to my ML project's configuration management. Much cleaner than dictionaries!

## Rating

⭐⭐⭐⭐⭐ 5/5 - Essential for serious Python developers
```

---

#### 5.3.5 Milestone File

```markdown
---
id: milestone-001
title: First ML Model in Production
type: goal-level
referenceType: goal
referenceId: goal-001
status: pending
created: 2024-12-13T10:00:00Z
targetDate: 2025-03-31
---

# Milestone: First ML Model in Production

## Definition of Done

- [ ] Model trained and evaluated (>85% accuracy on test set)
- [ ] Deployed to cloud infrastructure (AWS/GCP)
- [ ] API endpoint serving predictions
- [ ] Monitoring and logging in place
- [ ] Handling >100 requests/second
- [ ] Documentation written

## Success Metrics

- **Performance**: Model accuracy >85%, latency <100ms
- **Reliability**: 99.9% uptime
- **Scale**: Handle production traffic
- **Quality**: Code reviewed, tests written

## Importance

This milestone represents the transition from learning to doing. It proves I can take a model from notebook to production system.

## Evidence/Proof

- GitHub repository: (to be added)
- Production URL: (to be added)
- Blog post writeup: (to be added)

## Notes

Planning to use my fintech background to build a fraud detection system. This combines domain expertise with new ML skills.
```

---

#### 5.3.6 Progress Log File

```markdown
---
id: progress-2024-w50
weekOf: 2024-12-09
hoursInvested: 12
skillsWorked:
  - skill-002  # Python
  - skill-003  # Machine Learning
  - skill-007  # Math for ML
resourcesUsed:
  - resource-001  # Fluent Python
  - resource-002  # Fast.ai Course
milestonesAchieved: []
mood: frustrated
---

# Progress Log: Week 50, 2024

## Summary

Mixed week. Made progress on Fast.ai course but hit a wall with linear algebra. Math is harder than expected.

## Accomplishments

- ✅ Completed Fast.ai Lesson 3 (Data Ethics)
- ✅ Built digit classifier with 94% accuracy (MNIST)
- ✅ Read Fluent Python Chapter 5
- ✅ Set up Jupyter environment properly

## Challenges

- ❌ Linear algebra refresher taking longer than planned
- ❌ Matrix multiplication concepts still confusing
- ❌ Struggled with backpropagation math

## Time Breakdown

- Fast.ai course: 6 hours
- Math review: 4 hours
- Python study: 2 hours

## What I Learned

**Technical**:
- Transfer learning is incredibly powerful
- Data augmentation techniques for images
- Python data classes are great for model configs

**Meta-Learning**:
- Need to slow down on math, can't rush fundamentals
- Project-first approach is motivating but creates knowledge gaps
- Should schedule dedicated math study time

## Next Week Plan

- [ ] Complete 3Blue1Brown Linear Algebra series
- [ ] Start Fast.ai Lesson 4
- [ ] Build second image classifier project
- [ ] Review matrix operations daily

## Reflections

Feeling a bit frustrated with math, but reminding myself that I'm only 3 weeks in. Backend systems took years to master; ML will too. The key is consistent progress.

## Energy Level

🔋🔋🔋 (3/5) - Moderate energy, balancing with full-time job
```

---

### 5.4 Configuration File

```yaml
# ~/growth/config.yml

version: 1.0

user:
  name: "Kostiantyn"
  email: "kostiantyn@example.com"
  timezone: "America/New_York"

ai:
  provider: "anthropic"  # Options: openai, anthropic, local
  apiKeyEnvVar: "ANTHROPIC_API_KEY"
  model: "claude-opus-4-5"
  pathGenerationPrompt: "custom"  # Optional: override default prompt

git:
  autoCommit: true
  commitOnUpdate: true
  commitMessageTemplate: "Update {entityType}: {title}"

progress:
  defaultPeriod: "week"  # Options: day, week, month
  reminderEnabled: true
  reminderSchedule: "friday 5pm"

display:
  dateFormat: "2006-01-02"
  skillLevelLabels:
    beginner: "Beginner"
    intermediate: "Intermediate"
    advanced: "Advanced"
    expert: "Expert"

mcp:
  enabled: true
  serverPort: 8765
```

---

## 6. CLI Interface Design

### 6.1 Command Structure

```bash
growth <entity> <action> [flags] [args]
```

**Entities**: `skill`, `goal`, `path`, `resource`, `milestone`, `progress`
**Actions**: `create`, `list`, `view`, `edit`, `delete`, `search`

### 6.2 Core Commands

#### 6.2.1 Initialization

```bash
# Initialize new growth repository
growth init [directory]

# Example
growth init ~/my-growth
```

**Behavior**:
- Creates directory structure
- Initializes Git repository
- Creates default `config.yml`
- Optionally sets up MCP server
- Commits initial structure

---

#### 6.2.2 Skill Commands

```bash
# Create a new skill
growth skill create <title> [flags]
growth skill create "Python" --category programming --level intermediate

# List all skills
growth skill list [flags]
growth skill list --category ml --level beginner
growth skill list --status learning

# View skill details
growth skill view <id-or-slug>
growth skill view skill-002
growth skill view python

# Edit skill
growth skill edit <id-or-slug> [flags]
growth skill edit python --level advanced --status mastered

# Delete skill (with confirmation)
growth skill delete <id-or-slug>

# Search skills
growth skill search <query>
growth skill search "distributed"
```

---

#### 6.2.3 Goal Commands

```bash
# Create a new goal
growth goal create <title> [flags]
growth goal create "Become ML Engineer" --priority high --target-date 2025-12-31

# List goals
growth goal list [flags]
growth goal list --status active
growth goal list --priority high

# View goal details (shows associated paths, milestones, progress)
growth goal view <id-or-slug>
growth goal view goal-001
growth goal view ml-engineer

# Edit goal
growth goal edit <id-or-slug> [flags]
growth goal edit ml-engineer --status completed

# Add/remove learning paths
growth goal add-path <goal-id> <path-id>
growth goal remove-path <goal-id> <path-id>
```

---

#### 6.2.4 Learning Path Commands

```bash
# Generate AI-powered learning path for a goal
growth path generate <goal-id> [flags]
growth path generate goal-001 --approach "fast.ai top-down"
growth path generate goal-001 --manual  # Create blank template

# List paths
growth path list [flags]
growth path list --type ai-generated
growth path list --status active

# View path details (shows phases, skills, resources)
growth path view <id-or-slug>
growth path view path-001

# Update/regenerate AI path based on current progress
growth path update <id-or-slug>
growth path update path-001

# Edit path manually
growth path edit <id-or-slug>
```

---

#### 6.2.5 Resource Commands

```bash
# Add a resource to a skill
growth resource create <title> --skill <skill-id> --type <type> [flags]
growth resource create "Fluent Python" --skill python --type book --url "https://..."

# List resources
growth resource list [flags]
growth resource list --skill python
growth resource list --type course --status not-started

# Mark resource as started/completed
growth resource start <id>
growth resource complete <id>

# Search resources
growth resource search <query>
growth resource search "neural networks"
```

---

#### 6.2.6 Progress Commands

```bash
# Log progress (creates/updates current week's log)
growth progress log [message]
growth progress log "Completed Fast.ai Lesson 3"
growth progress log --skill python --resource fluent-python

# Log milestone achievement
growth progress log --milestone milestone-001 "Deployed first ML model!"

# View progress logs
growth progress view [flags]
growth progress view --week
growth progress view --month
growth progress view --date 2024-12-09

# Progress statistics
growth progress stats [flags]
growth progress stats --weekly
growth progress stats --skill python
```

---

#### 6.2.7 Milestone Commands

```bash
# Create milestone
growth milestone create <title> --type <type> --ref <id> [flags]
growth milestone create "First ML Model" --type goal --ref goal-001

# List milestones
growth milestone list [flags]
growth milestone list --status pending
growth milestone list --type goal-level

# Mark milestone achieved
growth milestone achieve <id> [flags]
growth milestone achieve milestone-001 --proof "https://github.com/user/project"
```

---

#### 6.2.8 Overview Commands

```bash
# Interactive TUI dashboard
growth board

# Text-based overview
growth overview

# Statistics summary
growth stats
```

---

### 6.3 Global Flags

```bash
--config <path>      # Custom config file location
--repo <path>        # Growth repository path (default: ~/growth)
--format <format>    # Output format: table, json, yaml
--verbose, -v        # Verbose output
--help, -h           # Show help
--version            # Show version
```

---

### 6.4 Example Workflows

#### Workflow 1: Starting a New Learning Journey

```bash
# 1. Initialize growth repository
growth init ~/my-growth

# 2. Create a goal
growth goal create "Become Machine Learning Engineer" \
  --priority high \
  --target-date 2025-12-31

# 3. Generate AI-powered learning path
growth path generate goal-001

# 4. Review the generated path
growth path view path-001

# 5. Start tracking progress
growth progress log "Started ML journey today!"
```

---

#### Workflow 2: Daily Learning Session

```bash
# 1. See what to work on
growth board

# 2. Mark resource as started
growth resource start resource-002  # Fast.ai course

# 3. After session, log progress
growth progress log --skill ml --resource resource-002 \
  "Completed Lesson 3 on data ethics. Built digit classifier."

# 4. Check weekly stats
growth progress stats --week
```

---

#### Workflow 3: Achieving a Milestone

```bash
# 1. Complete milestone
growth milestone achieve milestone-001 \
  --proof "https://github.com/user/fraud-detection"

# 2. Log in progress journal
growth progress log --milestone milestone-001 \
  "Deployed first ML model! Fraud detection system live."

# 3. Update skill level
growth skill edit ml --level intermediate

# 4. Regenerate learning path based on progress
growth path update path-001
```

---

## 7. AI Integration Strategy

### 7.1 Use Cases for AI

1. **Learning Path Generation**
   - Input: Goal, current skills, background, time constraints
   - Output: Structured, phased learning path with resources

2. **Resource Recommendations**
   - Input: Skill, current level, learning style
   - Output: Ranked list of books, courses, projects

3. **Progress Analysis**
   - Input: Progress logs, skill updates
   - Output: Insights, suggestions, pattern detection

4. **Skill Gap Analysis**
   - Input: Current skills, target role/goal
   - Output: Skills to learn, prioritization

### 7.2 Path Generation Prompt Template

```
You are an expert career coach for software engineers. Generate a personalized learning path.

CONTEXT:
- Goal: {goal_title}
- Current Skills: {skill_list_with_levels}
- Background: {user_background}
- Time Availability: {hours_per_week}
- Learning Style: {learning_style_preference}

REQUIREMENTS:
- Create 3-4 phases with clear progression
- Each phase should have:
  - Title and duration estimate
  - 2-4 required skills with target levels
  - Concrete project milestones
  - Specific resource recommendations
- Leverage existing strong skills
- Address skill gaps explicitly
- Provide rationale for sequencing

OUTPUT FORMAT:
Structured markdown following the Learning Path schema.
```

### 7.3 MCP Integration

**MCP Server Capabilities**:

```json
{
  "name": "growth-md",
  "version": "1.0.0",
  "capabilities": {
    "tools": [
      {
        "name": "list_skills",
        "description": "List all skills with optional filters",
        "parameters": {
          "category": "string (optional)",
          "level": "string (optional)",
          "status": "string (optional)"
        }
      },
      {
        "name": "get_skill",
        "description": "Get detailed information about a skill",
        "parameters": {
          "id": "string (required)"
        }
      },
      {
        "name": "create_goal",
        "description": "Create a new career goal",
        "parameters": {
          "title": "string (required)",
          "priority": "string (optional)",
          "targetDate": "string (optional)"
        }
      },
      {
        "name": "generate_path",
        "description": "Generate AI-powered learning path for a goal",
        "parameters": {
          "goalId": "string (required)",
          "approach": "string (optional)"
        }
      },
      {
        "name": "log_progress",
        "description": "Log learning progress",
        "parameters": {
          "message": "string (required)",
          "skills": "array (optional)",
          "resources": "array (optional)"
        }
      },
      {
        "name": "search",
        "description": "Search across all entities",
        "parameters": {
          "query": "string (required)",
          "entityType": "string (optional)"
        }
      }
    ]
  }
}
```

**Example AI Assistant Interaction**:

```
User: "Claude, what should I learn next to become an ML engineer?"

Claude (via MCP):
1. Calls list_skills to see current skill levels
2. Calls get_skill("goal-001") to understand the goal
3. Calls get_skill for each related skill to assess gaps
4. Responds with personalized recommendation
```

---

## 8. Development Roadmap

### 8.1 Phase 1: MVP (Weeks 1-6)

**Goal**: Working CLI with core entity management and basic AI integration

**Deliverables**:
- [x] Project setup (Go modules, Cobra CLI skeleton)
- [ ] Core domain models (Skill, Goal, Path, Phase, Resource, Milestone, Progress)
- [ ] Storage layer (Markdown parsing/writing, filesystem operations)
- [ ] CLI commands:
  - [ ] `init`, `skill`, `goal`, `path`, `resource`, `progress`
  - [ ] CRUD operations for each entity
- [ ] Git integration (auto-commit on changes)
- [ ] Basic AI path generation (OpenAI integration)
- [ ] Configuration management (YAML)
- [ ] Unit tests for core logic

**Success Criteria**:
- User can create goals, skills, and resources
- AI can generate a basic learning path
- All entities persist to Markdown files
- Git tracks all changes

---

### 8.2 Phase 2: Enhanced Features (Weeks 7-10)

**Goal**: Improve UX, add TUI, enhance AI capabilities

**Deliverables**:
- [ ] TUI dashboard (`growth board`) using Bubble Tea
- [ ] Advanced AI features:
  - [ ] Path regeneration based on progress
  - [ ] Resource recommendations
  - [ ] Skill gap analysis
- [ ] Search functionality (full-text search across entities)
- [ ] Progress analytics (velocity, consistency metrics)
- [ ] Milestone tracking
- [ ] Multiple AI provider support (Anthropic, OpenAI)

**Success Criteria**:
- Beautiful, interactive dashboard
- AI provides meaningful, personalized suggestions
- Search is fast and accurate
- Analytics provide actionable insights

---

### 8.3 Phase 3: MCP & Community (Weeks 11-14)

**Goal**: Enable AI assistant integration and community sharing

**Deliverables**:
- [ ] MCP server implementation
- [ ] Claude Desktop integration guide
- [ ] Example learning paths (ML, Backend, Frontend, etc.)
- [ ] Path templates for common career transitions
- [ ] Documentation:
  - [ ] Getting Started guide
  - [ ] CLI Reference
  - [ ] AI Integration guide
  - [ ] Contributing guide

**Success Criteria**:
- Users can interact with growth.md via Claude
- Example paths provide value out of the box
- Documentation is comprehensive

---

### 8.4 Phase 4: V1.0 Launch (Weeks 15-16)

**Goal**: Polish, package, and launch

**Deliverables**:
- [ ] Cross-platform binaries (Linux, macOS, Windows)
- [ ] Installation scripts (Homebrew, apt, etc.)
- [ ] Website/landing page
- [ ] Demo video
- [ ] Launch on Product Hunt, Hacker News, Reddit

**Success Metrics**:
- 1000 GitHub stars in first month
- 100+ active users
- Positive community feedback

---

### 8.5 Future Enhancements (V2.0+)

**Potential Features** (based on user feedback):
- [ ] Web UI (browser-based interface)
- [ ] Mobile app (read-only progress tracking)
- [ ] Team features (shared learning paths, mentorship)
- [ ] Integration with GitHub profile (auto-import skills from repos)
- [ ] Job market insights (trending skills, salary data)
- [ ] Learning recommendations from job postings
- [ ] Spaced repetition system for retention
- [ ] Skill hierarchy (parent/child skills)
- [ ] Multi-language support

---

## 9. Competitive Analysis

### 9.1 Comparison Matrix

| Feature | growth.md | LinkedIn Learning | Roadmap.sh | Pluralsight | Notion Career Tracker |
|---------|-----------|-------------------|------------|-------------|-----------------------|
| **Personalized AI Paths** | ✅ Full | ❌ None | ❌ None | 🟡 Basic | ❌ None |
| **Git-Native** | ✅ Yes | ❌ No | ❌ No | ❌ No | ❌ No |
| **Local-First** | ✅ Yes | ❌ No | ❌ No | ❌ No | ❌ No |
| **Data Ownership** | ✅ Full | ❌ Platform-locked | 🟡 Partial | ❌ Platform-locked | 🟡 Partial |
| **Open Source** | ✅ MIT | ❌ Proprietary | 🟡 Content only | ❌ Proprietary | ❌ Proprietary |
| **Cost** | ✅ Free | ❌ $40/mo | ✅ Free | ❌ $29/mo | ✅ Free tier |
| **CLI Tool** | ✅ Yes | ❌ No | ❌ No | ❌ No | ❌ No |
| **MCP Integration** | ✅ Yes | ❌ No | ❌ No | ❌ No | ❌ No |
| **Progress Tracking** | ✅ Detailed | 🟡 Basic | ❌ No | ✅ Yes | ✅ Manual |
| **Community Paths** | ✅ Via Git | ❌ No | ✅ Yes | ❌ No | ❌ No |
| **Offline Access** | ✅ Full | ❌ Limited | ✅ Static | ❌ Limited | 🟡 Limited |

### 9.2 Unique Value Propositions

1. **Developer-Native Workflow**
   - Uses tools developers already know (Git, Markdown, CLI)
   - No context-switching to web platforms
   - Scriptable and automatable

2. **True Data Ownership**
   - Your career data lives in plain text files you control
   - No vendor lock-in, no platform shutdown risk
   - Forkable, shareable via standard Git workflows

3. **AI-Powered Personalization**
   - Not generic courses, but paths tailored to YOUR background
   - Adapts as you learn and grow
   - Leverages latest LLM capabilities

4. **Privacy by Design**
   - Everything local, nothing uploaded
   - Share only what you choose
   - No tracking, no analytics, no surveillance

5. **Community Without Platform**
   - Share learning paths via GitHub, no platform required
   - Fork successful paths from others
   - Contribute to resource database

---

## 10. Open Source Strategy

### 10.1 Licensing & Philosophy

**License**: MIT License

**Philosophy**:
- Core tool is 100% free and open source forever
- Community-driven development
- Transparent roadmap and decision-making
- Documentation-first approach
- No bait-and-switch: features won't move to paid tier

### 10.2 Contribution Model

**Contribution Areas**:
1. **Code**: CLI features, AI integrations, bug fixes
2. **Learning Paths**: Example paths for common career transitions
3. **Resources**: Curated resource database
4. **Documentation**: Guides, tutorials, translations
5. **Community Support**: Help users in Discussions

**Contributor Recognition**:
- Contributors list in README
- Highlight community-contributed paths
- Regular "Contributor Spotlight" posts

### 10.3 Community Building

**Channels**:
- **GitHub Discussions**: Q&A, feature requests, path sharing
- **Discord** (optional): Real-time chat, community support
- **Twitter/X**: Updates, tips, user success stories
- **Blog**: Deep dives, case studies, technical posts

**Engagement**:
- Monthly community calls
- Showcase user success stories
- Feature high-quality community paths
- Quarterly contributor awards

### 10.4 Sustainability Model

**Current (MVP)**: Fully open source, no monetization

**Future (Optional, if needed)**:
- **Cloud Sync Service** ($5/month): Optional cloud backup and sync
- **Premium AI Features** ($10/month): Access to Claude Opus, GPT-4o for path generation
- **Team Features** ($50/month): Shared paths, mentorship tools, team analytics
- **Enterprise** (Custom): On-premise deployment, SSO, compliance features

**Key Principle**: Core CLI remains free forever. Monetization only for optional value-adds.

---

## 11. Conclusion

### 11.1 Summary

growth.md represents a paradigm shift in career development tooling for software engineers. By combining:

- **Git-native design** for version control of career progression
- **Local-first architecture** for data ownership and privacy
- **AI-powered personalization** for context-aware learning paths
- **Developer-friendly workflows** through CLI-first tooling
- **Community sharing** via standard Git mechanisms

...we create a system that feels natural to engineers while providing unprecedented insight and guidance for career growth.

### 11.2 Core Innovation

The key innovation is treating **career development as code**:
- Skills are modules
- Goals are releases
- Learning paths are roadmaps
- Progress logs are commit history
- Milestones are semantic versions

This mental model resonates with engineers and leverages existing tooling (Git) rather than requiring new platforms.

### 11.3 Success Vision

**6 Months**:
- 1000+ GitHub stars
- 500+ active users
- 20+ community-contributed learning paths
- Integration with Claude, ChatGPT via MCP

**12 Months**:
- 5000+ GitHub stars
- De facto tool for engineer career tracking
- Thriving community sharing paths and resources
- Success stories: users landing dream jobs using growth.md

**18+ Months**:
- Expand beyond software engineering (data science, design, product, etc.)
- Mobile companion app
- Team/mentorship features
- Integration with job market data

### 11.4 Call to Action

growth.md is more than a tool—it's a movement toward **developer-owned career data** and **AI-powered personalized learning**.

We invite:
- **Engineers**: Use growth.md to structure your career growth
- **Contributors**: Help build the platform
- **Community**: Share your learning paths and help others
- **Companies**: Sponsor development or offer integration opportunities

---

## Appendix A: Entity Relationship Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         GOAL (Root)                         │
│  - id, title, status, priority, targetDate                 │
│  - learningPaths: [Path IDs]                               │
│  - milestones: [Milestone IDs]                             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ 1:N (unidirectional)
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                     LEARNING PATH                           │
│  - id, title, type, status                                 │
│  - phases: [Phase IDs] (ordered)                           │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ 1:N
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                        PHASE                                │
│  - id, title, order, estimatedDuration                     │
│  - requiredSkills: [Skill IDs + target levels]             │
│  - milestones: [Milestone IDs]                             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ M:N (reference)
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                    SKILL (Global)                           │
│  - id, title, category, level, status                      │
│  - resources: [Resource IDs]                               │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        │ 1:N
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                      RESOURCE                               │
│  - id, title, type, skillId, status                        │
│  - url, author, estimatedHours                             │
└─────────────────────────────────────────────────────────────┘

CROSS-CUTTING ENTITIES:

┌─────────────────────────────────────────────────────────────┐
│                      MILESTONE                              │
│  - id, title, type, referenceType, referenceId             │
│  - status, achievedDate                                    │
│  - Can reference: Goal, Path, or Skill                     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    PROGRESS LOG                             │
│  - id, weekOf, hoursInvested                               │
│  - skillsWorked: [Skill IDs]                               │
│  - resourcesUsed: [Resource IDs]                           │
│  - milestonesAchieved: [Milestone IDs]                     │
│  - notes (free-form)                                       │
└─────────────────────────────────────────────────────────────┘
```

---

## Appendix B: Technology Stack Details

### Language & Runtime
- **Go 1.21+**: Modern language features, excellent stdlib
- **Modules**: Dependency management

### Key Dependencies
- **github.com/spf13/cobra**: CLI framework
- **github.com/spf13/viper**: Configuration management
- **github.com/charmbracelet/bubbletea**: TUI framework
- **github.com/charmbracelet/lipgloss**: TUI styling
- **gopkg.in/yaml.v3**: YAML parsing
- **github.com/yuin/goldmark**: Markdown parsing
- **github.com/sashabaranov/go-openai**: OpenAI client
- **github.com/anthropics/anthropic-sdk-go**: Anthropic client

### Development Tools
- **make**: Build automation
- **golangci-lint**: Linting
- **go test**: Testing
- **goreleaser**: Multi-platform releases

---

## Appendix C: Glossary

- **Entity**: Core data model (Goal, Skill, Path, etc.)
- **Frontmatter**: YAML metadata at top of Markdown files
- **MCP**: Model Context Protocol, standard for AI tool integration
- **Phase**: Sequential step within a Learning Path
- **Proficiency Level**: Beginner, Intermediate, Advanced, Expert
- **Slug**: URL-safe, human-readable identifier (e.g., `python-programming`)
- **Unidirectional Relationship**: Reference from parent to child only, not bidirectional

---