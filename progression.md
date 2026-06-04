### **day 1**: 23/5/2026
I added a network connection demo to see how the API might look like in go,
for now I'm just going to be reading the guide provided with the git repo to get a general understanding of the architecture of DNS
so when I start porting it I understand what its trying to do [this is the last part i stopped at](http://www.tcpipguide.com/free/t_DNSLabelsNamesandSyntaxRules.htm)

### **day 2**: 24/5/26
I started some of the work on it by adding the message type which represents the shape of the message to be sent and received 
I worked on the header for now added a type for it and functions to work with it (encode ,decode , print) 
the guide i was following was still stuck at the explaining the system point and i think i know enough to start working so I jumped to 
the section that specifies protocol details and conventions on how to store and manage DNS data [here is the link where I stopped](http://www.tcpipguide.com/free/t_DNSMessageHeaderandQuestionSectionFormat.htm)

### **day 3**: 29/5/26 
I'm getting the hang of it just implementing the spec finished the questions and at least the reading and writing seems to work
correctly I need to implement the other sections and make a successful response then think of a way to store actual records and maybe 
even load them from master files and after that a fun idea could be trying to optimize it to work well for a certain number of 
requests that I will decide at some point. Anyways [here is where I stopped](http://www.tcpipguide.com/free/t_DNSMessageResourceRecordFieldFormats-2.htm)
