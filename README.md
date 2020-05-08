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
