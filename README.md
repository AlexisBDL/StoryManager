# StoryManager

Make and manage decentralized stories
It use Noms (https://github.com/attic-labs/noms) to manage databases

Model of datas :
    struct Story {
      Title string        // init title command
	    Description string  // init ""
	    Effort int          // init 0
	    State string        // Open or Close
        Tasks               // list of tasks
            Goal  string    // init goal command
            Maker string    // init current user
            State string    // init ""
    }

The goal of this programme is to syncronize databases of stories in local machine.

!!!
The StorieManager need a file ".dbconfig" in the path of the executable. You can configure the programme with this file like default database and user.
Use the syntax of the example .dbconfig in this repository.
!!!

To run tests, use Robot Framework and run tests.robot

Commands :
______________________________
list

List ID and title of stories in my db, you can use an other db with flag -d <dbTarget>. Also, you can filtre closed and opened stories
______________________________
user

Show current user in .dbconfig
______________________________
update

Add stories that are not present in my BDD. The imported stories provide of the dbTarget
______________________________
log

Show all of the historic about commits in a story
______________________________
story create

Create a new story with random ID
______________________________
story edit

Change a field value in a story except "Tasks" and "State"
______________________________
story show

Show the last state (commit) of a story
______________________________
story close

Change the state of a story with value "Close"
______________________________
story Tadd

Add a task in the list "Tasks" of a story
______________________________
story Tedit

Edit a task in the list "Tasks" of a story
______________________________
story Tsearch

Found tasks by value of "State" or "Maker" in the list "Tasks" of a story
______________________________
story merge

Merge two stories that have common references
______________________________
story copy

Copy a story in an other database or add duplicate (other ID) in my database with option -d
______________________________
story sync



__________________________________________________________________________

If you need more informations about commands, use --help after the command
__________________________________________________________________________

Demo :

*** Create stories ***

./StoryManager story create test        // Exec cmd cmd title

test edited
ID : ehin6t1k6cojac1k70qkgc0830vjia0k

*** Edit it ***

./StoryManager story edit ehin6t1k6cojac1k70qkgc0830vjia0k -e 5         // Exec cmd cmd ID [-e -d -t] value

test edited
ID : ehin6t1k6cojac1k70qkgc0830vjia0k

*** Show it ***

./StoryManager story show ehin6t1k6cojac1k70qkgc0830vjia0k              // Exec cmd cmd ID [-t]

struct Commit {
  meta: struct Meta {
    author: "Author config --> PO",
    date: "2020-06-22T10:08:37Z",
    message: "Edit value effort in story test with ID ehin6t1k6cojac1k70qkgc0830vjia0k",
  },
  parents: set {
    #ideku8d93ekku555tjkf79ghi4gsavj2,
  },
  value: struct Story {
    Author: "Author config --> PO",
    Description: "",
    Effort: 5,
    State: "Open",
    Tasks: [],
    Title: "test",
  },
}

*** List stories ***

/StoryManager list      // Exec cmd [-d -c -o]

ehin6t1k6cojac1k70qkgc0830vjia0k		test

*** Add a task ***

./StoryManager story Tadd ehin6t1j74s34c9m70o1ejgk5s6gu1oh newGoal          // Exec cmd cmd ID goal

Task newGoal (0) added in story test
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh

*** Edit a task ***

./StoryManager story Tedit ehin6t1j74s34c9m70o1ejgk5s6gu1oh 0 -s Complete         // Exec cmd cmd ID IDtask [-s -g -m] value

Task 0 edited in story test
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh

*** Search task by ***

./StoryManager story Tsearch ehin6t1j74s34c9m70o1ejgk5s6gu1oh -s Complete         // Exec cmd cmd ID [-s -m] value

struct Task {
  Goal: "newGoal",
  ID: 0,
  Maker: "Alexis Bredel --> PO",
  State: "Complete",
}

*** Close story ***

./StoryManager story close ehin6t1j74s34c9m70o1ejgk5s6gu1oh         // Exec cmd cmd ID

test closed
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh

!!! You can't modify the story after close and it's not possible to reopen it