// You can query JSON data and store large amounts of data, but you don't support transactions.

// Redis data is stored in memory and is written to disk regularly.
// When there is insufficient memory, you can select the specified LRU algorithm to delete the data.
// All the data of mongodb is actually stored on the hard disk,
// and all the data to be operated on is mapped to an area of memory by mmap.

// If the table structure changes frequently, you don't need to store complex data structures,
// and you need to store document-type data in real time, and you need to extends The amount of data aways, you need mongoDB.\

// Bson(Binary Json) is Json's extends data format for mongoDB.

package mongo

import "gopkg.in/mgo.v2/bson"
