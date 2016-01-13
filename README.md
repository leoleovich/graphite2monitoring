# Description

This util is getting metric **m** with timerange from **range1From** till **range1Until** and compare it with timerange from **range2From** till **range2Until**
If there is a increasing more, than **-wi** it will exit with code 1 and if more, than **-ci** - with exit code 2.
Same for decreasing with arguments **-wd** and **-cd**
Also exit message will be percentage of the difference between timeranges.
# Usage
Usage:
- -a string  
        AuthToken to access the graphite-API. For example 'qqq'  
- -ci int  
        Metrics above this threshold will be marked as critical (default 40)
- -wi int
        Metrics above this threshold will be marked as warning (default 20)
- -cd int
        Metrics below this threshold will be marked as critical (default 40)
- -wd int
        Metrics below this threshold will be marked as warning (default 20)
- -m string  
        Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems

