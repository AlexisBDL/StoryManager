# StoryManager

Make and manage decentralized stories
It use Noms (https://github.com/attic-labs/noms) to manage databases

Model of datas :
    struct Story {
        Title string        // init ""
	    Description string  // init ""
	    Effort int          // init 0
	    State string        // Open or Close
    }

The goal of this programme is to syncronize databases of stories in local machine.

The StorieManager need a file ".dbconfig" in the path of the executable. You can configure the programme with this file like default database and user.
Use the syntax of the example .dbconfig in this repository.