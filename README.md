# Description

This util is getting **metric** with timerange from **range1From** till **range1Until** and compare it with timerange from **range2From** till **range2Until**
If there is a difference (increasing/decreasing) more, than **-w** it will exit with code 1 and if more, than **-c** - with exit code 2.
Also exit message will be percentage of the difference between timeranges.
# Usage
Usage of graphite2nagios:  
- -a string  
        AuthToken to access the graphite-API. For example 'qqq'  
- -c int  
        Metrics above the threshold will be marked as critical (default 40)  
-  -metric string  
        Name of metric or metric filter e.g. Character.* (default "qqqq.test.leoleovich.currentProblems")  
- -range1From string  
        e.g. 2015-10-23 10:00 (default "2016-01-04 13:28")  
- -range1Until string  
        e.g. 2015-10-23 11:00 (default "2016-01-04 14:28")  
- -range2From string  
        e.g. 2015-10-23 10:00 (default "2016-01-05 13:28")  
- -range2Until string  
        e.g. 2015-10-23 11:00 (default "2016-01-05 14:28")  
- -w int  
        Metrics above the threshold will be marked as warning (default 20)