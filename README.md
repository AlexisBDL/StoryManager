# StoryManager

Make and manage decentralized stories.
It use Noms (https://github.com/attic-labs/noms) to manage databases

```
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
```

The goal of this programme is to syncronize databases of stories in local machine.

!!!
The StorieManager need a file ".dbconfig" in the path of the executable. You can configure the programme with this file like default database and user.
Use the syntax of the example .dbconfig in this repository.
!!!

To run tests, use Robot Framework and run tests.robot

Commands :
______________________________
```
list
list [-c -o]
list [-d] <value>
list [-c -o -d] <value>
```

List ID and title of stories in my db, you can use another db with flag -d <dbTarget>. Also, you can filter closed and opened stories with flags -c for "close" and -o "open"
______________________________
```
user
```

Show current user in .dbconfig
______________________________
```
update <dbTarget>
```

Add stories that are not present in my BDD. The imported stories provide of the dbTarget
______________________________
```
log <ID>
```

Show all of the historic about commits in a story ID
______________________________
```
story create <title>
```

Create a new story with random ID and title 
______________________________
```
story edit <ID> [-t -d -e] <value>
```

Change a field value in the story ID except "Tasks" and "State". You can modify title with -t or description with -d or effort with -e
______________________________
```
story show <ID>
story show <ID> [-t]
```

Show the last state (commit) of the story ID. You can just show the tasks with -t
______________________________
```
story close <ID>
```

Change the state of the story ID with value "Close". You can't modify the story after close and it's not possible to reopen it. It is always possible to show/search-task/copy the story
______________________________
```
story Tadd <ID> <goal>
story Tadd <ID> <goal> <maker>
```

Add a task in the list "Tasks" of a story ID with goal and maker. If you don't give a maker, it will be the current maker in .dbconfig
______________________________
```
story Tdel <ID> <lastIDtask>
```

Delete a task in the list "Tasks" of a story ID. This command can now delete only the last task
______________________________
```
story Tedit <ID> <IDtask> [-g -m -s] <value>
```

Edit a task IDtask in the list "Tasks" of a story ID. You can modify the goal with -g, the maker with -m, the state with -s
______________________________
```
story Tsearch <ID> [-s -m] <value>
```

Found tasks by value of "State" with -s or "Maker" with -m in the list "Tasks" of a story ID
______________________________
```
story merge <ID1> <ID2> [-c -l -r]
```

Merge two stories that have common references (command copy duplicate for exemple). You need to resolve conflicts in CLI if you use option -c. Or force update with left -l or right -r. This command create a new ID for the merged story and the two lastes stories will be replaced by it
______________________________
```
story copy <ID> <value>
story copy <ID> [-d]
```

Copy a story in an other database or add duplicate (other ID) in my database with option -d. Value is the path of the database. Don't forget the name of the database in the path : ./your/path/Stories
______________________________
```
story sync <ID> <value>
```

Synchronize two databases (same ID) about the story ID. Value is the path of the database. Don't forget the name of the database in the path : ./home/user/Documents/Stories
The CLI ask you for conflict and finaly you and the database at value have the same story ID
______________________________

__________________________________________________________________________

If you need more informations about commands, use --help after the command
__________________________________________________________________________

Demo :

*** Create stories ***
```
./StoryManager story create test        // Exec cmd cmd title

test edited
ID : ehin6t1k6cojac1k70qkgc0830vjia0k
```
*** Edit it ***
```
./StoryManager story edit ehin6t1k6cojac1k70qkgc0830vjia0k -e 5         // Exec cmd cmd ID [-e -d -t] value

test edited
ID : ehin6t1k6cojac1k70qkgc0830vjia0k
```
*** Show it ***
```
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
```
*** List stories ***
```
/StoryManager list      // Exec cmd [-d -c -o]

ehin6t1k6cojac1k70qkgc0830vjia0k		test
```
*** Add a task ***
```
./StoryManager story Tadd ehin6t1j74s34c9m70o1ejgk5s6gu1oh newGoal          // Exec cmd cmd ID goal

Task newGoal (0) added in story test
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh
```
*** Edit a task ***
```
./StoryManager story Tedit ehin6t1j74s34c9m70o1ejgk5s6gu1oh 0 -s Complete         // Exec cmd cmd ID IDtask [-s -g -m] value

Task 0 edited in story test
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh
```
*** Search task by ***
```
./StoryManager story Tsearch ehin6t1j74s34c9m70o1ejgk5s6gu1oh -s Complete         // Exec cmd cmd ID [-s -m] value

struct Task {
  Goal: "newGoal",
  ID: 0,
  Maker: "Author config --> PO",
  State: "Complete",
}
```
*** Copy story ***
```
./StoryManager story copy ehin6t1j74s34c9m70o1ejgk5s6gu1oh ~/your/path/Usb/Stories 

Done - Synced 0 B in 0s (0 B/s)
```
If you want to check, whith the command : 
```
./StoryManager story show ~/your/path/Usb/Stories::ehin6t1j74s34c9m70o1ejgk5s6gu1oh
```

*** Merge stories ***
Create a copy, edit, merge it
```
./StoryManager story copy ehin6t1j74s34c9m70o1ejgk5s6gu1oh -d

All chunks already exist at destination! Created new dataset Stories.
Duplicate ID is : ehin6t1g70pjedpi68pka5017olje61o

./StoryManager story edit ehin6t1j74s34c9m70o1ejgk5s6gu1oh -e 4
./StoryManager story edit ehin6t1g70pjedpi68pka5017olje61o -e 5

./StoryManager story merge ehin6t1j74s34c9m70o1ejgk5s6gu1oh ehin6t1g70pjedpi68pka5017olje61o -c

Done - Synced 0 B in 0s (0 B/s)

Conflict at: .Effort
Left:  5
Right: 4

Enter 'l' to accept the left value, 'r' to accept the right value
l
Applied 1 changes...ed783dvklqs4cn8mstmjcge45qm8oatl
New ID : ehin6t1o6os34dhi6gr1ejgk5s6gu1oh
```

*** Sync stories ***
You create a story, this one will be update by an other people and he modify it. You want to synchronise your modifications with his modifications
```
./StoryManager story sync ehin6t1j74s34c9m70o1ejgk5s6gu1oh ~/your/path/Usb/Stories -c
Done - Synced 0 B in 0s (0 B/s)

Conflict at: .Effort
Left:  5
Right: 4

Enter 'l' to accept the left value, 'r' to accept the right value
l
Applied 1 changes...ed783dvklqs4cn8mstmjcge45qm8oatl
```

*** Close story ***
```
./StoryManager story close ehin6t1j74s34c9m70o1ejgk5s6gu1oh         // Exec cmd cmd ID

test closed
ID : ehin6t1j74s34c9m70o1ejgk5s6gu1oh
```
!!! Remember, you can't modify the story after close and it's not possible to reopen it
