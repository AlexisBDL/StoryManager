*** Settings ***
Documentation	Tests of StoryManager
Library    OperatingSystem
Library    String

*** Test Cases ***
Create story
	${stdout}=	Run	./StoryManager story create hello
	Should Contain	${stdout}	was created
	${id}=	Get Substring	${stdout}	0	-12
	${stdout}=	Run	./StoryManager list
	Log	${id}
	Should Contain	${stdout}	${id}
	
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


Remove db and no more db
	${files}=	Count Directories In Directory	${CURDIR}
	Should Be Equal As Integers	${files}	1
        Remove Directory        ${CURDIR}/Stories       recursive=True
