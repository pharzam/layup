# F-0002. Second sample fact (fixture)

A second fact file, so that the facts list handed to awk holds two names. The
PRD linter once passed that list newline-separated, which macOS awk rejects
("newline in string"); with one fact the defect was invisible. Kept as the
regression fixture of task T-wjq4.

1. Sample fact.
