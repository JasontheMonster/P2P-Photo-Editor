package main

const MAX_NUM_SERVERS 	= 10
const MAX_LOG_ENTRY		= 15

const LOG_PATH = "/logs/"

const INVITE	= 0
const PUBLIC	= 1
const HEARTBEAT	= 2
const ACCEPT	= 3
// const DECLINE	= 4
const ACK		= 5
const COMMIT	= 6

const UPDATEINFO 	= 7
const UPDATEREQUEST = 8
const IMAGE 		= 9

const QUIT = 10

const TFAIL = 5 * 1000000000
const TCLEANUP = 20 * 1000000000