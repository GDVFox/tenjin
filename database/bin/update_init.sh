#!/bin/bash

tables=(
    initial.sql
    department.sql
    appointement.sql
    person.sql
    employee.sql
    employee_post.sql
    comment.sql
    task.sql
    solution.sql
    vote.sql
    attachment.sql
    skill.sql
    permission.sql
    vacancy.sql
    interview.sql
    task_skill_requirement.sql
    vacancy_skill_requirement.sql
    requirement_check.sql
    has_permission.sql
    works_in.sql
)

# clear file content
> init/init.sql

for t in ${tables[*]}; do (cat "${t}"; echo; echo;) >> init/init.sql; done