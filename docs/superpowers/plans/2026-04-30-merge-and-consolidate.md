# Merge and Consolidate Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Consolidate all working branches, verify system integrity through tests, update documentation, and synchronize with the remote repository.

**Architecture:** This is an operational consolidation task. We will move all feature work from `granite-adobe` and any other relevant branches into `main`, resolve conflicts, and verify stability.

**Tech Stack:** Git, Go (testing), Vitest (frontend testing).

---

### Task 1: Commit Current Changes on main

**Files:**
- Modify: Various files in `backend/`, `frontend/`, `db/`

- [ ] **Step 1: Check status of local changes**

Run: `git status`

- [ ] **Step 2: Add and commit local changes**

Run: `git add . && git commit -m "feat: complete Lark integration and UI enhancements"`

### Task 2: Merge Feature Branches

- [ ] **Step 1: Merge granite-adobe into main**

Run: `git merge granite-adobe`

- [ ] **Step 2: Resolve any conflicts (if any)**

- [ ] **Step 3: Check claude/fervent-hofstadter-71400d**

Run: `git log main..claude/fervent-hofstadter-71400d`
If there are unique commits, merge it: `git merge claude/fervent-hofstadter-71400d`

### Task 3: Run Backend Tests

- [ ] **Step 1: Run all backend tests**

Run: `cd backend && go test ./...`
Expected: All tests PASS.

### Task 4: Run Frontend Tests

- [ ] **Step 1: Run frontend tests**

Run: `cd frontend && npm test`
Expected: All tests PASS.

### Task 5: Update README.md

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add mention of multi-user Lark bot configuration**

Update the "飞书机器人集成" section to mention that users can now configure their own bot credentials in the system.

### Task 6: Push to GitHub

- [ ] **Step 1: Push main branch to origin**

Run: `git push origin main`

- [ ] **Step 2: Cleanup merged branches (optional but recommended)**

Run: `git branch -d granite-adobe jackrabbit-mirage`
