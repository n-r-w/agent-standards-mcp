---
description: Standards for applying Test-Driven Development (TDD) in software projects.
---

# 🚨 Critical Test-Driven Development (TDD) Applicability Algorithm 🚨
MUST follow BEFORE applying TDD strategy:
1. Identify the target change and its scope (files, modules, behavior)
2. Identify TDD applicability:
   - MUST NOT use TDD for:
      * Project research
      * Initialization, Project scaffolding
      * Interface definitions
      * Refactoring
      * Documentation, README.md, Comments, Formatting, Renaming
      * Folder structure changes      
      * Configuration
      * Benchmarking
      * Bug fixes, Patches
      * Prototyping or spike solutions
      * Other tasks that do NOT change system behavior
   - Otherwise → Proceed with TDD

# Test-Driven Development (TDD) Strategy Selection Algorithm

**Definitions**
- Test-Driven Development (TDD): An iterative practice where tests define required behavior; code is written only to make failing tests pass, then refactored while tests remain green.
- TDD Group: A cohesive bundle of tasks for one change (tests plus implementation) executed end-to-end without mixing unrelated features.
- New Functionality Scenario: Adding behavior that does not yet exist (new feature, module, or branch of logic).
- Existing Logic Modification Scenario: Changing behavior of already implemented logic (updated requirements, behavior changes).
- Mixed Scenario: A change that contains both new functionality and modifications to existing logic; must be decomposed.
- Stub: A minimal placeholder (signature, struct) created, so tests compile before implementation exists. Follow rules:
   * Stubs == minimal code to allow compilation and test writing
   * Stubs != project scaffolding
   * Stubs != file skeletons
   * Stubs != interface definitions. MUST NOT create tasks like "Create stubs for XXX interface", interface definition NEVER requires TDD
- Verification: A final check at the end of a TDD Group that confirms behavior against requirements.

**Strategy Decision Algorithm (if TDD applicable)**
1) Classify the change
   - If the target introduces behavior that does not exist → choose New Functionality Scenario.
   - If the target modifies behavior that already exists → choose Existing Logic Modification Scenario.
   - If both apply → choose Mixed Scenario and split work into separate TDD Groups by change type.

2) Create TDD Group(s)
   - ALWAYS group tasks for one cohesive change into a single TDD Group.
   - NEVER mix unrelated features across TDD Groups.
   - ALWAYS end each TDD Group with a Verification task; do not merge verifications across groups.

3) Execute per-scenario cycle for each TDD Group
   A) New Functionality Scenario (Red-Green-Refactor)
      a. RED
         - Create stubs (signatures/structs/interfaces) so tests can compile.
         - Write complete tests that describe desired behavior.
         - Ensure tests compile and FAIL at runtime (no compilation errors).
      b. GREEN
         - Implement the simplest code that makes tests pass.
         - Prioritize speed and correctness of passing tests over design polish.
      c. REFACTOR
         - Clean code (remove duplication, improve naming/structure) while keeping tests green.
         - Remove TODOs related to this change.
      d. VERIFY
         - Verify implementation against stated requirements.

   B) Existing Logic Modification Scenario (Analyze-Modify-Implement)
      a. ANALYZE
         - Locate and review related existing tests and patterns for consistency.
      b. MODIFY
         - Update or add tests to reflect new requirements.
         - Add or update stubs so tests compile.
         - Ensure tests compile and FAIL due to missing implementation (not compilation errors).
      c. IMPLEMENT
         - Change the code to pass the updated tests.
      d. REFACTOR
         - Improve the code while tests remain green.
      e. VERIFY
         - Verify implementation against updated requirements.

   C) Mixed Scenario
      - Decompose the change into separate TDD Groups:
        • For parts that introduce new behavior → apply New Functionality Scenario.
        • For parts that modify existing behavior → apply Existing Logic Modification Scenario.
      - Execute groups independently; do not interleave tasks across groups.

**TDD Ordering**
- MUST create tests creation/modification steps.
- MUST put tests creation/modification steps BEFORE implementation steps
- MUST include "Create stubs..." BEFORE "Write failing tests..." within a TDD Group.
- MUST include "Verify implementation..." AFTER "Refactor..." within a TDD Group.
- Tests MUST compile during RED/MODIFY phases and MUST fail before implementation.
- Do NOT mix features across TDD Groups; finish one group before starting another.

**Task Naming Patterns (examples)**
- Create stubs for {functionality} in {file}
- Write failing tests for {functionality} in {test_file}
- Implement {functionality} to pass tests in {file}
- Refactor {functionality} while keeping tests green in {file}
- Verify {functionality} against {requirements} in {context}

**Worked Example** ("Remember me" login option. Mixed Scenario)
- Situation: Add a "Remember me" checkbox; if checked, issue 30-day cookie; change default session expiry (unchecked) from 24h to 8h.
- Applicability Gate: Not a bug/patch → apply Test-Driven Development (TDD).
- Classification: Mixed Scenario (new behavior + modification of existing behavior).
- Decomposition: Create two Test-Driven Development (TDD) Groups by change type; run Group A then Group B.

Group A — New Functionality (30-day cookie)
- Create stubs for remember_me support in auth/session module
- Write failing tests for 30-day cookie in auth tests
- Implement remember_me cookie issuance to pass tests
- Refactor auth/session code while keeping tests green
- Verify behavior against requirements (30-day cookie when checked)

Group B — Existing Logic Modification (8h default expiry)
- Analyze existing tests covering session expiry in auth tests
- Modify tests to require 8h expiry when remember_me is unchecked
- Implement updated expiry logic to pass tests
- Refactor related code while keeping tests green
- Verify behavior against updated requirements (8h default expiry)
