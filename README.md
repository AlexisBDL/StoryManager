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

list

List ID and title of stories in my db, you can use an other db with flag -d <dbTarget>. Also, you can filtre closed and opened stories

user
update
log
story create
story edit
story show
story close
story Tadd
story Tedit
story Tsearch
story merge
story copy
story sync


Demo :

if you need informations about commands, use --help after the command

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
