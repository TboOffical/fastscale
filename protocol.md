## FS Protocol
Operates over **TCP**

### Steps for forming a connection
1. Client connects to server via TCP
2. Client sends an "OK" to signify that it is ready to receive data
3. Server sends "KEY?" to request the client's public key
4. Client sends its public key (via CRYSTALS-Kyber encryption)
5. Server sends "OK"
6. Client sends "KEY?" to request the server's public key
7. Server sends its public key (via CRYSTALS-Kyber encryption)
8. One last header is sent in JSON format to provide information about the type of terminal that is being used
   ```json
   {
       "device": "terminal | headless",
       "features": ["termCommands", "guiForms", "fileTransfers", "editors"]
   }
   ```

Data can be exchanged after this point

### Special terminal features

#### termCommands
termCommands allow the server to ajust the terminal's settings remotely. 
Here are some examples of term commands

*the fastscale terminal dosen't allow users to type the escape characters used by term commands, and the 
characters below should not be used in data sent back to terminals as it might mess them up*
>`{{ command, arg1, arg2 }}`

here are some term commands
- >`{{ clear }}` clear the terminal
- >`{{ notification, "task finished" }} send a notification`