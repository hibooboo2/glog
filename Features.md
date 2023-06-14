# Write to multiple writers so that you can have logs go to different places
# Writer that sends logs to different sources based on level
# Request logger middleware 
# Add structured logging 
# Add flags for loggers similar to std logger
# Try amd make it 0 allocation

# Add http handler to get logs for a specific logger 
Make a http middleware that creates a logger for the given request.
The logger can produce a link that goes to the logs for that request
Rotate out saved logs for requests so that memory is not over used
Store requests in memory for some cache time with a max for the logs before it is saved to disk
Store logs with max memory time rotate as size gets to large

# Add handler to change log level for requests
Add ability to set log level per endpoint maybe even per request using a header

# add structured logging
Eg with fields 

See http://github.com/ssgreg/logf
