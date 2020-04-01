// tuyue_serv project doc.go

/*
tuyue_serv document
*/
package main

// thrift -out .. -gen go example.thrift
// thrift -out .. -gen py example.thrift

// bool         	布尔值, 占1个字节
// i8            	有符号整数，占1个字节
// i16           	有符号整数，占2个字节
// i32           	有符号整数，占4个字节
// i64           	有符号整数，占8个字节
// double     		浮点数，占8个字节
// string       	字符串
// binary       	二进制大对象（python里对应的是string类型，java里对应的是ByteBuffer等）
// map<t1,t2>  	mapparam 键值对，t1,t2 代指其他类型
// list<t1>  		listparam 列表，t1代指其他类型
// set<t1>  		setparam 散列集，t1代指其他类型
