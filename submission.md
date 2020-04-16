Griddy Engineering Take Home (Payment Service) Submission
=====


Open Questions:
List of unresolved implementation questions. Please only fill in (in the submission response file) if relevant.

Response:
I made the following assumptions in this project:
* I assumed the format of "amount" sent from the client was in the format of "10.00"
and was in USD. 

The following are issues that I would fix/implement if this were for production
* I would have liked to add 100% test coverage given the time
* I'd add a proper logger that integrates with AWS/GCP etc
* I'd obviously have more security (DB password, tokens, authentication on making 
sure client not requesting data for another customer etc)
* I have a bug with the MongoDB client where if passed a context, it cancels the 
context prematurely and does not continue with the db request. As a quick fix I 
passed nil instead of the context passed from the request, however it may be 
leaking connections as a result but requires further debugging to resolve.
* I'd add more validation on user input, especially the "amount" field as the
instructions required it to be a string so there are many possibilities of
what the value could be and perhaps would need to support.  

Feedback:
Let us know what you liked, or didn’t like about this experience? Was the take-home in line with your expectations? 
Would you ask that we would have done anything differently?

Response:
I felt the assessment was fair. I like how you put future concerns to think about to
discuss in subsequent rounds so I can think about it without the pressure to
respond immediately and so I can keep those future concerns in mind when I design
the software. 