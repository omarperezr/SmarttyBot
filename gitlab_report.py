import json
import gitlab
import sys
from datetime import datetime


gl = gitlab.Gitlab("https://gitlab.com/", private_token="")

def get_issues(project_name):
    project = gl.projects.list(search=project_name)
 
    issues_todo = [i.attributes['title'].replace("\"", "'") for i in project[0].issues.list(labels=['To Do'])]
    if not issues_todo:
        issues_todo = ["No issues todo"]
    issues_in_progress = [i.attributes['title'].replace("\"", "'") for i in project[0].issues.list(labels=['Doing'])]
    if not issues_in_progress:
        issues_in_progress = ["No issues in progress"]
    issues_done = [i.attributes['title'].replace("\"", "'") for i in project[0].issues.list(labels=['Done'])]
    if not issues_done:
        issues_done = ["No issues dode"]

    project_report = {"issues_todo": issues_todo, "issues_in_progress": issues_in_progress, "issues_done": issues_done}

    return project_report

def get_older_issue(project_name):
    project = gl.projects.list(search=project_name)

    issues_todo = [(i.attributes['title'].replace("\"", "'"), i.attributes['created_at']) for i in project[0].issues.list(labels=['To Do'])]
    issues_in_progress = [(i.attributes['title'].replace("\"", "'"), i.attributes['created_at']) for i in project[0].issues.list(labels=['Doing'])]

    if issues_todo:
        oldest_todo = min(map(lambda x: (x[0], datetime.strptime(x[1], "%Y-%m-%dT%H:%M:%S.%fZ")), issues_todo), key=lambda x: x[1])
        oldest_todo = f"{oldest_todo[0]} created at {oldest_todo[1].strftime('%Y-%m-%d %H:%M:%S')}"
    else:
        oldest_todo = "No issue todo"
    if issues_in_progress:
        oldest_in_progress = min(map(lambda x: (x[0], datetime.strptime(x[1], "%Y-%m-%dT%H:%M:%S.%fZ")), issues_in_progress), key=lambda x: x[1])
        oldest_in_progress = f"{oldest_in_progress[0]} created at {oldest_in_progress[1].strftime('%Y-%m-%d %H:%M:%S')}"
    else:
        oldest_in_progress = "No issue in progress"

    oldest_issues = {"oldest_todo": oldest_todo, "oldest_in_progress": oldest_in_progress}

    return oldest_issues

def get_status(project_name):
    project = gl.projects.list(search=project_name)

    issues_todo = len(project[0].issues.list(labels=['To Do']))
    issues_in_progress = len(project[0].issues.list(labels=['Doing']))
    issues_revision = len(project[0].issues.list(labels=['Code Revision']))
    issues_qa = len(project[0].issues.list(labels=['QA']))
    issues_deployment = len(project[0].issues.list(labels=['Deployment']))
    issues_done = len(project[0].issues.list(labels=['Done']))

    all_issues = (issues_todo + issues_in_progress + issues_revision + issues_qa + issues_deployment + issues_done)
    all_not_done = (issues_todo + issues_in_progress + issues_revision + issues_qa + issues_deployment)

    if all_issues == 0 or all_not_done == 0:
        total = 100
    else:
        total = (all_not_done/all_issues) * 100

    status = {"status": total}

    return status

def get_report(project_name):

    report = {**get_issues(project_name), **get_older_issue(project_name), **get_status(project_name)}

    print(json.dumps(report, indent=2))

def get_recent_mr(project_name):
    project = gl.projects.list(search=project_name)
    print(project.mergerequests.list())


if __name__ == "__main__":
    if sys.argv[1] == "--report":
        get_report(sys.argv[2])
    elif sys.argv[1] == "--oldest_issue":
        get_older_issue(sys.argv[2])
    elif sys.argv[1] == "--status":
        get_status(sys.argv[2])
