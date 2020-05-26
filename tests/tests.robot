*** Settings ***
Documentation	Tests of StoryManager commands
Library    OperatingSystem
Library    String

*** Test Cases ***
Create story
	${stdout}=	Run	./StoryManager story create test
	Should Contain	${stdout}	was created
	${id}=	Get Substring	${stdout}	0	-12
	${stdout}=	Run	./StoryManager list
	Log	${id}
	Should Contain	${stdout}	${id}
        ${stdout}=      Run     ./StoryManager story show ${id}
        Should Contain  ${stdout}       value: struct Story { 
        Should Contain  ${stdout}       Author: "Alexis Bredel --> PO",
        Should Contain  ${stdout}       Description: "",
        Should Contain  ${stdout}       Effort: 0,
        Should Contain  ${stdout}       State: "Open",
        Should Contain  ${stdout}       Title: "test",
	
Edit story effort
	${stdout}=      Run     ./StoryManager story create test
        Should Contain  ${stdout}       was created
        ${id}=  Get Substring   ${stdout}       0       -12
	${effort}=	Run	./StoryManager story show ${id}.value.Effort
	Should Be Equal	${effort}	0
        ${stdout}=      Run     ./StoryManager story edit ${id} -e 5
        ${effort}=	Run	./StoryManager story show ${id}.value.Effort
	Should Be Equal	${effort}	5

Edit story description
        ${stdout}=      Run     ./StoryManager story create test
        Should Contain  ${stdout}       was created
        ${id}=	Get Substring   ${stdout}       0       -12
        ${desc}=	Run     ./StoryManager story show ${id}.value.Description
        Should Be Equal	${desc}	""
        ${stdout}=      Run     ./StoryManager story edit ${id} -d "new description"
        ${desc}=	Run     ./StoryManager story show ${id}.value.Description
        Should Be Equal	${desc}	"new description"

Edit story title
        ${stdout}=      Run     ./StoryManager story create test
        Should Contain  ${stdout}       was created
        ${id}=	Get Substring   ${stdout}       0       -12
        ${title}=       Run     ./StoryManager story show ${id}.value.Title
        Should Be Equal         ${title}        "test"
        ${stdout}=      Run     ./StoryManager story edit ${id} -t "new title"
        ${title}=	Run     ./StoryManager story show ${id}.value.Title
        Should Be Equal	${title}        "new title"

Close story
        ${stdout}=      Run     ./StoryManager story create test
        Should Contain	${stdout}	was created
	${id}=	Get Substring	${stdout}	0	-12
        ${state}=	Run     ./StoryManager story show ${id}.value.State
        Should Be Equal         ${state}        "Open"
        ${stdout}=      Run     ./StoryManager story close ${id}
        ${state}=	Run     ./StoryManager story show ${id}.value.State
        Should Be Equal         ${state}        "Close"

User
        ${stdout}=      Run     ./StoryManager user
        Should Be Equal         ${stdout}       Alexis Bredel --> PO

Remove db and no more db
	${files}=	Count Directories In Directory	${CURDIR}
	Should Be Equal As Integers	${files}	1
        Remove Directory        ${CURDIR}/Stories       recursive=True
