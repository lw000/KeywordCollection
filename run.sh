#!/bin/sh

DIR="./"

typeset -l COMMAND
COMMAND=${1}

typeset -l PROCESS_NAME
PROCESS_NAME="tquery_server"

startProcess() {
    echo "server starting ..." 
	${DIR}${PROCESS_NAME} -d=true
	
	sleep 1s
    
	Pid=$(queryProcessPid ${PROCESS_NAME})
    if [[ ${Pid} ]]; then
        echo pid:[${Pid}]
		echo start success
	else
		echo "server start fail"
    fi
}

stopProcess() {
    Pid=${1}
    if [[ ${Pid} ]]; then
        for id in ${Pid}
        do
            kill -INT ${id}
        done
    fi
	
	sleep 2s
		
	Pid=$(queryProcessPid ${PROCESS_NAME})
    if [[ ${Pid} ]]; then
        echo pid:[${Pid}] stop fail
	else
		echo "stop success"
    fi
}

queryProcessPid() {
    if [[ ${1} ]]; then
        Pid=`ps -ef | grep ${1} | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
        if [[ "$Pid" ]]; then
            echo ${Pid}
        fi
    fi
}

if [[ ${COMMAND} ]]; then
    Pid=$(queryProcessPid ${PROCESS_NAME})
    if [[ ${Pid} ]]; then
        echo pid:[${Pid}]
    fi

    if [[ ${COMMAND} = "restart" ]]; then
        stopProcess ${Pid}
        sleep 1s
        startProcess
    elif [[ ${COMMAND} = "start" ]]; then
        if [[ "$Pid" ]]; then
            echo "server running"
        else
            startProcess
        fi
    elif [[ ${COMMAND} = "stop" ]]; then
        stopProcess ${Pid}
    else
        echo "please input [start/restart/stop] command"
    fi
else
    echo "please input [start/restart/stop] command"
fi
