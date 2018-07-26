# gopi-finder

Gopi-Finder is a peer to peer microservice location tool written in go.

## Server Installation

Download the release file for your operating system from https://github.com/Brumawen/gopi-finder/releases 

Extract the 2 files to a folder.

To install the server, run the following from the command line

        finderserver -server install
        finderserver -server run

This will install and run the server as a background service on your machine.

## Service Discovery

The finderclient program is used to search the network for any machine running the server software and will return the Name and IP address of each server found.  

        $ .\finderclient

        machineA  [192.168.1.10]
        machineB  [192.168.1.11]

To get a list of services running on a machine, run the following on the commandline with the machine IP address

        $ .\finderclient -services -ip 192.168.1.10

This will return a list of strings in the format {machine name}  {service name}  {ip address}  {service port number}

        machineA  Service1  192.168.1.10  12345

