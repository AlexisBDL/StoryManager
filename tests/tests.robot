*** Settings ***
Documentation	Tests of StoryManager commands
Library    OperatingSystem
Library    String
Library    Process

*** Test Cases ***
Create story
	${stdout}=	Run	./StoryManager story create test
        ${title}=       Get Line        ${stdout}       0
	${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        Should Contain  ${stdout}   Story test created
        ${stdout}=      Run     ./StoryManager story show ${id}
        Should Contain  ${stdout}       value: struct Story { 
        Should Contain  ${stdout}       Author: "Alexis Bredel --> PO",
        Should Contain  ${stdout}       Description: "",
        Should Contain  ${stdout}       Effort: 0,
        Should Contain  ${stdout}       State: "Open",
        Should Contain  ${stdout}       Title: "test",
        Should Contain  ${stdout}       Tasks: [],
	
Edit story effort
	${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
	${effort}=	Run	./StoryManager story show ${id}.value.Effort
	Should Be Equal	${effort}	0
        ${stdout}=      Run     ./StoryManager story edit ${id} -e 5
        ${effort}=	Run	./StoryManager story show ${id}.value.Effort
	Should Be Equal	${effort}	5

Edit story description
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        ${desc}=	Run     ./StoryManager story show ${id}.value.Description
        Should Be Equal	${desc}	""
        ${stdout}=      Run     ./StoryManager story edit ${id} -d "new description"
        ${desc}=	Run     ./StoryManager story show ${id}.value.Description
        Should Be Equal	${desc}	"new description"

Edit story title
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        ${title}=       Run     ./StoryManager story show ${id}.value.Title
        Should Be Equal         ${title}        "test"
        ${stdout}=      Run     ./StoryManager story edit ${id} -t "new title"
        ${title}=	Run     ./StoryManager story show ${id}.value.Title
        Should Be Equal	${title}        "new title"

Close story
        ${stdout}=      Run     ./StoryManager story create test
	${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        ${state}=	Run     ./StoryManager story show ${id}.value.State
        Should Be Equal         ${state}        "Open"
        ${stdout}=      Run     ./StoryManager story close ${id}
        ${state}=	Run     ./StoryManager story show ${id}.value.State
        Should Be Equal         ${state}        "Close"
        ${stdout}=      Run     ./StoryManager story edit ${id} -e 4
        Should Contain         ${stdout}       The story ${id} is close

User
        ${stdout}=      Run     ./StoryManager user
        Should Be Equal         ${stdout}       Alexis Bredel --> PO

Add task
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
	${stdout}=	Run	./StoryManager story Tadd ${id} task
        Should Contain  ${stdout}       Task task (0) added in story
        Should Contain  ${stdout}       ID : ${id}
        ${stdout}=      Run     ./StoryManager story show ${id} -t
        Should Contain  ${stdout}       Goal: "task",
        Should Contain  ${stdout}       Maker: "Alexis Bredel --> PO",
        Should Contain  ${stdout}       State: "",

Edit task
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
	${stdout}=	Run	./StoryManager story Tadd ${id} task
        ${stdout}=	Run	./StoryManager story Tedit ${id} 0 -s Complete
        Should Contain  ${stdout}       Task 0 edited in story test
        Should Contain  ${stdout}       ID : ${id}
        ${stdout}=      Run     ./StoryManager story show ${id} -t
        Should Contain  ${stdout}       Goal: "task",
        Should Contain  ${stdout}       ID: 0,
        Should Contain  ${stdout}       Maker: "Alexis Bredel --> PO",
        Should Contain  ${stdout}       State: "Complete",

Search task
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
	${stdout}=	Run	./StoryManager story Tadd ${id} task
        ${stdout}=	Run	./StoryManager story Tadd ${id} task2
        ${stdout}=	Run	./StoryManager story Tedit ${id} 1 -s Complete
        ${stdout}=	Run	./StoryManager story Tsearch ${id} -s Complete
        Should Contain  ${stdout}       Goal: "task2",
        Should Contain  ${stdout}       ID: 1,
        Should Contain  ${stdout}       Maker: "Alexis Bredel --> PO",
        Should Contain  ${stdout}       State: "Complete",
        Should Not Contain       ${stdout}       ID: 0,

List story
        ${stdout}=      Run     ./StoryManager story create test
        ${id1}=  Get Line	${stdout}	1
        ${id1}=  Get Substring   ${id1}   5
        ${stdout}=      Run     ./StoryManager story create test
        ${id2}=  Get Line	${stdout}	1
        ${id2}=  Get Substring   ${id2}   5
        ${stdout}=      Run     ./StoryManager story create test
        ${id3}=  Get Line	${stdout}	1
        ${id3}=  Get Substring   ${id3}   5
        ${stdout}=      Run     ./StoryManager list
        Should Contain  ${stdout}       ${id1}  test
        Should Contain  ${stdout}       ${id2}  test
        Should Contain  ${stdout}       ${id3}  test

Copy story
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        ${stdout}=      Run     ./StoryManager story copy ${id} Stories -d
        ${idd}=  Get Line	${stdout}	1
        ${idd}=  Get Substring   ${idd}   18
        ${stdout}=      Run     ./StoryManager list
        Should Contain  ${stdout}       ${id}   test
        Should Contain  ${stdout}       ${idd}   test
        ${stdout}=      Run     ./StoryManager story show ${id}.value
        ${stdoutd}=     Run     ./StoryManager story show ${idd}.value
        Should Be Equal         ${stdout}       ${stdoutd}

Merge story
        ${stdout}=      Run     ./StoryManager story create test
        ${id}=  Get Line	${stdout}	1
        ${id}=  Get Substring   ${id}   5
        ${stdout}=      Run     ./StoryManager story copy ${id} Stories -d
        ${idd}=  Get Line	${stdout}	1
        ${idd}=  Get Substring   ${idd}   18
        ${stdout}=      Run     ./StoryManager story edit ${id} -e 4
        ${stdout}=      Run     ./StoryManager story edit ${idd} -e 5
        ${stdout}=      Run     ./StoryManager story merge ${id} ${idd} -l
        ${idm}=  Get Line	${stdout}	5
        ${idm}=  Get Substring   ${idm}   9
        ${stdout}=      Run     ./StoryManager list
        Should Contain  ${stdout}       ${idm}   test
        Should Not Contain      ${stdout}       ${id}
        Should Not Contain      ${stdout}       ${idd}

Remove db and no more db
        ${files}=	Count Directories In Directory	${CURDIR}
	Should Be Equal As Integers	${files}	1
        Remove Directory        ${CURDIR}/Stories       recursive=True

        