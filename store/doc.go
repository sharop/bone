// Package badger  for BoneDB 2019
// SHAROP (sharop@gmail.com)
// Use of this source code is governed by a BSD-style
/*
MODEL Struct
KEYDATA
-------------------------------------------------------------------------------
|	BYTE 0 	|	BYTE 1-2	| 	BYTE[3:METALEN]	|	8 Bytes		|	8 Bytes   |
|	KeyType	|	MetaLen		| 	KeyName			|	 UID		|    Reserved |
-------------------------------------------------------------------------------

KEYINDEX


Namespace is a byte that specifies kind of key.
KEYIREVERSE

The kind of values in the database are for holding values of the record, index reference and relations.



*/
package store
