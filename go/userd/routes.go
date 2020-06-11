package main

import (
	commentcreate "github.com/GDVFox/tenjin/userd/comment/create"
	commentdelete "github.com/GDVFox/tenjin/userd/comment/delete"
	commentupdate "github.com/GDVFox/tenjin/userd/comment/update"
	interviewcreate "github.com/GDVFox/tenjin/userd/interview/create"
	interviewdelete "github.com/GDVFox/tenjin/userd/interview/delete"
	interviewread "github.com/GDVFox/tenjin/userd/interview/read"
	interviewupdatecommon "github.com/GDVFox/tenjin/userd/interview/update/common"
	interviewupdatecomplete "github.com/GDVFox/tenjin/userd/interview/update/complete"
	personcreate "github.com/GDVFox/tenjin/userd/person/create"
	persondelete "github.com/GDVFox/tenjin/userd/person/delete"
	personreadbase "github.com/GDVFox/tenjin/userd/person/read/base"
	personreadfull "github.com/GDVFox/tenjin/userd/person/read/full"
	personupdate "github.com/GDVFox/tenjin/userd/person/update"
	solutioncreate "github.com/GDVFox/tenjin/userd/solution/create"
	solutiondelete "github.com/GDVFox/tenjin/userd/solution/delete"
	solutionupdate "github.com/GDVFox/tenjin/userd/solution/update"
	taskcreate "github.com/GDVFox/tenjin/userd/task/create"
	taskdelete "github.com/GDVFox/tenjin/userd/task/delete"
	taskreadbase "github.com/GDVFox/tenjin/userd/task/read/base"
	taskreadfull "github.com/GDVFox/tenjin/userd/task/read/full"
	taskupdate "github.com/GDVFox/tenjin/userd/task/update"
	vacancycreate "github.com/GDVFox/tenjin/userd/vacancy/create"
	vacancydelete "github.com/GDVFox/tenjin/userd/vacancy/delete"
	vacancyreadbase "github.com/GDVFox/tenjin/userd/vacancy/read/base"
	vacancyreadfull "github.com/GDVFox/tenjin/userd/vacancy/read/full"
	vacancyupdate "github.com/GDVFox/tenjin/userd/vacancy/update"
	votecreate "github.com/GDVFox/tenjin/userd/vote/create"
	workappointmentcreate "github.com/GDVFox/tenjin/userd/work/appointment/create"
	workappointmentread "github.com/GDVFox/tenjin/userd/work/appointment/read"
	workdepartmentcreate "github.com/GDVFox/tenjin/userd/work/department/create"
	workdepartmentread "github.com/GDVFox/tenjin/userd/work/department/read"
	"github.com/GDVFox/tenjin/utils/server"
)

func routes() []*server.Route {
	return []*server.Route{
		{
			Method:     "POST",
			Pattern:    "/person.json",
			Handler:    server.RequestProcessor(personcreate.Process),
			ArgFactory: personcreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/person.json",
			Handler:    server.RequestProcessor(personreadbase.Process),
			ArgFactory: personreadbase.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/person/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(personreadfull.Process),
			ArgFactory: personreadfull.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/person/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(personupdate.Process),
			ArgFactory: personupdate.NewArguments,
		},
		{
			Method:     "DELETE",
			Pattern:    "/person/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(persondelete.Process),
			ArgFactory: persondelete.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/department.json",
			Handler:    server.RequestProcessor(workdepartmentcreate.Process),
			ArgFactory: workdepartmentcreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/department.json",
			Handler:    server.RequestProcessor(workdepartmentread.Process),
			ArgFactory: nil,
		},
		{
			Method:     "POST",
			Pattern:    "/appointment.json",
			Handler:    server.RequestProcessor(workappointmentcreate.Process),
			ArgFactory: workappointmentcreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/appointment.json",
			Handler:    server.RequestProcessor(workappointmentread.Process),
			ArgFactory: nil,
		},
		{
			Method:     "POST",
			Pattern:    "/task.json",
			Handler:    server.RequestProcessor(taskcreate.Process),
			ArgFactory: taskcreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/task.json",
			Handler:    server.RequestProcessor(taskreadbase.Process),
			ArgFactory: taskreadbase.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/task/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(taskreadfull.Process),
			ArgFactory: taskreadfull.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/task/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(taskupdate.Process),
			ArgFactory: taskupdate.NewArguments,
		},
		{
			Method:     "DELETE",
			Pattern:    "/task/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(taskdelete.Process),
			ArgFactory: taskdelete.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/task/{task_id:[0-9]+}/solve.json",
			Handler:    server.RequestProcessor(solutioncreate.Process),
			ArgFactory: solutioncreate.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/solution/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(solutionupdate.Process),
			ArgFactory: solutionupdate.NewArguments,
		},
		{
			Method:     "DELETE",
			Pattern:    "/solution/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(solutiondelete.Process),
			ArgFactory: solutiondelete.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/task/{post_id:[0-9]+}/comment.json",
			Handler:    server.RequestProcessor(commentcreate.Process),
			ArgFactory: commentcreate.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/solution/{post_id:[0-9]+}/comment.json",
			Handler:    server.RequestProcessor(commentcreate.Process),
			ArgFactory: commentcreate.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/comment/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(commentupdate.Process),
			ArgFactory: commentupdate.NewArguments,
		},
		{
			Method:     "DELETE",
			Pattern:    "/comment/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(commentdelete.Process),
			ArgFactory: commentdelete.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/task/{post_id:[0-9]+}/vote.json",
			Handler:    server.RequestProcessor(votecreate.Process),
			ArgFactory: votecreate.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/solution/{post_id:[0-9]+}/vote.json",
			Handler:    server.RequestProcessor(votecreate.Process),
			ArgFactory: votecreate.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/comment/{comment_id:[0-9]+}/vote.json",
			Handler:    server.RequestProcessor(votecreate.Process),
			ArgFactory: votecreate.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/vacancy.json",
			Handler:    server.RequestProcessor(vacancycreate.Process),
			ArgFactory: vacancycreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/vacancy.json",
			Handler:    server.RequestProcessor(vacancyreadbase.Process),
			ArgFactory: vacancyreadbase.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/vacancy/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(vacancyreadfull.Process),
			ArgFactory: vacancyreadfull.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/vacancy/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(vacancyupdate.Process),
			ArgFactory: vacancyupdate.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/vacancy/{id:[0-9]+}/close.json",
			Handler:    server.RequestProcessor(vacancydelete.Process),
			ArgFactory: vacancydelete.NewArguments,
		},
		{
			Method:     "POST",
			Pattern:    "/interview.json",
			Handler:    server.RequestProcessor(interviewcreate.Process),
			ArgFactory: interviewcreate.NewArguments,
		},
		{
			Method:     "GET",
			Pattern:    "/interview/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(interviewread.Process),
			ArgFactory: interviewread.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/interview/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(interviewupdatecommon.Process),
			ArgFactory: interviewupdatecommon.NewArguments,
		},
		{
			Method:     "PUT",
			Pattern:    "/interview/{id:[0-9]+}/complete.json",
			Handler:    server.RequestProcessor(interviewupdatecomplete.Process),
			ArgFactory: interviewupdatecomplete.NewArguments,
		},
		{
			Method:     "DELETE",
			Pattern:    "/interview/{id:[0-9]+}.json",
			Handler:    server.RequestProcessor(interviewdelete.Process),
			ArgFactory: interviewdelete.NewArguments,
		},
	}
}
